package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysDictTypeModel = (*customSysDictTypeModel)(nil)

type (
	// SysDictTypeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysDictTypeModel.
	SysDictTypeModel interface {
		sysDictTypeModel
		withSession(session sqlx.Session) SysDictTypeModel
	}

	customSysDictTypeModel struct {
		*defaultSysDictTypeModel
	}
)

// NewSysDictTypeModel returns a model for the database table.
func NewSysDictTypeModel(conn sqlx.SqlConn) SysDictTypeModel {
	return &customSysDictTypeModel{
		defaultSysDictTypeModel: newSysDictTypeModel(conn),
	}
}

func (m *customSysDictTypeModel) withSession(session sqlx.Session) SysDictTypeModel {
	return NewSysDictTypeModel(sqlx.NewSqlConnFromSession(session))
}
