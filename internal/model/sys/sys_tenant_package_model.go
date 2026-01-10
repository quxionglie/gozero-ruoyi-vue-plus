package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysTenantPackageModel = (*customSysTenantPackageModel)(nil)

type (
	// SysTenantPackageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysTenantPackageModel.
	SysTenantPackageModel interface {
		sysTenantPackageModel
		withSession(session sqlx.Session) SysTenantPackageModel
	}

	customSysTenantPackageModel struct {
		*defaultSysTenantPackageModel
	}
)

// NewSysTenantPackageModel returns a model for the database table.
func NewSysTenantPackageModel(conn sqlx.SqlConn) SysTenantPackageModel {
	return &customSysTenantPackageModel{
		defaultSysTenantPackageModel: newSysTenantPackageModel(conn),
	}
}

func (m *customSysTenantPackageModel) withSession(session sqlx.Session) SysTenantPackageModel {
	return NewSysTenantPackageModel(sqlx.NewSqlConnFromSession(session))
}
