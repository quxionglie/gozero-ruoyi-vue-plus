package sys

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// DictTypeQuery 字典类型查询条件
type DictTypeQuery struct {
	DictName string // 字典名称（模糊查询）
	DictType string // 字典类型（模糊查询）
	Status   string // 状态（0正常 1停用）
}

var _ SysDictTypeModel = (*customSysDictTypeModel)(nil)

type (
	// SysDictTypeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysDictTypeModel.
	SysDictTypeModel interface {
		sysDictTypeModel
		withSession(session sqlx.Session) SysDictTypeModel
		FindAll(ctx context.Context) ([]*SysDictType, error)
		CheckDictTypeUnique(ctx context.Context, dictType string, excludeDictId int64) (bool, error)
		CountByDictType(ctx context.Context, dictType string) (int64, error)
		FindPage(ctx context.Context, query *DictTypeQuery, pageQuery *PageQuery) ([]*SysDictType, int64, error)
	}

	customSysDictTypeModel struct {
		*defaultSysDictTypeModel
	}
)

// NewSysDictTypeModel returns a model for the database table.
func NewSysDictTypeModel(conn sqlx.SqlConn) SysDictTypeModel {
	return &customSysDictTypeModel{
		defaultSysDictTypeModel: newSysDictTypeModel(conn),
	}
}

func (m *customSysDictTypeModel) withSession(session sqlx.Session) SysDictTypeModel {
	return NewSysDictTypeModel(sqlx.NewSqlConnFromSession(session))
}

// FindAll 查询所有字典类型
func (m *customSysDictTypeModel) FindAll(ctx context.Context) ([]*SysDictType, error) {
	query := fmt.Sprintf("select %s from %s order by dict_id asc", sysDictTypeRows, m.table)
	var resp []*SysDictType
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return resp, nil
}

// CheckDictTypeUnique 校验字典类型唯一性
func (m *customSysDictTypeModel) CheckDictTypeUnique(ctx context.Context, dictType string, excludeDictId int64) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where dict_type = ?", m.table)
	if excludeDictId > 0 {
		query += " and dict_id != ?"
	}

	var count int64
	var err error
	if excludeDictId > 0 {
		err = m.conn.QueryRowPartialCtx(ctx, &count, query, dictType, excludeDictId)
	} else {
		err = m.conn.QueryRowPartialCtx(ctx, &count, query, dictType)
	}
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// CountByDictType 统计字典类型数量（用于检查是否存在）
func (m *customSysDictTypeModel) CountByDictType(ctx context.Context, dictType string) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where dict_type = ?", m.table)
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, dictType)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindPage 分页查询字典类型（支持条件查询和分页）
func (m *customSysDictTypeModel) FindPage(ctx context.Context, query *DictTypeQuery, pageQuery *PageQuery) ([]*SysDictType, int64, error) {
	if query == nil {
		query = &DictTypeQuery{}
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{}
	}

	// 构建 WHERE 条件
	whereClause := "1=1"
	var args []interface{}

	if query.DictName != "" {
		whereClause += " and dict_name like ?"
		args = append(args, "%"+query.DictName+"%")
	}
	if query.DictType != "" {
		whereClause += " and dict_type like ?"
		args = append(args, "%"+query.DictType+"%")
	}
	// 注意：sys_dict_type 表可能没有 status 字段，这里先不处理

	// 查询总数
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)
	var total int64
	err := m.conn.QueryRowPartialCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 构建排序
	orderBy := "dict_id"
	if pageQuery.OrderByColumn != "" {
		orderBy = pageQuery.OrderByColumn
	}
	orderDir := "asc"
	if pageQuery.IsAsc == "desc" {
		orderDir = "desc"
	}

	// 构建分页查询
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by %s %s", sysDictTypeRows, m.table, whereClause, orderBy, orderDir)
	if pageQuery.PageSize > 0 {
		offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
		if offset < 0 {
			offset = 0
		}
		sqlQuery += fmt.Sprintf(" limit %d offset %d", pageQuery.PageSize, offset)
	}

	var resp []*SysDictType
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}
