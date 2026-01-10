// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClientRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除客户端管理
func NewClientRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClientRemoveLogic {
	return &ClientRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClientRemoveLogic) ClientRemove(req *types.ClientRemoveReq) (resp *types.BaseResp, err error) {
	if req.Ids == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "客户端ID不能为空",
		}, nil
	}

	// 解析客户端ID列表（逗号分隔）
	idStrs := strings.Split(req.Ids, ",")
	var ids []int64
	for _, idStr := range idStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return &types.BaseResp{
				Code: 400,
				Msg:  fmt.Sprintf("客户端ID格式错误: %s", idStr),
			}, nil
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "客户端ID不能为空",
		}, nil
	}

	// 删除客户端管理
	for _, id := range ids {
		err = l.svcCtx.SysClientModel.Delete(l.ctx, id)
		if err != nil {
			l.Errorf("删除客户端管理失败: id=%d, err=%v", id, err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除客户端管理失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
