package sys

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysSocialModel = (*customSysSocialModel)(nil)

type (
	// SysSocialModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysSocialModel.
	SysSocialModel interface {
		sysSocialModel
		withSession(session sqlx.Session) SysSocialModel
		// FindByUserId 根据用户ID查询社会化关系列表
		FindByUserId(ctx context.Context, userId int64) ([]*SysSocial, error)
	}

	customSysSocialModel struct {
		*defaultSysSocialModel
	}
)

// NewSysSocialModel returns a model for the database table.
func NewSysSocialModel(conn sqlx.SqlConn) SysSocialModel {
	return &customSysSocialModel{
		defaultSysSocialModel: newSysSocialModel(conn),
	}
}

func (m *customSysSocialModel) withSession(session sqlx.Session) SysSocialModel {
	return NewSysSocialModel(sqlx.NewSqlConnFromSession(session))
}

// FindByUserId 根据用户ID查询社会化关系列表
func (m *customSysSocialModel) FindByUserId(ctx context.Context, userId int64) ([]*SysSocial, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `del_flag` = '0'", sysSocialRows, m.table)
	var resp []*SysSocial
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
