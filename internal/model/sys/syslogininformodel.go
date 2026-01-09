package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysLogininforModel = (*customSysLogininforModel)(nil)

type (
	// SysLogininforModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysLogininforModel.
	SysLogininforModel interface {
		sysLogininforModel
		withSession(session sqlx.Session) SysLogininforModel
	}

	customSysLogininforModel struct {
		*defaultSysLogininforModel
	}
)

// NewSysLogininforModel returns a model for the database table.
func NewSysLogininforModel(conn sqlx.SqlConn) SysLogininforModel {
	return &customSysLogininforModel{
		defaultSysLogininforModel: newSysLogininforModel(conn),
	}
}

func (m *customSysLogininforModel) withSession(session sqlx.Session) SysLogininforModel {
	return NewSysLogininforModel(sqlx.NewSqlConnFromSession(session))
}
