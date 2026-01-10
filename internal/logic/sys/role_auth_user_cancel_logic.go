// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAuthUserCancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消授权用户
func NewRoleAuthUserCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAuthUserCancelLogic {
	return &RoleAuthUserCancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleAuthUserCancelLogic) RoleAuthUserCancel(req *types.RoleAuthUserCancelReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.RoleId == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "角色ID不能为空",
		}, nil
	}
	if req.UserId == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "用户ID不能为空",
		}, nil
	}

	// 2. 获取当前用户ID
	currentUserId, _ := util.GetUserIdFromContext(l.ctx)

	// 3. 不允许修改当前用户角色
	if currentUserId == req.UserId {
		return &types.BaseResp{
			Code: 500,
			Msg:  "不允许修改当前用户角色!",
		}, nil
	}

	// 4. 删除用户角色关联
	err = l.svcCtx.SysUserRoleModel.Delete(l.ctx, req.UserId, req.RoleId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BaseResp{
				Code: 404,
				Msg:  "用户角色关联不存在",
			}, nil
		}
		l.Errorf("删除用户角色关联失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "删除用户角色关联失败",
		}, err
	}

	// TODO: 清除在线用户缓存（如果用户在线）

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
