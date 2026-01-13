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
		UpdateById(ctx context.Context, data *SysConfig) error
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
	allowedOrderColumns := buildAllowedOrderColumns(sysConfigFieldNames)
	orderBy := pageQuery.GetOrderBy("config_id", allowedOrderColumns)

	// 获取排序方向（默认升序）
	orderDir := pageQuery.GetOrderDir("asc")

	// 构建分页查询
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by %s %s", sysConfigRows, m.table, whereClause, orderBy, orderDir)
	if pageQuery.PageSize > 0 {
		offset := pageQuery.GetOffset()
		sqlQuery += fmt.Sprintf(" limit %d, %d", offset, pageQuery.PageSize)
	}

	var resp []*SysConfig
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}

// UpdateById 根据ID更新配置，只更新非零值字段
func (m *customSysConfigModel) UpdateById(ctx context.Context, data *SysConfig) error {
	if data.ConfigId == 0 {
		return fmt.Errorf("config_id cannot be zero")
	}

	var setParts []string
	var args []interface{}

	// 检查每个字段是否为非零值，如果是则加入更新列表
	if data.TenantId != "" {
		setParts = append(setParts, "`tenant_id` = ?")
		args = append(args, data.TenantId)
	}
	if data.ConfigName != "" {
		setParts = append(setParts, "`config_name` = ?")
		args = append(args, data.ConfigName)
	}
	if data.ConfigKey != "" {
		setParts = append(setParts, "`config_key` = ?")
		args = append(args, data.ConfigKey)
	}
	if data.ConfigValue != "" {
		setParts = append(setParts, "`config_value` = ?")
		args = append(args, data.ConfigValue)
	}
	if data.ConfigType != "" {
		setParts = append(setParts, "`config_type` = ?")
		args = append(args, data.ConfigType)
	}
	if data.CreateDept.Valid {
		setParts = append(setParts, "`create_dept` = ?")
		args = append(args, data.CreateDept.Int64)
	}
	if data.CreateBy.Valid {
		setParts = append(setParts, "`create_by` = ?")
		args = append(args, data.CreateBy.Int64)
	}
	if data.CreateTime.Valid {
		setParts = append(setParts, "`create_time` = ?")
		args = append(args, data.CreateTime.Time)
	}
	if data.UpdateBy.Valid {
		setParts = append(setParts, "`update_by` = ?")
		args = append(args, data.UpdateBy.Int64)
	}
	if data.UpdateTime.Valid {
		setParts = append(setParts, "`update_time` = ?")
		args = append(args, data.UpdateTime.Time)
	}
	if data.Remark.Valid {
		setParts = append(setParts, "`remark` = ?")
		args = append(args, data.Remark.String)
	}

	if len(setParts) == 0 {
		return nil // 没有需要更新的字段
	}

	// 构建更新SQL
	setClause := strings.Join(setParts, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE `config_id` = ?", m.table, setClause)
	args = append(args, data.ConfigId)

	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}
