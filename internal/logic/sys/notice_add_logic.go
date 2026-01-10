// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"database/sql"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增通知公告
func NewNoticeAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeAddLogic {
	return &NoticeAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeAddLogic) NoticeAdd(req *types.NoticeReq) (resp *types.BaseResp, err error) {
	// 1. 参数长度校验
	if err := util.ValidateStringLength(req.NoticeTitle, "公告标题", 50); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}

	// 2. 获取当前用户信息
	userId, _ := util.GetUserIdFromContext(l.ctx)
	tenantId, _ := util.GetTenantIdFromContext(l.ctx)

	// 获取用户的部门ID
	var deptId int64
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err == nil && user.DeptId.Valid {
		deptId = user.DeptId.Int64
	}

	// 3. 生成主键ID（使用雪花算法）
	noticeId, err := util.GenerateID()
	if err != nil {
		l.Errorf("生成公告ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "生成公告ID失败",
		}, err
	}

	// 4. 构建通知公告实体
	notice := &model.SysNotice{
		NoticeId:      noticeId,
		TenantId:      tenantId,
		NoticeTitle:   req.NoticeTitle,
		NoticeType:    req.NoticeType,
		NoticeContent: sql.NullString{String: req.NoticeContent, Valid: req.NoticeContent != ""},
		Status:        req.Status,
		Remark:        sql.NullString{String: req.Remark, Valid: req.Remark != ""},
		CreateDept:    sql.NullInt64{Int64: deptId, Valid: deptId > 0},
		CreateBy:      sql.NullInt64{Int64: userId, Valid: userId > 0},
	}

	// 5. 插入数据库
	_, err = l.svcCtx.SysNoticeModel.Insert(l.ctx, notice)
	if err != nil {
		l.Errorf("新增通知公告失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "新增通知公告失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
