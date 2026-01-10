package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UserResetPwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置密码
func NewUserResetPwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserResetPwdLogic {
	return &UserResetPwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserResetPwdLogic) UserResetPwd(req *types.UserResetPwdReq) (resp *types.BaseResp, err error) {
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

	// 3. 验证新密码
	if req.Password == "" {
		return &types.BaseResp{
			Code: 500,
			Msg:  "新密码不能为空",
		}, nil
	}

	// 4. 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("密码加密失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "密码加密失败",
		}, err
	}

	// 5. 重置用户密码
	err = l.svcCtx.SysUserModel.ResetUserPwd(l.ctx, req.UserId, string(hashedPassword))
	if err != nil {
		l.Errorf("重置用户密码失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "重置用户密码失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "重置成功",
	}, nil
}
