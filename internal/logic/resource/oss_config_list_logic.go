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

type OssConfigListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询对象存储配置列表
func NewOssConfigListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OssConfigListLogic {
	return &OssConfigListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OssConfigListLogic) OssConfigList(req *types.OssConfigListReq) (resp *types.OssConfigListResp, err error) {
	// 1. 构建查询条件
	ossConfigQuery := &model.OssConfigQuery{
		ConfigKey: req.ConfigKey,
		Status:    req.Status,
	}

	// 2. 构建分页查询条件
	pageQuery := &model.PageQuery{
		PageNum:       req.PageNum,
		PageSize:      req.PageSize,
		OrderByColumn: req.OrderByColumn,
		IsAsc:         req.IsAsc,
	}

	// 3. 查询数据
	ossConfigList, total, err := l.svcCtx.SysOssConfigModel.FindPage(l.ctx, ossConfigQuery, pageQuery)
	if err != nil {
		l.Errorf("查询对象存储配置列表失败: %v", err)
		return &types.OssConfigListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询对象存储配置列表失败",
			},
		}, err
	}

	// 4. 转换为响应格式
	rows := make([]types.OssConfigVo, 0, len(ossConfigList))
	for _, ossConfig := range ossConfigList {
		ossConfigVo := l.convertToOssConfigVo(ossConfig)
		rows = append(rows, ossConfigVo)
	}

	return &types.OssConfigListResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Total: total,
		Rows:  rows,
	}, nil
}

// convertToOssConfigVo 转换OSS配置实体为响应格式
func (l *OssConfigListLogic) convertToOssConfigVo(ossConfig *model.SysOssConfig) types.OssConfigVo {
	ossConfigVo := types.OssConfigVo{
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
		ossConfigVo.Remark = ossConfig.Remark.String
	}
	if ossConfig.CreateTime.Valid {
		ossConfigVo.CreateTime = ossConfig.CreateTime.Time.Format("2006-01-02 15:04:05")
	}
	return ossConfigVo
}
