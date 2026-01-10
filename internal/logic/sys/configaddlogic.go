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
	// 1. 校验参数键名唯一性
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

	// 2. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 3. 构建参数配置实体
	config := &model.SysConfig{
		TenantId:    "",
		ConfigName:  req.ConfigName,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		ConfigType:  req.ConfigType,
		Remark:      sql.NullString{String: req.Remark, Valid: req.Remark != ""},
		CreateBy:    sql.NullInt64{Int64: userId, Valid: userId > 0},
	}

	// 4. 插入数据库
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
