package sys

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysRoleDeptModel = (*customSysRoleDeptModel)(nil)

type (
	// SysRoleDeptModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysRoleDeptModel.
	SysRoleDeptModel interface {
		sysRoleDeptModel
		withSession(session sqlx.Session) SysRoleDeptModel
	}

	customSysRoleDeptModel struct {
		*defaultSysRoleDeptModel
	}
)

// NewSysRoleDeptModel returns a model for the database table.
func NewSysRoleDeptModel(conn sqlx.SqlConn) SysRoleDeptModel {
	return &customSysRoleDeptModel{
		defaultSysRoleDeptModel: newSysRoleDeptModel(conn),
	}
}

func (m *customSysRoleDeptModel) withSession(session sqlx.Session) SysRoleDeptModel {
	return NewSysRoleDeptModel(sqlx.NewSqlConnFromSession(session))
}
