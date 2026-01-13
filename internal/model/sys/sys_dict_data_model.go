package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// DictDataQuery 字典数据查询条件
type DictDataQuery struct {
	DictLabel string // 字典标签（模糊查询）
	DictType  string // 字典类型
	Status    string // 状态（0正常 1停用）
}

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
		FindPage(ctx context.Context, query *DictDataQuery, pageQuery *PageQuery) ([]*SysDictData, int64, error)
		UpdateById(ctx context.Context, data *SysDictData) error
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
func (m *customSysDictDataModel) FindPage(ctx context.Context, query *DictDataQuery, pageQuery *PageQuery) ([]*SysDictData, int64, error) {
	if query == nil {
		query = &DictDataQuery{}
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{}
	}
	// 初始化分页参数的非合规值
	pageQuery.Normalize()

	// 构建 WHERE 条件
	whereClause := "1=1"
	var args []interface{}

	if query.DictLabel != "" {
		whereClause += " and dict_label like ?"
		args = append(args, "%"+query.DictLabel+"%")
	}
	if query.DictType != "" {
		whereClause += " and dict_type = ?"
		args = append(args, query.DictType)
	}
	// 注意：sys_dict_data 表可能没有 status 字段，这里先不处理

	// 查询总数
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)
	var total int64
	err := m.conn.QueryRowPartialCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 构建排序（防止 SQL 注入）
	// 允许的排序列（支持 snake_case 和 camelCase，支持多列排序，用逗号分隔）
	allowedOrderColumns := buildAllowedOrderColumns(sysDictDataFieldNames)
	// 支持多列排序
	allowedOrderColumns["dict_type, dict_sort"] = true
	allowedOrderColumns["dictType, dictSort"] = true

	orderBy := "dict_type, dict_sort"
	if pageQuery.OrderByColumn != "" {
		originalColumn := strings.TrimSpace(pageQuery.OrderByColumn)
		// 检查是否为单列或多列排序（多列用逗号分隔）
		if allowedOrderColumns[originalColumn] {
			orderBy = originalColumn
		} else {
			// 将 camelCase 转换为 snake_case
			columnName := camelToSnake(originalColumn)
			if allowedOrderColumns[columnName] {
				orderBy = columnName
			} else {
				// 如果是多列排序，检查每一列是否在白名单中
				columns := strings.Split(originalColumn, ",")
				allValid := true
				convertedColumns := make([]string, 0, len(columns))
				for _, col := range columns {
					col = strings.TrimSpace(col)
					colSnake := camelToSnake(col)
					if allowedOrderColumns[col] || allowedOrderColumns[colSnake] {
						convertedColumns = append(convertedColumns, colSnake)
					} else {
						allValid = false
						break
					}
				}
				if allValid && len(convertedColumns) > 0 {
					orderBy = strings.Join(convertedColumns, ", ")
				}
			}
		}
	}

	// 获取排序方向（默认升序）
	orderDir := pageQuery.GetOrderDir("asc")

	// 构建分页查询
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by %s %s", sysDictDataRows, m.table, whereClause, orderBy, orderDir)
	if pageQuery.PageSize > 0 {
		offset := pageQuery.GetOffset()
		sqlQuery += fmt.Sprintf(" limit %d, %d", offset, pageQuery.PageSize)
	}

	var resp []*SysDictData
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}

// UpdateById 根据ID更新字典数据，只更新非零值字段
func (m *customSysDictDataModel) UpdateById(ctx context.Context, data *SysDictData) error {
	if data.DictCode == 0 {
		return fmt.Errorf("dict_code cannot be zero")
	}

	var setParts []string
	var args []interface{}

	// 检查每个字段是否为非零值，如果是则加入更新列表
	if data.TenantId != "" {
		setParts = append(setParts, "`tenant_id` = ?")
		args = append(args, data.TenantId)
	}
	if data.DictSort > 0 {
		setParts = append(setParts, "`dict_sort` = ?")
		args = append(args, data.DictSort)
	}
	if data.DictLabel != "" {
		setParts = append(setParts, "`dict_label` = ?")
		args = append(args, data.DictLabel)
	}
	if data.DictValue != "" {
		setParts = append(setParts, "`dict_value` = ?")
		args = append(args, data.DictValue)
	}
	if data.DictType != "" {
		setParts = append(setParts, "`dict_type` = ?")
		args = append(args, data.DictType)
	}
	if data.IsDefault != "" {
		setParts = append(setParts, "`is_default` = ?")
		args = append(args, data.IsDefault)
	}
	if data.CssClass.Valid {
		setParts = append(setParts, "`css_class` = ?")
		args = append(args, data.CssClass.String)
	}
	if data.ListClass.Valid {
		setParts = append(setParts, "`list_class` = ?")
		args = append(args, data.ListClass.String)
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
		return nil // 没有需要更新的字段
	}

	// 构建更新SQL
	setClause := strings.Join(setParts, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE `dict_code` = ?", m.table, setClause)
	args = append(args, data.DictCode)

	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}
