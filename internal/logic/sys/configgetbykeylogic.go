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

type ConfigGetByKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据参数键名查询参数值
func NewConfigGetByKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigGetByKeyLogic {
	return &ConfigGetByKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigGetByKeyLogic) ConfigGetByKey(req *types.ConfigGetByKeyReq) (resp *types.ConfigValueResp, err error) {
	if req.ConfigKey == "" {
		return &types.ConfigValueResp{
			BaseResp: types.BaseResp{
				Code: 400,
				Msg:  "参数键名不能为空",
			},
			Data: "",
		}, nil
	}

	config, err := l.svcCtx.SysConfigModel.FindByConfigKey(l.ctx, req.ConfigKey)
	if err != nil {
		if err == model.ErrNotFound {
			// 配置不存在时返回空字符串
			return &types.ConfigValueResp{
				BaseResp: types.BaseResp{
					Code: 200,
					Msg:  "操作成功",
				},
				Data: "",
			}, nil
		}
		l.Errorf("查询参数配置失败: %v", err)
		return &types.ConfigValueResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询参数配置失败",
			},
			Data: "",
		}, err
	}

	return &types.ConfigValueResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: config.ConfigValue,
	}, nil
}
