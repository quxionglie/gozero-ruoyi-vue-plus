// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictDataByTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据字典类型查询字典数据
func NewDictDataByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDataByTypeLogic {
	return &DictDataByTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictDataByTypeLogic) DictDataByType(req *types.DictDataByTypeReq) (resp *types.DictDataListResp, err error) {
	if req.DictType == "" {
		return &types.DictDataListResp{
			BaseResp: types.BaseResp{
				Code: 400,
				Msg:  "字典类型不能为空",
			},
			Data: []types.DictDataVo{},
		}, nil
	}

	rows, err := l.svcCtx.SysDictDataModel.FindByDictType(l.ctx, req.DictType)
	if err != nil {
		l.Errorf("查询字典数据失败: %v", err)
		return &types.DictDataListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询字典数据失败",
			},
			Data: []types.DictDataVo{},
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

	return &types.DictDataListResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: voList,
	}, nil
}
