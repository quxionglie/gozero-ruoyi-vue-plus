// Code generated manually for joint primary key table. DO NOT EDIT.
package sys

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	sysUserRoleRows = "`user_id`,`role_id`"
)

type (
	sysUserRoleModel interface {
		Insert(ctx context.Context, data *SysUserRole) (sql.Result, error)
		FindOne(ctx context.Context, userId, roleId int64) (*SysUserRole, error)
		Delete(ctx context.Context, userId, roleId int64) error
		DeleteByUserId(ctx context.Context, userId int64) error
		DeleteByRoleId(ctx context.Context, roleId int64) error
		FindByUserId(ctx context.Context, userId int64) ([]*SysUserRole, error)
		FindByRoleId(ctx context.Context, roleId int64) ([]*SysUserRole, error)
	}

	defaultSysUserRoleModel struct {
		conn  sqlx.SqlConn
		table string
	}

	SysUserRole struct {
		UserId int64 `db:"user_id"` // 用户ID
		RoleId int64 `db:"role_id"` // 角色ID
	}
)

func newSysUserRoleModel(conn sqlx.SqlConn) *defaultSysUserRoleModel {
	return &defaultSysUserRoleModel{
		conn:  conn,
		table: "`sys_user_role`",
	}
}

func (m *defaultSysUserRoleModel) Insert(ctx context.Context, data *SysUserRole) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, sysUserRoleRows)
	return m.conn.ExecCtx(ctx, query, data.UserId, data.RoleId)
}

func (m *defaultSysUserRoleModel) FindOne(ctx context.Context, userId, roleId int64) (*SysUserRole, error) {
	var resp SysUserRole
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `role_id` = ? limit 1", sysUserRoleRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, roleId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSysUserRoleModel) Delete(ctx context.Context, userId, roleId int64) error {
	query := fmt.Sprintf("delete from %s where `user_id` = ? and `role_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId, roleId)
	return err
}

func (m *defaultSysUserRoleModel) DeleteByUserId(ctx context.Context, userId int64) error {
	query := fmt.Sprintf("delete from %s where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId)
	return err
}

func (m *defaultSysUserRoleModel) DeleteByRoleId(ctx context.Context, roleId int64) error {
	query := fmt.Sprintf("delete from %s where `role_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, roleId)
	return err
}

func (m *defaultSysUserRoleModel) FindByUserId(ctx context.Context, userId int64) ([]*SysUserRole, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", sysUserRoleRows, m.table)
	var resp []*SysUserRole
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultSysUserRoleModel) FindByRoleId(ctx context.Context, roleId int64) ([]*SysUserRole, error) {
	query := fmt.Sprintf("select %s from %s where `role_id` = ?", sysUserRoleRows, m.table)
	var resp []*SysUserRole
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, roleId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultSysUserRoleModel) tableName() string {
	return m.table
}
