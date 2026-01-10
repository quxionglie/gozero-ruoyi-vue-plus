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

type LogininforRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogininforRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogininforRemoveLogic {
	return &LogininforRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogininforRemoveLogic) LogininforRemove(req *types.LogininforRemoveReq) (resp *types.BaseResp, err error) {
	// 解析 infoIds（逗号分隔的字符串）
	infoIdStrs := strings.Split(req.InfoIds, ",")
	infoIds := make([]int64, 0, len(infoIdStrs))
	for _, idStr := range infoIdStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr != "" {
			id, parseErr := strconv.ParseInt(idStr, 10, 64)
			if parseErr == nil {
				infoIds = append(infoIds, id)
			}
		}
	}

	if len(infoIds) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "参数错误：未提供有效的日志ID",
		}, nil
	}

	// 批量删除
	err = l.svcCtx.SysLogininforModel.DeleteByIds(l.ctx, infoIds)
	if err != nil {
		l.Errorf("批量删除登录日志失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "批量删除登录日志失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}
