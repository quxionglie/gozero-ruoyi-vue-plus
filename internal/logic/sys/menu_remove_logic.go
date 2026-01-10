// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除菜单
func NewMenuRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuRemoveLogic {
	return &MenuRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuRemoveLogic) MenuRemove(req *types.MenuRemoveReq) (resp *types.BaseResp, err error) {
	// 1. 检查是否存在子菜单
	hasChild, err := l.svcCtx.SysMenuModel.HasChildByMenuId(l.ctx, req.MenuId)
	if err != nil {
		l.Errorf("检查是否存在子菜单失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "检查是否存在子菜单失败",
		}, err
	}
	if hasChild {
		return &types.BaseResp{
			Code: 500,
			Msg:  "存在子菜单,不允许删除",
		}, nil
	}

	// 2. 检查菜单是否分配给角色
	existRole, err := l.svcCtx.SysMenuModel.CheckMenuExistRole(l.ctx, req.MenuId)
	if err != nil {
		l.Errorf("检查菜单是否分配给角色失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "检查菜单是否分配给角色失败",
		}, err
	}
	if existRole {
		return &types.BaseResp{
			Code: 500,
			Msg:  "菜单已分配,不允许删除",
		}, nil
	}

	// 3. 删除菜单
	err = l.svcCtx.SysMenuModel.Delete(l.ctx, req.MenuId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BaseResp{
				Code: 500,
				Msg:  "菜单不存在",
			}, nil
		}
		l.Errorf("删除菜单失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "删除菜单失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
