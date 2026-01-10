// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package monitor

import (
	"context"
	"strconv"
	"strings"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除操作日志记录
func NewOperLogRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogRemoveLogic {
	return &OperLogRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperLogRemoveLogic) OperLogRemove(req *types.OperLogRemoveReq) (resp *types.BaseResp, err error) {
	// 解析 operIds（逗号分隔的字符串）
	operIdStrs := strings.Split(req.OperIds, ",")
	operIds := make([]int64, 0, len(operIdStrs))
	for _, idStr := range operIdStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr != "" {
			id, parseErr := strconv.ParseInt(idStr, 10, 64)
			if parseErr == nil {
				operIds = append(operIds, id)
			}
		}
	}

	if len(operIds) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "参数错误：未提供有效的日志ID",
		}, nil
	}

	// 批量删除
	err = l.svcCtx.SysOperLogModel.DeleteByIds(l.ctx, operIds)
	if err != nil {
		l.Errorf("批量删除操作日志失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "批量删除操作日志失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}
