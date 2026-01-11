package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// NoticeQuery 通知公告查询条件
type NoticeQuery struct {
	NoticeTitle string // 公告标题（模糊查询）
	NoticeType  string // 公告类型（1通知 2公告）
	CreateBy    int64  // 创建者ID（通过用户名查询）
}

var _ SysNoticeModel = (*customSysNoticeModel)(nil)

type (
	// SysNoticeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysNoticeModel.
	SysNoticeModel interface {
		sysNoticeModel
		withSession(session sqlx.Session) SysNoticeModel
		FindPage(ctx context.Context, query *NoticeQuery, pageQuery *PageQuery) ([]*SysNotice, int64, error)
	}

	customSysNoticeModel struct {
		*defaultSysNoticeModel
	}
)

// NewSysNoticeModel returns a model for the database table.
func NewSysNoticeModel(conn sqlx.SqlConn) SysNoticeModel {
	return &customSysNoticeModel{
		defaultSysNoticeModel: newSysNoticeModel(conn),
	}
}

func (m *customSysNoticeModel) withSession(session sqlx.Session) SysNoticeModel {
	return NewSysNoticeModel(sqlx.NewSqlConnFromSession(session))
}

// FindPage 分页查询通知公告（支持条件查询和分页）
func (m *customSysNoticeModel) FindPage(ctx context.Context, query *NoticeQuery, pageQuery *PageQuery) ([]*SysNotice, int64, error) {
	if query == nil {
		query = &NoticeQuery{}
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{}
	}

	// 构建 WHERE 条件
	whereClause := "1=1"
	var args []interface{}

	if query.NoticeTitle != "" {
		whereClause += " and notice_title like ?"
		args = append(args, "%"+query.NoticeTitle+"%")
	}
	if query.NoticeType != "" {
		whereClause += " and notice_type = ?"
		args = append(args, query.NoticeType)
	}
	if query.CreateBy > 0 {
		whereClause += " and create_by = ?"
		args = append(args, query.CreateBy)
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
		"notice_id":    true,
		"noticeId":     true,
		"notice_title": true,
		"noticeTitle":  true,
		"notice_type":  true,
		"noticeType":   true,
		"status":       true,
		"create_time":  true,
		"createTime":   true,
		"update_time":  true,
		"updateTime":   true,
	}

	orderBy := "notice_id"
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
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by %s %s", sysNoticeRows, m.table, whereClause, orderBy, orderDir)
	if pageQuery.PageSize > 0 {
		offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
		if offset < 0 {
			offset = 0
		}
		sqlQuery += fmt.Sprintf(" limit %d, %d", offset, pageQuery.PageSize)
	}

	var resp []*SysNotice
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}
