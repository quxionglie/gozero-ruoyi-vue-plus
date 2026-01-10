package sys

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ConfigQuery 参数配置查询条件
type ConfigQuery struct {
	ConfigName string // 参数名称（模糊查询）
	ConfigKey  string // 参数键名（模糊查询）
	ConfigType string // 系统内置（Y是 N否）
}

// PageQuery 分页查询参数
type PageQuery struct {
	PageNum       int32  // 当前页数
	PageSize      int32  // 分页大小
	OrderByColumn string // 排序列
	IsAsc         string // 排序方向（desc 或 asc）
}

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
		FindPage(ctx context.Context, query *ConfigQuery, pageQuery *PageQuery) ([]*SysConfig, int64, error)
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

// FindPage 分页查询参数配置（支持条件查询和分页）
func (m *customSysConfigModel) FindPage(ctx context.Context, query *ConfigQuery, pageQuery *PageQuery) ([]*SysConfig, int64, error) {
	if query == nil {
		query = &ConfigQuery{}
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{}
	}

	// 构建 WHERE 条件
	whereClause := "1=1"
	var args []interface{}

	if query.ConfigName != "" {
		whereClause += " and config_name like ?"
		args = append(args, "%"+query.ConfigName+"%")
	}
	if query.ConfigKey != "" {
		whereClause += " and config_key like ?"
		args = append(args, "%"+query.ConfigKey+"%")
	}
	if query.ConfigType != "" {
		whereClause += " and config_type = ?"
		args = append(args, query.ConfigType)
	}

	// 查询总数
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)
	var total int64
	err := m.conn.QueryRowPartialCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 构建排序
	orderBy := "config_id"
	if pageQuery.OrderByColumn != "" {
		orderBy = pageQuery.OrderByColumn
	}
	orderDir := "asc"
	if pageQuery.IsAsc == "desc" {
		orderDir = "desc"
	}

	// 构建分页查询
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by %s %s", sysConfigRows, m.table, whereClause, orderBy, orderDir)
	if pageQuery.PageSize > 0 {
		offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
		if offset < 0 {
			offset = 0
		}
		sqlQuery += fmt.Sprintf(" limit %d offset %d", pageQuery.PageSize, offset)
	}

	var resp []*SysConfig
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}
