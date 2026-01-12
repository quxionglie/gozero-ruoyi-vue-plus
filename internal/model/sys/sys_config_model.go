package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ConfigQuery 参数配置查询条件
type ConfigQuery struct {
	ConfigName string // 参数名称（模糊查询）
	ConfigKey  string // 参数键名（模糊查询）
	ConfigType string // 系统内置（Y是 N否）
}

// PageQuery 已移至 vars.go，这里保留类型别名以保持兼容性
// 实际使用请使用 vars.go 中的 PageQuery 和 Normalize 方法

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
	// 初始化分页参数的非合规值
	pageQuery.Normalize()

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

	// 构建排序（防止 SQL 注入）
	// 允许的排序列（支持 snake_case 和 camelCase）
	allowedOrderColumns := map[string]bool{
		"config_id":   true,
		"configId":    true,
		"config_name": true,
		"configName":  true,
		"config_key":  true,
		"configKey":   true,
		"config_type": true,
		"configType":  true,
		"create_time": true,
		"createTime":  true,
		"update_time": true,
		"updateTime":  true,
	}

	orderBy := "config_id"
	if pageQuery.OrderByColumn != "" {
		// 将 camelCase 转换为 snake_case
		columnName := camelToSnake(strings.TrimSpace(pageQuery.OrderByColumn))
		// 检查原始字段名和转换后的字段名是否在允许列表中
		originalColumn := strings.TrimSpace(pageQuery.OrderByColumn)
		if allowedOrderColumns[originalColumn] || allowedOrderColumns[columnName] {
			// 使用转换后的 snake_case 字段名
			orderBy = columnName
		}
	}

	// 处理排序方向（兼容 asc、desc、descending 等）
	orderDir := "asc"
	isAscStr := strings.ToLower(strings.TrimSpace(pageQuery.IsAsc))
	if isAscStr == "asc" || isAscStr == "ascending" {
		orderDir = "asc"
	} else if isAscStr == "desc" || isAscStr == "descending" {
		orderDir = "desc"
	}

	// 构建分页查询
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by %s %s", sysConfigRows, m.table, whereClause, orderBy, orderDir)
	if pageQuery.PageSize > 0 {
		offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
		if offset < 0 {
			offset = 0
		}
		sqlQuery += fmt.Sprintf(" limit %d, %d", offset, pageQuery.PageSize)
	}

	var resp []*SysConfig
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}
