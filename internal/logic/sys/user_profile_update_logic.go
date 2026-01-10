package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户信息
func NewUserProfileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileUpdateLogic {
	return &UserProfileUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserProfileUpdateLogic) UserProfileUpdate(req *types.UserProfileReq) (resp *types.BaseResp, err error) {
	// 1. 获取当前用户ID
	userId, err := util.GetUserIdFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "获取用户ID失败",
		}, err
	}

	// 2. 校验手机号是否唯一（排除当前用户）
	if req.Phonenumber != "" {
		unique, err := l.svcCtx.SysUserModel.CheckPhoneUnique(l.ctx, req.Phonenumber, userId)
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
				Msg:  "修改个人信息失败，手机号码已存在",
			}, nil
		}
	}

	// 3. 校验邮箱是否唯一（排除当前用户）
	if req.Email != "" {
		unique, err := l.svcCtx.SysUserModel.CheckEmailUnique(l.ctx, req.Email, userId)
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
				Msg:  "修改个人信息失败，邮箱账号已存在",
			}, nil
		}
	}

	// 4. 更新用户基本信息
	err = l.svcCtx.SysUserModel.UpdateUserProfile(l.ctx, userId, req.NickName, req.Email, req.Phonenumber, req.Sex)
	if err != nil {
		l.Errorf("修改个人信息失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改个人信息失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "修改成功",
	}, nil
}
