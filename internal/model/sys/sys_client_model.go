package sys

import (
	"context"
	"database/sql"
	"fmt"

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
	// 允许排序的列名白名单
	allowedOrderColumns := map[string]bool{
		"id":          true,
		"client_id":   true,
		"client_key":  true,
		"status":      true,
		"create_time": true,
		"update_time": true,
	}
	orderBy := "id"
	if pageQuery.OrderByColumn != "" && allowedOrderColumns[pageQuery.OrderByColumn] {
		orderBy = pageQuery.OrderByColumn
	}
	// 只允许 asc 或 desc
	orderDir := "asc"
	if pageQuery.IsAsc == "desc" {
		orderDir = "desc"
	}

	// 构建分页查询
	sqlQuery := fmt.Sprintf("select %s from %s where %s order by %s %s", sysClientRows, m.table, whereClause, orderBy, orderDir)
	if pageQuery.PageSize > 0 {
		offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
		if offset < 0 {
			offset = 0
		}
		sqlQuery += fmt.Sprintf(" limit %d offset %d", pageQuery.PageSize, offset)
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
