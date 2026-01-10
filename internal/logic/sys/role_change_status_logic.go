// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"strings"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleChangeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 状态修改
func NewRoleChangeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleChangeStatusLogic {
	return &RoleChangeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleChangeStatusLogic) RoleChangeStatus(req *types.RoleChangeStatusReq) (resp *types.BaseResp, err error) {
	// 1. 查询角色是否存在
	role, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.RoleId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BaseResp{
				Code: 404,
				Msg:  "角色不存在",
			}, nil
		}
		l.Errorf("查询角色信息失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询角色信息失败",
		}, err
	}

	// 2. 校验角色是否允许操作（超级管理员角色不能操作）
	if role.RoleId == 1 || strings.ToLower(role.RoleKey) == "superadmin" {
		return &types.BaseResp{
			Code: 500,
			Msg:  "不允许操作超级管理员角色",
		}, nil
	}

	// 3. 如果角色状态改为停用，需要检查是否有用户使用该角色
	if req.Status == "1" {
		count, err := l.svcCtx.SysRoleModel.CountUserRoleByRoleId(l.ctx, req.RoleId)
		if err != nil {
			l.Errorf("统计角色使用数量失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "统计角色使用数量失败",
			}, err
		}
		if count > 0 {
			return &types.BaseResp{
				Code: 500,
				Msg:  "角色已分配，不能禁用!",
			}, nil
		}
	}

	// 4. 更新角色状态
	err = l.svcCtx.SysRoleModel.UpdateRoleStatus(l.ctx, req.RoleId, req.Status)
	if err != nil {
		l.Errorf("更新角色状态失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "更新角色状态失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
