// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package resource

import (
	"context"
	"database/sql"
	"time"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type OssConfigEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改对象存储配置
func NewOssConfigEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OssConfigEditLogic {
	return &OssConfigEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OssConfigEditLogic) OssConfigEdit(req *types.OssConfigReq) (resp *types.BaseResp, err error) {
	if req.OssConfigId == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "配置ID不能为空",
		}, nil
	}

	// 1. 参数校验
	if req.ConfigKey == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "配置key不能为空",
		}, nil
	}
	if req.AccessKey == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "accessKey不能为空",
		}, nil
	}
	if req.SecretKey == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "秘钥不能为空",
		}, nil
	}
	if req.BucketName == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "桶名称不能为空",
		}, nil
	}
	if req.Endpoint == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "访问站点不能为空",
		}, nil
	}

	// 2. 查询原配置
	oldOssConfig, err := l.svcCtx.SysOssConfigModel.FindOne(l.ctx, req.OssConfigId)
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

	// 3. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 4. 更新OSS配置
	ossConfig := &model.SysOssConfig{
		OssConfigId:  req.OssConfigId,
		TenantId:     oldOssConfig.TenantId, // 保持原租户ID
		ConfigKey:    req.ConfigKey,
		AccessKey:    req.AccessKey,
		SecretKey:    req.SecretKey,
		BucketName:   req.BucketName,
		Prefix:       req.Prefix,
		Endpoint:     req.Endpoint,
		Domain:       req.Domain,
		IsHttps:      req.IsHttps,
		Region:       req.Region,
		AccessPolicy: req.AccessPolicy,
		Status:       req.Status,
		Ext1:         req.Ext1,
		CreateDept:   oldOssConfig.CreateDept, // 保持原部门ID
		CreateBy:     oldOssConfig.CreateBy,   // 保持原创建者
		CreateTime:   oldOssConfig.CreateTime, // 保持原创建时间
		UpdateBy:     sql.NullInt64{Int64: userId, Valid: userId > 0},
		UpdateTime:   sql.NullTime{Time: time.Now(), Valid: true},
		Remark:       sql.NullString{String: req.Remark, Valid: req.Remark != ""},
	}

	err = l.svcCtx.SysOssConfigModel.Update(l.ctx, ossConfig)
	if err != nil {
		l.Errorf("修改对象存储配置失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改对象存储配置失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
