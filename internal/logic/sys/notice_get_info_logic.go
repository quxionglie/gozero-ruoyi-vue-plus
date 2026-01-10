// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询通知公告详细
func NewNoticeGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeGetInfoLogic {
	return &NoticeGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeGetInfoLogic) NoticeGetInfo(req *types.NoticeGetInfoReq) (resp *types.NoticeResp, err error) {
	// 1. 查询通知公告信息
	notice, err := l.svcCtx.SysNoticeModel.FindOne(l.ctx, req.NoticeId)
	if err != nil {
		if err == sys.ErrNotFound {
			return &types.NoticeResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "通知公告不存在",
				},
			}, nil
		}
		l.Errorf("查询通知公告信息失败: %v", err)
		return &types.NoticeResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询通知公告信息失败",
			},
		}, err
	}

	// 2. 转换为响应格式
	noticeVo := types.NoticeVo{
		NoticeId:      notice.NoticeId,
		NoticeTitle:   notice.NoticeTitle,
		NoticeType:    notice.NoticeType,
		NoticeContent: "",
		Status:        notice.Status,
		Remark:        "",
		CreateTime:    "",
	}
	if notice.NoticeContent.Valid {
		noticeVo.NoticeContent = notice.NoticeContent.String
	}
	if notice.Remark.Valid {
		noticeVo.Remark = notice.Remark.String
	}
	if notice.CreateTime.Valid {
		noticeVo.CreateTime = notice.CreateTime.Time.Format("2006-01-02 15:04:05")
	}

	return &types.NoticeResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: noticeVo,
	}, nil
}
