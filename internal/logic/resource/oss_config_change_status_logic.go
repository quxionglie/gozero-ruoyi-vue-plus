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

type OssConfigChangeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 状态修改
func NewOssConfigChangeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OssConfigChangeStatusLogic {
	return &OssConfigChangeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OssConfigChangeStatusLogic) OssConfigChangeStatus(req *types.OssConfigChangeStatusReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.OssConfigId == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "配置ID不能为空",
		}, nil
	}
	if req.Status == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "状态不能为空",
		}, nil
	}

	// 2. 查询配置是否存在
	_, err = l.svcCtx.SysOssConfigModel.FindOne(l.ctx, req.OssConfigId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BaseResp{
				Code: 500,
				Msg:  "对象存储配置不存在",
			}, nil
		}
		l.Errorf("查询对象存储配置失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询对象存储配置失败",
		}, err
	}

	// 3. 更新状态
	err = l.svcCtx.SysOssConfigModel.UpdateStatus(l.ctx, req.OssConfigId, req.Status)
	if err != nil {
		l.Errorf("修改对象存储配置状态失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改对象存储配置状态失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
