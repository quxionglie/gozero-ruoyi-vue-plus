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

type DictTypeAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增字典类型
func NewDictTypeAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeAddLogic {
	return &DictTypeAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictTypeAddLogic) DictTypeAdd(req *types.DictTypeReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验（长度和格式）
	if err := util.ValidateStringLength(req.DictName, "字典名称", 100); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}
	if err := util.ValidateStringLength(req.DictType, "字典类型", 100); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}
	if err := util.ValidateDictType(req.DictType); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}

	// 2. 校验字典类型唯一性
	unique, err := l.svcCtx.SysDictTypeModel.CheckDictTypeUnique(l.ctx, req.DictType, 0)
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
			Msg:  fmt.Sprintf("新增字典'%s'失败，字典类型已存在", req.DictName),
		}, nil
	}

	// 3. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 4. 构建字典类型实体
	dictType := &model.SysDictType{
		TenantId: "",
		DictName: req.DictName,
		DictType: req.DictType,
		Remark:   sql.NullString{String: req.Remark, Valid: req.Remark != ""},
		CreateBy: sql.NullInt64{Int64: userId, Valid: userId > 0},
	}

	// 5. 插入数据库
	_, err = l.svcCtx.SysDictTypeModel.Insert(l.ctx, dictType)
	if err != nil {
		l.Errorf("新增字典类型失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "新增字典类型失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
