package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysSocialModel = (*customSysSocialModel)(nil)

type (
	// SysSocialModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysSocialModel.
	SysSocialModel interface {
		sysSocialModel
		withSession(session sqlx.Session) SysSocialModel
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
