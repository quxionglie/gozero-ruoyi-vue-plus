// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询参数配置列表
func NewConfigListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigListLogic {
	return &ConfigListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigListLogic) ConfigList(req *types.ConfigListReq) (resp *types.TableDataInfoResp, err error) {
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
	configQuery := &model.ConfigQuery{
		ConfigName: req.ConfigName,
		ConfigKey:  req.ConfigKey,
		ConfigType: req.ConfigType,
	}
	pageQuery := &model.PageQuery{
		PageNum:       pageNum,
		PageSize:      pageSize,
		OrderByColumn: req.OrderByColumn,
		IsAsc:         req.IsAsc,
	}

	// 使用 SQL 分页查询
	rows, total, err := l.svcCtx.SysConfigModel.FindPage(l.ctx, configQuery, pageQuery)
	if err != nil {
		l.Errorf("查询参数配置列表失败: %v", err)
		return &types.TableDataInfoResp{
			Total: 0,
			Rows:  []types.ConfigVo{},
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询参数配置列表失败",
			},
		}, err
	}

	// 转换为 VO
	voList := make([]types.ConfigVo, 0, len(rows))
	for _, row := range rows {
		vo := types.ConfigVo{
			ConfigId:    row.ConfigId,
			ConfigName:  row.ConfigName,
			ConfigKey:   row.ConfigKey,
			ConfigValue: row.ConfigValue,
			ConfigType:  row.ConfigType,
			Remark:      "",
			CreateTime:  "",
		}
		if row.Remark.Valid {
			vo.Remark = row.Remark.String
		}
		if row.CreateTime.Valid {
			vo.CreateTime = row.CreateTime.Time.Format("2006-01-02 15:04:05")
		}
		voList = append(voList, vo)
	}

	return &types.TableDataInfoResp{
		Total: total,
		Rows:  voList,
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "查询成功",
		},
	}, nil
}
