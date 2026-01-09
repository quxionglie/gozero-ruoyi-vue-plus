package sys

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysUserModel = (*customSysUserModel)(nil)

type (
	// SysUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysUserModel.
	SysUserModel interface {
		sysUserModel
		withSession(session sqlx.Session) SysUserModel
		FindOneByUserName(ctx context.Context, userName, tenantId string) (*SysUser, error)
	}

	customSysUserModel struct {
		*defaultSysUserModel
	}
)

// NewSysUserModel returns a model for the database table.
func NewSysUserModel(conn sqlx.SqlConn) SysUserModel {
	return &customSysUserModel{
		defaultSysUserModel: newSysUserModel(conn),
	}
}

func (m *customSysUserModel) withSession(session sqlx.Session) SysUserModel {
	return NewSysUserModel(sqlx.NewSqlConnFromSession(session))
}

// FindOneByUserName 根据用户名和租户ID查询用户
func (m *customSysUserModel) FindOneByUserName(ctx context.Context, userName, tenantId string) (*SysUser, error) {
	query := "select `user_id`,`tenant_id`,`dept_id`,`user_name`,`nick_name`,`user_type`,`email`,`phonenumber`,`sex`,`avatar`,`password`,`status`,`del_flag`,`login_ip`,`login_date`,`create_dept`,`create_by`,`create_time`,`update_by`,`update_time`,`remark` from `sys_user` where `user_name` = ? and `tenant_id` = ? and `del_flag` = '0' limit 1"
	var resp SysUser
	err := m.conn.QueryRowCtx(ctx, &resp, query, userName, tenantId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
