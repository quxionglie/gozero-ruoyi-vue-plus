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

type LogininforListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取系统访问记录列表
func NewLogininforListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogininforListLogic {
	return &LogininforListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogininforListLogic) LogininforList(req *types.LogininforListReq) (resp *types.LogininforListResp, err error) {
	// 设置默认分页参数
	pageNum := req.PageNum
	pageSize := req.PageSize
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 构建查询条件
	query := &model.LogininforQuery{
		Ipaddr:    req.Ipaddr,
		Status:    req.Status,
		UserName:  req.UserName,
		BeginTime: req.BeginTime,
		EndTime:   req.EndTime,
	}
	pageQuery := &model.PageQuery{
		PageNum:       pageNum,
		PageSize:      pageSize,
		OrderByColumn: req.OrderByColumn,
		IsAsc:         req.IsAsc,
	}

	// 使用 SQL 分页查询
	rows, total, err := l.svcCtx.SysLogininforModel.FindPage(l.ctx, query, pageQuery)
	if err != nil {
		l.Errorf("查询系统访问记录列表失败: %v", err)
		return &types.LogininforListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询系统访问记录列表失败",
			},
		}, err
	}

	// 转换为 VO
	voList := make([]types.LogininforVo, 0, len(rows))
	for _, row := range rows {
		vo := types.LogininforVo{
			InfoId:        row.InfoId,
			TenantId:      row.TenantId,
			UserName:      row.UserName,
			ClientKey:     row.ClientKey,
			DeviceType:    row.DeviceType,
			Ipaddr:        row.Ipaddr,
			LoginLocation: row.LoginLocation,
			Browser:       row.Browser,
			Os:            row.Os,
			Status:        row.Status,
			Msg:           row.Msg,
		}
		if row.LoginTime.Valid {
			vo.LoginTime = row.LoginTime.Time.Format("2006-01-02 15:04:05")
		}
		voList = append(voList, vo)
	}

	return &types.LogininforListResp{
		Total: total,
		Rows:  voList,
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "查询成功",
		},
	}, nil
}
