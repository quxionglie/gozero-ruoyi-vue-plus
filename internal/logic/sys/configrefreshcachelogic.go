// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigRefreshCacheLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 刷新参数缓存
func NewConfigRefreshCacheLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigRefreshCacheLogic {
	return &ConfigRefreshCacheLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigRefreshCacheLogic) ConfigRefreshCache() (resp *types.BaseResp, err error) {
	// 刷新缓存：清除 Redis 中的配置缓存
	keys, err := l.svcCtx.RedisConn.KeysCtx(l.ctx, "sys_config:*")
	if err == nil {
		for _, key := range keys {
			_, _ = l.svcCtx.RedisConn.DelCtx(l.ctx, key)
		}
	}

	l.Infof("已刷新参数配置缓存")
	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
