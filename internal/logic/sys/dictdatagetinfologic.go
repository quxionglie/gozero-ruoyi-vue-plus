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

type DictDataGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询字典数据详细
func NewDictDataGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDataGetInfoLogic {
	return &DictDataGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictDataGetInfoLogic) DictDataGetInfo(req *types.DictDataGetInfoReq) (resp *types.DictDataResp, err error) {
	if req.DictCode == 0 {
		return &types.DictDataResp{
			BaseResp: types.BaseResp{
				Code: 400,
				Msg:  "字典编码不能为空",
			},
		}, nil
	}

	dictData, err := l.svcCtx.SysDictDataModel.FindOne(l.ctx, req.DictCode)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.DictDataResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "字典数据不存在",
				},
			}, nil
		}
		l.Errorf("查询字典数据失败: %v", err)
		return &types.DictDataResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询字典数据失败",
			},
		}, err
	}

	vo := types.DictDataVo{
		DictCode:   dictData.DictCode,
		DictSort:   int32(dictData.DictSort),
		DictLabel:  dictData.DictLabel,
		DictValue:  dictData.DictValue,
		DictType:   dictData.DictType,
		CssClass:   "",
		ListClass:  "",
		IsDefault:  dictData.IsDefault,
		Remark:     "",
		CreateTime: "",
	}
	if dictData.CssClass.Valid {
		vo.CssClass = dictData.CssClass.String
	}
	if dictData.ListClass.Valid {
		vo.ListClass = dictData.ListClass.String
	}
	if dictData.Remark.Valid {
		vo.Remark = dictData.Remark.String
	}
	if dictData.CreateTime.Valid {
		vo.CreateTime = dictData.CreateTime.Time.Format("2006-01-02 15:04:05")
	}

	return &types.DictDataResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: vo,
	}, nil
}
