package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// FindPage 分页查询用户列表
func (m *customSysUserModel) FindPage(ctx context.Context, query *UserQuery, pageQuery *PageQuery) ([]*SysUser, int64, error) {
	if query == nil {
		query = &UserQuery{}
	}
	if pageQuery == nil {
		pageQuery = &PageQuery{
			PageNum:  1,
			PageSize: 10,
		}
	} else {
		// 初始化分页参数的非合规值
		pageQuery.Normalize()
	}

	// 构建 WHERE 条件
	whereClause := "u.del_flag = '0'"
	var args []interface{}

	if query.UserId > 0 {
		whereClause += " and u.user_id = ?"
		args = append(args, query.UserId)
	}
	if query.UserIds != "" {
		// 解析用户ID串
		userIds := parseIds(query.UserIds)
		if len(userIds) > 0 {
			placeholders := ""
			for i, id := range userIds {
				if i > 0 {
					placeholders += ","
				}
				placeholders += "?"
				args = append(args, id)
			}
			whereClause += fmt.Sprintf(" and u.user_id IN (%s)", placeholders)
		}
	}
	if query.UserName != "" {
		whereClause += " and u.user_name LIKE ?"
		args = append(args, "%"+query.UserName+"%")
	}
	if query.NickName != "" {
		whereClause += " and u.nick_name LIKE ?"
		args = append(args, "%"+query.NickName+"%")
	}
	if query.Status != "" {
		whereClause += " and u.status = ?"
		args = append(args, query.Status)
	}
	if query.Phonenumber != "" {
		whereClause += " and u.phonenumber LIKE ?"
		args = append(args, "%"+query.Phonenumber+"%")
	}
	if query.DeptId > 0 {
		// TODO: 需要查询部门及其所有子部门
		// 这里简化处理，只查询直接关联的部门
		whereClause += " and u.dept_id = ?"
		args = append(args, query.DeptId)
	}
	if query.BeginTime != "" && query.EndTime != "" {
		whereClause += " and u.create_time BETWEEN ? AND ?"
		args = append(args, query.BeginTime, query.EndTime)
	}
	if query.ExcludeUserIds != "" {
		excludeIds := parseIds(query.ExcludeUserIds)
		if len(excludeIds) > 0 {
			placeholders := ""
			for i, id := range excludeIds {
				if i > 0 {
					placeholders += ","
				}
				placeholders += "?"
				args = append(args, id)
			}
			whereClause += fmt.Sprintf(" and u.user_id NOT IN (%s)", placeholders)
		}
	}

	// 构建排序（防止 SQL 注入）
	// 允许的排序列（支持 snake_case 和 camelCase，需要包含表别名）
	allowedOrderColumns := buildAllowedOrderColumnsWithPrefix(sysUserFieldNames, "u.")
	orderBy := pageQuery.GetOrderByWithDirAndPrefix("u.user_id ASC", allowedOrderColumns, "u.", "asc")

	// 计算总数
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT u.user_id)
		FROM sys_user u
		LEFT JOIN sys_dept d ON u.dept_id = d.dept_id
		WHERE %s
	`, whereClause)

	var total int64
	err := m.conn.QueryRowPartialCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 计算分页参数
	offset, limit := pageQuery.GetOffsetAndLimit()

	// 查询数据
	// 将 sysUserRows 转换为别名形式（u.xxx）
	fields := strings.Split(sysUserRows, ",")
	aliasFields := make([]string, len(fields))
	for i, field := range fields {
		field = strings.TrimSpace(field)
		field = strings.Trim(field, "`")
		aliasFields[i] = "u." + field
	}
	userRows := strings.Join(aliasFields, ",")
	querySQL := fmt.Sprintf(`
		SELECT DISTINCT %s
		FROM sys_user u
		LEFT JOIN sys_dept d ON u.dept_id = d.dept_id
		WHERE %s
		ORDER BY %s
		LIMIT ?, ?
	`, userRows, whereClause, orderBy)

	args = append(args, offset, limit)

	var userList []*SysUser
	err = m.conn.QueryRowsPartialCtx(ctx, &userList, querySQL, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	return userList, total, nil
}

// FindByIds 根据用户ID串查询用户基础信息
func (m *customSysUserModel) FindByIds(ctx context.Context, userIds []int64, deptId int64) ([]*SysUser, error) {
	if len(userIds) == 0 {
		return []*SysUser{}, nil
	}

	whereClause := "del_flag = '0' and status = '0'"
	var args []interface{}

	if deptId > 0 {
		whereClause += " and dept_id = ?"
		args = append(args, deptId)
	}

	placeholders := ""
	for i, id := range userIds {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, id)
	}
	whereClause += fmt.Sprintf(" and user_id IN (%s)", placeholders)

	userRows := "user_id,user_name,nick_name"
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE %s
		ORDER BY user_id ASC
	`, userRows, m.table, whereClause)

	var userList []*SysUser
	err := m.conn.QueryRowsPartialCtx(ctx, &userList, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return userList, nil
}

// FindByPhonenumber 根据手机号查询用户
func (m *customSysUserModel) FindByPhonenumber(ctx context.Context, phonenumber string) (*SysUser, error) {
	query := fmt.Sprintf("select %s from %s where `phonenumber` = ? and `del_flag` = '0' limit 1", sysUserRows, m.table)
	var resp SysUser
	err := m.conn.QueryRowCtx(ctx, &resp, query, phonenumber)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindByEmail 根据邮箱查询用户
func (m *customSysUserModel) FindByEmail(ctx context.Context, email string) (*SysUser, error) {
	query := fmt.Sprintf("select %s from %s where `email` = ? and `del_flag` = '0' limit 1", sysUserRows, m.table)
	var resp SysUser
	err := m.conn.QueryRowCtx(ctx, &resp, query, email)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// CheckUserNameUnique 校验用户名称是否唯一
func (m *customSysUserModel) CheckUserNameUnique(ctx context.Context, userName string, excludeUserId int64) (bool, error) {
	whereClause := "user_name = ? and del_flag = '0'"
	args := []interface{}{userName}
	if excludeUserId > 0 {
		whereClause += " and user_id != ?"
		args = append(args, excludeUserId)
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", m.table, whereClause)
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// CheckPhoneUnique 校验手机号码是否唯一
func (m *customSysUserModel) CheckPhoneUnique(ctx context.Context, phonenumber string, excludeUserId int64) (bool, error) {
	if phonenumber == "" {
		return true, nil
	}

	whereClause := "phonenumber = ? and del_flag = '0'"
	args := []interface{}{phonenumber}
	if excludeUserId > 0 {
		whereClause += " and user_id != ?"
		args = append(args, excludeUserId)
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", m.table, whereClause)
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// CheckEmailUnique 校验邮箱是否唯一
func (m *customSysUserModel) CheckEmailUnique(ctx context.Context, email string, excludeUserId int64) (bool, error) {
	if email == "" {
		return true, nil
	}

	whereClause := "email = ? and del_flag = '0'"
	args := []interface{}{email}
	if excludeUserId > 0 {
		whereClause += " and user_id != ?"
		args = append(args, excludeUserId)
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", m.table, whereClause)
	var count int64
	err := m.conn.QueryRowPartialCtx(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// UpdateUserStatus 修改用户状态
func (m *customSysUserModel) UpdateUserStatus(ctx context.Context, userId int64, status string) error {
	query := fmt.Sprintf("UPDATE %s SET status = ? WHERE user_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, status, userId)
	return err
}

// ResetUserPwd 重置用户密码
func (m *customSysUserModel) ResetUserPwd(ctx context.Context, userId int64, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password = ? WHERE user_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, password, userId)
	return err
}

// UpdateUserProfile 修改用户基本信息
func (m *customSysUserModel) UpdateUserProfile(ctx context.Context, userId int64, nickName, email, phonenumber, sex string) error {
	setClause := ""
	var args []interface{}

	if nickName != "" {
		setClause += "nick_name = ?, "
		args = append(args, nickName)
	}
	if email != "" {
		setClause += "email = ?, "
		args = append(args, email)
	}
	if phonenumber != "" {
		setClause += "phonenumber = ?, "
		args = append(args, phonenumber)
	}
	if sex != "" {
		setClause += "sex = ?, "
		args = append(args, sex)
	}

	if len(args) == 0 {
		return nil
	}

	// 移除最后的逗号和空格
	setClause = setClause[:len(setClause)-2]
	args = append(args, userId)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE user_id = ?", m.table, setClause)
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

// UpdateUserAvatar 修改用户头像
func (m *customSysUserModel) UpdateUserAvatar(ctx context.Context, userId int64, avatar int64) error {
	query := fmt.Sprintf("UPDATE %s SET avatar = ? WHERE user_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, avatar, userId)
	return err
}

// SelectUserListByDept 根据部门查询用户列表
func (m *customSysUserModel) SelectUserListByDept(ctx context.Context, deptId int64) ([]*SysUser, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE dept_id = ? AND del_flag = '0' AND status = '0'
		ORDER BY user_id ASC
	`, sysUserRows, m.table)

	var userList []*SysUser
	err := m.conn.QueryRowsPartialCtx(ctx, &userList, query, deptId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return userList, nil
}

// UpdateById 根据ID更新用户，只更新非零值字段
func (m *customSysUserModel) UpdateById(ctx context.Context, data *SysUser) error {
	if data.UserId == 0 {
		return fmt.Errorf("user_id cannot be zero")
	}

	var setParts []string
	var args []interface{}

	// 检查每个字段是否为非零值，如果是则加入更新列表
	if data.TenantId != "" {
		setParts = append(setParts, "`tenant_id` = ?")
		args = append(args, data.TenantId)
	}
	if data.UserName != "" {
		setParts = append(setParts, "`user_name` = ?")
		args = append(args, data.UserName)
	}
	if data.NickName != "" {
		setParts = append(setParts, "`nick_name` = ?")
		args = append(args, data.NickName)
	}
	if data.UserType != "" {
		setParts = append(setParts, "`user_type` = ?")
		args = append(args, data.UserType)
	}
	if data.Email != "" {
		setParts = append(setParts, "`email` = ?")
		args = append(args, data.Email)
	}
	if data.Phonenumber != "" {
		setParts = append(setParts, "`phonenumber` = ?")
		args = append(args, data.Phonenumber)
	}
	if data.Sex != "" {
		setParts = append(setParts, "`sex` = ?")
		args = append(args, data.Sex)
	}
	if data.Password != "" {
		setParts = append(setParts, "`password` = ?")
		args = append(args, data.Password)
	}
	if data.Status != "" {
		setParts = append(setParts, "`status` = ?")
		args = append(args, data.Status)
	}
	if data.DelFlag != "" {
		setParts = append(setParts, "`del_flag` = ?")
		args = append(args, data.DelFlag)
	}
	if data.LoginIp != "" {
		setParts = append(setParts, "`login_ip` = ?")
		args = append(args, data.LoginIp)
	}
	if data.DeptId.Valid {
		setParts = append(setParts, "`dept_id` = ?")
		args = append(args, data.DeptId.Int64)
	}
	if data.Avatar.Valid {
		setParts = append(setParts, "`avatar` = ?")
		args = append(args, data.Avatar.Int64)
	}
	if data.LoginDate.Valid {
		setParts = append(setParts, "`login_date` = ?")
		args = append(args, data.LoginDate.Time)
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
	query := fmt.Sprintf("UPDATE %s SET %s WHERE `user_id` = ?", m.table, setClause)
	args = append(args, data.UserId)

	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

// parseIds 解析ID串（逗号分隔）为int64数组
func parseIds(idStr string) []int64 {
	if idStr == "" {
		return []int64{}
	}

	var ids []int64
	parts := strings.Split(idStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}
