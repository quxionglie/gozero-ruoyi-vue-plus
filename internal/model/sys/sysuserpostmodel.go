package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysUserPostModel = (*customSysUserPostModel)(nil)

type (
	// SysUserPostModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysUserPostModel.
	SysUserPostModel interface {
		sysUserPostModel
		withSession(session sqlx.Session) SysUserPostModel
	}

	customSysUserPostModel struct {
		*defaultSysUserPostModel
	}
)

// NewSysUserPostModel returns a model for the database table.
func NewSysUserPostModel(conn sqlx.SqlConn) SysUserPostModel {
	return &customSysUserPostModel{
		defaultSysUserPostModel: newSysUserPostModel(conn),
	}
}

func (m *customSysUserPostModel) withSession(session sqlx.Session) SysUserPostModel {
	return NewSysUserPostModel(sqlx.NewSqlConnFromSession(session))
}
