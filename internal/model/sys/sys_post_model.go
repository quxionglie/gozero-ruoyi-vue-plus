package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// PostQuery 岗位查询条件
type PostQuery struct {
	PostCode     string // 岗位编码（模糊查询）
	PostCategory string // 岗位类别编码（模糊查询）
	PostName     string // 岗位名称（模糊查询）
	Status       string // 状态（0正常 1停用）
	DeptId       int64  // 部门id（单部门）
	BelongDeptId int64  // 归属部门id（部门树，需要查询子部门）
}

var _ SysPostModel = (*customSysPostModel)(nil)

type (
	// SysPostModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysPostModel.
	SysPostModel interface {
		sysPostModel
		withSession(session sqlx.Session) SysPostModel
		FindPage(ctx context.Context, query *PostQuery, pageQuery *PageQuery) ([]*SysPost, int64, error)
		FindAll(ctx context.Context, query *PostQuery) ([]*SysPost, error)
		FindByIds(ctx context.Context, postIds []int64) ([]*SysPost, error)
		CheckPostNameUnique(ctx context.Context, postName string, deptId int64, excludePostId int64) (bool, error)
		CheckPostCodeUnique(ctx context.Context, postCode string, excludePostId int64) (bool, error)
		CountUserPostById(ctx context.Context, postId int64) (int64, error)
		CountPostByDeptId(ctx context.Context, deptId int64) (int64, error)
	}

	customSysPostModel struct {
		*defaultSysPostModel
	}
)

// NewSysPostModel returns a model for the database table.
func NewSysPostModel(conn sqlx.SqlConn) SysPostModel {
	return &customSysPostModel{
		defaultSysPostModel: newSysPostModel(conn),
	}
}

func (m *customSysPostModel) withSession(session sqlx.Session) SysPostModel {
	return NewSysPostModel(sqlx.NewSqlConnFromSession(session))
}

// FindPage 分页查询岗位（支持条件查询和分页）
func (m *customSysPostModel) FindPage(ctx context.Context, query *PostQuery, pageQuery *PageQuery) ([]*SysPost, int64, error) {
	if query == nil {
		query = &PostQuery{}
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{}
	}
	// 初始化分页参数的非合规值
	pageQuery.Normalize()

	// 构建 WHERE 条件
	whereClause := "1=1"
	var args []interface{}

	if query.PostCode != "" {
		whereClause += " and post_code like ?"
		args = append(args, "%"+query.PostCode+"%")
	}
	if query.PostCategory != "" {
		whereClause += " and post_category like ?"
		args = append(args, "%"+query.PostCategory+"%")
	}
	if query.PostName != "" {
		whereClause += " and post_name like ?"
		args = append(args, "%"+query.PostName+"%")
	}
	if query.Status != "" {
		whereClause += " and status = ?"
		args = append(args, query.Status)
	}
	if query.DeptId > 0 {
		whereClause += " and dept_id = ?"
		args = append(args, query.DeptId)
	} else if query.BelongDeptId > 0 {
		// 部门树查询：需要查询所有子部门（这里简化处理，实际应该递归查询子部门）
		// TODO: 如果需要完整的部门树查询，需要调用 SysDeptModel 的方法
		whereClause += " and dept_id = ?"
		args = append(args, query.BelongDeptId)
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
		"post_id":     true,
		"postId":      true,
		"post_code":   true,
		"postCode":    true,
		"post_name":   true,
		"postName":    true,
		"post_sort":   true,
		"postSort":    true,
		"status":      true,
		"create_time": true,
		"createTime":  true,
		"update_time": true,
		"updateTime":  true,
	}

	orderBy := "post_sort"
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
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by %s %s", sysPostRows, m.table, whereClause, orderBy, orderDir)
	if pageQuery.PageSize > 0 {
		offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
		if offset < 0 {
			offset = 0
		}
		sqlQuery += fmt.Sprintf(" limit %d, %d", offset, pageQuery.PageSize)
	}

	var resp []*SysPost
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}

// FindAll 查询所有岗位（不分页）
func (m *customSysPostModel) FindAll(ctx context.Context, query *PostQuery) ([]*SysPost, error) {
	if query == nil {
		query = &PostQuery{}
	}

	// 构建 WHERE 条件（与 FindPage 相同逻辑）
	whereClause := "1=1"
	var args []interface{}

	if query.PostCode != "" {
		whereClause += " and post_code like ?"
		args = append(args, "%"+query.PostCode+"%")
	}
	if query.PostCategory != "" {
		whereClause += " and post_category like ?"
		args = append(args, "%"+query.PostCategory+"%")
	}
	if query.PostName != "" {
		whereClause += " and post_name like ?"
		args = append(args, "%"+query.PostName+"%")
	}
	if query.Status != "" {
		whereClause += " and status = ?"
		args = append(args, query.Status)
	}
	if query.DeptId > 0 {
		whereClause += " and dept_id = ?"
		args = append(args, query.DeptId)
	} else if query.BelongDeptId > 0 {
		whereClause += " and dept_id = ?"
		args = append(args, query.BelongDeptId)
	}

	sqlQuery := fmt.Sprintf("select %s from %s where %s order by post_sort asc", sysPostRows, m.table, whereClause)
	var resp []*SysPost
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return resp, nil
}

// FindByIds 根据岗位ID列表查询岗位
func (m *customSysPostModel) FindByIds(ctx context.Context, postIds []int64) ([]*SysPost, error) {
	if len(postIds) == 0 {
		return []*SysPost{}, nil
	}

	// 构建 IN 查询
	placeholders := ""
	for i := 0; i < len(postIds); i++ {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
	}

	query := fmt.Sprintf("select %s from %s where post_id in (%s) and status = '0' order by post_sort asc", sysPostRows, m.table, placeholders)
	var args []interface{}
	for _, id := range postIds {
		args = append(args, id)
	}

	var resp []*SysPost
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return resp, nil
}

// CheckPostNameUnique 校验岗位名称唯一性（同部门内唯一）
func (m *customSysPostModel) CheckPostNameUnique(ctx context.Context, postName string, deptId int64, excludePostId int64) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where post_name = ? and dept_id = ?", m.table)
	var args []interface{}
	args = append(args, postName, deptId)

	if excludePostId > 0 {
		query += " and post_id != ?"
		args = append(args, excludePostId)
	}

	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// CheckPostCodeUnique 校验岗位编码唯一性（全局唯一）
func (m *customSysPostModel) CheckPostCodeUnique(ctx context.Context, postCode string, excludePostId int64) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where post_code = ?", m.table)
	var args []interface{}
	args = append(args, postCode)

	if excludePostId > 0 {
		query += " and post_id != ?"
		args = append(args, excludePostId)
	}

	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// CountUserPostById 统计岗位使用数量
func (m *customSysPostModel) CountUserPostById(ctx context.Context, postId int64) (int64, error) {
	query := "select count(*) from `sys_user_post` where post_id = ?"
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, postId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountPostByDeptId 根据部门ID统计岗位数量
func (m *customSysPostModel) CountPostByDeptId(ctx context.Context, deptId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where dept_id = ?", m.table)
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, deptId)
	if err != nil {
		return 0, err
	}
	return count, nil
}
