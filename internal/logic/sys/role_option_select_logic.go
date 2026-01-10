// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"strconv"
	"strings"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleOptionSelectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取角色选择框列表
func NewRoleOptionSelectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleOptionSelectLogic {
	return &RoleOptionSelectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleOptionSelectLogic) RoleOptionSelect(req *types.RoleOptionSelectReq) (resp *types.RoleOptionSelectResp, err error) {
	var roleIds []int64

	// 1. 解析角色ID列表（如果有）
	if req.RoleIds != "" {
		roleIdStrs := strings.Split(req.RoleIds, ",")
		for _, idStr := range roleIdStrs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				continue
			}
			roleIds = append(roleIds, id)
		}
	}

	// 2. 查询角色列表
	var roles []*model.SysRole
	if len(roleIds) > 0 {
		// 根据角色ID列表查询
		roles, err = l.svcCtx.SysRoleModel.FindByIds(l.ctx, roleIds)
	} else {
		// 查询所有正常状态的角色
		roles, err = l.svcCtx.SysRoleModel.FindAll(l.ctx, &model.RoleQuery{Status: "0"})
	}
	if err != nil {
		l.Errorf("查询角色列表失败: %v", err)
		return &types.RoleOptionSelectResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询角色列表失败",
			},
		}, err
	}

	// 3. 转换为响应格式
	data := make([]types.RoleVo, 0, len(roles))
	for _, role := range roles {
		roleVo := convertRoleToVo(role)
		data = append(data, roleVo)
	}

	return &types.RoleOptionSelectResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: data,
	}, nil
}
