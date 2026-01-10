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

type RoleAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增角色
func NewRoleAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAddLogic {
	return &RoleAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleAddLogic) RoleAdd(req *types.RoleReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
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

	// 2. 参数长度校验
	if err := util.ValidateStringLength(req.RoleName, "角色名称", 30); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}
	if err := util.ValidateStringLength(req.RoleKey, "角色权限标识", 100); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}

	// 3. 校验角色是否允许操作（不允许使用系统内置管理员角色标识符）
	adminRoleKeys := []string{"superadmin", "tenantadmin"}
	for _, adminKey := range adminRoleKeys {
		if strings.EqualFold(req.RoleKey, adminKey) {
			return &types.BaseResp{
				Code: 500,
				Msg:  "不允许使用系统内置管理员角色标识符!",
			}, nil
		}
	}

	// 4. 校验角色名称唯一性
	unique, err := l.svcCtx.SysRoleModel.CheckRoleNameUnique(l.ctx, req.RoleName, 0)
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
			Msg:  fmt.Sprintf("新增角色'%s'失败，角色名称已存在", req.RoleName),
		}, nil
	}

	// 5. 校验角色权限标识唯一性
	unique, err = l.svcCtx.SysRoleModel.CheckRoleKeyUnique(l.ctx, req.RoleKey, 0)
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
			Msg:  fmt.Sprintf("新增角色'%s'失败，角色权限已存在", req.RoleName),
		}, nil
	}

	// 6. 获取当前用户信息
	userId, _ := util.GetUserIdFromContext(l.ctx)
	tenantId, _ := util.GetTenantIdFromContext(l.ctx)

	// 获取用户的部门ID
	var deptId int64
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err == nil && user.DeptId.Valid {
		deptId = user.DeptId.Int64
	}

	// 7. 生成角色ID（使用雪花算法）
	newRoleId, err := util.GenerateID()
	if err != nil {
		l.Errorf("生成角色ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "生成角色ID失败",
		}, err
	}

	// 8. 构建角色实体
	// MenuCheckStrictly 和 DeptCheckStrictly 默认为 0（false）
	menuCheckStrictly := int64(0)
	deptCheckStrictly := int64(0)

	roleSort := int64(req.RoleSort)
	if roleSort == 0 {
		roleSort = 0
	}

	if req.Status == "" {
		req.Status = "0"
	}
	if req.DataScope == "" {
		req.DataScope = "1"
	}

	role := &model.SysRole{
		RoleId:            newRoleId,
		TenantId:          tenantId,
		RoleName:          req.RoleName,
		RoleKey:           req.RoleKey,
		RoleSort:          roleSort,
		DataScope:         req.DataScope,
		MenuCheckStrictly: menuCheckStrictly,
		DeptCheckStrictly: deptCheckStrictly,
		Status:            req.Status,
		DelFlag:           "0",
		CreateDept:        sql.NullInt64{Int64: deptId, Valid: deptId > 0},
		CreateBy:          sql.NullInt64{Int64: userId, Valid: userId > 0},
		Remark:            sql.NullString{String: req.Remark, Valid: req.Remark != ""},
	}

	// 9. 插入角色
	_, err = l.svcCtx.SysRoleModel.Insert(l.ctx, role)
	if err != nil {
		l.Errorf("插入角色失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("新增角色'%s'失败", req.RoleName),
		}, err
	}

	// 10. 批量插入角色菜单关联
	if len(req.MenuIds) > 0 {
		err = l.svcCtx.SysRoleMenuModel.InsertBatch(l.ctx, newRoleId, req.MenuIds)
		if err != nil {
			l.Errorf("插入角色菜单关联失败: %v", err)
			// 回滚：删除刚插入的角色
			_ = l.svcCtx.SysRoleModel.Delete(l.ctx, newRoleId)
			return &types.BaseResp{
				Code: 500,
				Msg:  "插入角色菜单关联失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
