// Code generated manually for joint primary key table. DO NOT EDIT.
package sys

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	sysUserPostRows = "`user_id`,`post_id`"
)

type (
	sysUserPostModel interface {
		Insert(ctx context.Context, data *SysUserPost) (sql.Result, error)
		FindOne(ctx context.Context, userId, postId int64) (*SysUserPost, error)
		Delete(ctx context.Context, userId, postId int64) error
		DeleteByUserId(ctx context.Context, userId int64) error
		DeleteByPostId(ctx context.Context, postId int64) error
		FindByUserId(ctx context.Context, userId int64) ([]*SysUserPost, error)
		FindByPostId(ctx context.Context, postId int64) ([]*SysUserPost, error)
	}

	defaultSysUserPostModel struct {
		conn  sqlx.SqlConn
		table string
	}

	SysUserPost struct {
		UserId int64 `db:"user_id"` // 用户ID
		PostId int64 `db:"post_id"` // 岗位ID
	}
)

func newSysUserPostModel(conn sqlx.SqlConn) *defaultSysUserPostModel {
	return &defaultSysUserPostModel{
		conn:  conn,
		table: "`sys_user_post`",
	}
}

func (m *defaultSysUserPostModel) Insert(ctx context.Context, data *SysUserPost) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, sysUserPostRows)
	return m.conn.ExecCtx(ctx, query, data.UserId, data.PostId)
}

func (m *defaultSysUserPostModel) FindOne(ctx context.Context, userId, postId int64) (*SysUserPost, error) {
	var resp SysUserPost
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `post_id` = ? limit 1", sysUserPostRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, postId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSysUserPostModel) Delete(ctx context.Context, userId, postId int64) error {
	query := fmt.Sprintf("delete from %s where `user_id` = ? and `post_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId, postId)
	return err
}

func (m *defaultSysUserPostModel) DeleteByUserId(ctx context.Context, userId int64) error {
	query := fmt.Sprintf("delete from %s where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId)
	return err
}

func (m *defaultSysUserPostModel) DeleteByPostId(ctx context.Context, postId int64) error {
	query := fmt.Sprintf("delete from %s where `post_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, postId)
	return err
}

func (m *defaultSysUserPostModel) FindByUserId(ctx context.Context, userId int64) ([]*SysUserPost, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", sysUserPostRows, m.table)
	var resp []*SysUserPost
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultSysUserPostModel) FindByPostId(ctx context.Context, postId int64) ([]*SysUserPost, error) {
	query := fmt.Sprintf("select %s from %s where `post_id` = ?", sysUserPostRows, m.table)
	var resp []*SysUserPost
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, postId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultSysUserPostModel) tableName() string {
	return m.table
}
