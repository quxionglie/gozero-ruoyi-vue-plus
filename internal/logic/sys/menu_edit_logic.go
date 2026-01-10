// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改菜单
func NewMenuEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuEditLogic {
	return &MenuEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuEditLogic) MenuEdit(req *types.MenuReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.MenuId <= 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "菜单ID不能为空",
		}, nil
	}
	if req.MenuName == "" {
		return &types.BaseResp{Code: 400, Msg: "菜单名称不能为空"}, nil
	}
	if req.MenuType == "" {
		return &types.BaseResp{Code: 400, Msg: "菜单类型不能为空"}, nil
	}

	// 2. 参数长度校验
	if err := util.ValidateStringLength(req.MenuName, "菜单名称", 50); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}
	if req.Path != "" {
		if err := util.ValidateStringLength(req.Path, "路由地址", 200); err != nil {
			return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
		}
	}
	if req.Component != "" {
		if err := util.ValidateStringLength(req.Component, "组件路径", 200); err != nil {
			return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
		}
	}
	if req.Perms != "" {
		if err := util.ValidateStringLength(req.Perms, "权限标识", 100); err != nil {
			return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
		}
	}

	// 3. 检查菜单是否存在
	menu, err := l.svcCtx.SysMenuModel.FindOne(l.ctx, req.MenuId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BaseResp{
				Code: 500,
				Msg:  "菜单不存在",
			}, nil
		}
		l.Errorf("查询菜单失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询菜单失败",
		}, err
	}

	// 4. 默认父菜单ID为0（根菜单）
	if req.ParentId == 0 {
		req.ParentId = 0
	}

	// 5. 上级菜单不能是自己
	if req.ParentId == req.MenuId {
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("修改菜单'%s'失败，上级菜单不能选择自己", req.MenuName),
		}, nil
	}

	// 6. 校验菜单名称唯一性（同父菜单下唯一）
	unique, err := l.svcCtx.SysMenuModel.CheckMenuNameUnique(l.ctx, req.MenuName, req.ParentId, req.MenuId)
	if err != nil {
		l.Errorf("校验菜单名称唯一性失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "校验菜单名称唯一性失败",
		}, err
	}
	if !unique {
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("修改菜单'%s'失败，菜单名称已存在", req.MenuName),
		}, nil
	}

	// 7. 如果 isFrame 为 0（外链），path 必须以 http(s):// 开头
	if req.IsFrame == "0" && req.Path != "" {
		if !strings.HasPrefix(req.Path, "http://") && !strings.HasPrefix(req.Path, "https://") {
			return &types.BaseResp{
				Code: 500,
				Msg:  fmt.Sprintf("修改菜单'%s'失败，地址必须以http(s)://开头", req.MenuName),
			}, nil
		}
	}

	// 8. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 9. 转换 isFrame 和 isCache
	var isFrame int64 = 1
	if req.IsFrame == "0" {
		isFrame = 0
	}
	var isCache int64 = 0
	if req.IsCache == "0" {
		isCache = 0
	}

	// 10. 更新菜单信息
	menu.ParentId = req.ParentId
	menu.MenuName = req.MenuName
	menu.OrderNum = int64(req.OrderNum)
	menu.Path = req.Path
	menu.Component = sql.NullString{String: req.Component, Valid: req.Component != ""}
	menu.QueryParam = sql.NullString{String: req.QueryParam, Valid: req.QueryParam != ""}
	menu.IsFrame = isFrame
	menu.IsCache = isCache
	menu.MenuType = req.MenuType
	if req.Visible != "" {
		menu.Visible = req.Visible
	}
	if req.Status != "" {
		menu.Status = req.Status
	}
	menu.Perms = sql.NullString{String: req.Perms, Valid: req.Perms != ""}
	menu.Icon = req.Icon
	menu.Remark = req.Remark
	menu.UpdateBy = sql.NullInt64{Int64: userId, Valid: userId > 0}

	// 11. 更新数据库
	err = l.svcCtx.SysMenuModel.Update(l.ctx, menu)
	if err != nil {
		l.Errorf("修改菜单失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改菜单失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
