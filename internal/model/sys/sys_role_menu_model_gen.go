// Code generated manually for joint primary key table. DO NOT EDIT.
package sys

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	sysRoleMenuRows = "`role_id`,`menu_id`"
)

type (
	sysRoleMenuModel interface {
		Insert(ctx context.Context, data *SysRoleMenu) (sql.Result, error)
		FindOne(ctx context.Context, roleId, menuId int64) (*SysRoleMenu, error)
		Delete(ctx context.Context, roleId, menuId int64) error
		DeleteByRoleId(ctx context.Context, roleId int64) error
		DeleteByMenuId(ctx context.Context, menuId int64) error
		FindByRoleId(ctx context.Context, roleId int64) ([]*SysRoleMenu, error)
		FindByMenuId(ctx context.Context, menuId int64) ([]*SysRoleMenu, error)
	}

	defaultSysRoleMenuModel struct {
		conn  sqlx.SqlConn
		table string
	}

	SysRoleMenu struct {
		RoleId int64 `db:"role_id"` // 角色ID
		MenuId int64 `db:"menu_id"` // 菜单ID
	}
)

func newSysRoleMenuModel(conn sqlx.SqlConn) *defaultSysRoleMenuModel {
	return &defaultSysRoleMenuModel{
		conn:  conn,
		table: "`sys_role_menu`",
	}
}

func (m *defaultSysRoleMenuModel) Insert(ctx context.Context, data *SysRoleMenu) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, sysRoleMenuRows)
	return m.conn.ExecCtx(ctx, query, data.RoleId, data.MenuId)
}

func (m *defaultSysRoleMenuModel) FindOne(ctx context.Context, roleId, menuId int64) (*SysRoleMenu, error) {
	var resp SysRoleMenu
	query := fmt.Sprintf("select %s from %s where `role_id` = ? and `menu_id` = ? limit 1", sysRoleMenuRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, roleId, menuId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSysRoleMenuModel) Delete(ctx context.Context, roleId, menuId int64) error {
	query := fmt.Sprintf("delete from %s where `role_id` = ? and `menu_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, roleId, menuId)
	return err
}

func (m *defaultSysRoleMenuModel) DeleteByRoleId(ctx context.Context, roleId int64) error {
	query := fmt.Sprintf("delete from %s where `role_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, roleId)
	return err
}

func (m *defaultSysRoleMenuModel) DeleteByMenuId(ctx context.Context, menuId int64) error {
	query := fmt.Sprintf("delete from %s where `menu_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, menuId)
	return err
}

func (m *defaultSysRoleMenuModel) FindByRoleId(ctx context.Context, roleId int64) ([]*SysRoleMenu, error) {
	query := fmt.Sprintf("select %s from %s where `role_id` = ?", sysRoleMenuRows, m.table)
	var resp []*SysRoleMenu
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, roleId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultSysRoleMenuModel) FindByMenuId(ctx context.Context, menuId int64) ([]*SysRoleMenu, error) {
	query := fmt.Sprintf("select %s from %s where `menu_id` = ?", sysRoleMenuRows, m.table)
	var resp []*SysRoleMenu
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, menuId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultSysRoleMenuModel) tableName() string {
	return m.table
}
