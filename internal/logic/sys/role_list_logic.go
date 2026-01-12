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

type RoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取角色信息列表
func NewRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleListLogic {
	return &RoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleListLogic) RoleList(req *types.RoleListReq) (resp *types.RoleListResp, err error) {
	// 1. 构建查询条件
	roleQuery := &model.RoleQuery{
		RoleId:   req.RoleId,
		RoleName: req.RoleName,
		RoleKey:  req.RoleKey,
		Status:   req.Status,
	}

	// 2. 构建分页查询条件
	pageQuery := &model.PageQuery{
		PageNum:       req.PageNum,
		PageSize:      req.PageSize,
		OrderByColumn: req.OrderByColumn,
		IsAsc:         req.IsAsc,
	}

	// 3. 查询数据
	roleList, total, err := l.svcCtx.SysRoleModel.FindPage(l.ctx, roleQuery, pageQuery)
	if err != nil {
		l.Errorf("查询角色列表失败: %v", err)
		return &types.RoleListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询角色列表失败",
			},
		}, err
	}

	// 4. 转换为响应格式
	rows := make([]types.RoleVo, 0, len(roleList))
	for _, role := range roleList {
		roleVo := l.convertToRoleVo(role)
		rows = append(rows, roleVo)
	}

	return &types.RoleListResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Total: total,
		Rows:  rows,
	}, nil
}

// convertToRoleVo 转换角色实体为响应格式（复用函数）
func (l *RoleListLogic) convertToRoleVo(role *model.SysRole) types.RoleVo {
	return convertRoleToVo(role)
}
