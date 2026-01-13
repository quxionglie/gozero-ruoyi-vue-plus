package sys

import (
	"context"
	"database/sql"
	"time"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户
func NewUserEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserEditLogic {
	return &UserEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserEditLogic) UserEdit(req *types.UserReq) (resp *types.BaseResp, err error) {
	// 1. 校验用户ID
	if req.UserId <= 0 {
		return &types.BaseResp{
			Code: 500,
			Msg:  "用户ID不能为空",
		}, nil
	}

	// 2. 检查用户是否允许操作（不能操作超级管理员）
	isSuperAdmin, err := l.svcCtx.SysRoleModel.CheckIsSuperAdmin(l.ctx, req.UserId)
	if err != nil {
		l.Errorf("检查超级管理员失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "检查超级管理员失败",
		}, err
	}
	if isSuperAdmin {
		return &types.BaseResp{
			Code: 500,
			Msg:  "不允许操作超级管理员用户",
		}, nil
	}

	// 3. 校验用户名是否唯一
	unique, err := l.svcCtx.SysUserModel.CheckUserNameUnique(l.ctx, req.UserName, req.UserId)
	if err != nil {
		l.Errorf("校验用户名唯一性失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "校验用户名唯一性失败",
		}, err
	}
	if !unique {
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改用户'" + req.UserName + "'失败，登录账号已存在",
		}, nil
	}

	// 4. 校验手机号是否唯一
	if req.Phonenumber != "" {
		unique, err := l.svcCtx.SysUserModel.CheckPhoneUnique(l.ctx, req.Phonenumber, req.UserId)
		if err != nil {
			l.Errorf("校验手机号唯一性失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "校验手机号唯一性失败",
			}, err
		}
		if !unique {
			return &types.BaseResp{
				Code: 500,
				Msg:  "修改用户'" + req.UserName + "'失败，手机号码已存在",
			}, nil
		}
	}

	// 5. 校验邮箱是否唯一
	if req.Email != "" {
		unique, err := l.svcCtx.SysUserModel.CheckEmailUnique(l.ctx, req.Email, req.UserId)
		if err != nil {
			l.Errorf("校验邮箱唯一性失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "校验邮箱唯一性失败",
			}, err
		}
		if !unique {
			return &types.BaseResp{
				Code: 500,
				Msg:  "修改用户'" + req.UserName + "'失败，邮箱账号已存在",
			}, nil
		}
	}

	// 6. 获取当前用户ID
	userId, err := util.GetUserIdFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "获取用户ID失败",
		}, err
	}

	// 7. 更新用户角色关联
	if len(req.RoleIds) > 0 {
		// 先删除原有关联
		err = l.svcCtx.SysUserRoleModel.DeleteByUserId(l.ctx, req.UserId)
		if err != nil {
			l.Errorf("删除用户角色关联失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除用户角色关联失败",
			}, err
		}

		// 非超级管理员，禁止包含超级管理员角色
		currentUserId, _ := util.GetUserIdFromContext(l.ctx)
		isCurrentSuperAdmin, err := l.svcCtx.SysRoleModel.CheckIsSuperAdmin(l.ctx, currentUserId)
		if err != nil {
			l.Errorf("检查超级管理员失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "检查超级管理员失败",
			}, err
		}

		var roleIds []int64
		for _, roleId := range req.RoleIds {
			// 非超级管理员，禁止包含超级管理员角色（roleId=1）
			if !isCurrentSuperAdmin && roleId == 1 {
				continue
			}
			roleIds = append(roleIds, roleId)
		}

		if len(roleIds) > 0 {
			err = l.svcCtx.SysUserRoleModel.InsertBatchByUserId(l.ctx, req.UserId, roleIds)
			if err != nil {
				l.Errorf("新增用户角色关联失败: %v", err)
				return &types.BaseResp{
					Code: 500,
					Msg:  "新增用户角色关联失败",
				}, err
			}
		}
	}

	// 9. 更新用户岗位关联
	if len(req.PostIds) > 0 {
		// 先删除原有关联
		err = l.svcCtx.SysUserPostModel.DeleteByUserId(l.ctx, req.UserId)
		if err != nil {
			l.Errorf("删除用户岗位关联失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除用户岗位关联失败",
			}, err
		}

		// 批量插入新的关联
		err = l.svcCtx.SysUserPostModel.InsertBatch(l.ctx, req.UserId, req.PostIds)
		if err != nil {
			l.Errorf("新增用户岗位关联失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "新增用户岗位关联失败",
			}, err
		}
	}

	// 10. 构建更新用户实体（只设置表单输入的字段）
	updateUser := &model.SysUser{
		UserId:     req.UserId,
		UserName:   req.UserName,
		NickName:   req.NickName,
		UpdateBy:   sql.NullInt64{Int64: userId, Valid: true},
		UpdateTime: sql.NullTime{Time: time.Now(), Valid: true},
	}
	if req.Email != "" {
		updateUser.Email = req.Email
	}
	if req.Phonenumber != "" {
		updateUser.Phonenumber = req.Phonenumber
	}
	if req.Sex != "" {
		updateUser.Sex = req.Sex
	}
	if req.Status != "" {
		updateUser.Status = req.Status
	}
	if req.UserType != "" {
		updateUser.UserType = req.UserType
	}
	if req.DeptId > 0 {
		updateUser.DeptId = sql.NullInt64{Int64: req.DeptId, Valid: true}
	}
	if req.Remark != "" {
		updateUser.Remark = sql.NullString{String: req.Remark, Valid: true}
	}
	// 密码处理：如果传了密码则更新（需要加密），否则不更新
	if req.Password != "" {
		// TODO: 这里应该对密码进行加密
		updateUser.Password = req.Password
	}

	// 11. 更新用户信息
	_, err = l.svcCtx.SysUserModel.UpdateById(l.ctx, updateUser)
	if err != nil {
		l.Errorf("修改用户信息失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改用户'" + req.UserName + "'信息失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "修改成功",
	}, nil
}
