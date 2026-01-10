// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除角色
func NewRoleRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleRemoveLogic {
	return &RoleRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleRemoveLogic) RoleRemove(req *types.RoleRemoveReq) (resp *types.BaseResp, err error) {
	// 1. 解析角色ID列表
	roleIdStrs := strings.Split(req.RoleIds, ",")
	var roleIds []int64
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

	if len(roleIds) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "角色ID不能为空",
		}, nil
	}

	// 2. 查询角色信息，检查是否允许删除
	roles, err := l.svcCtx.SysRoleModel.FindByIds(l.ctx, roleIds)
	if err != nil {
		l.Errorf("查询角色信息失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询角色信息失败",
		}, err
	}

	// 3. 检查每个角色是否允许删除
	for _, role := range roles {
		// 检查是否为超级管理员角色
		if role.RoleId == 1 || strings.ToLower(role.RoleKey) == "superadmin" {
			return &types.BaseResp{
				Code: 500,
				Msg:  "不允许删除超级管理员角色",
			}, nil
		}

		// 检查角色是否被用户使用
		count, err := l.svcCtx.SysRoleModel.CountUserRoleByRoleId(l.ctx, role.RoleId)
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
				Msg:  fmt.Sprintf("%s已分配，不能删除!", role.RoleName),
			}, nil
		}
	}

	// 4. 批量删除角色与菜单关联
	err = l.svcCtx.SysRoleMenuModel.DeleteByRoleIds(l.ctx, roleIds)
	if err != nil {
		l.Errorf("删除角色菜单关联失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "删除角色菜单关联失败",
		}, err
	}

	// 5. 批量删除角色与部门关联
	err = l.svcCtx.SysRoleDeptModel.DeleteByRoleIds(l.ctx, roleIds)
	if err != nil {
		l.Errorf("删除角色部门关联失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "删除角色部门关联失败",
		}, err
	}

	// 6. 删除角色
	for _, roleId := range roleIds {
		err = l.svcCtx.SysRoleModel.Delete(l.ctx, roleId)
		if err != nil {
			l.Errorf("删除角色失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除角色失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
