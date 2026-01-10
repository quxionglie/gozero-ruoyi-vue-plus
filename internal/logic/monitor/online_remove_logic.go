// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package monitor

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnlineRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 强退用户
func NewOnlineRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnlineRemoveLogic {
	return &OnlineRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnlineRemoveLogic) OnlineRemove(req *types.OnlineRemoveReq) (resp *types.BaseResp, err error) {
	// 删除在线 token 缓存
	onlineTokenKey := "online_tokens:" + req.TokenId
	_, delErr := l.svcCtx.RedisConn.DelCtx(l.ctx, onlineTokenKey)
	if delErr != nil {
		l.Errorf("强退用户失败: %v", delErr)
		return &types.BaseResp{
			Code: 500,
			Msg:  "强退用户失败",
		}, delErr
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "强退成功",
	}, nil
}
