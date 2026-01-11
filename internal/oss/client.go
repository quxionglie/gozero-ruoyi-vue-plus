package oss

import (
	"io"
	"time"
)

// AccessPolicyType 访问策略类型
type AccessPolicyType int

const (
	AccessPolicyPrivate AccessPolicyType = 0 // private
	AccessPolicyPublic  AccessPolicyType = 1 // public
	AccessPolicyCustom  AccessPolicyType = 2 // custom
)

// UploadResult 上传结果
type UploadResult struct {
	URL      string // 文件URL
	Filename string // 文件名（对象键）
	ETag     string // ETag
}

// OssProperties OSS配置属性
type OssProperties struct {
	ConfigKey    string // 配置键
	AccessKey    string // accessKey
	SecretKey    string // 秘钥
	BucketName   string // 桶名称
	Prefix       string // 前缀
	Endpoint     string // 访问站点
	Domain       string // 自定义域名
	IsHttps      string // 是否https（Y=是,N=否）
	Region       string // 域
	AccessPolicy string // 桶权限类型(0=private 1=public 2=custom)
	TenantId     string // 租户编号
}

// OssClient OSS客户端接口
type OssClient interface {
	// Upload 上传文件（指定文件路径和对象键）
	// filePath: 本地文件路径
	// key: 对象键
	// md5Digest: MD5哈希值（可选，传空字符串表示不验证）
	// contentType: 文件内容类型
	Upload(filePath string, key string, md5Digest string, contentType string) (*UploadResult, error)

	// UploadStream 上传输入流
	// inputStream: 输入流
	// key: 对象键
	// length: 输入流长度
	// contentType: 文件内容类型
	UploadStream(inputStream io.Reader, key string, length int64, contentType string) (*UploadResult, error)

	// UploadSuffix 上传文件（使用后缀构造对象键）
	// data: 文件数据
	// suffix: 文件后缀（如 .jpg）
	// contentType: 文件内容类型（如 image/jpeg）
	UploadSuffix(data []byte, suffix string, contentType string) (*UploadResult, error)

	// UploadSuffixStream 上传输入流（使用后缀构造对象键）
	// inputStream: 输入流
	// suffix: 文件后缀
	// length: 输入流长度
	// contentType: 文件内容类型
	UploadSuffixStream(inputStream io.Reader, suffix string, length int64, contentType string) (*UploadResult, error)

	// UploadSuffixFile 上传文件（使用后缀构造对象键）
	// filePath: 本地文件路径
	// suffix: 文件后缀
	UploadSuffixFile(filePath string, suffix string) (*UploadResult, error)

	// FileDownload 下载文件到临时目录
	// path: 对象键或完整路径
	// 返回临时文件路径
	FileDownload(path string) (string, error)

	// Download 下载文件到输出流
	// key: 对象键（文件名）
	// out: 输出流
	Download(key string, out io.Writer) error

	// GetObjectContent 获取对象内容输入流
	// path: 对象键或完整路径
	// 返回输入流（调用者需要关闭）
	GetObjectContent(path string) (io.ReadCloser, error)

	// Delete 删除文件
	// path: 文件路径或对象键
	Delete(path string) error

	// CreatePresignedGetUrl 创建预签名下载URL
	// objectKey: 对象键
	// expiredTime: 过期时间
	CreatePresignedGetUrl(objectKey string, expiredTime time.Duration) (string, error)

	// CreatePresignedPutUrl 创建预签名上传URL
	// objectKey: 对象键
	// expiredTime: 过期时间
	// metadata: 元数据（可选）
	CreatePresignedPutUrl(objectKey string, expiredTime time.Duration, metadata map[string]string) (string, error)

	// GetConfigKey 获取配置键
	GetConfigKey() string

	// GetAccessPolicy 获取访问策略类型
	GetAccessPolicy() AccessPolicyType

	// GetUrl 获取文件URL（基础URL）
	GetUrl() string

	// GetEndpoint 获取端点URL
	GetEndpoint() string

	// GetDomain 获取域名URL
	GetDomain() string

	// GetPath 生成文件路径（对象键）
	// prefix: 前缀
	// suffix: 后缀
	GetPath(prefix string, suffix string) string

	// RemoveBaseUrl 移除基础URL
	// path: 完整路径
	RemoveBaseUrl(path string) string

	// GetIsHttps 获取是否使用HTTPS
	GetIsHttps() string

	// CheckPropertiesSame 检查配置是否相同
	CheckPropertiesSame(properties *OssProperties) bool
}
