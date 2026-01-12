// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改保存角色
func NewRoleEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleEditLogic {
	return &RoleEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleEditLogic) RoleEdit(req *types.RoleReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.RoleId == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "角色ID不能为空",
		}, nil
	}
	if req.RoleName == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "角色名称不能为空",
		}, nil
	}
	if req.RoleKey == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "角色权限标识不能为空",
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

	// 4. 校验不允许修改系统内置管理员角色标识符
	adminRoleKeys := []string{"superadmin", "tenantadmin"}
	if strings.EqualFold(role.RoleKey, "superadmin") || strings.EqualFold(role.RoleKey, "tenantadmin") {
		if !strings.EqualFold(role.RoleKey, req.RoleKey) {
			return &types.BaseResp{
				Code: 500,
				Msg:  "不允许修改系统内置管理员角色标识符!",
			}, nil
		}
	}

	// 5. 校验不允许使用系统内置管理员角色标识符
	for _, adminKey := range adminRoleKeys {
		if strings.EqualFold(req.RoleKey, adminKey) && !strings.EqualFold(role.RoleKey, adminKey) {
			return &types.BaseResp{
				Code: 500,
				Msg:  "不允许使用系统内置管理员角色标识符!",
			}, nil
		}
	}

	// 6. 如果角色状态改为停用，需要检查是否有用户使用该角色
	if req.Status == "1" && role.Status == "0" {
		count, err := l.svcCtx.SysRoleModel.CountUserRoleByRoleId(l.ctx, req.RoleId)
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
				Msg:  "角色已分配，不能禁用!",
			}, nil
		}
	}

	// 7. 参数长度校验
	if err := util.ValidateStringLength(req.RoleName, "角色名称", 30); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}
	if err := util.ValidateStringLength(req.RoleKey, "角色权限标识", 100); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}

	// 8. 校验角色名称唯一性
	unique, err := l.svcCtx.SysRoleModel.CheckRoleNameUnique(l.ctx, req.RoleName, req.RoleId)
	if err != nil {
		l.Errorf("校验角色名称唯一性失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "校验角色名称唯一性失败",
		}, err
	}
	if !unique {
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("修改角色'%s'失败，角色名称已存在", req.RoleName),
		}, nil
	}

	// 9. 校验角色权限标识唯一性
	unique, err = l.svcCtx.SysRoleModel.CheckRoleKeyUnique(l.ctx, req.RoleKey, req.RoleId)
	if err != nil {
		l.Errorf("校验角色权限标识唯一性失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "校验角色权限标识唯一性失败",
		}, err
	}
	if !unique {
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("修改角色'%s'失败，角色权限已存在", req.RoleName),
		}, nil
	}

	// 10. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 11. 构建角色实体（更新字段）
	roleSort := int64(req.RoleSort)
	if req.Status == "" {
		req.Status = role.Status
	}
	if req.DataScope == "" {
		req.DataScope = role.DataScope
	}

	updateRole := &model.SysRole{
		RoleId:            req.RoleId,
		TenantId:          role.TenantId, // 保持原租户ID
		RoleName:          req.RoleName,
		RoleKey:           req.RoleKey,
		RoleSort:          roleSort,
		DataScope:         req.DataScope,
		MenuCheckStrictly: role.MenuCheckStrictly, // 保持原值
		DeptCheckStrictly: role.DeptCheckStrictly, // 保持原值
		Status:            req.Status,
		DelFlag:           role.DelFlag,    // 保持原值
		CreateDept:        role.CreateDept, // 保持原值
		CreateBy:          role.CreateBy,   // 保持原值
		CreateTime:        role.CreateTime, // 保持原创建时间
		UpdateBy:          sql.NullInt64{Int64: userId, Valid: userId > 0},
		UpdateTime:        sql.NullTime{Time: time.Now(), Valid: true},
		Remark:            sql.NullString{String: req.Remark, Valid: req.Remark != ""},
	}

	// 12. 更新角色
	err = l.svcCtx.SysRoleModel.Update(l.ctx, updateRole)
	if err != nil {
		l.Errorf("更新角色失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("修改角色'%s'失败", req.RoleName),
		}, err
	}

	// 13. 删除角色与菜单关联
	err = l.svcCtx.SysRoleMenuModel.DeleteByRoleId(l.ctx, req.RoleId)
	if err != nil {
		l.Errorf("删除角色菜单关联失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "删除角色菜单关联失败",
		}, err
	}

	// 14. 批量插入角色菜单关联
	if len(req.MenuIds) > 0 {
		err = l.svcCtx.SysRoleMenuModel.InsertBatch(l.ctx, req.RoleId, req.MenuIds)
		if err != nil {
			l.Errorf("插入角色菜单关联失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "插入角色菜单关联失败",
			}, err
		}
	}

	// TODO: 清除在线用户缓存（如果角色状态或权限发生变化）

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
