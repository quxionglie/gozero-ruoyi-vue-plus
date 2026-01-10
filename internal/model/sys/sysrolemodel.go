package sys

import (
	"context"
	"database/sql"
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
