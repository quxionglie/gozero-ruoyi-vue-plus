package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 头像上传
func NewUserProfileAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileAvatarLogic {
	return &UserProfileAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserProfileAvatarLogic) UserProfileAvatar() (resp *types.BaseResp, err error) {
	// TODO: 实际头像上传功能需要处理文件上传
	// 这里先返回成功响应，后续需要：
	// 1. 接收multipart文件
	// 2. 上传到OSS
	// 3. 获取OSS文件ID
	// 4. 更新用户头像字段
	return &types.BaseResp{
		Code: 200,
		Msg:  "头像上传功能待实现",
	}, nil
}
