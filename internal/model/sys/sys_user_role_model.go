package sys

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysUserRoleModel = (*customSysUserRoleModel)(nil)

type (
	// SysUserRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysUserRoleModel.
	SysUserRoleModel interface {
		sysUserRoleModel
		withSession(session sqlx.Session) SysUserRoleModel
		// InsertBatch 批量插入用户角色关联（根据角色ID和用户ID列表）
		InsertBatch(ctx context.Context, roleId int64, userIds []int64) error
		// InsertBatchByUserId 批量插入用户角色关联（根据用户ID和角色ID列表）
		InsertBatchByUserId(ctx context.Context, userId int64, roleIds []int64) error
		// DeleteByUserId 根据用户ID删除用户角色关联
		DeleteByUserId(ctx context.Context, userId int64) error
		// DeleteByRoleIdAndUserIds 根据角色ID和用户ID列表删除用户角色关联
		DeleteByRoleIdAndUserIds(ctx context.Context, roleId int64, userIds []int64) error
		// FindByUserId 根据用户ID查询用户角色关联列表
		FindByUserId(ctx context.Context, userId int64) ([]*SysUserRole, error)
		// SelectUserIdsByRoleId 根据角色ID查询用户ID列表
		SelectUserIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error)
	}

	customSysUserRoleModel struct {
		*defaultSysUserRoleModel
	}
)

// NewSysUserRoleModel returns a model for the database table.
func NewSysUserRoleModel(conn sqlx.SqlConn) SysUserRoleModel {
	return &customSysUserRoleModel{
		defaultSysUserRoleModel: newSysUserRoleModel(conn),
	}
}

func (m *customSysUserRoleModel) withSession(session sqlx.Session) SysUserRoleModel {
	return NewSysUserRoleModel(sqlx.NewSqlConnFromSession(session))
}

// InsertBatch 批量插入用户角色关联
func (m *customSysUserRoleModel) InsertBatch(ctx context.Context, roleId int64, userIds []int64) error {
	if len(userIds) == 0 {
		return nil
	}

	// 构建批量插入SQL
	placeholders := ""
	args := make([]interface{}, 0, len(userIds)*2)
	for i, userId := range userIds {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "(?, ?)"
		args = append(args, userId, roleId)
	}

	query := fmt.Sprintf("INSERT INTO %s (user_id, role_id) VALUES %s", m.table, placeholders)
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

// DeleteByRoleIdAndUserIds 根据角色ID和用户ID列表删除用户角色关联
func (m *customSysUserRoleModel) DeleteByRoleIdAndUserIds(ctx context.Context, roleId int64, userIds []int64) error {
	if len(userIds) == 0 {
		return nil
	}

	// 构建 IN 查询
	placeholders := ""
	for i := 0; i < len(userIds); i++ {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE role_id = ? AND user_id IN (%s)", m.table, placeholders)
	args := make([]interface{}, len(userIds)+1)
	args[0] = roleId
	for i, userId := range userIds {
		args[i+1] = userId
	}
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

// InsertBatchByUserId 批量插入用户角色关联（根据用户ID和角色ID列表）
func (m *customSysUserRoleModel) InsertBatchByUserId(ctx context.Context, userId int64, roleIds []int64) error {
	if len(roleIds) == 0 {
		return nil
	}

	// 构建批量插入SQL
	placeholders := ""
	args := make([]interface{}, 0, len(roleIds)*2)
	for i, roleId := range roleIds {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "(?, ?)"
		args = append(args, userId, roleId)
	}

	query := fmt.Sprintf("INSERT INTO %s (user_id, role_id) VALUES %s", m.table, placeholders)
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

// DeleteByUserId 根据用户ID删除用户角色关联
func (m *customSysUserRoleModel) DeleteByUserId(ctx context.Context, userId int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId)
	return err
}

// FindByUserId 根据用户ID查询用户角色关联列表
func (m *customSysUserRoleModel) FindByUserId(ctx context.Context, userId int64) ([]*SysUserRole, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE user_id = ?", sysUserRoleRows, m.table)
	var resp []*SysUserRole
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	if err != nil && err != sqlx.ErrNotFound {
		return nil, err
	}
	return resp, nil
}

// SelectUserIdsByRoleId 根据角色ID查询用户ID列表
func (m *customSysUserRoleModel) SelectUserIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	query := fmt.Sprintf("SELECT user_id FROM %s WHERE role_id = ?", m.table)
	var userIds []int64
	err := m.conn.QueryRowsPartialCtx(ctx, &userIds, query, roleId)
	if err != nil && err != sqlx.ErrNotFound {
		return nil, err
	}
	return userIds, nil
}
