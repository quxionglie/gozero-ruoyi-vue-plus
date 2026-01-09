package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysOperLogModel = (*customSysOperLogModel)(nil)

type (
	// SysOperLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysOperLogModel.
	SysOperLogModel interface {
		sysOperLogModel
		withSession(session sqlx.Session) SysOperLogModel
	}

	customSysOperLogModel struct {
		*defaultSysOperLogModel
	}
)

// NewSysOperLogModel returns a model for the database table.
func NewSysOperLogModel(conn sqlx.SqlConn) SysOperLogModel {
	return &customSysOperLogModel{
		defaultSysOperLogModel: newSysOperLogModel(conn),
	}
}

func (m *customSysOperLogModel) withSession(session sqlx.Session) SysOperLogModel {
	return NewSysOperLogModel(sqlx.NewSqlConnFromSession(session))
}
