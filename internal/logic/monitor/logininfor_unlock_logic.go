// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package monitor

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogininforUnlockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogininforUnlockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogininforUnlockLogic {
	return &LogininforUnlockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogininforUnlockLogic) LogininforUnlock(req *types.LogininforUnlockReq) (resp *types.BaseResp, err error) {
	// 删除密码错误计数缓存
	// key格式: pwd_err_cnt:{userName}
	key := "pwd_err_cnt:" + req.UserName

	// 检查key是否存在
	exists, err := l.svcCtx.RedisConn.ExistsCtx(l.ctx, key)
	if err != nil {
		l.Errorf("检查密码错误计数缓存失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "解锁用户失败",
		}, err
	}

	if exists {
		// 删除key
		_, delErr := l.svcCtx.RedisConn.DelCtx(l.ctx, key)
		if delErr != nil {
			l.Errorf("删除密码错误计数缓存失败: %v", delErr)
			return &types.BaseResp{
				Code: 500,
				Msg:  "解锁用户失败",
			}, delErr
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "解锁成功",
	}, nil
}
