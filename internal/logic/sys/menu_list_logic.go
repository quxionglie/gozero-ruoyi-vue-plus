// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"fmt"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询菜单列表
func NewMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuListLogic {
	return &MenuListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuListLogic) MenuList(req *types.MenuListReq) (resp *types.MenuListResp, err error) {
	// 1. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 2. 检查是否是超级管理员
	isSuperAdmin, err := l.svcCtx.SysRoleModel.CheckIsSuperAdmin(l.ctx, userId)
	if err != nil {
		l.Errorf("检查是否是超级管理员失败: %v", err)
		// 默认按非超级管理员处理
		isSuperAdmin = false
	}

	// 3. 构建查询条件
	menuQuery := &model.MenuQuery{
		MenuName: req.MenuName,
		Visible:  req.Visible,
		Status:   req.Status,
		MenuType: req.MenuType,
		ParentId: req.ParentId,
	}

	// 4. 查询菜单列表（超级管理员传入 userId=0 表示不过滤）
	queryUserId := userId
	if isSuperAdmin {
		queryUserId = 0
	}
	menus, err := l.svcCtx.SysMenuModel.FindAll(l.ctx, menuQuery, queryUserId)
	if err != nil {
		l.Errorf("查询菜单列表失败: %v", err)
		return &types.MenuListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询菜单列表失败",
			},
		}, err
	}

	// 4. 转换为响应格式
	rows := make([]types.MenuVo, 0, len(menus))
	for _, menu := range menus {
		menuVo := l.convertToMenuVo(menu)
		rows = append(rows, menuVo)
	}

	return &types.MenuListResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: rows,
	}, nil
}

// convertToMenuVo 转换菜单实体为响应格式
func (l *MenuListLogic) convertToMenuVo(menu *model.SysMenu) types.MenuVo {
	menuVo := types.MenuVo{
		MenuId:     menu.MenuId,
		ParentId:   menu.ParentId,
		MenuName:   menu.MenuName,
		OrderNum:   int32(menu.OrderNum),
		Path:       menu.Path,
		Component:  "",
		QueryParam: "",
		IsFrame:    "",
		IsCache:    "",
		MenuType:   menu.MenuType,
		Visible:    menu.Visible,
		Status:     menu.Status,
		Perms:      "",
		Icon:       menu.Icon,
		CreateTime: "",
		Remark:     menu.Remark,
		Children:   []types.MenuVo{},
	}

	if menu.Component.Valid {
		menuVo.Component = menu.Component.String
	}
	if menu.QueryParam.Valid {
		menuVo.QueryParam = menu.QueryParam.String
	}
	if menu.Perms.Valid {
		menuVo.Perms = menu.Perms.String
	}
	menuVo.IsFrame = fmt.Sprintf("%d", menu.IsFrame)
	menuVo.IsCache = fmt.Sprintf("%d", menu.IsCache)
	if menu.CreateTime.Valid {
		menuVo.CreateTime = menu.CreateTime.Time.Format("2006-01-02 15:04:05")
	}

	return menuVo
}
