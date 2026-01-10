package sys

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysUserPostModel = (*customSysUserPostModel)(nil)

type (
	// SysUserPostModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysUserPostModel.
	SysUserPostModel interface {
		sysUserPostModel
		withSession(session sqlx.Session) SysUserPostModel
		// InsertBatch 批量插入用户岗位关联
		InsertBatch(ctx context.Context, userId int64, postIds []int64) error
		// DeleteByUserId 根据用户ID删除用户岗位关联
		DeleteByUserId(ctx context.Context, userId int64) error
		// SelectPostIdsByUserId 根据用户ID查询岗位ID列表
		SelectPostIdsByUserId(ctx context.Context, userId int64) ([]int64, error)
	}

	customSysUserPostModel struct {
		*defaultSysUserPostModel
	}
)

// NewSysUserPostModel returns a model for the database table.
func NewSysUserPostModel(conn sqlx.SqlConn) SysUserPostModel {
	return &customSysUserPostModel{
		defaultSysUserPostModel: newSysUserPostModel(conn),
	}
}

func (m *customSysUserPostModel) withSession(session sqlx.Session) SysUserPostModel {
	return NewSysUserPostModel(sqlx.NewSqlConnFromSession(session))
}

// InsertBatch 批量插入用户岗位关联
func (m *customSysUserPostModel) InsertBatch(ctx context.Context, userId int64, postIds []int64) error {
	if len(postIds) == 0 {
		return nil
	}

	// 构建批量插入SQL
	placeholders := ""
	args := make([]interface{}, 0, len(postIds)*2)
	for i, postId := range postIds {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "(?, ?)"
		args = append(args, userId, postId)
	}

	query := fmt.Sprintf("INSERT INTO %s (user_id, post_id) VALUES %s", m.table, placeholders)
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

// DeleteByUserId 根据用户ID删除用户岗位关联
func (m *customSysUserPostModel) DeleteByUserId(ctx context.Context, userId int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId)
	return err
}

// SelectPostIdsByUserId 根据用户ID查询岗位ID列表
func (m *customSysUserPostModel) SelectPostIdsByUserId(ctx context.Context, userId int64) ([]int64, error) {
	query := fmt.Sprintf("SELECT post_id FROM %s WHERE user_id = ?", m.table)
	var postIds []int64
	err := m.conn.QueryRowsPartialCtx(ctx, &postIds, query, userId)
	if err != nil && err != sqlx.ErrNotFound {
		return nil, err
	}
	return postIds, nil
}
