package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysRoleModel = (*customSysRoleModel)(nil)

type (
	// SysRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysRoleModel.
	SysRoleModel interface {
		sysRoleModel
		withSession(session sqlx.Session) SysRoleModel
		// SelectRolesByUserId 根据用户ID查询角色列表
		SelectRolesByUserId(ctx context.Context, userId int64) ([]*RoleVo, error)
		// CheckIsSuperAdmin 检查用户是否为超级管理员
		CheckIsSuperAdmin(ctx context.Context, userId int64) (bool, error)
		// SelectRoleKeysByUserId 根据用户ID查询角色权限标识
		SelectRoleKeysByUserId(ctx context.Context, userId int64) ([]string, error)
		// FindPage 分页查询角色列表
		FindPage(ctx context.Context, query *RoleQuery, pageQuery *PageQuery) ([]*SysRole, int64, error)
		// FindByIds 根据角色ID列表查询角色列表
		FindByIds(ctx context.Context, roleIds []int64) ([]*SysRole, error)
		// FindAll 查询所有角色
		FindAll(ctx context.Context, query *RoleQuery) ([]*SysRole, error)
		// CheckRoleNameUnique 检查角色名称是否唯一
		CheckRoleNameUnique(ctx context.Context, roleName string, excludeRoleId int64) (bool, error)
		// CheckRoleKeyUnique 检查角色权限标识是否唯一
		CheckRoleKeyUnique(ctx context.Context, roleKey string, excludeRoleId int64) (bool, error)
		// CountUserRoleByRoleId 统计角色使用数量
		CountUserRoleByRoleId(ctx context.Context, roleId int64) (int64, error)
		// UpdateRoleStatus 更新角色状态
		UpdateRoleStatus(ctx context.Context, roleId int64, status string) error
	}

	// RoleQuery 角色查询条件
	RoleQuery struct {
		RoleId   int64  // 角色ID
		RoleName string // 角色名称（模糊查询）
		RoleKey  string // 角色权限标识（模糊查询）
		Status   string // 角色状态（0正常 1停用）
	}

	customSysRoleModel struct {
		*defaultSysRoleModel
	}

	// RoleVo 角色视图对象（用于 SelectRolesByUserId 返回）
	RoleVo struct {
		RoleId     int64          `db:"role_id"`
		RoleName   string         `db:"role_name"`
		RoleKey    string         `db:"role_key"`
		RoleSort   int64          `db:"role_sort"`
		DataScope  string         `db:"data_scope"`
		Status     string         `db:"status"`
		Remark     sql.NullString `db:"remark"`
		CreateTime sql.NullTime   `db:"create_time"`
	}
)

// NewSysRoleModel returns a model for the database table.
func NewSysRoleModel(conn sqlx.SqlConn) SysRoleModel {
	return &customSysRoleModel{
		defaultSysRoleModel: newSysRoleModel(conn),
	}
}

func (m *customSysRoleModel) withSession(session sqlx.Session) SysRoleModel {
	return NewSysRoleModel(sqlx.NewSqlConnFromSession(session))
}

// SelectRolesByUserId 根据用户ID查询角色列表
func (m *customSysRoleModel) SelectRolesByUserId(ctx context.Context, userId int64) ([]*RoleVo, error) {
	query := `
		SELECT r.role_id, r.role_name, r.role_key, r.role_sort, r.data_scope, r.status, r.remark, r.create_time
		FROM sys_role r
		INNER JOIN sys_user_role ur ON r.role_id = ur.role_id
		WHERE ur.user_id = ? 
		  AND r.status = '0'
		  AND r.del_flag = '0'
		ORDER BY r.role_sort ASC
	`

	var rows []*RoleVo
	err := m.conn.QueryRowsPartialCtx(ctx, &rows, query, userId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return rows, nil
}

// CheckIsSuperAdmin 检查用户是否为超级管理员
func (m *customSysRoleModel) CheckIsSuperAdmin(ctx context.Context, userId int64) (bool, error) {
	query := `
		SELECT DISTINCT r.role_id, r.role_key
		FROM sys_role r
		INNER JOIN sys_user_role ur ON r.role_id = ur.role_id
		WHERE ur.user_id = ? 
		  AND r.status = '0'
		  AND r.del_flag = '0'
	`

	type roleRow struct {
		RoleId  int64  `db:"role_id"`
		RoleKey string `db:"role_key"`
	}

	var rows []roleRow
	err := m.conn.QueryRowsPartialCtx(ctx, &rows, query, userId)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	for _, row := range rows {
		if row.RoleId == 1 || strings.ToLower(row.RoleKey) == "superadmin" {
			return true, nil
		}
	}
	return false, nil
}

// SelectRoleKeysByUserId 根据用户ID查询角色权限标识
func (m *customSysRoleModel) SelectRoleKeysByUserId(ctx context.Context, userId int64) ([]string, error) {
	query := `
		SELECT DISTINCT r.role_key
		FROM sys_role r
		INNER JOIN sys_user_role ur ON r.role_id = ur.role_id
		WHERE ur.user_id = ? 
		  AND r.role_key IS NOT NULL 
		  AND r.role_key != ''
		  AND r.status = '0'
		  AND r.del_flag = '0'
	`

	type roleRow struct {
		RoleKey string `db:"role_key"`
	}

	var rows []roleRow
	err := m.conn.QueryRowsPartialCtx(ctx, &rows, query, userId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 处理角色权限标识（可能包含多个，用逗号分隔）
	result := make([]string, 0)
	for _, row := range rows {
		roleKey := strings.TrimSpace(row.RoleKey)
		if roleKey != "" {
			// 如果角色标识包含逗号，拆分成多个
			keys := strings.Split(roleKey, ",")
			for _, key := range keys {
				key = strings.TrimSpace(key)
				if key != "" {
					result = append(result, key)
				}
			}
		}
	}

	return result, nil
}

// FindPage 分页查询角色列表
func (m *customSysRoleModel) FindPage(ctx context.Context, query *RoleQuery, pageQuery *PageQuery) ([]*SysRole, int64, error) {
	if query == nil {
		query = &RoleQuery{}
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{
			PageNum:  1,
			PageSize: 10,
		}
	} else {
		// 初始化分页参数的非合规值
		pageQuery.Normalize()
	}

	// 构建 WHERE 条件
	whereClause := "del_flag = '0'"
	var args []interface{}

	if query.RoleId > 0 {
		whereClause += " and role_id = ?"
		args = append(args, query.RoleId)
	}
	if query.RoleName != "" {
		whereClause += " and role_name LIKE ?"
		args = append(args, "%"+query.RoleName+"%")
	}
	if query.RoleKey != "" {
		whereClause += " and role_key LIKE ?"
		args = append(args, "%"+query.RoleKey+"%")
	}
	if query.Status != "" {
		whereClause += " and status = ?"
		args = append(args, query.Status)
	}

	// 计算总数
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)
	var total int64
	err := m.conn.QueryRowPartialCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 构建 ORDER BY 子句（防止 SQL 注入）
	// 允许的排序列（支持 snake_case 和 camelCase）
	allowedOrderColumns := map[string]bool{
		"role_id":     true,
		"roleId":      true,
		"role_sort":   true,
		"roleSort":    true,
		"create_time": true,
		"createTime":  true,
		"role_name":   true,
		"roleName":    true,
	}

	orderBy := "role_sort ASC, create_time ASC"
	if pageQuery.OrderByColumn != "" {
		// 将 camelCase 转换为 snake_case
		columnName := camelToSnake(strings.TrimSpace(pageQuery.OrderByColumn))
		// 检查原始字段名和转换后的字段名是否在允许列表中
		originalColumn := strings.TrimSpace(pageQuery.OrderByColumn)
		if allowedOrderColumns[originalColumn] || allowedOrderColumns[columnName] {
			// 使用转换后的 snake_case 字段名
			orderBy = columnName + " "
			// 处理排序方向（兼容 asc、desc、descending 等）
			isAscStr := strings.ToLower(strings.TrimSpace(pageQuery.IsAsc))
			if isAscStr == "asc" || isAscStr == "ascending" {
				orderBy += "ASC"
			} else {
				orderBy += "DESC"
			}
		}
	}

	// 计算分页参数
	offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
	limit := pageQuery.PageSize

	// 查询数据（引用gen文件中的常量）
	rows := "role_id,tenant_id,role_name,role_key,role_sort,data_scope,menu_check_strictly,dept_check_strictly,status,del_flag,create_dept,create_by,create_time,update_by,update_time,remark"
	querySQL := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE %s
		ORDER BY %s
		LIMIT ?, ?
	`, rows, m.table, whereClause, orderBy)

	args = append(args, offset, limit)

	var roleList []*SysRole
	err = m.conn.QueryRowsPartialCtx(ctx, &roleList, querySQL, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	return roleList, total, nil
}

// FindByIds 根据角色ID列表查询角色列表
func (m *customSysRoleModel) FindByIds(ctx context.Context, roleIds []int64) ([]*SysRole, error) {
	if len(roleIds) == 0 {
		return []*SysRole{}, nil
	}

	// 构建 IN 查询
	placeholders := ""
	for i := 0; i < len(roleIds); i++ {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
	}

	// 引用gen文件中的常量
	rows := "role_id,tenant_id,role_name,role_key,role_sort,data_scope,menu_check_strictly,dept_check_strictly,status,del_flag,create_dept,create_by,create_time,update_by,update_time,remark"
	query := fmt.Sprintf(`
		SELECT %s 
		FROM %s 
		WHERE role_id IN (%s) 
		  AND status = '0' 
		  AND del_flag = '0'
		ORDER BY role_sort ASC
	`, rows, m.table, placeholders)

	var args []interface{}
	for _, id := range roleIds {
		args = append(args, id)
	}

	var roleList []*SysRole
	err := m.conn.QueryRowsPartialCtx(ctx, &roleList, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return roleList, nil
}

// FindAll 查询所有角色
func (m *customSysRoleModel) FindAll(ctx context.Context, query *RoleQuery) ([]*SysRole, error) {
	if query == nil {
		query = &RoleQuery{}
	}

	// 构建 WHERE 条件
	whereClause := "del_flag = '0'"
	var args []interface{}

	if query.RoleId > 0 {
		whereClause += " and role_id = ?"
		args = append(args, query.RoleId)
	}
	if query.RoleName != "" {
		whereClause += " and role_name LIKE ?"
		args = append(args, "%"+query.RoleName+"%")
	}
	if query.RoleKey != "" {
		whereClause += " and role_key LIKE ?"
		args = append(args, "%"+query.RoleKey+"%")
	}
	if query.Status != "" {
		whereClause += " and status = ?"
		args = append(args, query.Status)
	}

	// 引用gen文件中的常量
	rows := "role_id,tenant_id,role_name,role_key,role_sort,data_scope,menu_check_strictly,dept_check_strictly,status,del_flag,create_dept,create_by,create_time,update_by,update_time,remark"
	querySQL := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE %s
		ORDER BY role_sort ASC, create_time ASC
	`, rows, m.table, whereClause)

	var roleList []*SysRole
	err := m.conn.QueryRowsPartialCtx(ctx, &roleList, querySQL, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return roleList, nil
}

// CheckRoleNameUnique 检查角色名称是否唯一
func (m *customSysRoleModel) CheckRoleNameUnique(ctx context.Context, roleName string, excludeRoleId int64) (bool, error) {
	query := "SELECT COUNT(*) FROM sys_role WHERE role_name = ? AND del_flag = '0'"
	args := []interface{}{roleName}
	if excludeRoleId > 0 {
		query += " AND role_id != ?"
		args = append(args, excludeRoleId)
	}

	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// CheckRoleKeyUnique 检查角色权限标识是否唯一
func (m *customSysRoleModel) CheckRoleKeyUnique(ctx context.Context, roleKey string, excludeRoleId int64) (bool, error) {
	query := "SELECT COUNT(*) FROM sys_role WHERE role_key = ? AND del_flag = '0'"
	args := []interface{}{roleKey}
	if excludeRoleId > 0 {
		query += " AND role_id != ?"
		args = append(args, excludeRoleId)
	}

	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// CountUserRoleByRoleId 统计角色使用数量
func (m *customSysRoleModel) CountUserRoleByRoleId(ctx context.Context, roleId int64) (int64, error) {
	query := "SELECT COUNT(*) FROM sys_user_role WHERE role_id = ?"
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, roleId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateRoleStatus 更新角色状态
func (m *customSysRoleModel) UpdateRoleStatus(ctx context.Context, roleId int64, status string) error {
	query := fmt.Sprintf("UPDATE %s SET status = ? WHERE role_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, status, roleId)
	return err
}
