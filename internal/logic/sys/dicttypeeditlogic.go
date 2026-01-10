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

type DictTypeEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改字典类型
func NewDictTypeEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeEditLogic {
	return &DictTypeEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictTypeEditLogic) DictTypeEdit(req *types.DictTypeReq) (resp *types.BaseResp, err error) {
	if req.DictId == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "字典ID不能为空",
		}, nil
	}

	// 1. 校验字典类型唯一性
	unique, err := l.svcCtx.SysDictTypeModel.CheckDictTypeUnique(l.ctx, req.DictType, req.DictId)
	if err != nil {
		l.Errorf("校验字典类型唯一性失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "校验字典类型唯一性失败",
		}, err
	}
	if !unique {
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("修改字典'%s'失败，字典类型已存在", req.DictName),
		}, nil
	}

	// 2. 查询原字典类型
	oldDictType, err := l.svcCtx.SysDictTypeModel.FindOne(l.ctx, req.DictId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BaseResp{
				Code: 500,
				Msg:  "字典类型不存在",
			}, nil
		}
		l.Errorf("查询字典类型失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询字典类型失败",
		}, err
	}

	// 3. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 4. 更新字典类型
	dictType := &model.SysDictType{
		DictId:   req.DictId,
		DictName: req.DictName,
		DictType: req.DictType,
		Remark:   sql.NullString{String: req.Remark, Valid: req.Remark != ""},
		UpdateBy: sql.NullInt64{Int64: userId, Valid: userId > 0},
	}

	err = l.svcCtx.SysDictTypeModel.Update(l.ctx, dictType)
	if err != nil {
		l.Errorf("修改字典类型失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改字典类型失败",
		}, err
	}

	// 5. 如果字典类型改变，需要更新关联的字典数据
	if oldDictType.DictType != req.DictType {
		err = l.svcCtx.SysDictDataModel.UpdateDictTypeByOldDictType(l.ctx, oldDictType.DictType, req.DictType)
		if err != nil {
			l.Errorf("更新关联字典数据失败: %v", err)
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
