package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysDeptModel = (*customSysDeptModel)(nil)

type (
	// SysDeptModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysDeptModel.
	SysDeptModel interface {
		sysDeptModel
		withSession(session sqlx.Session) SysDeptModel
	}

	customSysDeptModel struct {
		*defaultSysDeptModel
	}
)

// NewSysDeptModel returns a model for the database table.
func NewSysDeptModel(conn sqlx.SqlConn) SysDeptModel {
	return &customSysDeptModel{
		defaultSysDeptModel: newSysDeptModel(conn),
	}
}

func (m *customSysDeptModel) withSession(session sqlx.Session) SysDeptModel {
	return NewSysDeptModel(sqlx.NewSqlConnFromSession(session))
}
