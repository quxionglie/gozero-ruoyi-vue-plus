// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package monitor

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogCleanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 清理操作日志记录
func NewOperLogCleanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogCleanLogic {
	return &OperLogCleanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperLogCleanLogic) OperLogClean() (resp *types.BaseResp, err error) {
	// 清空所有操作日志
	err = l.svcCtx.SysOperLogModel.Clean(l.ctx)
	if err != nil {
		l.Errorf("清理操作日志失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "清理操作日志失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "清理成功",
	}, nil
}
