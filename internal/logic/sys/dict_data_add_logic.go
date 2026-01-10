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

type DictDataAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增字典数据
func NewDictDataAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDataAddLogic {
	return &DictDataAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictDataAddLogic) DictDataAdd(req *types.DictDataReq) (resp *types.BaseResp, err error) {
	// 1. 参数长度校验
	if err := util.ValidateStringLength(req.DictLabel, "字典标签", 100); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}
	if err := util.ValidateStringLength(req.DictValue, "字典键值", 100); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}
	if err := util.ValidateStringLength(req.DictType, "字典类型", 100); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}
	if req.CssClass != "" {
		if err := util.ValidateStringLength(req.CssClass, "样式属性", 100); err != nil {
			return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
		}
	}

	// 2. 校验字典数据唯一性（同一字典类型下，字典键值唯一）
	unique, err := l.svcCtx.SysDictDataModel.CheckDictDataUnique(l.ctx, req.DictType, req.DictValue, 0)
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
			Msg:  fmt.Sprintf("新增字典数据'%s'失败，字典键值已存在", req.DictValue),
		}, nil
	}

	// 3. 获取当前用户信息
	userId, _ := util.GetUserIdFromContext(l.ctx)
	tenantId, _ := util.GetTenantIdFromContext(l.ctx)

	// 获取用户的部门ID
	var deptId int64
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err == nil && user.DeptId.Valid {
		deptId = user.DeptId.Int64
	}

	// 生成主键ID（使用雪花算法）
	dictCode, err := util.GenerateID()
	if err != nil {
		l.Errorf("生成字典数据ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "生成字典数据ID失败",
		}, err
	}

	// 4. 构建字典数据实体
	dictData := &model.SysDictData{
		DictCode:   dictCode,
		TenantId:   tenantId,
		DictSort:   int64(req.DictSort),
		DictLabel:  req.DictLabel,
		DictValue:  req.DictValue,
		DictType:   req.DictType,
		CssClass:   sql.NullString{String: req.CssClass, Valid: req.CssClass != ""},
		ListClass:  sql.NullString{String: req.ListClass, Valid: req.ListClass != ""},
		IsDefault:  req.IsDefault,
		Remark:     sql.NullString{String: req.Remark, Valid: req.Remark != ""},
		CreateDept: sql.NullInt64{Int64: deptId, Valid: deptId > 0},
		CreateBy:   sql.NullInt64{Int64: userId, Valid: userId > 0},
	}

	// 5. 插入数据库
	_, err = l.svcCtx.SysDictDataModel.Insert(l.ctx, dictData)
	if err != nil {
		l.Errorf("新增字典数据失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "新增字典数据失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
