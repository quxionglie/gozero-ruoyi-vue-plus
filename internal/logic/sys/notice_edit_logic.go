// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"database/sql"
	"time"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改通知公告
func NewNoticeEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeEditLogic {
	return &NoticeEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeEditLogic) NoticeEdit(req *types.NoticeReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.NoticeId <= 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "公告ID不能为空",
		}, nil
	}

	// 2. 参数长度校验
	if err := util.ValidateStringLength(req.NoticeTitle, "公告标题", 50); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}

	// 3. 检查公告是否存在
	notice, err := l.svcCtx.SysNoticeModel.FindOne(l.ctx, req.NoticeId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BaseResp{
				Code: 500,
				Msg:  "通知公告不存在",
			}, nil
		}
		l.Errorf("查询通知公告失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询通知公告失败",
		}, err
	}

	// 4. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 5. 更新通知公告信息
	notice.NoticeTitle = req.NoticeTitle
	if req.NoticeType != "" {
		notice.NoticeType = req.NoticeType
	}
	notice.NoticeContent = sql.NullString{String: req.NoticeContent, Valid: req.NoticeContent != ""}
	if req.Status != "" {
		notice.Status = req.Status
	}
	notice.Remark = sql.NullString{String: req.Remark, Valid: req.Remark != ""}
	notice.UpdateBy = sql.NullInt64{Int64: userId, Valid: userId > 0}
	notice.UpdateTime = sql.NullTime{Time: time.Now(), Valid: true}

	// 6. 更新数据库
	err = l.svcCtx.SysNoticeModel.Update(l.ctx, notice)
	if err != nil {
		l.Errorf("修改通知公告失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改通知公告失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
