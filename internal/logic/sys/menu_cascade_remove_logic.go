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

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuCascadeRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量级联删除菜单
func NewMenuCascadeRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuCascadeRemoveLogic {
	return &MenuCascadeRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuCascadeRemoveLogic) MenuCascadeRemove(req *types.MenuCascadeRemoveReq) (resp *types.BaseResp, err error) {
	// 1. 解析菜单ID列表
	menuIdStrs := strings.Split(req.MenuIds, ",")
	var menuIds []int64
	for _, idStr := range menuIdStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		menuIds = append(menuIds, id)
	}

	if len(menuIds) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "菜单ID不能为空",
		}, nil
	}

	// 2. 检查是否存在子菜单（批量）
	hasChild, err := l.svcCtx.SysMenuModel.HasChildByMenuIds(l.ctx, menuIds)
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

	// 3. 批量删除菜单（级联删除，包括角色菜单关联）
	for _, menuId := range menuIds {
		// 删除菜单
		err = l.svcCtx.SysMenuModel.Delete(l.ctx, menuId)
		if err != nil && err != sys.ErrNotFound {
			l.Errorf("删除菜单失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除菜单失败",
			}, err
		}
		// TODO: 删除角色菜单关联（需要在 SysRoleMenuModel 中添加 DeleteByMenuIds 方法）
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
