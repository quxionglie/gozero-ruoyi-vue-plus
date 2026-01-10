// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictTypeRefreshCacheLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 刷新字典缓存
func NewDictTypeRefreshCacheLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeRefreshCacheLogic {
	return &DictTypeRefreshCacheLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictTypeRefreshCacheLogic) DictTypeRefreshCache() (resp *types.BaseResp, err error) {
	// 刷新缓存：清除 Redis 中的字典缓存
	keys, err := l.svcCtx.RedisConn.KeysCtx(l.ctx, "sys_dict:*")
	if err == nil {
		for _, key := range keys {
			_, _ = l.svcCtx.RedisConn.DelCtx(l.ctx, key)
		}
	}

	keys, err = l.svcCtx.RedisConn.KeysCtx(l.ctx, "sys_dict_type:*")
	if err == nil {
		for _, key := range keys {
			_, _ = l.svcCtx.RedisConn.DelCtx(l.ctx, key)
		}
	}

	l.Infof("已刷新字典缓存")
	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
