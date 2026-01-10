package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserChangeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 状态修改
func NewUserChangeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserChangeStatusLogic {
	return &UserChangeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserChangeStatusLogic) UserChangeStatus(req *types.UserChangeStatusReq) (resp *types.BaseResp, err error) {
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

	// 3. 验证状态值
	if req.Status != "0" && req.Status != "1" {
		return &types.BaseResp{
			Code: 500,
			Msg:  "状态值无效",
		}, nil
	}

	// 4. 更新用户状态
	err = l.svcCtx.SysUserModel.UpdateUserStatus(l.ctx, req.UserId, req.Status)
	if err != nil {
		l.Errorf("修改用户状态失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改用户状态失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
