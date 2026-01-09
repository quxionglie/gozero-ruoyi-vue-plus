package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysNoticeModel = (*customSysNoticeModel)(nil)

type (
	// SysNoticeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysNoticeModel.
	SysNoticeModel interface {
		sysNoticeModel
		withSession(session sqlx.Session) SysNoticeModel
	}

	customSysNoticeModel struct {
		*defaultSysNoticeModel
	}
)

// NewSysNoticeModel returns a model for the database table.
func NewSysNoticeModel(conn sqlx.SqlConn) SysNoticeModel {
	return &customSysNoticeModel{
		defaultSysNoticeModel: newSysNoticeModel(conn),
	}
}

func (m *customSysNoticeModel) withSession(session sqlx.Session) SysNoticeModel {
	return NewSysNoticeModel(sqlx.NewSqlConnFromSession(session))
}
