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

type ConfigGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询参数配置详细
func NewConfigGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigGetInfoLogic {
	return &ConfigGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigGetInfoLogic) ConfigGetInfo(req *types.ConfigGetInfoReq) (resp *types.ConfigResp, err error) {
	if req.ConfigId == 0 {
		return &types.ConfigResp{
			BaseResp: types.BaseResp{
				Code: 400,
				Msg:  "参数ID不能为空",
			},
		}, nil
	}

	config, err := l.svcCtx.SysConfigModel.FindOne(l.ctx, req.ConfigId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.ConfigResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "参数配置不存在",
				},
			}, nil
		}
		l.Errorf("查询参数配置失败: %v", err)
		return &types.ConfigResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询参数配置失败",
			},
		}, err
	}

	vo := types.ConfigVo{
		ConfigId:    config.ConfigId,
		ConfigName:  config.ConfigName,
		ConfigKey:   config.ConfigKey,
		ConfigValue: config.ConfigValue,
		ConfigType:  config.ConfigType,
		Remark:      "",
		CreateTime:  "",
	}
	if config.Remark.Valid {
		vo.Remark = config.Remark.String
	}
	if config.CreateTime.Valid {
		vo.CreateTime = config.CreateTime.Time.Format("2006-01-02 15:04:05")
	}

	return &types.ConfigResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: vo,
	}, nil
}
