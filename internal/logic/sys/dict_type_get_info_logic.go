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

type DictTypeGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询字典类型详细
func NewDictTypeGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeGetInfoLogic {
	return &DictTypeGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictTypeGetInfoLogic) DictTypeGetInfo(req *types.DictTypeGetInfoReq) (resp *types.DictTypeResp, err error) {
	if req.DictId == 0 {
		return &types.DictTypeResp{
			BaseResp: types.BaseResp{
				Code: 400,
				Msg:  "字典ID不能为空",
			},
		}, nil
	}

	// 查询字典类型
	dictType, err := l.svcCtx.SysDictTypeModel.FindOne(l.ctx, req.DictId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.DictTypeResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "字典类型不存在",
				},
			}, nil
		}
		l.Errorf("查询字典类型失败: %v", err)
		return &types.DictTypeResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询字典类型失败",
			},
		}, err
	}

	// 转换为 VO
	vo := types.DictTypeVo{
		DictId:     dictType.DictId,
		DictName:   dictType.DictName,
		DictType:   dictType.DictType,
		Remark:     "",
		CreateTime: "",
	}
	if dictType.Remark.Valid {
		vo.Remark = dictType.Remark.String
	}
	if dictType.CreateTime.Valid {
		vo.CreateTime = dictType.CreateTime.Time.Format("2006-01-02 15:04:05")
	}

	return &types.DictTypeResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: vo,
	}, nil
}
