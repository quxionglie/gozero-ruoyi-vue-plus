package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysPostModel = (*customSysPostModel)(nil)

type (
	// SysPostModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysPostModel.
	SysPostModel interface {
		sysPostModel
		withSession(session sqlx.Session) SysPostModel
	}

	customSysPostModel struct {
		*defaultSysPostModel
	}
)

// NewSysPostModel returns a model for the database table.
func NewSysPostModel(conn sqlx.SqlConn) SysPostModel {
	return &customSysPostModel{
		defaultSysPostModel: newSysPostModel(conn),
	}
}

func (m *customSysPostModel) withSession(session sqlx.Session) SysPostModel {
	return NewSysPostModel(sqlx.NewSqlConnFromSession(session))
}
