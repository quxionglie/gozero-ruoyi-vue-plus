// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClientChangeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 状态修改
func NewClientChangeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClientChangeStatusLogic {
	return &ClientChangeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClientChangeStatusLogic) ClientChangeStatus(req *types.ClientReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.ClientId == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "客户端ID不能为空",
		}, nil
	}
	if req.Status == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "状态不能为空",
		}, nil
	}

	// 2. 更新客户端状态（根据 clientId 更新）
	err = l.svcCtx.SysClientModel.UpdateClientStatus(l.ctx, req.ClientId, req.Status)
	if err != nil {
		l.Errorf("修改客户端状态失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改客户端状态失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
