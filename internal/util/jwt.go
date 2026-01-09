package util

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const claimsKey contextKey = "claims"

// jwtSecret 已移除，改为从配置中读取

// Claims JWT 声明
type Claims struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
	TenantId string `json:"tenantId"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token（使用 go-zero 标准格式）
// secret: JWT 密钥
// userId: 用户ID
// username: 用户名
// tenantId: 租户ID
// expireTime: 过期时间（秒）
func GenerateToken(secret string, userId int64, username, tenantId string, expireTime int64) (string, error) {
	nowTime := time.Now()
	expireDuration := time.Duration(expireTime) * time.Second
	if expireTime == 0 {
		// 默认 30 分钟
		expireDuration = 30 * time.Minute
	}

	claims := Claims{
		UserId:   userId,
		Username: username,
		TenantId: tenantId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(nowTime.Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(secret))
	return token, err
}

// ParseToken 解析 JWT token（如果需要手动解析，传入 secret）
func ParseToken(token, secret string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// WithClaims 将 Claims 存储到 context 中
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// GetClaims 从 context 中获取 Claims
func GetClaims(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	return claims, ok
}

// GetUserIdFromContext 从 context 中获取 userId（go-zero JWT 会将 claims 字段存储到 context 中）
// 支持多种类型：int64、json.Number、float64、string
func GetUserIdFromContext(ctx context.Context) (int64, error) {
	userIdValue := ctx.Value("userId")
	if userIdValue == nil {
		return 0, fmt.Errorf("userId not found in context")
	}

	return convertToInt64(userIdValue)
}

// GetUsernameFromContext 从 context 中获取 username
func GetUsernameFromContext(ctx context.Context) (string, error) {
	usernameValue := ctx.Value("username")
	if usernameValue == nil {
		return "", fmt.Errorf("username not found in context")
	}

	username, ok := usernameValue.(string)
	if !ok {
		return "", fmt.Errorf("username is not a string type")
	}

	return username, nil
}

// GetTenantIdFromContext 从 context 中获取 tenantId
func GetTenantIdFromContext(ctx context.Context) (string, error) {
	tenantIdValue := ctx.Value("tenantId")
	if tenantIdValue == nil {
		return "", fmt.Errorf("tenantId not found in context")
	}

	tenantId, ok := tenantIdValue.(string)
	if !ok {
		return "", fmt.Errorf("tenantId is not a string type")
	}

	return tenantId, nil
}

// GetJwtInfo 从 context 中获取完整的 JWT 信息
// 返回 userId、username、tenantId
func GetJwtInfo(ctx context.Context) (userId int64, username, tenantId string, err error) {
	userId, err = GetUserIdFromContext(ctx)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to get userId: %w", err)
	}

	username, err = GetUsernameFromContext(ctx)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to get username: %w", err)
	}

	tenantId, err = GetTenantIdFromContext(ctx)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to get tenantId: %w", err)
	}

	return userId, username, tenantId, nil
}

// convertToInt64 将各种数字类型转换为 int64
// 支持：int64、int、int32、json.Number、float64、string
func convertToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case json.Number:
		// json.Number 类型，需要转换为 int64
		result, err := v.Int64()
		if err != nil {
			return 0, fmt.Errorf("failed to convert json.Number to int64: %w", err)
		}
		return result, nil
	case float64:
		// JSON 数字可能被解析为 float64
		return int64(v), nil
	case float32:
		return int64(v), nil
	case string:
		// 字符串类型，尝试解析
		result, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to convert string to int64: %w", err)
		}
		return result, nil
	default:
		return 0, fmt.Errorf("unsupported type for userId: %T", value)
	}
}
