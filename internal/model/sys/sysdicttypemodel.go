package sys

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

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
		FindPage(ctx context.Context, dictName, dictType, status string, pageNum, pageSize int32, orderByColumn, isAsc string) ([]*SysDictType, int64, error)
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
func (m *customSysDictTypeModel) FindPage(ctx context.Context, dictName, dictType, status string, pageNum, pageSize int32, orderByColumn, isAsc string) ([]*SysDictType, int64, error) {
	// 构建 WHERE 条件
	whereClause := "1=1"
	var args []interface{}

	if dictName != "" {
		whereClause += " and dict_name like ?"
		args = append(args, "%"+dictName+"%")
	}
	if dictType != "" {
		whereClause += " and dict_type like ?"
		args = append(args, "%"+dictType+"%")
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
	if orderByColumn != "" {
		orderBy = orderByColumn
	}
	orderDir := "asc"
	if isAsc == "desc" {
		orderDir = "desc"
	}

	// 构建分页查询
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", sysDictTypeRows, m.table, whereClause, orderBy, orderDir)
	if pageSize > 0 {
		offset := (pageNum - 1) * pageSize
		query += fmt.Sprintf(" limit %d offset %d", pageSize, offset)
	}

	var resp []*SysDictType
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}
