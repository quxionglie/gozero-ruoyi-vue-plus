// Code generated manually for joint primary key table. DO NOT EDIT.
package sys

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	sysRoleDeptRows = "`role_id`,`dept_id`"
)

type (
	sysRoleDeptModel interface {
		Insert(ctx context.Context, data *SysRoleDept) (sql.Result, error)
		FindOne(ctx context.Context, roleId, deptId int64) (*SysRoleDept, error)
		Delete(ctx context.Context, roleId, deptId int64) error
		DeleteByRoleId(ctx context.Context, roleId int64) error
		DeleteByDeptId(ctx context.Context, deptId int64) error
		FindByRoleId(ctx context.Context, roleId int64) ([]*SysRoleDept, error)
		FindByDeptId(ctx context.Context, deptId int64) ([]*SysRoleDept, error)
	}

	defaultSysRoleDeptModel struct {
		conn  sqlx.SqlConn
		table string
	}

	SysRoleDept struct {
		RoleId int64 `db:"role_id"` // 角色ID
		DeptId int64 `db:"dept_id"` // 部门ID
	}
)

func newSysRoleDeptModel(conn sqlx.SqlConn) *defaultSysRoleDeptModel {
	return &defaultSysRoleDeptModel{
		conn:  conn,
		table: "`sys_role_dept`",
	}
}

func (m *defaultSysRoleDeptModel) Insert(ctx context.Context, data *SysRoleDept) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, sysRoleDeptRows)
	return m.conn.ExecCtx(ctx, query, data.RoleId, data.DeptId)
}

func (m *defaultSysRoleDeptModel) FindOne(ctx context.Context, roleId, deptId int64) (*SysRoleDept, error) {
	var resp SysRoleDept
	query := fmt.Sprintf("select %s from %s where `role_id` = ? and `dept_id` = ? limit 1", sysRoleDeptRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, roleId, deptId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSysRoleDeptModel) Delete(ctx context.Context, roleId, deptId int64) error {
	query := fmt.Sprintf("delete from %s where `role_id` = ? and `dept_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, roleId, deptId)
	return err
}

func (m *defaultSysRoleDeptModel) DeleteByRoleId(ctx context.Context, roleId int64) error {
	query := fmt.Sprintf("delete from %s where `role_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, roleId)
	return err
}

func (m *defaultSysRoleDeptModel) DeleteByDeptId(ctx context.Context, deptId int64) error {
	query := fmt.Sprintf("delete from %s where `dept_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, deptId)
	return err
}

func (m *defaultSysRoleDeptModel) FindByRoleId(ctx context.Context, roleId int64) ([]*SysRoleDept, error) {
	query := fmt.Sprintf("select %s from %s where `role_id` = ?", sysRoleDeptRows, m.table)
	var resp []*SysRoleDept
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, roleId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultSysRoleDeptModel) FindByDeptId(ctx context.Context, deptId int64) ([]*SysRoleDept, error) {
	query := fmt.Sprintf("select %s from %s where `dept_id` = ?", sysRoleDeptRows, m.table)
	var resp []*SysRoleDept
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, deptId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultSysRoleDeptModel) tableName() string {
	return m.table
}
