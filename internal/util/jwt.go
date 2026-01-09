package util

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("ruoyi-vue-plus-secret-key") // TODO: 从配置文件读取

// Claims JWT 声明
type Claims struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
	TenantId string `json:"tenantId"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(userId int64, username, tenantId string, expireTime int64) (string, error) {
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
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseToken 解析 JWT token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
