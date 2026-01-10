// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictTypeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询字典类型列表
func NewDictTypeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeListLogic {
	return &DictTypeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictTypeListLogic) DictTypeList() (resp *types.TableDataInfoResp, err error) {
	// 查询字典类型列表
	rows, err := l.svcCtx.SysDictTypeModel.FindAll(l.ctx)
	if err != nil {
		l.Errorf("查询字典类型列表失败: %v", err)
		return &types.TableDataInfoResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询字典类型列表失败",
			},
			Data: types.TableDataInfo{
				Total: 0,
				Rows:  []types.DictTypeVo{},
			},
		}, err
	}

	// 转换为 VO
	voList := make([]types.DictTypeVo, 0, len(rows))
	for _, row := range rows {
		vo := types.DictTypeVo{
			DictId:     row.DictId,
			DictName:   row.DictName,
			DictType:   row.DictType,
			Remark:     "",
			CreateTime: "",
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
