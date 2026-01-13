package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysMenuModel = (*customSysMenuModel)(nil)

type (
	// SysMenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysMenuModel.
	SysMenuModel interface {
		sysMenuModel
		withSession(session sqlx.Session) SysMenuModel
		// SelectMenuTreeAll 查询所有菜单（超级管理员）
		SelectMenuTreeAll(ctx context.Context) ([]*SysMenu, error)
		// SelectMenuListByUserId 根据用户ID查询菜单列表（非超级管理员）
		SelectMenuListByUserId(ctx context.Context, userId int64) ([]*SysMenu, error)
		// SelectMenuPermissionsByUserId 根据用户ID查询菜单权限
		SelectMenuPermissionsByUserId(ctx context.Context, userId int64) ([]string, error)
		// FindAll 查询菜单列表（根据条件）
		FindAll(ctx context.Context, query *MenuQuery, userId int64) ([]*SysMenu, error)
		// CheckMenuNameUnique 检查菜单名称唯一性（同父菜单下唯一）
		CheckMenuNameUnique(ctx context.Context, menuName string, parentId int64, excludeMenuId int64) (bool, error)
		// HasChildByMenuId 是否存在子菜单
		HasChildByMenuId(ctx context.Context, menuId int64) (bool, error)
		// HasChildByMenuIds 是否存在子菜单（批量）
		HasChildByMenuIds(ctx context.Context, menuIds []int64) (bool, error)
		// CheckMenuExistRole 检查菜单是否分配给角色
		CheckMenuExistRole(ctx context.Context, menuId int64) (bool, error)
		// SelectMenuListByRoleId 根据角色ID查询菜单ID列表
		SelectMenuListByRoleId(ctx context.Context, roleId int64) ([]int64, error)
		// UpdateById 根据ID更新菜单，只更新非零值字段
		UpdateById(ctx context.Context, data *SysMenu) error
	}

	// MenuQuery 菜单查询条件
	MenuQuery struct {
		MenuName string // 菜单名称（模糊查询）
		Visible  string // 显示状态（0显示 1隐藏）
		Status   string // 菜单状态（0正常 1停用）
		MenuType string // 菜单类型（M目录 C菜单 F按钮）
		ParentId int64  // 父菜单ID
	}

	customSysMenuModel struct {
		*defaultSysMenuModel
	}
)

// NewSysMenuModel returns a model for the database table.
func NewSysMenuModel(conn sqlx.SqlConn) SysMenuModel {
	return &customSysMenuModel{
		defaultSysMenuModel: newSysMenuModel(conn),
	}
}

func (m *customSysMenuModel) withSession(session sqlx.Session) SysMenuModel {
	return NewSysMenuModel(sqlx.NewSqlConnFromSession(session))
}

// SelectMenuTreeAll 查询所有菜单（超级管理员）
func (m *customSysMenuModel) SelectMenuTreeAll(ctx context.Context) ([]*SysMenu, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE menu_type IN ('M', 'C')
		  AND status = '0'
		ORDER BY parent_id ASC, order_num ASC
	`, sysMenuRows, m.table)

	var menus []*SysMenu
	err := m.conn.QueryRowsPartialCtx(ctx, &menus, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return menus, nil
}

// SelectMenuListByUserId 根据用户ID查询菜单列表（非超级管理员）
func (m *customSysMenuModel) SelectMenuListByUserId(ctx context.Context, userId int64) ([]*SysMenu, error) {
	// 构建字段列表，添加 m. 前缀
	fields := strings.Split(sysMenuRows, ",")
	fieldList := make([]string, 0, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(field)
		field = strings.Trim(field, "`")
		fieldList = append(fieldList, fmt.Sprintf("m.`%s`", field))
	}
	selectFields := strings.Join(fieldList, ",")

	query := fmt.Sprintf(`
		SELECT DISTINCT %s
		FROM %s m
		INNER JOIN sys_role_menu rm ON m.menu_id = rm.menu_id
		INNER JOIN sys_user_role ur ON rm.role_id = ur.role_id
		INNER JOIN sys_role r ON ur.role_id = r.role_id
		WHERE ur.user_id = ?
		  AND m.menu_type IN ('M', 'C')
		  AND m.status = '0'
		  AND r.status = '0'
		ORDER BY m.parent_id ASC, m.order_num ASC
	`, selectFields, m.table)

	var menus []*SysMenu
	err := m.conn.QueryRowsPartialCtx(ctx, &menus, query, userId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return menus, nil
}

// SelectMenuPermissionsByUserId 根据用户ID查询菜单权限
func (m *customSysMenuModel) SelectMenuPermissionsByUserId(ctx context.Context, userId int64) ([]string, error) {
	query := `
		SELECT DISTINCT m.perms
		FROM sys_menu m
		INNER JOIN sys_role_menu rm ON m.menu_id = rm.menu_id
		INNER JOIN sys_user_role ur ON rm.role_id = ur.role_id
		WHERE ur.user_id = ? 
		  AND m.perms IS NOT NULL 
		  AND m.perms != ''
		  AND m.status = '0'
		  AND m.del_flag = '0'
	`

	type permRow struct {
		Perms sql.NullString `db:"perms"`
	}

	var rows []permRow
	err := m.conn.QueryRowsPartialCtx(ctx, &rows, query, userId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 过滤空字符串
	result := make([]string, 0)
	for _, row := range rows {
		if row.Perms.Valid {
			perm := strings.TrimSpace(row.Perms.String)
			if perm != "" {
				result = append(result, perm)
			}
		}
	}

	return result, nil
}

// FindAll 查询菜单列表（根据条件）
// 注意：该方法需要调用者判断是否是超级管理员，如果是超级管理员则传入 userId=0
func (m *customSysMenuModel) FindAll(ctx context.Context, query *MenuQuery, userId int64) ([]*SysMenu, error) {
	if query == nil {
		query = &MenuQuery{}
	}

	// 构建 WHERE 条件
	var whereClause []string
	var args []interface{}

	// 如果不是超级管理员（userId > 0），需要根据用户ID过滤菜单
	if userId > 0 {
		// 使用子查询过滤：通过用户角色关联查询菜单
		// 子查询：SELECT DISTINCT m.menu_id FROM sys_menu m
		//   INNER JOIN sys_role_menu rm ON m.menu_id = rm.menu_id
		//   INNER JOIN sys_user_role ur ON rm.role_id = ur.role_id
		//   WHERE ur.user_id = ?
		whereClause = append(whereClause, fmt.Sprintf(`
			menu_id IN (
				SELECT DISTINCT m.menu_id 
				FROM %s m
				INNER JOIN sys_role_menu rm ON m.menu_id = rm.menu_id
				INNER JOIN sys_user_role ur ON rm.role_id = ur.role_id
				INNER JOIN sys_role r ON ur.role_id = r.role_id
				WHERE ur.user_id = ? AND r.status = '0'
			)
		`, m.table))
		args = append(args, userId)
	}

	// 构建查询条件
	if query.MenuName != "" {
		whereClause = append(whereClause, "menu_name LIKE ?")
		args = append(args, "%"+query.MenuName+"%")
	}
	if query.Visible != "" {
		whereClause = append(whereClause, "visible = ?")
		args = append(args, query.Visible)
	}
	if query.Status != "" {
		whereClause = append(whereClause, "status = ?")
		args = append(args, query.Status)
	}
	if query.MenuType != "" {
		whereClause = append(whereClause, "menu_type = ?")
		args = append(args, query.MenuType)
	}
	if query.ParentId > 0 {
		whereClause = append(whereClause, "parent_id = ?")
		args = append(args, query.ParentId)
	}

	whereSQL := ""
	if len(whereClause) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClause, " AND ")
	}

	querySQL := fmt.Sprintf(`
		SELECT %s
		FROM %s
		%s
		ORDER BY parent_id ASC, order_num ASC
	`, sysMenuRows, m.table, whereSQL)

	var menus []*SysMenu
	err := m.conn.QueryRowsPartialCtx(ctx, &menus, querySQL, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return menus, nil
}

// CheckMenuNameUnique 检查菜单名称唯一性（同父菜单下唯一）
func (m *customSysMenuModel) CheckMenuNameUnique(ctx context.Context, menuName string, parentId int64, excludeMenuId int64) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where menu_name = ? and parent_id = ?", m.table)
	var args []interface{}
	args = append(args, menuName, parentId)

	if excludeMenuId > 0 {
		query += " and menu_id != ?"
		args = append(args, excludeMenuId)
	}

	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// HasChildByMenuId 是否存在子菜单
func (m *customSysMenuModel) HasChildByMenuId(ctx context.Context, menuId int64) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where parent_id = ?", m.table)
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, menuId)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasChildByMenuIds 是否存在子菜单（批量）
func (m *customSysMenuModel) HasChildByMenuIds(ctx context.Context, menuIds []int64) (bool, error) {
	if len(menuIds) == 0 {
		return false, nil
	}

	// 构建 IN 查询
	placeholders := ""
	for i := 0; i < len(menuIds); i++ {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
	}

	// 查询是否有子菜单（parent_id 在 menuIds 中，但 menu_id 不在 menuIds 中）
	query := fmt.Sprintf(`
		select count(*) 
		from %s 
		where parent_id in (%s) 
		  and menu_id not in (%s)
	`, m.table, placeholders, placeholders)

	var args []interface{}
	for _, id := range menuIds {
		args = append(args, id)
	}
	for _, id := range menuIds {
		args = append(args, id)
	}

	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckMenuExistRole 检查菜单是否分配给角色
func (m *customSysMenuModel) CheckMenuExistRole(ctx context.Context, menuId int64) (bool, error) {
	query := "select count(*) from `sys_role_menu` where menu_id = ?"
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, menuId)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SelectMenuListByRoleId 根据角色ID查询菜单ID列表
func (m *customSysMenuModel) SelectMenuListByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	// 查询角色信息以获取 menuCheckStrictly
	roleQuery := "select menu_check_strictly from `sys_role` where role_id = ?"
	var menuCheckStrictly bool
	err := m.conn.QueryRowPartialCtx(ctx, &menuCheckStrictly, roleQuery, roleId)
	if err != nil {
		return nil, err
	}

	var query string
	if menuCheckStrictly {
		// 严格模式：只返回直接分配的菜单ID
		query = "select menu_id from `sys_role_menu` where role_id = ?"
	} else {
		// 非严格模式：返回菜单及其所有父菜单的ID
		// 这里简化处理，直接返回分配的菜单ID
		// 实际应该递归查询父菜单
		query = "select menu_id from `sys_role_menu` where role_id = ?"
	}

	var menuIds []int64
	err = m.conn.QueryRowsPartialCtx(ctx, &menuIds, query, roleId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return menuIds, nil
}

// UpdateById 根据ID更新菜单，只更新非零值字段
func (m *customSysMenuModel) UpdateById(ctx context.Context, data *SysMenu) error {
	if data.MenuId == 0 {
		return fmt.Errorf("menu_id cannot be zero")
	}

	var setParts []string
	var args []interface{}

	// 检查每个字段是否为非零值，如果是则加入更新列表
	if data.MenuName != "" {
		setParts = append(setParts, "`menu_name` = ?")
		args = append(args, data.MenuName)
	}
	if data.ParentId > 0 {
		setParts = append(setParts, "`parent_id` = ?")
		args = append(args, data.ParentId)
	}
	if data.OrderNum > 0 {
		setParts = append(setParts, "`order_num` = ?")
		args = append(args, data.OrderNum)
	}
	if data.Path != "" {
		setParts = append(setParts, "`path` = ?")
		args = append(args, data.Path)
	}
	if data.Component.Valid {
		setParts = append(setParts, "`component` = ?")
		args = append(args, data.Component.String)
	}
	if data.QueryParam.Valid {
		setParts = append(setParts, "`query_param` = ?")
		args = append(args, data.QueryParam.String)
	}
	if data.IsFrame > 0 {
		setParts = append(setParts, "`is_frame` = ?")
		args = append(args, data.IsFrame)
	}
	if data.IsCache > 0 {
		setParts = append(setParts, "`is_cache` = ?")
		args = append(args, data.IsCache)
	}
	if data.MenuType != "" {
		setParts = append(setParts, "`menu_type` = ?")
		args = append(args, data.MenuType)
	}
	if data.Visible != "" {
		setParts = append(setParts, "`visible` = ?")
		args = append(args, data.Visible)
	}
	if data.Status != "" {
		setParts = append(setParts, "`status` = ?")
		args = append(args, data.Status)
	}
	if data.Perms.Valid {
		setParts = append(setParts, "`perms` = ?")
		args = append(args, data.Perms.String)
	}
	if data.Icon != "" {
		setParts = append(setParts, "`icon` = ?")
		args = append(args, data.Icon)
	}
	if data.Remark != "" {
		setParts = append(setParts, "`remark` = ?")
		args = append(args, data.Remark)
	}
	if data.CreateDept.Valid {
		setParts = append(setParts, "`create_dept` = ?")
		args = append(args, data.CreateDept.Int64)
	}
	if data.CreateBy.Valid {
		setParts = append(setParts, "`create_by` = ?")
		args = append(args, data.CreateBy.Int64)
	}
	if data.CreateTime.Valid {
		setParts = append(setParts, "`create_time` = ?")
		args = append(args, data.CreateTime.Time)
	}
	if data.UpdateBy.Valid {
		setParts = append(setParts, "`update_by` = ?")
		args = append(args, data.UpdateBy.Int64)
	}
	if data.UpdateTime.Valid {
		setParts = append(setParts, "`update_time` = ?")
		args = append(args, data.UpdateTime.Time)
	}

	if len(setParts) == 0 {
		return nil // 没有需要更新的字段
	}

	// 构建更新SQL
	setClause := strings.Join(setParts, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE `menu_id` = ?", m.table, setClause)
	args = append(args, data.MenuId)

	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}
