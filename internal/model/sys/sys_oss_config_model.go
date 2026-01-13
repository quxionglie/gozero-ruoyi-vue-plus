package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// OssConfigQuery OSS配置查询条件
type OssConfigQuery struct {
	ConfigKey string // 配置key（模糊查询）
	Status    string // 状态（0是 1否）
}

var _ SysOssConfigModel = (*customSysOssConfigModel)(nil)

type (
	// SysOssConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysOssConfigModel.
	SysOssConfigModel interface {
		sysOssConfigModel
		withSession(session sqlx.Session) SysOssConfigModel
		FindPage(ctx context.Context, query *OssConfigQuery, pageQuery *PageQuery) ([]*SysOssConfig, int64, error)
		UpdateStatus(ctx context.Context, ossConfigId int64, status string) error
		FindDefault(ctx context.Context, tenantId string) (*SysOssConfig, error)
		FindByConfigKey(ctx context.Context, configKey string, tenantId string) (*SysOssConfig, error)
		UpdateById(ctx context.Context, data *SysOssConfig) (int64, error)
	}

	customSysOssConfigModel struct {
		*defaultSysOssConfigModel
	}
)

// NewSysOssConfigModel returns a model for the database table.
func NewSysOssConfigModel(conn sqlx.SqlConn) SysOssConfigModel {
	return &customSysOssConfigModel{
		defaultSysOssConfigModel: newSysOssConfigModel(conn),
	}
}

func (m *customSysOssConfigModel) withSession(session sqlx.Session) SysOssConfigModel {
	return NewSysOssConfigModel(sqlx.NewSqlConnFromSession(session))
}

// FindPage 分页查询OSS对象存储配置列表（支持条件查询和分页）
func (m *customSysOssConfigModel) FindPage(ctx context.Context, query *OssConfigQuery, pageQuery *PageQuery) ([]*SysOssConfig, int64, error) {
	if query == nil {
		query = &OssConfigQuery{}
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{
			PageNum:  1,
			PageSize: 10,
		}
	} else {
		// 初始化分页参数的非合规值
		pageQuery.Normalize()
	}

	// 构建 WHERE 条件
	whereClause := "1=1"
	var args []interface{}

	if query.ConfigKey != "" {
		whereClause += " and config_key LIKE ?"
		args = append(args, "%"+query.ConfigKey+"%")
	}
	if query.Status != "" {
		whereClause += " and status = ?"
		args = append(args, query.Status)
	}

	// 计算总数
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)
	var total int64
	err := m.conn.QueryRowPartialCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 构建 ORDER BY 子句（防止 SQL 注入）
	// 允许的排序列（支持 snake_case 和 camelCase）
	allowedOrderColumns := buildAllowedOrderColumns(sysOssConfigFieldNames)
	orderBy := pageQuery.GetOrderBy("oss_config_id", allowedOrderColumns)

	// 获取排序方向（默认降序）
	orderDir := pageQuery.GetOrderDir("desc")

	// 计算分页参数
	offset := pageQuery.GetOffset()

	// 查询数据
	querySQL := fmt.Sprintf("select %s from %s where %s order by %s %s limit %d, %d",
		sysOssConfigRows, m.table, whereClause, orderBy, orderDir, offset, pageQuery.PageSize)

	var ossConfigList []*SysOssConfig
	err = m.conn.QueryRowsPartialCtx(ctx, &ossConfigList, querySQL, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	return ossConfigList, total, nil
}

// UpdateStatus 更新OSS配置状态
func (m *customSysOssConfigModel) UpdateStatus(ctx context.Context, ossConfigId int64, status string) error {
	query := fmt.Sprintf("update %s set `status` = ? where `oss_config_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, status, ossConfigId)
	return err
}

// FindDefault 查找默认OSS配置（status=0）
func (m *customSysOssConfigModel) FindDefault(ctx context.Context, tenantId string) (*SysOssConfig, error) {
	whereClause := "status = '0'"
	args := []interface{}{}

	if tenantId != "" {
		whereClause += " and tenant_id = ?"
		args = append(args, tenantId)
	}

	query := fmt.Sprintf("select %s from %s where %s limit 1", sysOssConfigRows, m.table, whereClause)
	var resp SysOssConfig
	err := m.conn.QueryRowCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindByConfigKey 根据配置键查询OSS配置
func (m *customSysOssConfigModel) FindByConfigKey(ctx context.Context, configKey string, tenantId string) (*SysOssConfig, error) {
	whereClause := "config_key = ?"
	args := []interface{}{configKey}

	if tenantId != "" {
		whereClause += " and tenant_id = ?"
		args = append(args, tenantId)
	}

	query := fmt.Sprintf("select %s from %s where %s limit 1", sysOssConfigRows, m.table, whereClause)
	var resp SysOssConfig
	err := m.conn.QueryRowCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// UpdateById 根据ID更新OSS配置，只更新非零值字段
func (m *customSysOssConfigModel) UpdateById(ctx context.Context, data *SysOssConfig) (int64, error) {
	if data.OssConfigId == 0 {
		return 0, fmt.Errorf("oss_config_id cannot be zero")
	}

	var setParts []string
	var args []interface{}

	// 检查每个字段是否为非零值，如果是则加入更新列表
	if data.TenantId != "" {
		setParts = append(setParts, "`tenant_id` = ?")
		args = append(args, data.TenantId)
	}
	if data.ConfigKey != "" {
		setParts = append(setParts, "`config_key` = ?")
		args = append(args, data.ConfigKey)
	}
	if data.AccessKey != "" {
		setParts = append(setParts, "`access_key` = ?")
		args = append(args, data.AccessKey)
	}
	if data.SecretKey != "" {
		setParts = append(setParts, "`secret_key` = ?")
		args = append(args, data.SecretKey)
	}
	if data.BucketName != "" {
		setParts = append(setParts, "`bucket_name` = ?")
		args = append(args, data.BucketName)
	}
	if data.Prefix != "" {
		setParts = append(setParts, "`prefix` = ?")
		args = append(args, data.Prefix)
	}
	if data.Endpoint != "" {
		setParts = append(setParts, "`endpoint` = ?")
		args = append(args, data.Endpoint)
	}
	if data.Domain != "" {
		setParts = append(setParts, "`domain` = ?")
		args = append(args, data.Domain)
	}
	if data.IsHttps != "" {
		setParts = append(setParts, "`is_https` = ?")
		args = append(args, data.IsHttps)
	}
	if data.Region != "" {
		setParts = append(setParts, "`region` = ?")
		args = append(args, data.Region)
	}
	if data.AccessPolicy != "" {
		setParts = append(setParts, "`access_policy` = ?")
		args = append(args, data.AccessPolicy)
	}
	if data.Status != "" {
		setParts = append(setParts, "`status` = ?")
		args = append(args, data.Status)
	}
	if data.Ext1 != "" {
		setParts = append(setParts, "`ext1` = ?")
		args = append(args, data.Ext1)
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
		return 0, nil // 没有需要更新的字段
	}

	// 构建更新SQL
	setClause := strings.Join(setParts, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE `oss_config_id` = ?", m.table, setClause)
	args = append(args, data.OssConfigId)

	result, err := m.conn.ExecCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
