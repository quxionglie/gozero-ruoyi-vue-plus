package oss

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/zeromicro/go-zero/core/logx"
)

// MinioClient MinIO客户端实现
type MinioClient struct {
	properties *OssProperties
	client     *minio.Client
	bucketName string
	configKey  string
}

// NewMinioClient 创建MinIO客户端
func NewMinioClient(properties *OssProperties) (*MinioClient, error) {
	if properties == nil {
		return nil, fmt.Errorf("OSS配置不能为空")
	}

	// 确定端点URL
	endpoint := properties.Endpoint
	if endpoint == "" {
		return nil, fmt.Errorf("端点不能为空")
	}

	// 确定是否使用HTTPS
	useSSL := properties.IsHttps == "Y"

	// 如果端点包含协议，提取主机名
	if strings.HasPrefix(endpoint, "http://") {
		endpoint = strings.TrimPrefix(endpoint, "http://")
		useSSL = false
	} else if strings.HasPrefix(endpoint, "https://") {
		endpoint = strings.TrimPrefix(endpoint, "https://")
		useSSL = true
	}

	// 创建MinIO客户端
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(properties.AccessKey, properties.SecretKey, ""),
		Secure: useSSL,
		Region: properties.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("创建MinIO客户端失败: %w", err)
	}

	client := &MinioClient{
		properties: properties,
		client:     minioClient,
		bucketName: properties.BucketName,
		configKey:  properties.ConfigKey,
	}

	// 检查bucket是否存在，如果不存在则创建
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, properties.BucketName)
	if err != nil {
		return nil, fmt.Errorf("检查bucket是否存在失败: %w", err)
	}

	if !exists {
		// 创建bucket
		err = minioClient.MakeBucket(ctx, properties.BucketName, minio.MakeBucketOptions{
			Region: properties.Region,
		})
		if err != nil {
			return nil, fmt.Errorf("创建bucket失败: %w", err)
		}
		logx.Infof("成功创建bucket: %s", properties.BucketName)
	}

	return client, nil
}

// UploadSuffix 上传文件（使用后缀构造对象键）
func (c *MinioClient) UploadSuffix(data []byte, suffix string, contentType string) (*UploadResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("文件数据不能为空")
	}

	// 生成对象键（文件名）
	objectKey := c.GetPath(c.properties.Prefix, suffix)

	// 上传文件
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	reader := bytes.NewReader(data)
	_, err := c.client.PutObject(ctx, c.bucketName, objectKey, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	// 构建URL
	url := c.GetUrl() + "/" + objectKey

	return &UploadResult{
		URL:      url,
		Filename: objectKey,
		ETag:     "", // MinIO的ETag在上传时不会立即返回，如果需要可以重新获取
	}, nil
}

// Download 下载文件到输出流
func (c *MinioClient) Download(key string, out io.Writer) error {
	// 移除基础URL（如果包含）
	objectKey := c.removeBaseUrl(key)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	object, err := c.client.GetObject(ctx, c.bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("获取对象失败: %w", err)
	}
	defer object.Close()

	_, err = io.Copy(out, object)
	if err != nil {
		return fmt.Errorf("下载文件失败: %w", err)
	}

	return nil
}

// Delete 删除文件
func (c *MinioClient) Delete(path string) error {
	// 移除基础URL（如果包含）
	objectKey := c.removeBaseUrl(path)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.client.RemoveObject(ctx, c.bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// GetConfigKey 获取配置键
func (c *MinioClient) GetConfigKey() string {
	return c.configKey
}

// GetAccessPolicy 获取访问策略类型
func (c *MinioClient) GetAccessPolicy() AccessPolicyType {
	policy := strings.TrimSpace(c.properties.AccessPolicy)
	switch policy {
	case "1":
		return AccessPolicyPublic
	case "2":
		return AccessPolicyCustom
	default:
		return AccessPolicyPrivate
	}
}

// GetUrl 获取文件URL（基础URL）
func (c *MinioClient) GetUrl() string {
	domain := c.properties.Domain
	endpoint := c.properties.Endpoint
	header := c.getIsHttps()

	// 如果有自定义域名，使用域名
	if domain != "" {
		if strings.HasPrefix(domain, "http://") || strings.HasPrefix(domain, "https://") {
			return domain + "/" + c.bucketName
		}
		return header + domain + "/" + c.bucketName
	}

	// 否则使用端点
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint + "/" + c.bucketName
	}
	return header + endpoint + "/" + c.bucketName
}

// CreatePresignedGetUrl 创建预签名下载URL
func (c *MinioClient) CreatePresignedGetUrl(objectKey string, expiredTime time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url, err := c.client.PresignedGetObject(ctx, c.bucketName, objectKey, expiredTime, nil)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %w", err)
	}

	return url.String(), nil
}

// getPath 生成文件路径（对象键）
// 格式: prefix/yyyy/MM/dd/uuid + suffix
func (c *MinioClient) getPath(suffix string) string {
	// 生成UUID
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")

	// 生成日期路径
	now := time.Now()
	datePath := fmt.Sprintf("%04d/%02d/%02d", now.Year(), now.Month(), now.Day())

	// 拼接路径
	var path string
	if c.properties.Prefix != "" {
		path = c.properties.Prefix + "/" + datePath + "/" + id
	} else {
		path = datePath + "/" + id
	}

	return path + suffix
}

// removeBaseUrl 移除路径中的基础URL部分，得到相对路径
func (c *MinioClient) removeBaseUrl(path string) string {
	baseUrl := c.GetUrl() + "/"
	return strings.TrimPrefix(path, baseUrl)
}

// getIsHttps 获取是否使用HTTPS的协议头部
func (c *MinioClient) getIsHttps() string {
	if c.properties.IsHttps == "Y" {
		return "https://"
	}
	return "http://"
}
