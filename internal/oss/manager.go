package oss

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
)

// OssManager OSS客户端管理器
type OssManager struct {
	ossConfigModel model.SysOssConfigModel
}

// NewOssManager 创建OSS管理器
func NewOssManager(ossConfigModel model.SysOssConfigModel) *OssManager {
	return &OssManager{
		ossConfigModel: ossConfigModel,
	}
}

// GetDefaultClient 获取默认OSS客户端（status=0的配置）
func (m *OssManager) GetDefaultClient(ctx context.Context, tenantId string) (OssClient, error) {
	return GetDefaultClient(ctx, m.ossConfigModel, tenantId)
}

// GetClientByConfigKey 根据配置键获取OSS客户端
func (m *OssManager) GetClientByConfigKey(ctx context.Context, configKey string, tenantId string) (OssClient, error) {
	return GetClientByConfigKey(ctx, m.ossConfigModel, configKey, tenantId)
}
