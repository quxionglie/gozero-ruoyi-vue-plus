// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改参数配置
func NewConfigEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigEditLogic {
	return &ConfigEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigEditLogic) ConfigEdit(req *types.ConfigReq) (resp *types.BaseResp, err error) {
	if req.ConfigId == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "参数ID不能为空",
		}, nil
	}

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
	unique, err := l.svcCtx.SysConfigModel.CheckConfigKeyUnique(l.ctx, req.ConfigKey, req.ConfigId)
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
			Msg:  fmt.Sprintf("修改参数'%s'失败，参数键名已存在", req.ConfigName),
		}, nil
	}

	// 3. 查询原配置（用于清除旧缓存和保留创建信息）
	oldConfig, err := l.svcCtx.SysConfigModel.FindOne(l.ctx, req.ConfigId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BaseResp{
				Code: 500,
				Msg:  "参数配置不存在",
			}, nil
		}
		l.Errorf("查询参数配置失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询参数配置失败",
		}, err
	}

	// 4. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 5. 更新参数配置
	config := &model.SysConfig{
		ConfigId:    req.ConfigId,
		ConfigName:  req.ConfigName,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		ConfigType:  req.ConfigType,
		Remark:      sql.NullString{String: req.Remark, Valid: req.Remark != ""},
		CreateDept:  oldConfig.CreateDept, // 保持原部门ID
		CreateBy:    oldConfig.CreateBy,   // 保持原创建者
		CreateTime:  oldConfig.CreateTime, // 保持原创建时间
		UpdateBy:    sql.NullInt64{Int64: userId, Valid: userId > 0},
		UpdateTime:  sql.NullTime{Time: time.Now(), Valid: true},
	}

	err = l.svcCtx.SysConfigModel.Update(l.ctx, config)
	if err != nil {
		l.Errorf("修改参数配置失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改参数配置失败",
		}, err
	}

	// 5. 如果配置键改变，清除旧键的缓存
	if oldConfig != nil && oldConfig.ConfigKey != req.ConfigKey {
		oldCacheKey := "sys_config:" + oldConfig.ConfigKey
		_, _ = l.svcCtx.RedisConn.DelCtx(l.ctx, oldCacheKey)
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
