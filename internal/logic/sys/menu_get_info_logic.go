// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"fmt"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询菜单详细
func NewMenuGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuGetInfoLogic {
	return &MenuGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuGetInfoLogic) MenuGetInfo(req *types.MenuGetInfoReq) (resp *types.MenuResp, err error) {
	// 1. 查询菜单信息
	menu, err := l.svcCtx.SysMenuModel.FindOne(l.ctx, req.MenuId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.MenuResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "菜单不存在",
				},
			}, nil
		}
		l.Errorf("查询菜单信息失败: %v", err)
		return &types.MenuResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询菜单信息失败",
			},
		}, err
	}

	// 2. 转换为响应格式
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

	return &types.MenuResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: menuVo,
	}, nil
}
