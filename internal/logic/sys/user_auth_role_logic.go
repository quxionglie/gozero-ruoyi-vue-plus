package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAuthRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户授权角色
func NewUserAuthRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAuthRoleLogic {
	return &UserAuthRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserAuthRoleLogic) UserAuthRole(req *types.UserAuthRoleReq) (resp *types.BaseResp, err error) {
	// 1. 校验用户ID
	if req.UserId <= 0 {
		return &types.BaseResp{
			Code: 500,
			Msg:  "用户ID不能为空",
		}, nil
	}

	// 2. 获取当前用户ID，检查是否为超级管理员
	currentUserId, err := util.GetUserIdFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "获取用户ID失败",
		}, err
	}

	isSuperAdmin, err := l.svcCtx.SysRoleModel.CheckIsSuperAdmin(l.ctx, currentUserId)
	if err != nil {
		l.Errorf("检查超级管理员失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "检查超级管理员失败",
		}, err
	}

	// 3. 处理角色ID列表（非超级管理员，禁止包含超级管理员角色）
	var roleIds []int64
	for _, roleId := range req.RoleIds {
		// 非超级管理员，禁止包含超级管理员角色（roleId=1）
		if !isSuperAdmin && roleId == 1 {
			continue
		}
		roleIds = append(roleIds, roleId)
	}

	// 4. 先删除用户的所有角色关联
	err = l.svcCtx.SysUserRoleModel.DeleteByUserId(l.ctx, req.UserId)
	if err != nil {
		l.Errorf("删除用户角色关联失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "删除用户角色关联失败",
		}, err
	}

	// 5. 批量插入新的角色关联
	if len(roleIds) > 0 {
		err = l.svcCtx.SysUserRoleModel.InsertBatchByUserId(l.ctx, req.UserId, roleIds)
		if err != nil {
			l.Errorf("新增用户角色关联失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "新增用户角色关联失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "授权成功",
	}, nil
}
