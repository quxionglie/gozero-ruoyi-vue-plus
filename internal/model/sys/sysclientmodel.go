package sys

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysClientModel = (*customSysClientModel)(nil)

type (
	// SysClientModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysClientModel.
	SysClientModel interface {
		sysClientModel
		withSession(session sqlx.Session) SysClientModel
		FindOneByClientId(ctx context.Context, clientId string) (*SysClient, error)
	}

	customSysClientModel struct {
		*defaultSysClientModel
	}
)

// NewSysClientModel returns a model for the database table.
func NewSysClientModel(conn sqlx.SqlConn) SysClientModel {
	return &customSysClientModel{
		defaultSysClientModel: newSysClientModel(conn),
	}
}

func (m *customSysClientModel) withSession(session sqlx.Session) SysClientModel {
	return NewSysClientModel(sqlx.NewSqlConnFromSession(session))
}

// FindOneByClientId 根据客户端ID查询
func (m *customSysClientModel) FindOneByClientId(ctx context.Context, clientId string) (*SysClient, error) {
	query := "select `id`,`client_id`,`client_key`,`client_secret`,`grant_type`,`device_type`,`active_timeout`,`timeout`,`status`,`del_flag`,`create_dept`,`create_by`,`create_time`,`update_by`,`update_time` from `sys_client` where `client_id` = ? and `del_flag` = '0' limit 1"
	var resp SysClient
	err := m.conn.QueryRowCtx(ctx, &resp, query, clientId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
