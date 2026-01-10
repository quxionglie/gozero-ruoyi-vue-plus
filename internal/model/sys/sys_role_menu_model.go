package sys

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysRoleMenuModel = (*customSysRoleMenuModel)(nil)

type (
	// SysRoleMenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysRoleMenuModel.
	SysRoleMenuModel interface {
		sysRoleMenuModel
		withSession(session sqlx.Session) SysRoleMenuModel
		// InsertBatch 批量插入角色菜单关联
		InsertBatch(ctx context.Context, roleId int64, menuIds []int64) error
		// DeleteByRoleIds 批量删除角色菜单关联（根据角色ID列表）
		DeleteByRoleIds(ctx context.Context, roleIds []int64) error
	}

	customSysRoleMenuModel struct {
		*defaultSysRoleMenuModel
	}
)

// NewSysRoleMenuModel returns a model for the database table.
func NewSysRoleMenuModel(conn sqlx.SqlConn) SysRoleMenuModel {
	return &customSysRoleMenuModel{
		defaultSysRoleMenuModel: newSysRoleMenuModel(conn),
	}
}

func (m *customSysRoleMenuModel) withSession(session sqlx.Session) SysRoleMenuModel {
	return NewSysRoleMenuModel(sqlx.NewSqlConnFromSession(session))
}

// InsertBatch 批量插入角色菜单关联
func (m *customSysRoleMenuModel) InsertBatch(ctx context.Context, roleId int64, menuIds []int64) error {
	if len(menuIds) == 0 {
		return nil
	}

	// 构建批量插入SQL
	placeholders := ""
	args := make([]interface{}, 0, len(menuIds)*2)
	for i, menuId := range menuIds {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "(?, ?)"
		args = append(args, roleId, menuId)
	}

	query := fmt.Sprintf("INSERT INTO %s (role_id, menu_id) VALUES %s", m.table, placeholders)
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

// DeleteByRoleIds 批量删除角色菜单关联（根据角色ID列表）
func (m *customSysRoleMenuModel) DeleteByRoleIds(ctx context.Context, roleIds []int64) error {
	if len(roleIds) == 0 {
		return nil
	}

	// 构建 IN 查询
	placeholders := ""
	for i := 0; i < len(roleIds); i++ {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE role_id IN (%s)", m.table, placeholders)
	args := make([]interface{}, len(roleIds))
	for i, roleId := range roleIds {
		args[i] = roleId
	}
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}
