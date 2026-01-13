package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ClientQuery 客户端管理查询条件
type ClientQuery struct {
	ClientId     string // 客户端id
	ClientKey    string // 客户端key
	ClientSecret string // 客户端秘钥
	Status       string // 状态（0正常 1停用）
}

var _ SysClientModel = (*customSysClientModel)(nil)

type (
	// SysClientModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysClientModel.
	SysClientModel interface {
		sysClientModel
		withSession(session sqlx.Session) SysClientModel
		FindOneByClientId(ctx context.Context, clientId string) (*SysClient, error)
		CheckClientKeyUnique(ctx context.Context, clientKey string, excludeId int64) (bool, error)
		FindPage(ctx context.Context, query *ClientQuery, pageQuery *PageQuery) ([]*SysClient, int64, error)
		UpdateClientStatus(ctx context.Context, clientId string, status string) error
		UpdateById(ctx context.Context, data *SysClient) error
	}

	customSysClientModel struct {
		*defaultSysClientModel
	}
)

// NewSysClientModel returns a model for the database table.
func NewSysClientModel(conn sqlx.SqlConn) SysClientModel {
	return &customSysClientModel{
		defaultSysClientModel: newSysClientModel(conn),
	}
}

func (m *customSysClientModel) withSession(session sqlx.Session) SysClientModel {
	return NewSysClientModel(sqlx.NewSqlConnFromSession(session))
}

// FindOneByClientId 根据客户端ID查询
func (m *customSysClientModel) FindOneByClientId(ctx context.Context, clientId string) (*SysClient, error) {
	query := fmt.Sprintf("select %s from %s where `client_id` = ? and `del_flag` = '0' limit 1", sysClientRows, m.table)
	var resp SysClient
	err := m.conn.QueryRowCtx(ctx, &resp, query, clientId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// CheckClientKeyUnique 校验客户端key唯一性
func (m *customSysClientModel) CheckClientKeyUnique(ctx context.Context, clientKey string, excludeId int64) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where client_key = ? and del_flag = '0'", m.table)
	if excludeId > 0 {
		query += " and id != ?"
	}

	var count int64
	var err error
	if excludeId > 0 {
		err = m.conn.QueryRowPartialCtx(ctx, &count, query, clientKey, excludeId)
	} else {
		err = m.conn.QueryRowPartialCtx(ctx, &count, query, clientKey)
	}
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// FindPage 分页查询客户端管理（支持条件查询和分页）
func (m *customSysClientModel) FindPage(ctx context.Context, query *ClientQuery, pageQuery *PageQuery) ([]*SysClient, int64, error) {
	if query == nil {
		query = &ClientQuery{}
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{}
	}
	// 初始化分页参数的非合规值
	pageQuery.Normalize()

	// 构建 WHERE 条件
	whereClause := "del_flag = '0'"
	var args []interface{}

	if query.ClientId != "" {
		whereClause += " and client_id = ?"
		args = append(args, query.ClientId)
	}
	if query.ClientKey != "" {
		whereClause += " and client_key = ?"
		args = append(args, query.ClientKey)
	}
	if query.ClientSecret != "" {
		whereClause += " and client_secret = ?"
		args = append(args, query.ClientSecret)
	}
	if query.Status != "" {
		whereClause += " and status = ?"
		args = append(args, query.Status)
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
	allowedOrderColumns := buildAllowedOrderColumns(sysClientFieldNames)
	orderBy := pageQuery.GetOrderBy("id", allowedOrderColumns)

	// 获取排序方向（默认升序）
	orderDir := pageQuery.GetOrderDir("asc")

	// 构建分页查询
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by %s %s", sysClientRows, m.table, whereClause, orderBy, orderDir)
	if pageQuery.PageSize > 0 {
		offset := pageQuery.GetOffset()
		sqlQuery += fmt.Sprintf(" limit %d, %d", offset, pageQuery.PageSize)
	}

	var resp []*SysClient
	err = m.conn.QueryRowsPartialCtx(ctx, &resp, sqlQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return resp, total, nil
}

// UpdateClientStatus 更新客户端状态
func (m *customSysClientModel) UpdateClientStatus(ctx context.Context, clientId string, status string) error {
	query := fmt.Sprintf("update %s set status = ? where client_id = ? and del_flag = '0'", m.table)
	_, err := m.conn.ExecCtx(ctx, query, status, clientId)
	return err
}

// UpdateById 根据ID更新客户端，只更新非零值字段
func (m *customSysClientModel) UpdateById(ctx context.Context, data *SysClient) error {
	if data.Id == 0 {
		return fmt.Errorf("id cannot be zero")
	}

	var setParts []string
	var args []interface{}

	// 检查每个字段是否为非零值，如果是则加入更新列表
	if data.ClientId.Valid {
		setParts = append(setParts, "`client_id` = ?")
		args = append(args, data.ClientId.String)
	}
	if data.ClientKey.Valid {
		setParts = append(setParts, "`client_key` = ?")
		args = append(args, data.ClientKey.String)
	}
	if data.ClientSecret.Valid {
		setParts = append(setParts, "`client_secret` = ?")
		args = append(args, data.ClientSecret.String)
	}
	if data.GrantType.Valid {
		setParts = append(setParts, "`grant_type` = ?")
		args = append(args, data.GrantType.String)
	}
	if data.DeviceType.Valid {
		setParts = append(setParts, "`device_type` = ?")
		args = append(args, data.DeviceType.String)
	}
	if data.ActiveTimeout > 0 {
		setParts = append(setParts, "`active_timeout` = ?")
		args = append(args, data.ActiveTimeout)
	}
	if data.Timeout > 0 {
		setParts = append(setParts, "`timeout` = ?")
		args = append(args, data.Timeout)
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
	query := fmt.Sprintf("UPDATE %s SET %s WHERE `id` = ?", m.table, setClause)
	args = append(args, data.Id)

	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}
