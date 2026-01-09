package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysOssConfigModel = (*customSysOssConfigModel)(nil)

type (
	// SysOssConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysOssConfigModel.
	SysOssConfigModel interface {
		sysOssConfigModel
		withSession(session sqlx.Session) SysOssConfigModel
	}

	customSysOssConfigModel struct {
		*defaultSysOssConfigModel
	}
)

// NewSysOssConfigModel returns a model for the database table.
func NewSysOssConfigModel(conn sqlx.SqlConn) SysOssConfigModel {
	return &customSysOssConfigModel{
		defaultSysOssConfigModel: newSysOssConfigModel(conn),
	}
}

func (m *customSysOssConfigModel) withSession(session sqlx.Session) SysOssConfigModel {
	return NewSysOssConfigModel(sqlx.NewSqlConnFromSession(session))
}
