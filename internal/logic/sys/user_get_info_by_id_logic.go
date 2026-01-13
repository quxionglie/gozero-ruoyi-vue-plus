package sys

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetInfoByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据用户编号获取详细信息
func NewUserGetInfoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetInfoByIdLogic {
	return &UserGetInfoByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserGetInfoByIdLogic) UserGetInfoById(req *types.UserGetInfoReq) (resp *types.UserDetailResp, err error) {
	userInfoDataVo := types.UserInfoDataVo{}

	// 如果 userId 不为 0，查询用户相关信息
	if req.UserId > 0 {
		// 1. 查询用户信息
		user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, req.UserId)
		if err != nil {
			if err == model.ErrNotFound {
				return &types.UserDetailResp{
					BaseResp: types.BaseResp{
						Code: 404,
						Msg:  "用户不存在",
					},
				}, nil
			}
			l.Errorf("查询用户信息失败: %v", err)
			return &types.UserDetailResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "查询用户信息失败",
				},
			}, err
		}

		// 2. 转换为响应格式
		userVo := convertUserToVo(l.ctx, l.svcCtx, user)

		// 3. 查询部门名称
		if user.DeptId.Valid && user.DeptId.Int64 > 0 {
			dept, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, user.DeptId.Int64)
			if err == nil {
				userVo.DeptName = dept.DeptName
			}
		}

		// 4. 查询用户已拥有的角色ID列表
		userRoles, err := l.svcCtx.SysUserRoleModel.FindByUserId(l.ctx, req.UserId)
		if err != nil {
			l.Errorf("查询用户角色关联失败: %v", err)
		} else {
			roleIds := make([]int64, 0, len(userRoles))
			for _, ur := range userRoles {
				roleIds = append(roleIds, ur.RoleId)
			}
			userInfoDataVo.RoleIds = util.Int64SliceToStringSlice(roleIds)
		}

		// 5. 查询用户已拥有的角色（用于 user.roles）
		roles, err := l.svcCtx.SysRoleModel.SelectRolesByUserId(l.ctx, req.UserId)
		if err != nil {
			l.Errorf("查询用户角色失败: %v", err)
		} else {
			roleVos := make([]types.SysRoleVo, 0, len(roles))
			for _, role := range roles {
				roleVo := convertRoleVoToSysRoleVo(role)
				roleVos = append(roleVos, roleVo)
			}
			userVo.Roles = roleVos
		}

		// 6. 查询用户已拥有的岗位ID列表
		postIds, err := l.svcCtx.SysUserPostModel.SelectPostIdsByUserId(l.ctx, req.UserId)
		if err != nil {
			l.Errorf("查询用户岗位ID列表失败: %v", err)
		} else {
			userInfoDataVo.PostIds = util.Int64SliceToStringSlice(postIds)
		}

		// 7. 根据用户的部门ID查询岗位列表
		if user.DeptId.Valid && user.DeptId.Int64 > 0 {
			postQuery := &model.PostQuery{
				DeptId: user.DeptId.Int64,
			}
			posts, err := l.svcCtx.SysPostModel.FindAll(l.ctx, postQuery)
			if err != nil {
				l.Errorf("查询岗位列表失败: %v", err)
			} else {
				postVos := make([]types.PostVo, 0, len(posts))
				for _, post := range posts {
					postVo := convertPostToVo(l.ctx, l.svcCtx, post)
					postVos = append(postVos, postVo)
				}
				userInfoDataVo.Posts = postVos
			}
		}

		userInfoDataVo.User = userVo
	}

	// 8. 查询所有正常状态的角色列表（无论 userId 是否为 0）
	roleQuery := &model.RoleQuery{
		Status: "0", // 正常状态
	}
	allRoles, err := l.svcCtx.SysRoleModel.FindAll(l.ctx, roleQuery)
	if err != nil {
		l.Errorf("查询所有角色列表失败: %v", err)
	} else {
		// 如果是超级管理员，返回所有角色；否则过滤掉超级管理员角色
		isSuperAdmin := false
		if req.UserId > 0 {
			isSuperAdmin, _ = l.svcCtx.SysRoleModel.CheckIsSuperAdmin(l.ctx, req.UserId)
		}

		roleVos := make([]types.RoleVo, 0)
		for _, role := range allRoles {
			// 如果不是超级管理员，过滤掉超级管理员角色
			if !isSuperAdmin && (role.RoleId == 1 || role.RoleKey == "superadmin") {
				continue
			}
			roleVo := convertRoleToVo(role) // 使用 role_common.go 中的函数
			roleVos = append(roleVos, roleVo)
		}
		userInfoDataVo.Roles = roleVos
	}

	return &types.UserDetailResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: userInfoDataVo,
	}, nil
}

// convertRoleVoToSysRoleVo 转换 RoleVo 为 SysRoleVo（用于 user.roles）
func convertRoleVoToSysRoleVo(role *model.RoleVo) types.SysRoleVo {
	roleVo := types.SysRoleVo{
		RoleId:            role.RoleId,
		RoleName:          role.RoleName,
		RoleKey:           role.RoleKey,
		RoleSort:          int32(role.RoleSort),
		DataScope:         role.DataScope,
		MenuCheckStrictly: nil,
		DeptCheckStrictly: nil,
		Status:            role.Status,
		Remark:            nil,
		CreateTime:        nil,
		SuperAdmin:        role.RoleId == 1 || role.RoleKey == "superadmin",
		Flag:              false,
	}
	if role.Remark.Valid {
		remarkStr := role.Remark.String
		roleVo.Remark = &remarkStr
	}
	if role.CreateTime.Valid {
		createTimeStr := role.CreateTime.Time.Format("2006-01-02 15:04:05")
		roleVo.CreateTime = &createTimeStr
	}
	return roleVo
}

// convertPostToVo 转换岗位为 PostVo
func convertPostToVo(ctx context.Context, svcCtx *svc.ServiceContext, post *model.SysPost) types.PostVo {
	postVo := types.PostVo{
		PostId:       post.PostId,
		DeptId:       post.DeptId,
		PostCode:     post.PostCode,
		PostCategory: "",
		PostName:     post.PostName,
		PostSort:     int32(post.PostSort),
		Status:       post.Status,
		Remark:       "",
		CreateTime:   "",
		DeptName:     "",
	}
	// 查询部门名称
	if post.DeptId > 0 {
		dept, err := svcCtx.SysDeptModel.FindOne(ctx, post.DeptId)
		if err == nil {
			postVo.DeptName = dept.DeptName
		}
	}
	if post.PostCategory.Valid {
		postVo.PostCategory = post.PostCategory.String
	}
	if post.Remark.Valid {
		postVo.Remark = post.Remark.String
	}
	if post.CreateTime.Valid {
		postVo.CreateTime = post.CreateTime.Time.Format("2006-01-02 15:04:05")
	}
	return postVo
}
