// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictDataListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询字典数据列表
func NewDictDataListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDataListLogic {
	return &DictDataListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictDataListLogic) DictDataList() (resp *types.TableDataInfoResp, err error) {
	rows, err := l.svcCtx.SysDictDataModel.FindAll(l.ctx)
	if err != nil {
		l.Errorf("查询字典数据列表失败: %v", err)
		return &types.TableDataInfoResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询字典数据列表失败",
			},
			Data: types.TableDataInfo{
				Total: 0,
				Rows:  []types.DictDataVo{},
			},
		}, err
	}

	// 转换为 VO
	voList := make([]types.DictDataVo, 0, len(rows))
	for _, row := range rows {
		vo := types.DictDataVo{
			DictCode:   row.DictCode,
			DictSort:   int32(row.DictSort),
			DictLabel:  row.DictLabel,
			DictValue:  row.DictValue,
			DictType:   row.DictType,
			CssClass:   "",
			ListClass:  "",
			IsDefault:  row.IsDefault,
			Remark:     "",
			CreateTime: "",
		}
		if row.CssClass.Valid {
			vo.CssClass = row.CssClass.String
		}
		if row.ListClass.Valid {
			vo.ListClass = row.ListClass.String
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
