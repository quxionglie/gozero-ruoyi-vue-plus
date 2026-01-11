// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package resource

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OssConfigGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取对象存储配置详细信息
func NewOssConfigGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OssConfigGetInfoLogic {
	return &OssConfigGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OssConfigGetInfoLogic) OssConfigGetInfo(req *types.OssConfigGetInfoReq) (resp *types.OssConfigResp, err error) {
	if req.OssConfigId == 0 {
		return &types.OssConfigResp{
			BaseResp: types.BaseResp{
				Code: 400,
				Msg:  "配置ID不能为空",
			},
		}, nil
	}

	ossConfig, err := l.svcCtx.SysOssConfigModel.FindOne(l.ctx, req.OssConfigId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.OssConfigResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "对象存储配置不存在",
				},
			}, nil
		}
		l.Errorf("查询对象存储配置失败: %v", err)
		return &types.OssConfigResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询对象存储配置失败",
			},
		}, err
	}

	vo := types.OssConfigVo{
		OssConfigId:  ossConfig.OssConfigId,
		ConfigKey:    ossConfig.ConfigKey,
		AccessKey:    ossConfig.AccessKey,
		SecretKey:    ossConfig.SecretKey,
		BucketName:   ossConfig.BucketName,
		Prefix:       ossConfig.Prefix,
		Endpoint:     ossConfig.Endpoint,
		Domain:       ossConfig.Domain,
		IsHttps:      ossConfig.IsHttps,
		Region:       ossConfig.Region,
		AccessPolicy: ossConfig.AccessPolicy,
		Status:       ossConfig.Status,
		Ext1:         ossConfig.Ext1,
		Remark:       "",
		CreateTime:   "",
	}
	if ossConfig.Remark.Valid {
		vo.Remark = ossConfig.Remark.String
	}
	if ossConfig.CreateTime.Valid {
		vo.CreateTime = ossConfig.CreateTime.Time.Format("2006-01-02 15:04:05")
	}

	return &types.OssConfigResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: vo,
	}, nil
}
