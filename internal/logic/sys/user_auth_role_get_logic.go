package sys

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAuthRoleGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据用户编号获取授权角色
func NewUserAuthRoleGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAuthRoleGetLogic {
	return &UserAuthRoleGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserAuthRoleGetLogic) UserAuthRoleGet(req *types.UserAuthRoleGetReq) (resp *types.UserDetailResp, err error) {
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

	// 2. 查询所有角色
	allRoles, err := l.svcCtx.SysRoleModel.FindAll(l.ctx, &model.RoleQuery{})
	if err != nil {
		l.Errorf("查询所有角色失败: %v", err)
		return &types.UserDetailResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询所有角色失败",
			},
		}, err
	}

	// 3. 查询用户已分配的角色
	userRoles, err := l.svcCtx.SysRoleModel.SelectRolesByUserId(l.ctx, req.UserId)
	if err != nil {
		l.Errorf("查询用户角色失败: %v", err)
		return &types.UserDetailResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询用户角色失败",
			},
		}, err
	}

	// 4. 构建用户已分配的角色ID集合
	userRoleIds := make(map[int64]bool)
	for _, role := range userRoles {
		userRoleIds[role.RoleId] = true
	}

	// 5. 转换角色列表，设置flag字段
	roleVos := make([]types.SysRoleVo, 0, len(allRoles))
	currentUserId, _ := util.GetUserIdFromContext(l.ctx)
	isCurrentSuperAdmin, _ := l.svcCtx.SysRoleModel.CheckIsSuperAdmin(l.ctx, currentUserId)

	for _, role := range allRoles {
		// 非超级管理员，过滤掉超级管理员角色
		if !isCurrentSuperAdmin && (role.RoleId == 1 || role.RoleKey == "superadmin") {
			continue
		}

		roleVo := types.SysRoleVo{
			RoleId:    role.RoleId,
			RoleName:  role.RoleName,
			RoleKey:   role.RoleKey,
			RoleSort:  int32(role.RoleSort),
			DataScope: role.DataScope,
			Status:    role.Status,
			Flag:      userRoleIds[role.RoleId], // 用户是否存在此角色标识
		}

		// 判断是否为超级管理员
		roleVo.SuperAdmin = role.RoleId == 1 || role.RoleKey == "superadmin"

		if role.Remark.Valid {
			roleVo.Remark = role.Remark.String
		}
		if role.CreateTime.Valid {
			roleVo.CreateTime = role.CreateTime.Time.Format("2006-01-02 15:04:05")
		}
		roleVos = append(roleVos, roleVo)
	}

	// 6. 转换为用户响应格式
	userVo := convertUserToVo(user)
	userVo.Roles = roleVos

	return &types.UserDetailResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: types.UserInfoDataVo{
			User: userVo,
		},
	}, nil
}
