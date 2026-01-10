// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"database/sql"
	"fmt"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictDataEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改字典数据
func NewDictDataEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDataEditLogic {
	return &DictDataEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictDataEditLogic) DictDataEdit(req *types.DictDataReq) (resp *types.BaseResp, err error) {
	if req.DictCode == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "字典编码不能为空",
		}, nil
	}

	// 1. 校验字典数据唯一性
	unique, err := l.svcCtx.SysDictDataModel.CheckDictDataUnique(l.ctx, req.DictType, req.DictValue, req.DictCode)
	if err != nil {
		l.Errorf("校验字典数据唯一性失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "校验字典数据唯一性失败",
		}, err
	}
	if !unique {
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("修改字典数据'%s'失败，字典键值已存在", req.DictValue),
		}, nil
	}

	// 2. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 3. 更新字典数据
	dictData := &model.SysDictData{
		DictCode:  req.DictCode,
		DictSort:  int64(req.DictSort),
		DictLabel: req.DictLabel,
		DictValue: req.DictValue,
		DictType:  req.DictType,
		CssClass:  sql.NullString{String: req.CssClass, Valid: req.CssClass != ""},
		ListClass: sql.NullString{String: req.ListClass, Valid: req.ListClass != ""},
		IsDefault: req.IsDefault,
		Remark:    sql.NullString{String: req.Remark, Valid: req.Remark != ""},
		UpdateBy:  sql.NullInt64{Int64: userId, Valid: userId > 0},
	}

	err = l.svcCtx.SysDictDataModel.Update(l.ctx, dictData)
	if err != nil {
		l.Errorf("修改字典数据失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改字典数据失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
