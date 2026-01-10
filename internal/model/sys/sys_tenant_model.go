package sys

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysTenantModel = (*customSysTenantModel)(nil)

type (
	// SysTenantModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysTenantModel.
	SysTenantModel interface {
		sysTenantModel
		withSession(session sqlx.Session) SysTenantModel
		FindOneByTenantId(ctx context.Context, tenantId string) (*SysTenant, error)
		FindAllAvailable(ctx context.Context) ([]*SysTenant, error)
	}

	customSysTenantModel struct {
		*defaultSysTenantModel
	}
)

// NewSysTenantModel returns a model for the database table.
func NewSysTenantModel(conn sqlx.SqlConn) SysTenantModel {
	return &customSysTenantModel{
		defaultSysTenantModel: newSysTenantModel(conn),
	}
}

func (m *customSysTenantModel) withSession(session sqlx.Session) SysTenantModel {
	return NewSysTenantModel(sqlx.NewSqlConnFromSession(session))
}

// FindOneByTenantId 根据租户ID查询
func (m *customSysTenantModel) FindOneByTenantId(ctx context.Context, tenantId string) (*SysTenant, error) {
	query := fmt.Sprintf("select %s from %s where `tenant_id` = ? and `del_flag` = '0' limit 1", sysTenantRows, m.table)
	var resp SysTenant
	err := m.conn.QueryRowCtx(ctx, &resp, query, tenantId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindAllAvailable 查询所有可用的租户（status=0, del_flag=0）
func (m *customSysTenantModel) FindAllAvailable(ctx context.Context) ([]*SysTenant, error) {
	query := fmt.Sprintf("select %s from %s where `status` = '0' and `del_flag` = '0'", sysTenantRows, m.table)
	var resp []*SysTenant
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
