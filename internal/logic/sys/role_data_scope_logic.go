// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"database/sql"
	"strings"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleDataScopeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改保存数据权限
func NewRoleDataScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleDataScopeLogic {
	return &RoleDataScopeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleDataScopeLogic) RoleDataScope(req *types.RoleReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.RoleId == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "角色ID不能为空",
		}, nil
	}

	// 2. 查询角色是否存在
	role, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.RoleId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BaseResp{
				Code: 404,
				Msg:  "角色不存在",
			}, nil
		}
		l.Errorf("查询角色信息失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询角色信息失败",
		}, err
	}

	// 3. 校验角色是否允许操作（超级管理员角色不能操作）
	if role.RoleId == 1 || strings.ToLower(role.RoleKey) == "superadmin" {
		return &types.BaseResp{
			Code: 500,
			Msg:  "不允许操作超级管理员角色",
		}, nil
	}

	// 4. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 5. 更新角色数据权限信息
	if req.DataScope == "" {
		req.DataScope = role.DataScope
	}

	updateRole := &model.SysRole{
		RoleId:            req.RoleId,
		TenantId:          role.TenantId, // 保持原租户ID
		RoleName:          role.RoleName, // 保持原值
		RoleKey:           role.RoleKey,  // 保持原值
		RoleSort:          role.RoleSort, // 保持原值
		DataScope:         req.DataScope,
		MenuCheckStrictly: role.MenuCheckStrictly, // 保持原值
		DeptCheckStrictly: role.DeptCheckStrictly, // 保持原值
		Status:            role.Status,            // 保持原值
		DelFlag:           role.DelFlag,           // 保持原值
		CreateDept:        role.CreateDept,        // 保持原值
		CreateBy:          role.CreateBy,          // 保持原值
		UpdateBy:          sql.NullInt64{Int64: userId, Valid: userId > 0},
		Remark:            role.Remark, // 保持原值
	}

	// 6. 更新角色
	err = l.svcCtx.SysRoleModel.Update(l.ctx, updateRole)
	if err != nil {
		l.Errorf("更新角色数据权限失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "更新角色数据权限失败",
		}, err
	}

	// 7. 删除角色与部门关联
	err = l.svcCtx.SysRoleDeptModel.DeleteByRoleId(l.ctx, req.RoleId)
	if err != nil {
		l.Errorf("删除角色部门关联失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "删除角色部门关联失败",
		}, err
	}

	// 8. 批量插入角色部门关联（如果是自定义数据权限）
	if req.DataScope == "2" && len(req.DeptIds) > 0 {
		err = l.svcCtx.SysRoleDeptModel.InsertBatch(l.ctx, req.RoleId, req.DeptIds)
		if err != nil {
			l.Errorf("插入角色部门关联失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "插入角色部门关联失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
