package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// OperLogQuery 操作日志查询条件
type OperLogQuery struct {
	OperIp       string // 主机地址（模糊查询）
	Title        string // 模块标题（模糊查询）
	BusinessType int32  // 业务类型（0其它 1新增 2修改 3删除）
	Status       int32  // 操作状态（0正常 1异常）
	OperName     string // 操作人员（模糊查询）
	BeginTime    string // 开始时间
	EndTime      string // 结束时间
}

var _ SysOperLogModel = (*customSysOperLogModel)(nil)

type (
	// SysOperLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysOperLogModel.
	SysOperLogModel interface {
		sysOperLogModel
		withSession(session sqlx.Session) SysOperLogModel
		FindPage(ctx context.Context, query *OperLogQuery, pageQuery *PageQuery) ([]*SysOperLog, int64, error)
		DeleteByIds(ctx context.Context, operIds []int64) error
		Clean(ctx context.Context) error
	}

	customSysOperLogModel struct {
		*defaultSysOperLogModel
	}
)

// NewSysOperLogModel returns a model for the database table.
func NewSysOperLogModel(conn sqlx.SqlConn) SysOperLogModel {
	return &customSysOperLogModel{
		defaultSysOperLogModel: newSysOperLogModel(conn),
	}
}

func (m *customSysOperLogModel) withSession(session sqlx.Session) SysOperLogModel {
	return NewSysOperLogModel(sqlx.NewSqlConnFromSession(session))
}

// FindPage 分页查询操作日志（支持条件查询和分页）
func (m *customSysOperLogModel) FindPage(ctx context.Context, query *OperLogQuery, pageQuery *PageQuery) ([]*SysOperLog, int64, error) {
	if query == nil {
		query = &OperLogQuery{}
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{}
	}
	// 初始化分页参数的非合规值
	pageQuery.Normalize()

	// 构建 WHERE 条件
	whereClause := "1=1"
	var args []interface{}

	if query.OperIp != "" {
		whereClause += " and oper_ip like ?"
		args = append(args, "%"+query.OperIp+"%")
	}
	if query.Title != "" {
		whereClause += " and title like ?"
		args = append(args, "%"+query.Title+"%")
	}
	if query.BusinessType > 0 {
		whereClause += " and business_type = ?"
		args = append(args, query.BusinessType)
	}
	if query.Status >= 0 {
		whereClause += " and status = ?"
		args = append(args, query.Status)
	}
	if query.OperName != "" {
		whereClause += " and oper_name like ?"
		args = append(args, "%"+query.OperName+"%")
	}
	if query.BeginTime != "" && query.EndTime != "" {
		whereClause += " and oper_time between ? and ?"
		args = append(args, query.BeginTime, query.EndTime)
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
		"oper_id":   true,
		"operId":    true,
		"title":     true,
		"oper_ip":   true,
		"operIp":    true,
		"status":    true,
		"oper_time": true,
		"operTime":  true,
		"cost_time": true,
		"costTime":  true,
	}

	orderBy := "oper_id"
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
	orderDir := "desc"
	isAscStr := strings.ToLower(strings.TrimSpace(pageQuery.IsAsc))
	if isAscStr == "asc" || isAscStr == "ascending" {
		orderDir = "asc"
	} else if isAscStr == "desc" || isAscStr == "descending" {
		orderDir = "desc"
	}

	// 构建分页查询
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by %s %s", sysOperLogRows, m.table, whereClause, orderBy, orderDir)
	if pageQuery.PageSize > 0 {
		offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
		if offset < 0 {
			offset = 0
		}
		sqlQuery += fmt.Sprintf(" limit %d, %d", offset, pageQuery.PageSize)
	}

	var resp []*SysOperLog
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}

// DeleteByIds 批量删除操作日志
func (m *customSysOperLogModel) DeleteByIds(ctx context.Context, operIds []int64) error {
	if len(operIds) == 0 {
		return nil
	}
	placeholders := strings.Repeat("?,", len(operIds))
	placeholders = placeholders[:len(placeholders)-1]
	query := fmt.Sprintf("delete from %s where oper_id in (%s)", m.table, placeholders)
	args := make([]interface{}, len(operIds))
	for i, id := range operIds {
		args[i] = id
	}
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

// Clean 清空所有操作日志
func (m *customSysOperLogModel) Clean(ctx context.Context) error {
	query := fmt.Sprintf("delete from %s", m.table)
	_, err := m.conn.ExecCtx(ctx, query)
	return err
}
