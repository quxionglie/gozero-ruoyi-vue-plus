// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除参数配置
func NewConfigRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigRemoveLogic {
	return &ConfigRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigRemoveLogic) ConfigRemove(req *types.ConfigRemoveReq) (resp *types.BaseResp, err error) {
	if req.ConfigIds == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "参数ID不能为空",
		}, nil
	}

	// 解析参数ID列表（逗号分隔）
	configIdStrs := strings.Split(req.ConfigIds, ",")
	var configIds []int64
	for _, idStr := range configIdStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return &types.BaseResp{
				Code: 400,
				Msg:  fmt.Sprintf("参数ID格式错误: %s", idStr),
			}, nil
		}
		configIds = append(configIds, id)
	}

	if len(configIds) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "参数ID不能为空",
		}, nil
	}

	// 检查系统内置参数，不能删除
	for _, configId := range configIds {
		config, err := l.svcCtx.SysConfigModel.FindOne(l.ctx, configId)
		if err != nil {
			continue
		}

		// 系统内置参数（configType='Y'）不能删除
		if config.ConfigType == "Y" {
			return &types.BaseResp{
				Code: 500,
				Msg:  fmt.Sprintf("内置参数【%s】不能删除", config.ConfigKey),
			}, nil
		}

		// 清除缓存
		cacheKey := "sys_config:" + config.ConfigKey
		_, _ = l.svcCtx.RedisConn.DelCtx(l.ctx, cacheKey)
	}

	// 删除参数配置
	for _, configId := range configIds {
		err = l.svcCtx.SysConfigModel.Delete(l.ctx, configId)
		if err != nil {
			l.Errorf("删除参数配置失败: configId=%d, err=%v", configId, err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除参数配置失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
