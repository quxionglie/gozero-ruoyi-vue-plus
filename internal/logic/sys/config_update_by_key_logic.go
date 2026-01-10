// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"database/sql"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigUpdateByKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据参数键名修改参数配置
func NewConfigUpdateByKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigUpdateByKeyLogic {
	return &ConfigUpdateByKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigUpdateByKeyLogic) ConfigUpdateByKey(req *types.ConfigReq) (resp *types.BaseResp, err error) {
	if req.ConfigKey == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "参数键名不能为空",
		}, nil
	}

	// 1. 参数长度校验
	if err := util.ValidateStringLength(req.ConfigValue, "参数键值", 500); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}

	// 2. 根据配置键查询配置
	config, err := l.svcCtx.SysConfigModel.FindByConfigKey(l.ctx, req.ConfigKey)
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

	// 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 更新配置（只更新配置值和备注）
	updateConfig := &model.SysConfig{
		ConfigId:    config.ConfigId,
		ConfigKey:   config.ConfigKey, // 保持原键名
		ConfigValue: req.ConfigValue,
		Remark:      sql.NullString{String: req.Remark, Valid: req.Remark != ""},
		UpdateBy:    sql.NullInt64{Int64: userId, Valid: userId > 0},
	}

	err = l.svcCtx.SysConfigModel.Update(l.ctx, updateConfig)
	if err != nil {
		l.Errorf("根据参数键名修改参数配置失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改参数配置失败",
		}, err
	}

	// 清除缓存
	cacheKey := "sys_config:" + req.ConfigKey
	_, _ = l.svcCtx.RedisConn.DelCtx(l.ctx, cacheKey)

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
