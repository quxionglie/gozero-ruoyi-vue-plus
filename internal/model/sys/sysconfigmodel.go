package sys

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysConfigModel = (*customSysConfigModel)(nil)

type (
	// SysConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysConfigModel.
	SysConfigModel interface {
		sysConfigModel
		withSession(session sqlx.Session) SysConfigModel
		FindAll(ctx context.Context) ([]*SysConfig, error)
		FindByConfigKey(ctx context.Context, configKey string) (*SysConfig, error)
		CheckConfigKeyUnique(ctx context.Context, configKey string, excludeConfigId int64) (bool, error)
	}

	customSysConfigModel struct {
		*defaultSysConfigModel
	}
)

// NewSysConfigModel returns a model for the database table.
func NewSysConfigModel(conn sqlx.SqlConn) SysConfigModel {
	return &customSysConfigModel{
		defaultSysConfigModel: newSysConfigModel(conn),
	}
}

func (m *customSysConfigModel) withSession(session sqlx.Session) SysConfigModel {
	return NewSysConfigModel(sqlx.NewSqlConnFromSession(session))
}

// FindAll 查询所有参数配置
func (m *customSysConfigModel) FindAll(ctx context.Context) ([]*SysConfig, error) {
	query := fmt.Sprintf("select %s from %s order by config_id asc", sysConfigRows, m.table)
	var resp []*SysConfig
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return resp, nil
}

// FindByConfigKey 根据参数键名查询参数配置
func (m *customSysConfigModel) FindByConfigKey(ctx context.Context, configKey string) (*SysConfig, error) {
	query := fmt.Sprintf("select %s from %s where config_key = ? limit 1", sysConfigRows, m.table)
	var resp SysConfig
	err := m.conn.QueryRowCtx(ctx, &resp, query, configKey)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// CheckConfigKeyUnique 校验参数键名唯一性
func (m *customSysConfigModel) CheckConfigKeyUnique(ctx context.Context, configKey string, excludeConfigId int64) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where config_key = ?", m.table)
	if excludeConfigId > 0 {
		query += " and config_id != ?"
	}

	var count int64
	var err error
	if excludeConfigId > 0 {
		err = m.conn.QueryRowPartialCtx(ctx, &count, query, configKey, excludeConfigId)
	} else {
		err = m.conn.QueryRowPartialCtx(ctx, &count, query, configKey)
	}
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
