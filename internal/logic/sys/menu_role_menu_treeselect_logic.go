// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuRoleMenuTreeselectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取角色菜单树
func NewMenuRoleMenuTreeselectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuRoleMenuTreeselectLogic {
	return &MenuRoleMenuTreeselectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuRoleMenuTreeselectLogic) MenuRoleMenuTreeselect(req *types.MenuRoleMenuTreeselectReq) (resp *types.MenuTreeselectResp, err error) {
	// 1. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 2. 检查是否是超级管理员
	isSuperAdmin, err := l.svcCtx.SysRoleModel.CheckIsSuperAdmin(l.ctx, userId)
	if err != nil {
		l.Errorf("检查是否是超级管理员失败: %v", err)
		isSuperAdmin = false
	}

	// 3. 查询所有菜单（超级管理员传入 userId=0 表示不过滤）
	queryUserId := userId
	if isSuperAdmin {
		queryUserId = 0
	}
	menus, err := l.svcCtx.SysMenuModel.FindAll(l.ctx, &model.MenuQuery{}, queryUserId)
	if err != nil {
		l.Errorf("查询菜单列表失败: %v", err)
		return &types.MenuTreeselectResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询菜单列表失败",
			},
		}, err
	}

	// 4. 查询角色已分配的菜单ID列表
	menuIds, err := l.svcCtx.SysMenuModel.SelectMenuListByRoleId(l.ctx, req.RoleId)
	if err != nil {
		l.Errorf("查询角色菜单列表失败: %v", err)
		return &types.MenuTreeselectResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询角色菜单列表失败",
			},
		}, err
	}

	// 5. 构建菜单ID集合（用于快速查找）
	menuIdSet := make(map[int64]bool)
	for _, menuId := range menuIds {
		menuIdSet[menuId] = true
	}

	// 6. 构建树形结构
	treeList := l.buildMenuTreeSelect(menus, menuIdSet)

	return &types.MenuTreeselectResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: treeList,
	}, nil
}

// buildMenuTreeSelect 构建菜单树形结构
func (l *MenuRoleMenuTreeselectLogic) buildMenuTreeSelect(menus []*model.SysMenu, checkedMenuIds map[int64]bool) []types.MenuTreeVo {
	if len(menus) == 0 {
		return []types.MenuTreeVo{}
	}

	// 构建树形结构
	var buildTree func(parentId int64) []types.MenuTreeVo
	buildTree = func(parentId int64) []types.MenuTreeVo {
		var children []types.MenuTreeVo
		for _, menu := range menus {
			if menu.ParentId == parentId {
				treeNode := types.MenuTreeVo{
					Id:       menu.MenuId,
					Label:    menu.MenuName,
					ParentId: menu.ParentId,
					Children: buildTree(menu.MenuId),
				}
				children = append(children, treeNode)
			}
		}
		return children
	}

	return buildTree(0)
}
