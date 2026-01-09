package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysDictDataModel = (*customSysDictDataModel)(nil)

type (
	// SysDictDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysDictDataModel.
	SysDictDataModel interface {
		sysDictDataModel
		withSession(session sqlx.Session) SysDictDataModel
	}

	customSysDictDataModel struct {
		*defaultSysDictDataModel
	}
)

// NewSysDictDataModel returns a model for the database table.
func NewSysDictDataModel(conn sqlx.SqlConn) SysDictDataModel {
	return &customSysDictDataModel{
		defaultSysDictDataModel: newSysDictDataModel(conn),
	}
}

func (m *customSysDictDataModel) withSession(session sqlx.Session) SysDictDataModel {
	return NewSysDictDataModel(sqlx.NewSqlConnFromSession(session))
}
