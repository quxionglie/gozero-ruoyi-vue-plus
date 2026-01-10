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
