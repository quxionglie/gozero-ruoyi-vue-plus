package sys

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysRoleDeptModel = (*customSysRoleDeptModel)(nil)

type (
	// SysRoleDeptModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysRoleDeptModel.
	SysRoleDeptModel interface {
		sysRoleDeptModel
		withSession(session sqlx.Session) SysRoleDeptModel
		// InsertBatch 批量插入角色部门关联
		InsertBatch(ctx context.Context, roleId int64, deptIds []int64) error
		// DeleteByRoleId 根据角色ID删除角色部门关联
		DeleteByRoleId(ctx context.Context, roleId int64) error
		// DeleteByRoleIds 批量删除角色部门关联（根据角色ID列表）
		DeleteByRoleIds(ctx context.Context, roleIds []int64) error
		// SelectDeptIdsByRoleId 根据角色ID查询部门ID列表
		SelectDeptIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error)
	}

	customSysRoleDeptModel struct {
		*defaultSysRoleDeptModel
	}
)

// NewSysRoleDeptModel returns a model for the database table.
func NewSysRoleDeptModel(conn sqlx.SqlConn) SysRoleDeptModel {
	return &customSysRoleDeptModel{
		defaultSysRoleDeptModel: newSysRoleDeptModel(conn),
	}
}

func (m *customSysRoleDeptModel) withSession(session sqlx.Session) SysRoleDeptModel {
	return NewSysRoleDeptModel(sqlx.NewSqlConnFromSession(session))
}

// InsertBatch 批量插入角色部门关联
func (m *customSysRoleDeptModel) InsertBatch(ctx context.Context, roleId int64, deptIds []int64) error {
	if len(deptIds) == 0 {
		return nil
	}

	// 构建批量插入SQL
	placeholders := ""
	args := make([]interface{}, 0, len(deptIds)*2)
	for i, deptId := range deptIds {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "(?, ?)"
		args = append(args, roleId, deptId)
	}

	query := fmt.Sprintf("INSERT INTO %s (role_id, dept_id) VALUES %s", m.table, placeholders)
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

// DeleteByRoleId 根据角色ID删除角色部门关联
func (m *customSysRoleDeptModel) DeleteByRoleId(ctx context.Context, roleId int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE role_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, roleId)
	return err
}

// DeleteByRoleIds 批量删除角色部门关联（根据角色ID列表）
func (m *customSysRoleDeptModel) DeleteByRoleIds(ctx context.Context, roleIds []int64) error {
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

// SelectDeptIdsByRoleId 根据角色ID查询部门ID列表
func (m *customSysRoleDeptModel) SelectDeptIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	query := fmt.Sprintf("SELECT dept_id FROM %s WHERE role_id = ?", m.table)
	var deptIds []int64
	err := m.conn.QueryRowsPartialCtx(ctx, &deptIds, query, roleId)
	if err != nil && err != sqlx.ErrNotFound {
		return nil, err
	}
	return deptIds, nil
}
