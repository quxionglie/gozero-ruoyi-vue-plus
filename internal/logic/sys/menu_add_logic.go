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

type MenuAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增菜单
func NewMenuAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuAddLogic {
	return &MenuAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuAddLogic) MenuAdd(req *types.MenuReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
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

	// 3. 默认父菜单ID为0（根菜单）
	if req.ParentId == 0 {
		req.ParentId = 0
	}

	// 4. 校验菜单名称唯一性（同父菜单下唯一）
	unique, err := l.svcCtx.SysMenuModel.CheckMenuNameUnique(l.ctx, req.MenuName, req.ParentId, 0)
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
			Msg:  fmt.Sprintf("新增菜单'%s'失败，菜单名称已存在", req.MenuName),
		}, nil
	}

	// 5. 如果 isFrame 为 0（外链），path 必须以 http(s):// 开头
	if req.IsFrame == "0" && req.Path != "" {
		if !strings.HasPrefix(req.Path, "http://") && !strings.HasPrefix(req.Path, "https://") {
			return &types.BaseResp{
				Code: 500,
				Msg:  fmt.Sprintf("新增菜单'%s'失败，地址必须以http(s)://开头", req.MenuName),
			}, nil
		}
	}

	// 6. 获取当前用户信息
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 获取用户的部门ID
	var deptId int64
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err == nil && user.DeptId.Valid {
		deptId = user.DeptId.Int64
	}

	// 7. 生成主键ID（使用雪花算法）
	newMenuId, err := util.GenerateID()
	if err != nil {
		l.Errorf("生成菜单ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "生成菜单ID失败",
		}, err
	}

	// 8. 转换 isFrame 和 isCache
	var isFrame int64 = 1
	if req.IsFrame == "0" {
		isFrame = 0
	}
	var isCache int64 = 0
	if req.IsCache == "0" {
		isCache = 0
	}

	// 9. 构建菜单实体
	menu := &model.SysMenu{
		MenuId:     newMenuId,
		ParentId:   req.ParentId,
		MenuName:   req.MenuName,
		OrderNum:   int64(req.OrderNum),
		Path:       req.Path,
		Component:  sql.NullString{String: req.Component, Valid: req.Component != ""},
		QueryParam: sql.NullString{String: req.QueryParam, Valid: req.QueryParam != ""},
		IsFrame:    isFrame,
		IsCache:    isCache,
		MenuType:   req.MenuType,
		Visible:    req.Visible,
		Status:     req.Status,
		Perms:      sql.NullString{String: req.Perms, Valid: req.Perms != ""},
		Icon:       req.Icon,
		Remark:     req.Remark,
		CreateDept: sql.NullInt64{Int64: deptId, Valid: deptId > 0},
		CreateBy:   sql.NullInt64{Int64: userId, Valid: userId > 0},
	}
	if menu.Visible == "" {
		menu.Visible = "0"
	}
	if menu.Status == "" {
		menu.Status = "0"
	}
	if menu.OrderNum == 0 {
		menu.OrderNum = 0
	}

	// 10. 插入数据库
	_, err = l.svcCtx.SysMenuModel.Insert(l.ctx, menu)
	if err != nil {
		l.Errorf("新增菜单失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "新增菜单失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
