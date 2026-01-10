package sys

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysDictDataModel = (*customSysDictDataModel)(nil)

type (
	// SysDictDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysDictDataModel.
	SysDictDataModel interface {
		sysDictDataModel
		withSession(session sqlx.Session) SysDictDataModel
		FindAll(ctx context.Context) ([]*SysDictData, error)
		FindByDictType(ctx context.Context, dictType string) ([]*SysDictData, error)
		CheckDictDataUnique(ctx context.Context, dictType, dictValue string, excludeDictCode int64) (bool, error)
		CountByDictType(ctx context.Context, dictType string) (int64, error)
		UpdateDictTypeByOldDictType(ctx context.Context, oldDictType, newDictType string) error
		FindPage(ctx context.Context, dictLabel, dictType, status string, pageNum, pageSize int32, orderByColumn, isAsc string) ([]*SysDictData, int64, error)
	}

	customSysDictDataModel struct {
		*defaultSysDictDataModel
	}
)

// NewSysDictDataModel returns a model for the database table.
func NewSysDictDataModel(conn sqlx.SqlConn) SysDictDataModel {
	return &customSysDictDataModel{
		defaultSysDictDataModel: newSysDictDataModel(conn),
	}
}

func (m *customSysDictDataModel) withSession(session sqlx.Session) SysDictDataModel {
	return NewSysDictDataModel(sqlx.NewSqlConnFromSession(session))
}

// FindAll 查询所有字典数据
func (m *customSysDictDataModel) FindAll(ctx context.Context) ([]*SysDictData, error) {
	query := fmt.Sprintf("select %s from %s order by dict_type asc, dict_sort asc", sysDictDataRows, m.table)
	var resp []*SysDictData
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return resp, nil
}

// FindByDictType 根据字典类型查询字典数据
func (m *customSysDictDataModel) FindByDictType(ctx context.Context, dictType string) ([]*SysDictData, error) {
	query := fmt.Sprintf("select %s from %s where dict_type = ? order by dict_sort asc", sysDictDataRows, m.table)
	var resp []*SysDictData
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, dictType)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return resp, nil
}

// CheckDictDataUnique 校验字典数据唯一性（同一字典类型下，字典键值唯一）
func (m *customSysDictDataModel) CheckDictDataUnique(ctx context.Context, dictType, dictValue string, excludeDictCode int64) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where dict_type = ? and dict_value = ?", m.table)
	if excludeDictCode > 0 {
		query += " and dict_code != ?"
	}

	var count int64
	var err error
	if excludeDictCode > 0 {
		err = m.conn.QueryRowPartialCtx(ctx, &count, query, dictType, dictValue, excludeDictCode)
	} else {
		err = m.conn.QueryRowPartialCtx(ctx, &count, query, dictType, dictValue)
	}
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// CountByDictType 统计字典数据数量（用于检查是否已分配）
func (m *customSysDictDataModel) CountByDictType(ctx context.Context, dictType string) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where dict_type = ?", m.table)
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, dictType)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateDictTypeByOldDictType 更新关联的字典数据的字典类型
func (m *customSysDictDataModel) UpdateDictTypeByOldDictType(ctx context.Context, oldDictType, newDictType string) error {
	query := fmt.Sprintf("update %s set dict_type = ? where dict_type = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, newDictType, oldDictType)
	return err
}

// FindPage 分页查询字典数据（支持条件查询和分页）
func (m *customSysDictDataModel) FindPage(ctx context.Context, dictLabel, dictType, status string, pageNum, pageSize int32, orderByColumn, isAsc string) ([]*SysDictData, int64, error) {
	// 构建 WHERE 条件
	whereClause := "1=1"
	var args []interface{}

	if dictLabel != "" {
		whereClause += " and dict_label like ?"
		args = append(args, "%"+dictLabel+"%")
	}
	if dictType != "" {
		whereClause += " and dict_type = ?"
		args = append(args, dictType)
	}
	// 注意：sys_dict_data 表可能没有 status 字段，这里先不处理

	// 查询总数
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)
	var total int64
	err := m.conn.QueryRowPartialCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 构建排序
	orderBy := "dict_type, dict_sort"
	if orderByColumn != "" {
		orderBy = orderByColumn
	}
	orderDir := "asc"
	if isAsc == "desc" {
		orderDir = "desc"
	}

	// 构建分页查询
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", sysDictDataRows, m.table, whereClause, orderBy, orderDir)
	if pageSize > 0 {
		offset := (pageNum - 1) * pageSize
		query += fmt.Sprintf(" limit %d offset %d", pageSize, offset)
	}

	var resp []*SysDictData
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}
