// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"strconv"
	"strings"

	"gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResp, err error) {
	// 1. 从 JWT token 中获取用户信息（go-zero 会将 JWT claims 中的字段存储到 context 中）
	userId, err := util.GetUserIdFromContext(l.ctx)
	if err != nil {
		return &types.UserInfoResp{
			BaseResp: types.BaseResp{
				Code: 401,
				Msg:  "未授权，请先登录",
			},
		}, nil
	}

	// 2. 查询用户信息
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == sys.ErrNotFound {
			return &types.UserInfoResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "没有权限访问用户数据",
				},
			}, nil
		}
		l.Errorf("查询用户信息失败: %v", err)
		return &types.UserInfoResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询用户信息失败",
			},
		}, err
	}

	// 3. 查询部门名称
	deptName := ""
	if user.DeptId.Valid {
		dept, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, user.DeptId.Int64)
		if err == nil {
			deptName = dept.DeptName
		}
	}

	// 4. 构建用户信息响应
	userVo := types.SysUserVo{
		UserId:      user.UserId,
		TenantId:    user.TenantId,
		DeptId:      0,
		UserName:    user.UserName,
		NickName:    user.NickName,
		UserType:    user.UserType,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Sex:         user.Sex,
		Avatar:      "",
		Status:      user.Status,
		LoginIp:     user.LoginIp,
		LoginDate:   "",
		Remark:      "",
		CreateTime:  "",
		DeptName:    deptName,
	}

	if user.DeptId.Valid {
		userVo.DeptId = user.DeptId.Int64
	}
	if user.Avatar.Valid {
		userVo.Avatar = strconv.FormatInt(user.Avatar.Int64, 10)
	} else {
		userVo.Avatar = ""
	}
	if user.LoginDate.Valid {
		userVo.LoginDate = user.LoginDate.Time.Format("2006-01-02 15:04:05")
	}
	if user.Remark.Valid {
		userVo.Remark = user.Remark.String
	}
	if user.CreateTime.Valid {
		userVo.CreateTime = user.CreateTime.Time.Format("2006-01-02 15:04:05")
	}

	// 5. 查询用户的角色列表（填充到 SysUserVo.roles）
	userRoles, err := l.getUserRoles(userId)
	if err != nil {
		l.Errorf("查询用户角色列表失败: %v", err)
		userRoles = []types.SysRoleVo{}
	}
	userVo.Roles = userRoles

	// 6. 检查是否为超级管理员（role_id = 1 或 role_key = 'superadmin'）
	isSuperAdmin := l.isSuperAdmin(userId, userRoles)

	// 7. 查询菜单权限
	var permissions []string
	if isSuperAdmin {
		// 超级管理员拥有所有权限
		permissions = []string{"*:*:*"}
	} else {
		permissions, err = l.getMenuPermissions(userId)
		if err != nil {
			l.Errorf("查询菜单权限失败: %v", err)
			permissions = []string{}
		}
	}

	// 8. 查询角色权限（role_key 字符串数组）
	var roleKeys []string
	if isSuperAdmin {
		// 超级管理员角色标识
		roleKeys = []string{"superadmin"}
	} else {
		roleKeys, err = l.getRolePermissions(userId)
		if err != nil {
			l.Errorf("查询角色权限失败: %v", err)
			roleKeys = []string{}
		}
	}

	return &types.UserInfoResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: types.UserInfoVo{
			User:        userVo,
			Permissions: permissions,
			Roles:       roleKeys,
		},
	}, nil
}

// getUserRoles 查询用户的角色列表
func (l *GetUserInfoLogic) getUserRoles(userId int64) ([]types.SysRoleVo, error) {
	// 使用 model 方法查询角色列表
	rows, err := l.svcCtx.SysRoleModel.SelectRolesByUserId(l.ctx, userId)
	if err != nil {
		return nil, err
	}

	result := make([]types.SysRoleVo, 0, len(rows))
	for _, row := range rows {
		// 判断是否为超级管理员（role_id = 1 或 role_key = 'superadmin'）
		isSuperAdmin := row.RoleId == 1 || strings.ToLower(row.RoleKey) == "superadmin"

		roleVo := types.SysRoleVo{
			RoleId:     row.RoleId,
			RoleName:   row.RoleName,
			RoleKey:    row.RoleKey,
			RoleSort:   int32(row.RoleSort),
			DataScope:  row.DataScope,
			Status:     row.Status,
			Remark:     "",
			CreateTime: "",
			SuperAdmin: isSuperAdmin,
			Flag:       true, // 用户已拥有的角色，flag 为 true
		}
		if row.Remark.Valid {
			roleVo.Remark = row.Remark.String
		}
		if row.CreateTime.Valid {
			roleVo.CreateTime = row.CreateTime.Time.Format("2006-01-02 15:04:05")
		}
		result = append(result, roleVo)
	}

	return result, nil
}

// isSuperAdmin 检查用户是否为超级管理员
// 超级管理员：role_id = 1 或 role_key = 'superadmin'
func (l *GetUserInfoLogic) isSuperAdmin(userId int64, roles []types.SysRoleVo) bool {
	// 检查角色列表中是否有超级管理员角色
	for _, role := range roles {
		if role.RoleId == 1 || role.RoleKey == "superadmin" {
			return true
		}
	}
	return false
}

// getMenuPermissions 查询用户的菜单权限
func (l *GetUserInfoLogic) getMenuPermissions(userId int64) ([]string, error) {
	// 使用 model 方法查询菜单权限
	return l.svcCtx.SysMenuModel.SelectMenuPermissionsByUserId(l.ctx, userId)
}

// getRolePermissions 查询用户的角色权限
func (l *GetUserInfoLogic) getRolePermissions(userId int64) ([]string, error) {
	// 使用 model 方法查询角色权限标识
	return l.svcCtx.SysRoleModel.SelectRoleKeysByUserId(l.ctx, userId)
}
