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
		UpdateById(ctx context.Context, data *SysNotice) (int64, error)
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
	// 初始化分页参数的非合规值
	pageQuery.Normalize()

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
	allowedOrderColumns := buildAllowedOrderColumns(sysNoticeFieldNames)
	orderBy := pageQuery.GetOrderBy("notice_id", allowedOrderColumns)

	// 获取排序方向（默认升序）
	orderDir := pageQuery.GetOrderDir("asc")

	// 构建分页查询
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by %s %s", sysNoticeRows, m.table, whereClause, orderBy, orderDir)
	if pageQuery.PageSize > 0 {
		offset := pageQuery.GetOffset()
		sqlQuery += fmt.Sprintf(" limit %d, %d", offset, pageQuery.PageSize)
	}

	var resp []*SysNotice
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}

// UpdateById 根据ID更新通知公告，只更新非零值字段
func (m *customSysNoticeModel) UpdateById(ctx context.Context, data *SysNotice) (int64, error) {
	if data.NoticeId == 0 {
		return 0, fmt.Errorf("notice_id cannot be zero")
	}

	var setParts []string
	var args []interface{}

	// 检查每个字段是否为非零值，如果是则加入更新列表
	if data.TenantId != "" {
		setParts = append(setParts, "`tenant_id` = ?")
		args = append(args, data.TenantId)
	}
	if data.NoticeTitle != "" {
		setParts = append(setParts, "`notice_title` = ?")
		args = append(args, data.NoticeTitle)
	}
	if data.NoticeType != "" {
		setParts = append(setParts, "`notice_type` = ?")
		args = append(args, data.NoticeType)
	}
	if data.NoticeContent.Valid {
		setParts = append(setParts, "`notice_content` = ?")
		args = append(args, data.NoticeContent.String)
	}
	if data.Status != "" {
		setParts = append(setParts, "`status` = ?")
		args = append(args, data.Status)
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
	query := fmt.Sprintf("UPDATE %s SET %s WHERE `notice_id` = ?", m.table, setClause)
	args = append(args, data.NoticeId)

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
