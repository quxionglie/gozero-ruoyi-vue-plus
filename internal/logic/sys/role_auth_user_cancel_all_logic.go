// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAuthUserCancelAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量取消授权用户
func NewRoleAuthUserCancelAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAuthUserCancelAllLogic {
	return &RoleAuthUserCancelAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleAuthUserCancelAllLogic) RoleAuthUserCancelAll(req *types.RoleAuthUserCancelAllReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.RoleId == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "角色ID不能为空",
		}, nil
	}
	if len(req.UserIds) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "用户ID列表不能为空",
		}, nil
	}

	// 2. 获取当前用户ID
	currentUserId, _ := util.GetUserIdFromContext(l.ctx)

	// 3. 检查是否包含当前用户
	for _, userId := range req.UserIds {
		if currentUserId == userId {
			return &types.BaseResp{
				Code: 500,
				Msg:  "不允许修改当前用户角色!",
			}, nil
		}
	}

	// 4. 批量删除用户角色关联
	err = l.svcCtx.SysUserRoleModel.DeleteByRoleIdAndUserIds(l.ctx, req.RoleId, req.UserIds)
	if err != nil {
		l.Errorf("批量删除用户角色关联失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "批量删除用户角色关联失败",
		}, err
	}

	// TODO: 清除在线用户缓存（如果用户在线）

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
