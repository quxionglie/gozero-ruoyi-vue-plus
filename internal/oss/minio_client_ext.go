package oss

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// Upload 上传文件（指定文件路径和对象键）
func (c *MinioClient) Upload(filePath string, key string, md5Digest string, contentType string) (*UploadResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}
	if md5Digest != "" {
		// MinIO 支持 MD5 校验，但需要在 PutObjectOptions 中设置
		// 注意：MinIO Go SDK 的 PutObjectOptions 没有直接的 MD5 字段
		// 这里可以后续扩展
	}

	_, err = c.client.PutObject(ctx, c.bucketName, key, file, fileInfo.Size(), opts)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	url := c.GetUrl() + "/" + key
	return &UploadResult{
		URL:      url,
		Filename: key,
		ETag:     "",
	}, nil
}

// UploadStream 上传输入流
func (c *MinioClient) UploadStream(inputStream io.Reader, key string, length int64, contentType string) (*UploadResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := c.client.PutObject(ctx, c.bucketName, key, inputStream, length, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	url := c.GetUrl() + "/" + key
	return &UploadResult{
		URL:      url,
		Filename: key,
		ETag:     "",
	}, nil
}

// UploadSuffixStream 上传输入流（使用后缀构造对象键）
func (c *MinioClient) UploadSuffixStream(inputStream io.Reader, suffix string, length int64, contentType string) (*UploadResult, error) {
	objectKey := c.GetPath(c.properties.Prefix, suffix)
	return c.UploadStream(inputStream, objectKey, length, contentType)
}

// UploadSuffixFile 上传文件（使用后缀构造对象键）
func (c *MinioClient) UploadSuffixFile(filePath string, suffix string) (*UploadResult, error) {
	// 检测文件类型
	mtype, err := mimetype.DetectFile(filePath)
	contentType := "application/octet-stream"
	if err == nil && mtype != nil {
		contentType = mtype.String()
	}

	objectKey := c.GetPath(c.properties.Prefix, suffix)
	return c.Upload(filePath, objectKey, "", contentType)
}

// FileDownload 下载文件到临时目录
func (c *MinioClient) FileDownload(path string) (string, error) {
	objectKey := c.RemoveBaseUrl(path)

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "oss_download_*")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %w", err)
	}
	tmpFilePath := tmpFile.Name()
	tmpFile.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 下载文件
	err = c.client.FGetObject(ctx, c.bucketName, objectKey, tmpFilePath, minio.GetObjectOptions{})
	if err != nil {
		os.Remove(tmpFilePath)
		return "", fmt.Errorf("下载文件失败: %w", err)
	}

	return tmpFilePath, nil
}

// GetObjectContent 获取对象内容输入流
func (c *MinioClient) GetObjectContent(path string) (io.ReadCloser, error) {
	objectKey := c.RemoveBaseUrl(path)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	object, err := c.client.GetObject(ctx, c.bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取对象失败: %w", err)
	}

	return object, nil
}

// CreatePresignedPutUrl 创建预签名上传URL
func (c *MinioClient) CreatePresignedPutUrl(objectKey string, expiredTime time.Duration, metadata map[string]string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reqParams := make(map[string]string)
	for k, v := range metadata {
		reqParams[k] = v
	}

	url, err := c.client.PresignedPutObject(ctx, c.bucketName, objectKey, expiredTime)
	if err != nil {
		return "", fmt.Errorf("生成预签名上传URL失败: %w", err)
	}

	return url.String(), nil
}

// GetEndpoint 获取端点URL
func (c *MinioClient) GetEndpoint() string {
	header := c.GetIsHttps()
	endpoint := c.properties.Endpoint

	// 如果端点包含协议，直接返回
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint
	}

	return header + endpoint
}

// GetDomain 获取域名URL
func (c *MinioClient) GetDomain() string {
	domain := c.properties.Domain
	header := c.GetIsHttps()

	if domain != "" {
		if strings.HasPrefix(domain, "http://") || strings.HasPrefix(domain, "https://") {
			return domain
		}
		return header + domain
	}

	return c.GetEndpoint()
}

// GetPath 生成文件路径（对象键）
func (c *MinioClient) GetPath(prefix string, suffix string) string {
	// 生成UUID
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")

	// 生成日期路径
	now := time.Now()
	datePath := fmt.Sprintf("%04d/%02d/%02d", now.Year(), now.Month(), now.Day())

	// 拼接路径
	var path string
	if prefix != "" {
		path = prefix + "/" + datePath + "/" + id
	} else {
		path = datePath + "/" + id
	}

	return path + suffix
}

// RemoveBaseUrl 移除基础URL
func (c *MinioClient) RemoveBaseUrl(path string) string {
	return c.removeBaseUrl(path)
}

// GetIsHttps 获取是否使用HTTPS
func (c *MinioClient) GetIsHttps() string {
	return c.getIsHttps()
}

// CheckPropertiesSame 检查配置是否相同
func (c *MinioClient) CheckPropertiesSame(properties *OssProperties) bool {
	if properties == nil {
		return false
	}

	return c.properties.ConfigKey == properties.ConfigKey &&
		c.properties.AccessKey == properties.AccessKey &&
		c.properties.SecretKey == properties.SecretKey &&
		c.properties.BucketName == properties.BucketName &&
		c.properties.Prefix == properties.Prefix &&
		c.properties.Endpoint == properties.Endpoint &&
		c.properties.Domain == properties.Domain &&
		c.properties.IsHttps == properties.IsHttps &&
		c.properties.Region == properties.Region &&
		c.properties.AccessPolicy == properties.AccessPolicy &&
		c.properties.TenantId == properties.TenantId
}
