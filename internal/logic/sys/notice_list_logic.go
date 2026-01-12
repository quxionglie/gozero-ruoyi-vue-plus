// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询通知公告列表
func NewNoticeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeListLogic {
	return &NoticeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeListLogic) NoticeList(req *types.NoticeListReq) (resp *types.TableDataInfoResp, err error) {
	// 1. 构建查询条件
	query := &sys.NoticeQuery{
		NoticeTitle: req.NoticeTitle,
		NoticeType:  req.NoticeType,
	}

	// 2. 如果提供了创建人名称，需要通过用户名查询用户ID
	if req.CreateByName != "" {
		// 获取当前租户ID（从 context 或请求中）
		tenantId, _ := util.GetTenantIdFromContext(l.ctx)
		user, err := l.svcCtx.SysUserModel.FindOneByUserName(l.ctx, req.CreateByName, tenantId)
		if err == nil && user != nil {
			query.CreateBy = user.UserId
		}
	}

	// 3. 构建分页查询参数
	pageQuery := &sys.PageQuery{
		PageNum:       req.PageNum,
		PageSize:      req.PageSize,
		OrderByColumn: req.OrderByColumn,
		IsAsc:         req.IsAsc,
	}

	// 4. 查询数据
	notices, total, err := l.svcCtx.SysNoticeModel.FindPage(l.ctx, query, pageQuery)
	if err != nil {
		l.Errorf("查询通知公告列表失败: %v", err)
		return &types.TableDataInfoResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询通知公告列表失败",
			},
		}, err
	}

	// 5. 转换为响应格式
	rows := make([]types.NoticeVo, 0, len(notices))
	for _, notice := range notices {
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
		rows = append(rows, noticeVo)
	}

	return &types.TableDataInfoResp{
		Total: total,
		Rows:  rows,
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "查询成功",
		},
	}, nil
}
