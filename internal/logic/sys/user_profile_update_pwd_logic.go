package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UserProfileUpdatePwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置密码
func NewUserProfileUpdatePwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileUpdatePwdLogic {
	return &UserProfileUpdatePwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserProfileUpdatePwdLogic) UserProfileUpdatePwd(req *types.UserPasswordReq) (resp *types.BaseResp, err error) {
	// 1. 获取当前用户ID
	userId, err := util.GetUserIdFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "获取用户ID失败",
		}, err
	}

	// 2. 查询用户信息
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		l.Errorf("查询用户信息失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询用户信息失败",
		}, err
	}

	// 3. 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改密码失败，旧密码错误",
		}, nil
	}

	// 4. 验证新密码不能与旧密码相同
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.NewPassword))
	if err == nil {
		return &types.BaseResp{
			Code: 500,
			Msg:  "新密码不能与旧密码相同",
		}, nil
	}

	// 5. 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("密码加密失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "密码加密失败",
		}, err
	}

	// 6. 更新密码
	err = l.svcCtx.SysUserModel.ResetUserPwd(l.ctx, userId, string(hashedPassword))
	if err != nil {
		l.Errorf("修改密码失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改密码失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "修改成功",
	}, nil
}
