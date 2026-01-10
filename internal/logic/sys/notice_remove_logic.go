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

type NoticeRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除通知公告
func NewNoticeRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeRemoveLogic {
	return &NoticeRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeRemoveLogic) NoticeRemove(req *types.NoticeRemoveReq) (resp *types.BaseResp, err error) {
	if req.NoticeIds == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "公告ID不能为空",
		}, nil
	}

	// 解析公告ID列表（逗号分隔）
	noticeIdStrs := strings.Split(req.NoticeIds, ",")
	var noticeIds []int64
	for _, idStr := range noticeIdStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return &types.BaseResp{
				Code: 400,
				Msg:  fmt.Sprintf("公告ID格式错误: %s", idStr),
			}, nil
		}
		noticeIds = append(noticeIds, id)
	}

	if len(noticeIds) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "公告ID不能为空",
		}, nil
	}

	// 删除通知公告
	for _, noticeId := range noticeIds {
		err = l.svcCtx.SysNoticeModel.Delete(l.ctx, noticeId)
		if err != nil {
			l.Errorf("删除通知公告失败: noticeId=%d, err=%v", noticeId, err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除通知公告失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
