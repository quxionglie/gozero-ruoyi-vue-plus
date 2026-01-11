package sys

import (
	"context"
	"strconv"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 个人信息
func NewUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileLogic {
	return &UserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserProfileLogic) UserProfile() (resp *types.UserProfileResp, err error) {
	// 1. 获取当前用户ID
	userId, err := util.GetUserIdFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return &types.UserProfileResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "获取用户ID失败",
			},
		}, err
	}

	// 2. 查询用户信息
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		l.Errorf("查询用户信息失败: %v", err)
		return &types.UserProfileResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询用户信息失败",
			},
		}, err
	}

	// 3. 转换为响应格式（ProfileUserVo 只包含基础字段）
	profileVo := types.ProfileUserVo{
		UserId:      user.UserId,
		TenantId:    user.TenantId,
		DeptId:      0,
		UserName:    user.UserName,
		NickName:    user.NickName,
		UserType:    user.UserType,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Sex:         user.Sex,
		Avatar:      "",
		LoginIp:     user.LoginIp,
		LoginDate:   "",
		DeptName:    "",
	}
	if user.DeptId.Valid {
		profileVo.DeptId = user.DeptId.Int64
	}
	if user.Avatar.Valid {
		profileVo.Avatar = strconv.FormatInt(user.Avatar.Int64, 10)
	} else {
		profileVo.Avatar = ""
	}
	if user.LoginDate.Valid {
		profileVo.LoginDate = user.LoginDate.Time.Format("2006-01-02 15:04:05")
	}

	return &types.UserProfileResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: profileVo,
	}, nil
}
