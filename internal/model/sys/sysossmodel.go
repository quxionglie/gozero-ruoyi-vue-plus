package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// OssQuery OSS查询条件
type OssQuery struct {
	FileName     string // 文件名（模糊查询）
	OriginalName string // 原名（模糊查询）
	FileSuffix   string // 文件后缀名
	Url          string // URL地址
	CreateBy     int64  // 创建者
	Service      string // 服务商
}

var _ SysOssModel = (*customSysOssModel)(nil)

type (
	// SysOssModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysOssModel.
	SysOssModel interface {
		sysOssModel
		withSession(session sqlx.Session) SysOssModel
		FindPage(ctx context.Context, query *OssQuery, pageQuery *PageQuery) ([]*SysOss, int64, error)
		FindByIds(ctx context.Context, ossIds []int64) ([]*SysOss, error)
	}

	customSysOssModel struct {
		*defaultSysOssModel
	}
)

// NewSysOssModel returns a model for the database table.
func NewSysOssModel(conn sqlx.SqlConn) SysOssModel {
	return &customSysOssModel{
		defaultSysOssModel: newSysOssModel(conn),
	}
}

func (m *customSysOssModel) withSession(session sqlx.Session) SysOssModel {
	return NewSysOssModel(sqlx.NewSqlConnFromSession(session))
}

// FindPage 分页查询OSS对象存储列表
func (m *customSysOssModel) FindPage(ctx context.Context, query *OssQuery, pageQuery *PageQuery) ([]*SysOss, int64, error) {
	if query == nil {
		query = &OssQuery{}
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

	if query.FileName != "" {
		whereClause += " and file_name LIKE ?"
		args = append(args, "%"+query.FileName+"%")
	}
	if query.OriginalName != "" {
		whereClause += " and original_name LIKE ?"
		args = append(args, "%"+query.OriginalName+"%")
	}
	if query.FileSuffix != "" {
		whereClause += " and file_suffix = ?"
		args = append(args, query.FileSuffix)
	}
	if query.Url != "" {
		whereClause += " and url = ?"
		args = append(args, query.Url)
	}
	if query.CreateBy > 0 {
		whereClause += " and create_by = ?"
		args = append(args, query.CreateBy)
	}
	if query.Service != "" {
		whereClause += " and service = ?"
		args = append(args, query.Service)
	}

	// 计算总数
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)
	var total int64
	err := m.conn.QueryRowPartialCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 构建 ORDER BY 子句
	orderBy := "oss_id ASC"
	if pageQuery.OrderByColumn != "" {
		// 列名白名单验证（防止SQL注入）
		ossOrderColumns := map[string]bool{
			"oss_id": true, "create_time": true, "file_name": true, "original_name": true,
			"file_suffix": true, "service": true,
		}
		orderColumn := strings.ToLower(pageQuery.OrderByColumn)
		if ossOrderColumns[orderColumn] {
			orderBy = pageQuery.OrderByColumn + " "
			if pageQuery.IsAsc == "ascending" || pageQuery.IsAsc == "asc" {
				orderBy += "ASC"
			} else {
				orderBy += "DESC"
			}
		}
	}

	// 计算分页参数
	offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
	limit := pageQuery.PageSize

	// 查询数据
	querySQL := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, sysOssRows, m.table, whereClause, orderBy)

	args = append(args, limit, offset)

	var ossList []*SysOss
	err = m.conn.QueryRowsPartialCtx(ctx, &ossList, querySQL, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	return ossList, total, nil
}

// FindByIds 根据OSS ID列表查询OSS对象
func (m *customSysOssModel) FindByIds(ctx context.Context, ossIds []int64) ([]*SysOss, error) {
	if len(ossIds) == 0 {
		return []*SysOss{}, nil
	}

	// 构建 IN 查询
	placeholders := ""
	for i := 0; i < len(ossIds); i++ {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
	}

	query := fmt.Sprintf("select %s from %s where oss_id in (%s) order by oss_id asc", sysOssRows, m.table, placeholders)
	var args []interface{}
	for _, id := range ossIds {
		args = append(args, id)
	}

	var ossList []*SysOss
	err := m.conn.QueryRowsPartialCtx(ctx, &ossList, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return ossList, nil
}
