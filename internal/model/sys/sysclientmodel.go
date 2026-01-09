package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysClientModel = (*customSysClientModel)(nil)

type (
	// SysClientModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysClientModel.
	SysClientModel interface {
		sysClientModel
		withSession(session sqlx.Session) SysClientModel
	}

	customSysClientModel struct {
		*defaultSysClientModel
	}
)

// NewSysClientModel returns a model for the database table.
func NewSysClientModel(conn sqlx.SqlConn) SysClientModel {
	return &customSysClientModel{
		defaultSysClientModel: newSysClientModel(conn),
	}
}

func (m *customSysClientModel) withSession(session sqlx.Session) SysClientModel {
	return NewSysClientModel(sqlx.NewSqlConnFromSession(session))
}
