package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// LogininforQuery 登录日志查询条件
type LogininforQuery struct {
	Ipaddr    string // IP地址（模糊查询）
	Status    string // 登录状态（0成功 1失败）
	UserName  string // 用户账号（模糊查询）
	BeginTime string // 开始时间
	EndTime   string // 结束时间
}

var _ SysLogininforModel = (*customSysLogininforModel)(nil)

type (
	// SysLogininforModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysLogininforModel.
	SysLogininforModel interface {
		sysLogininforModel
		withSession(session sqlx.Session) SysLogininforModel
		FindPage(ctx context.Context, query *LogininforQuery, pageQuery *PageQuery) ([]*SysLogininfor, int64, error)
		DeleteByIds(ctx context.Context, infoIds []int64) error
		Clean(ctx context.Context) error
	}

	customSysLogininforModel struct {
		*defaultSysLogininforModel
	}
)

// NewSysLogininforModel returns a model for the database table.
func NewSysLogininforModel(conn sqlx.SqlConn) SysLogininforModel {
	return &customSysLogininforModel{
		defaultSysLogininforModel: newSysLogininforModel(conn),
	}
}

func (m *customSysLogininforModel) withSession(session sqlx.Session) SysLogininforModel {
	return NewSysLogininforModel(sqlx.NewSqlConnFromSession(session))
}

// FindPage 分页查询登录日志（支持条件查询和分页）
func (m *customSysLogininforModel) FindPage(ctx context.Context, query *LogininforQuery, pageQuery *PageQuery) ([]*SysLogininfor, int64, error) {
	if query == nil {
		query = &LogininforQuery{}
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{}
	}
	// 初始化分页参数的非合规值
	pageQuery.Normalize()

	// 构建 WHERE 条件
	whereClause := "1=1"
	var args []interface{}

	if query.Ipaddr != "" {
		whereClause += " and ipaddr like ?"
		args = append(args, "%"+query.Ipaddr+"%")
	}
	if query.Status != "" {
		whereClause += " and status = ?"
		args = append(args, query.Status)
	}
	if query.UserName != "" {
		whereClause += " and user_name like ?"
		args = append(args, "%"+query.UserName+"%")
	}
	if query.BeginTime != "" && query.EndTime != "" {
		whereClause += " and login_time between ? and ?"
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
		"info_id":    true,
		"infoId":     true,
		"user_name":  true,
		"userName":   true,
		"ipaddr":     true,
		"status":     true,
		"login_time": true,
		"loginTime":  true,
	}

	orderBy := "info_id"
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
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by %s %s", sysLogininforRows, m.table, whereClause, orderBy, orderDir)
	if pageQuery.PageSize > 0 {
		offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
		if offset < 0 {
			offset = 0
		}
		sqlQuery += fmt.Sprintf(" limit %d, %d", offset, pageQuery.PageSize)
	}

	var resp []*SysLogininfor
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}

// DeleteByIds 批量删除登录日志
func (m *customSysLogininforModel) DeleteByIds(ctx context.Context, infoIds []int64) error {
	if len(infoIds) == 0 {
		return nil
	}
	placeholders := strings.Repeat("?,", len(infoIds))
	placeholders = placeholders[:len(placeholders)-1]
	query := fmt.Sprintf("delete from %s where info_id in (%s)", m.table, placeholders)
	args := make([]interface{}, len(infoIds))
	for i, id := range infoIds {
		args[i] = id
	}
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

// Clean 清空所有登录日志
func (m *customSysLogininforModel) Clean(ctx context.Context) error {
	query := fmt.Sprintf("delete from %s", m.table)
	_, err := m.conn.ExecCtx(ctx, query)
	return err
}
