// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

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

func (l *ConfigListLogic) ConfigList() (resp *types.TableDataInfoResp, err error) {
	rows, err := l.svcCtx.SysConfigModel.FindAll(l.ctx)
	if err != nil {
		l.Errorf("查询参数配置列表失败: %v", err)
		return &types.TableDataInfoResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询参数配置列表失败",
			},
			Data: types.TableDataInfo{
				Total: 0,
				Rows:  []types.ConfigVo{},
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
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: types.TableDataInfo{
			Total: int64(len(voList)),
			Rows:  voList,
		},
	}, nil
}
