package sys

import (
	"context"
	"fmt"
	"strings"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除用户
func NewUserRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRemoveLogic {
	return &UserRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRemoveLogic) UserRemove(req *types.UserRemoveReq) (resp *types.BaseResp, err error) {
	// 1. 解析用户ID串
	userIdsStr := strings.Split(req.UserIds, ",")
	userIds := make([]int64, 0, len(userIdsStr))
	for _, idStr := range userIdsStr {
		if idStr == "" {
			continue
		}
		var id int64
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil {
			continue
		}
		userIds = append(userIds, id)
	}
	if len(userIds) == 0 {
		return &types.BaseResp{
			Code: 500,
			Msg:  "用户ID不能为空",
		}, nil
	}

	// 2. 获取当前用户ID
	currentUserId, err := util.GetUserIdFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "获取用户ID失败",
		}, err
	}

	// 3. 检查每个用户是否允许删除
	for _, userId := range userIds {
		// 不能删除自己
		if userId == currentUserId {
			return &types.BaseResp{
				Code: 500,
				Msg:  "当前用户不能删除",
			}, nil
		}

		// 不能删除超级管理员
		isSuperAdmin, err := l.svcCtx.SysRoleModel.CheckIsSuperAdmin(l.ctx, userId)
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
	}

	// 4. 批量删除用户
	for _, userId := range userIds {
		// 删除用户与角色关联
		err = l.svcCtx.SysUserRoleModel.DeleteByUserId(l.ctx, userId)
		if err != nil {
			l.Errorf("删除用户角色关联失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除用户角色关联失败",
			}, err
		}

		// 删除用户与岗位关联
		err = l.svcCtx.SysUserPostModel.DeleteByUserId(l.ctx, userId)
		if err != nil {
			l.Errorf("删除用户岗位关联失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除用户岗位关联失败",
			}, err
		}

		// 删除用户（软删除，设置del_flag=1）
		user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
		if err != nil {
			if err == model.ErrNotFound {
				continue // 用户不存在，跳过
			}
			l.Errorf("查询用户信息失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "查询用户信息失败",
			}, err
		}

		user.DelFlag = "1"
		err = l.svcCtx.SysUserModel.Update(l.ctx, user)
		if err != nil {
			l.Errorf("删除用户失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除用户失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}
