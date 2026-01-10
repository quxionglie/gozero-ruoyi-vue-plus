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

type ConfigAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增参数配置
func NewConfigAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigAddLogic {
	return &ConfigAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigAddLogic) ConfigAdd(req *types.ConfigReq) (resp *types.BaseResp, err error) {
	// 1. 参数长度校验
	if err := util.ValidateStringLength(req.ConfigName, "参数名称", 100); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}
	if err := util.ValidateStringLength(req.ConfigKey, "参数键名", 100); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}
	if err := util.ValidateStringLength(req.ConfigValue, "参数键值", 500); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}

	// 2. 校验参数键名唯一性
	unique, err := l.svcCtx.SysConfigModel.CheckConfigKeyUnique(l.ctx, req.ConfigKey, 0)
	if err != nil {
		l.Errorf("校验参数键名唯一性失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "校验参数键名唯一性失败",
		}, err
	}
	if !unique {
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("新增参数'%s'失败，参数键名已存在", req.ConfigName),
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
	configId, err := util.GenerateID()
	if err != nil {
		l.Errorf("生成配置ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "生成配置ID失败",
		}, err
	}

	// 4. 构建参数配置实体
	config := &model.SysConfig{
		ConfigId:    configId,
		TenantId:    tenantId,
		ConfigName:  req.ConfigName,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		ConfigType:  req.ConfigType,
		Remark:      sql.NullString{String: req.Remark, Valid: req.Remark != ""},
		CreateDept:  sql.NullInt64{Int64: deptId, Valid: deptId > 0},
		CreateBy:    sql.NullInt64{Int64: userId, Valid: userId > 0},
	}

	// 5. 插入数据库
	_, err = l.svcCtx.SysConfigModel.Insert(l.ctx, config)
	if err != nil {
		l.Errorf("新增参数配置失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "新增参数配置失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
