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
	allowedOrderColumns := map[string]bool{
		"oss_config_id": true, "ossConfigId": true,
		"config_key": true, "configKey": true,
		"status":      true,
		"create_time": true, "createTime": true,
		"update_time": true, "updateTime": true,
	}

	orderBy := "oss_config_id"
	if pageQuery.OrderByColumn != "" {
		columnName := camelToSnake(strings.TrimSpace(pageQuery.OrderByColumn))
		originalColumn := strings.TrimSpace(pageQuery.OrderByColumn)
		if allowedOrderColumns[originalColumn] || allowedOrderColumns[columnName] {
			orderBy = columnName
		}
	}

	orderDir := "desc"
	isAscStr := strings.ToLower(strings.TrimSpace(pageQuery.IsAsc))
	if isAscStr == "asc" || isAscStr == "ascending" {
		orderDir = "asc"
	} else if isAscStr == "desc" || isAscStr == "descending" {
		orderDir = "desc"
	}

	// 计算分页参数
	offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
	if offset < 0 {
		offset = 0
	}

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
