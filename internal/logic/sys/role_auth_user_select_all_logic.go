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

type RoleAuthUserSelectAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量选择用户授权
func NewRoleAuthUserSelectAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAuthUserSelectAllLogic {
	return &RoleAuthUserSelectAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleAuthUserSelectAllLogic) RoleAuthUserSelectAll(req *types.RoleAuthUserSelectAllReq) (resp *types.BaseResp, err error) {
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

	// 2. 获取当前用户ID（用于数据权限检查，这里简化处理）
	currentUserId, _ := util.GetUserIdFromContext(l.ctx)
	_ = currentUserId // 暂时不使用，后续可以根据需要实现数据权限检查

	// TODO: 检查数据权限（checkRoleDataScope）

	// 3. 批量插入用户角色关联
	err = l.svcCtx.SysUserRoleModel.InsertBatch(l.ctx, req.RoleId, req.UserIds)
	if err != nil {
		l.Errorf("批量插入用户角色关联失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "批量插入用户角色关联失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
