package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysOssModel = (*customSysOssModel)(nil)

type (
	// SysOssModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysOssModel.
	SysOssModel interface {
		sysOssModel
		withSession(session sqlx.Session) SysOssModel
	}

	customSysOssModel struct {
		*defaultSysOssModel
	}
)

// NewSysOssModel returns a model for the database table.
func NewSysOssModel(conn sqlx.SqlConn) SysOssModel {
	return &customSysOssModel{
		defaultSysOssModel: newSysOssModel(conn),
	}
}

func (m *customSysOssModel) withSession(session sqlx.Session) SysOssModel {
	return NewSysOssModel(sqlx.NewSqlConnFromSession(session))
}
