// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package monitor

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取操作日志记录列表
func NewOperLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogListLogic {
	return &OperLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperLogListLogic) OperLogList(req *types.OperLogListReq) (resp *types.OperLogListResp, err error) {
	// 构建查询条件
	query := &model.OperLogQuery{
		OperIp:       req.OperIp,
		Title:        req.Title,
		BusinessType: req.BusinessType,
		Status:       req.Status,
		OperName:     req.OperName,
		BeginTime:    req.BeginTime,
		EndTime:      req.EndTime,
	}
	pageQuery := &model.PageQuery{
		PageNum:       req.PageNum,
		PageSize:      req.PageSize,
		OrderByColumn: req.OrderByColumn,
		IsAsc:         req.IsAsc,
	}

	// 使用 SQL 分页查询
	rows, total, err := l.svcCtx.SysOperLogModel.FindPage(l.ctx, query, pageQuery)
	if err != nil {
		l.Errorf("查询操作日志记录列表失败: %v", err)
		return &types.OperLogListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询操作日志记录列表失败",
			},
		}, err
	}

	// 转换为 VO
	voList := make([]types.OperLogVo, 0, len(rows))
	for _, row := range rows {
		vo := types.OperLogVo{
			OperId:        row.OperId,
			TenantId:      row.TenantId,
			Title:         row.Title,
			BusinessType:  int32(row.BusinessType),
			Method:        row.Method,
			RequestMethod: row.RequestMethod,
			OperatorType:  int32(row.OperatorType),
			OperName:      row.OperName,
			DeptName:      row.DeptName,
			OperUrl:       row.OperUrl,
			OperIp:        row.OperIp,
			OperLocation:  row.OperLocation,
			OperParam:     row.OperParam,
			JsonResult:    row.JsonResult,
			Status:        int32(row.Status),
			ErrorMsg:      row.ErrorMsg,
			CostTime:      row.CostTime,
		}
		if row.OperTime.Valid {
			vo.OperTime = row.OperTime.Time.Format("2006-01-02 15:04:05")
		}
		voList = append(voList, vo)
	}

	return &types.OperLogListResp{
		Total: total,
		Rows:  voList,
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "查询成功",
		},
	}, nil
}
