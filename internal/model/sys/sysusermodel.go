package sys

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// UserQuery 用户查询条件
type UserQuery struct {
	UserId         int64  // 用户ID
	UserIds        string // 用户ID串（逗号分隔）
	UserName       string // 用户账号（模糊查询）
	NickName       string // 用户昵称（模糊查询）
	Status         string // 帐号状态（0正常 1停用）
	Phonenumber    string // 手机号码（模糊查询）
	DeptId         int64  // 部门ID（包含子部门）
	BeginTime      string // 开始时间
	EndTime        string // 结束时间
	ExcludeUserIds string // 排除不查询的用户ID串（逗号分隔）
	RoleId         int64  // 角色ID（用于查询已分配/未分配用户）
}

var _ SysUserModel = (*customSysUserModel)(nil)

type (
	// SysUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysUserModel.
	SysUserModel interface {
		sysUserModel
		withSession(session sqlx.Session) SysUserModel
		FindOneByUserName(ctx context.Context, userName, tenantId string) (*SysUser, error)
		// FindAllocatedPage 分页查询已分配用户角色列表
		FindAllocatedPage(ctx context.Context, query *UserQuery, pageQuery *PageQuery) ([]*SysUser, int64, error)
		// FindUnallocatedPage 分页查询未分配用户角色列表
		FindUnallocatedPage(ctx context.Context, query *UserQuery, pageQuery *PageQuery) ([]*SysUser, int64, error)
		// FindPage 分页查询用户列表
		FindPage(ctx context.Context, query *UserQuery, pageQuery *PageQuery) ([]*SysUser, int64, error)
		// FindByIds 根据用户ID串查询用户基础信息
		FindByIds(ctx context.Context, userIds []int64, deptId int64) ([]*SysUser, error)
		// FindByPhonenumber 根据手机号查询用户
		FindByPhonenumber(ctx context.Context, phonenumber string) (*SysUser, error)
		// FindByEmail 根据邮箱查询用户
		FindByEmail(ctx context.Context, email string) (*SysUser, error)
		// CheckUserNameUnique 校验用户名称是否唯一
		CheckUserNameUnique(ctx context.Context, userName string, excludeUserId int64) (bool, error)
		// CheckPhoneUnique 校验手机号码是否唯一
		CheckPhoneUnique(ctx context.Context, phonenumber string, excludeUserId int64) (bool, error)
		// CheckEmailUnique 校验邮箱是否唯一
		CheckEmailUnique(ctx context.Context, email string, excludeUserId int64) (bool, error)
		// UpdateUserStatus 修改用户状态
		UpdateUserStatus(ctx context.Context, userId int64, status string) error
		// ResetUserPwd 重置用户密码
		ResetUserPwd(ctx context.Context, userId int64, password string) error
		// UpdateUserProfile 修改用户基本信息
		UpdateUserProfile(ctx context.Context, userId int64, nickName, email, phonenumber, sex string) error
		// UpdateUserAvatar 修改用户头像
		UpdateUserAvatar(ctx context.Context, userId int64, avatar int64) error
		// SelectUserListByDept 根据部门查询用户列表
		SelectUserListByDept(ctx context.Context, deptId int64) ([]*SysUser, error)
	}

	customSysUserModel struct {
		*defaultSysUserModel
	}
)

// NewSysUserModel returns a model for the database table.
func NewSysUserModel(conn sqlx.SqlConn) SysUserModel {
	return &customSysUserModel{
		defaultSysUserModel: newSysUserModel(conn),
	}
}

func (m *customSysUserModel) withSession(session sqlx.Session) SysUserModel {
	return NewSysUserModel(sqlx.NewSqlConnFromSession(session))
}

// FindOneByUserName 根据用户名和租户ID查询用户
func (m *customSysUserModel) FindOneByUserName(ctx context.Context, userName, tenantId string) (*SysUser, error) {
	query := "select `user_id`,`tenant_id`,`dept_id`,`user_name`,`nick_name`,`user_type`,`email`,`phonenumber`,`sex`,`avatar`,`password`,`status`,`del_flag`,`login_ip`,`login_date`,`create_dept`,`create_by`,`create_time`,`update_by`,`update_time`,`remark` from `sys_user` where `user_name` = ? and `tenant_id` = ? and `del_flag` = '0' limit 1"
	var resp SysUser
	err := m.conn.QueryRowCtx(ctx, &resp, query, userName, tenantId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindAllocatedPage 分页查询已分配用户角色列表
func (m *customSysUserModel) FindAllocatedPage(ctx context.Context, query *UserQuery, pageQuery *PageQuery) ([]*SysUser, int64, error) {
	if query == nil {
		query = &UserQuery{}
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{
			PageNum:  1,
			PageSize: 10,
		}
	}

	// 构建 WHERE 条件
	whereClause := "u.del_flag = '0'"
	var args []interface{}

	if query.RoleId > 0 {
		whereClause += " and r.role_id = ?"
		args = append(args, query.RoleId)
	}
	if query.UserName != "" {
		whereClause += " and u.user_name LIKE ?"
		args = append(args, "%"+query.UserName+"%")
	}
	if query.Status != "" {
		whereClause += " and u.status = ?"
		args = append(args, query.Status)
	}
	if query.Phonenumber != "" {
		whereClause += " and u.phonenumber LIKE ?"
		args = append(args, "%"+query.Phonenumber+"%")
	}

	// 计算总数
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT u.user_id)
		FROM sys_user u
		LEFT JOIN sys_dept d ON u.dept_id = d.dept_id
		LEFT JOIN sys_user_role sur ON u.user_id = sur.user_id
		LEFT JOIN sys_role r ON r.role_id = sur.role_id
		WHERE %s
	`, whereClause)

	var total int64
	err := m.conn.QueryRowPartialCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 计算分页参数
	offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
	limit := pageQuery.PageSize

	// 查询数据
	userRows := "u.user_id,u.tenant_id,u.dept_id,u.user_name,u.nick_name,u.user_type,u.email,u.phonenumber,u.sex,u.avatar,u.password,u.status,u.del_flag,u.login_ip,u.login_date,u.create_dept,u.create_by,u.create_time,u.update_by,u.update_time,u.remark"
	querySQL := fmt.Sprintf(`
		SELECT DISTINCT %s
		FROM sys_user u
		LEFT JOIN sys_dept d ON u.dept_id = d.dept_id
		LEFT JOIN sys_user_role sur ON u.user_id = sur.user_id
		LEFT JOIN sys_role r ON r.role_id = sur.role_id
		WHERE %s
		ORDER BY u.user_id ASC
		LIMIT ? OFFSET ?
	`, userRows, whereClause)

	args = append(args, limit, offset)

	var userList []*SysUser
	err = m.conn.QueryRowsPartialCtx(ctx, &userList, querySQL, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	return userList, total, nil
}

// FindUnallocatedPage 分页查询未分配用户角色列表
func (m *customSysUserModel) FindUnallocatedPage(ctx context.Context, query *UserQuery, pageQuery *PageQuery) ([]*SysUser, int64, error) {
	if query == nil || query.RoleId == 0 {
		return []*SysUser{}, 0, fmt.Errorf("角色ID不能为空")
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{
			PageNum:  1,
			PageSize: 10,
		}
	}

	// 1. 查询已分配该角色的用户ID列表（通过sys_user_role表）
	userRoleModel := NewSysUserRoleModel(m.conn)
	userIds, err := userRoleModel.SelectUserIdsByRoleId(ctx, query.RoleId)
	if err != nil {
		return nil, 0, err
	}

	// 构建 WHERE 条件
	whereClause := "u.del_flag = '0'"
	var args []interface{}

	// 排除已分配的用户
	if len(userIds) > 0 {
		placeholders := ""
		for i := 0; i < len(userIds); i++ {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args = append(args, userIds[i])
		}
		whereClause += fmt.Sprintf(" AND u.user_id NOT IN (%s)", placeholders)
	}

	// 排除角色ID不等于当前角色的用户（如果用户有其他角色）
	whereClause += " AND (r.role_id != ? OR r.role_id IS NULL)"
	args = append(args, query.RoleId)

	if query.UserName != "" {
		whereClause += " and u.user_name LIKE ?"
		args = append(args, "%"+query.UserName+"%")
	}
	if query.Phonenumber != "" {
		whereClause += " and u.phonenumber LIKE ?"
		args = append(args, "%"+query.Phonenumber+"%")
	}

	// 计算总数
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT u.user_id)
		FROM sys_user u
		LEFT JOIN sys_dept d ON u.dept_id = d.dept_id
		LEFT JOIN sys_user_role sur ON u.user_id = sur.user_id
		LEFT JOIN sys_role r ON r.role_id = sur.role_id
		WHERE %s
	`, whereClause)

	var total int64
	err = m.conn.QueryRowPartialCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 计算分页参数
	offset := (pageQuery.PageNum - 1) * pageQuery.PageSize
	limit := pageQuery.PageSize

	// 查询数据
	userRows := "u.user_id,u.tenant_id,u.dept_id,u.user_name,u.nick_name,u.user_type,u.email,u.phonenumber,u.sex,u.avatar,u.password,u.status,u.del_flag,u.login_ip,u.login_date,u.create_dept,u.create_by,u.create_time,u.update_by,u.update_time,u.remark"
	querySQL := fmt.Sprintf(`
		SELECT DISTINCT %s
		FROM sys_user u
		LEFT JOIN sys_dept d ON u.dept_id = d.dept_id
		LEFT JOIN sys_user_role sur ON u.user_id = sur.user_id
		LEFT JOIN sys_role r ON r.role_id = sur.role_id
		WHERE %s
		ORDER BY u.user_id ASC
		LIMIT ? OFFSET ?
	`, userRows, whereClause)

	args = append(args, limit, offset)

	var userList []*SysUser
	err = m.conn.QueryRowsPartialCtx(ctx, &userList, querySQL, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	return userList, total, nil
}
