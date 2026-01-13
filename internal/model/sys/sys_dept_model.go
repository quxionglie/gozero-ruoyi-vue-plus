package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// DeptQuery 部门查询条件
type DeptQuery struct {
	DeptId       int64  // 部门id
	ParentId     int64  // 父部门id
	DeptName     string // 部门名称（模糊查询）
	DeptCategory string // 部门类别编码（模糊查询）
	Status       string // 状态（0正常 1停用）
	BelongDeptId int64  // 归属部门id（部门树）
}

var _ SysDeptModel = (*customSysDeptModel)(nil)

type (
	// SysDeptModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysDeptModel.
	SysDeptModel interface {
		sysDeptModel
		withSession(session sqlx.Session) SysDeptModel
		FindAll(ctx context.Context, query *DeptQuery) ([]*SysDept, error)
		FindByIds(ctx context.Context, deptIds []int64) ([]*SysDept, error)
		CheckDeptNameUnique(ctx context.Context, deptName string, parentId int64, excludeDeptId int64) (bool, error)
		HasChildByDeptId(ctx context.Context, deptId int64) (bool, error)
		CheckDeptExistUser(ctx context.Context, deptId int64) (bool, error)
		CountNormalChildrenDeptById(ctx context.Context, deptId int64) (int64, error)
		SelectDeptAndChildById(ctx context.Context, deptId int64) ([]int64, error)
		UpdateById(ctx context.Context, data *SysDept) error
	}

	customSysDeptModel struct {
		*defaultSysDeptModel
	}
)

// NewSysDeptModel returns a model for the database table.
func NewSysDeptModel(conn sqlx.SqlConn) SysDeptModel {
	return &customSysDeptModel{
		defaultSysDeptModel: newSysDeptModel(conn),
	}
}

func (m *customSysDeptModel) withSession(session sqlx.Session) SysDeptModel {
	return NewSysDeptModel(sqlx.NewSqlConnFromSession(session))
}

// FindAll 查询所有部门（不分页）
func (m *customSysDeptModel) FindAll(ctx context.Context, query *DeptQuery) ([]*SysDept, error) {
	if query == nil {
		query = &DeptQuery{}
	}

	// 构建 WHERE 条件
	whereClause := "del_flag = '0'"
	var args []interface{}

	if query.DeptId > 0 {
		whereClause += " and dept_id = ?"
		args = append(args, query.DeptId)
	}
	if query.ParentId > 0 {
		whereClause += " and parent_id = ?"
		args = append(args, query.ParentId)
	}
	if query.DeptName != "" {
		whereClause += " and dept_name like ?"
		args = append(args, "%"+query.DeptName+"%")
	}
	if query.DeptCategory != "" {
		whereClause += " and dept_category like ?"
		args = append(args, "%"+query.DeptCategory+"%")
	}
	if query.Status != "" {
		whereClause += " and status = ?"
		args = append(args, query.Status)
	}
	if query.BelongDeptId > 0 {
		// 部门树查询：需要查询所有子部门（这里简化处理，实际应该递归查询子部门）
		// TODO: 如果需要完整的部门树查询，需要调用 SelectDeptAndChildById
		whereClause += " and dept_id = ?"
		args = append(args, query.BelongDeptId)
	}

	// 排序：先按 ancestors，再按 parent_id，再按 order_num，最后按 dept_id
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by ancestors asc, parent_id asc, order_num asc, dept_id asc", sysDeptRows, m.table, whereClause)

	var resp []*SysDept
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return resp, nil
}

// FindByIds 根据部门ID列表查询部门
func (m *customSysDeptModel) FindByIds(ctx context.Context, deptIds []int64) ([]*SysDept, error) {
	if len(deptIds) == 0 {
		return []*SysDept{}, nil
	}

	// 构建 IN 查询
	placeholders := ""
	for i := 0; i < len(deptIds); i++ {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
	}

	query := fmt.Sprintf("select dept_id, dept_name, leader from %s where dept_id in (%s) and status = '0' and del_flag = '0' order by order_num asc", m.table, placeholders)
	var args []interface{}
	for _, id := range deptIds {
		args = append(args, id)
	}

	var resp []*SysDept
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return resp, nil
}

// CheckDeptNameUnique 校验部门名称唯一性（同父部门下唯一）
func (m *customSysDeptModel) CheckDeptNameUnique(ctx context.Context, deptName string, parentId int64, excludeDeptId int64) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where dept_name = ? and parent_id = ? and del_flag = '0'", m.table)
	var args []interface{}
	args = append(args, deptName, parentId)

	if excludeDeptId > 0 {
		query += " and dept_id != ?"
		args = append(args, excludeDeptId)
	}

	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// HasChildByDeptId 是否存在子节点
func (m *customSysDeptModel) HasChildByDeptId(ctx context.Context, deptId int64) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where parent_id = ? and del_flag = '0'", m.table)
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, deptId)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckDeptExistUser 查询部门是否存在用户
func (m *customSysDeptModel) CheckDeptExistUser(ctx context.Context, deptId int64) (bool, error) {
	query := "select count(*) from `sys_user` where dept_id = ? and del_flag = '0'"
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, deptId)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountNormalChildrenDeptById 根据ID查询所有子部门数（正常状态）
func (m *customSysDeptModel) CountNormalChildrenDeptById(ctx context.Context, deptId int64) (int64, error) {
	// 使用 FIND_IN_SET 查询 ancestors 中包含 deptId 的部门
	query := fmt.Sprintf("select count(*) from %s where status = '0' and del_flag = '0' and FIND_IN_SET(?, ancestors)", m.table)
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, deptId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SelectDeptAndChildById 查询部门及其所有子部门ID
func (m *customSysDeptModel) SelectDeptAndChildById(ctx context.Context, deptId int64) ([]int64, error) {
	// 查询包括自己和所有子部门
	query := fmt.Sprintf("select dept_id from %s where (dept_id = ? or FIND_IN_SET(?, ancestors)) and del_flag = '0'", m.table)
	var resp []*SysDept
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, deptId, deptId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var deptIds []int64
	for _, dept := range resp {
		deptIds = append(deptIds, dept.DeptId)
	}
	return deptIds, nil
}

// UpdateById 根据ID更新部门，只更新非零值字段
func (m *customSysDeptModel) UpdateById(ctx context.Context, data *SysDept) error {
	if data.DeptId == 0 {
		return fmt.Errorf("dept_id cannot be zero")
	}

	var setParts []string
	var args []interface{}

	// 检查每个字段是否为非零值，如果是则加入更新列表
	if data.TenantId != "" {
		setParts = append(setParts, "`tenant_id` = ?")
		args = append(args, data.TenantId)
	}
	if data.ParentId > 0 {
		setParts = append(setParts, "`parent_id` = ?")
		args = append(args, data.ParentId)
	}
	if data.Ancestors != "" {
		setParts = append(setParts, "`ancestors` = ?")
		args = append(args, data.Ancestors)
	}
	if data.DeptName != "" {
		setParts = append(setParts, "`dept_name` = ?")
		args = append(args, data.DeptName)
	}
	if data.DeptCategory.Valid {
		setParts = append(setParts, "`dept_category` = ?")
		args = append(args, data.DeptCategory.String)
	}
	if data.OrderNum > 0 {
		setParts = append(setParts, "`order_num` = ?")
		args = append(args, data.OrderNum)
	}
	if data.Leader.Valid {
		setParts = append(setParts, "`leader` = ?")
		args = append(args, data.Leader.Int64)
	}
	if data.Phone.Valid {
		setParts = append(setParts, "`phone` = ?")
		args = append(args, data.Phone.String)
	}
	if data.Email.Valid {
		setParts = append(setParts, "`email` = ?")
		args = append(args, data.Email.String)
	}
	if data.Status != "" {
		setParts = append(setParts, "`status` = ?")
		args = append(args, data.Status)
	}
	if data.DelFlag != "" {
		setParts = append(setParts, "`del_flag` = ?")
		args = append(args, data.DelFlag)
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

	if len(setParts) == 0 {
		return nil // 没有需要更新的字段
	}

	// 构建更新SQL
	setClause := strings.Join(setParts, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE `dept_id` = ?", m.table, setClause)
	args = append(args, data.DeptId)

	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}
