package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysTenantModel = (*customSysTenantModel)(nil)

type (
	// SysTenantModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysTenantModel.
	SysTenantModel interface {
		sysTenantModel
		withSession(session sqlx.Session) SysTenantModel
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
