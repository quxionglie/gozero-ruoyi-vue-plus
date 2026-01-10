// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package monitor

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogininforCleanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogininforCleanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogininforCleanLogic {
	return &LogininforCleanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogininforCleanLogic) LogininforClean() (resp *types.BaseResp, err error) {
	// 清空所有登录日志
	err = l.svcCtx.SysLogininforModel.Clean(l.ctx)
	if err != nil {
		l.Errorf("清理系统访问记录失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "清理系统访问记录失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "清理成功",
	}, nil
}
