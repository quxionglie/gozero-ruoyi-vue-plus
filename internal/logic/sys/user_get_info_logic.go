// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据用户编号获取详细信息
func NewUserGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetInfoLogic {
	return &UserGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserGetInfoLogic) UserGetInfo(req *types.UserGetInfoReq) (resp *types.UserDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
