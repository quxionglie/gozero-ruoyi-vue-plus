package oss

import (
	"context"
	"fmt"
	"sync"

	model "gozero-ruoyi-vue-plus/internal/model/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	clientCache = make(map[string]OssClient)
	cacheLock   sync.RWMutex
)

// GetDefaultClient 获取默认OSS客户端（status=0的配置）
func GetDefaultClient(ctx context.Context, ossConfigModel model.SysOssConfigModel, tenantId string) (OssClient, error) {
	// 查询默认配置
	ossConfig, err := ossConfigModel.FindDefault(ctx, tenantId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, fmt.Errorf("默认OSS配置不存在")
		}
		return nil, fmt.Errorf("查询默认OSS配置失败: %w", err)
	}

	return GetClientByConfigKey(ctx, ossConfigModel, ossConfig.ConfigKey, tenantId)
}

// GetClientByConfigKey 根据配置键获取OSS客户端
func GetClientByConfigKey(ctx context.Context, ossConfigModel model.SysOssConfigModel, configKey string, tenantId string) (OssClient, error) {
	if configKey == "" {
		return nil, fmt.Errorf("配置键不能为空")
	}

	// 查询配置
	ossConfig, err := ossConfigModel.FindByConfigKey(ctx, configKey, tenantId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, fmt.Errorf("OSS配置不存在: %s", configKey)
		}
		return nil, fmt.Errorf("查询OSS配置失败: %w", err)
	}

	// 构建缓存键
	cacheKey := configKey
	if tenantId != "" {
		cacheKey = tenantId + ":" + configKey
	}

	// 先从缓存获取
	cacheLock.RLock()
	client, exists := clientCache[cacheKey]
	cacheLock.RUnlock()

	if exists && client != nil {
		return client, nil
	}

	// 创建客户端
	cacheLock.Lock()
	defer cacheLock.Unlock()

	// 双重检查
	client, exists = clientCache[cacheKey]
	if exists && client != nil {
		return client, nil
	}

	// 构建属性
	properties := &OssProperties{
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
		TenantId:     ossConfig.TenantId,
	}

	// 目前只支持MinIO，后续可以扩展
	client, err = NewMinioClient(properties)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %w", err)
	}

	// 存入缓存
	clientCache[cacheKey] = client
	logx.Infof("创建OSS客户端: configKey=%s, tenantId=%s", configKey, tenantId)

	return client, nil
}

// ClearCache 清空客户端缓存（用于配置更新后）
func ClearCache() {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	clientCache = make(map[string]OssClient)
	logx.Info("清空OSS客户端缓存")
}
