// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"net/http"
	"strings"

	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取 token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "未提供认证令牌", http.StatusUnauthorized)
			return
		}

		// 支持 Bearer token 格式
		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 解析 JWT token
		claims, err := util.ParseToken(token)
		if err != nil {
			logx.Errorf("解析 JWT token 失败: %v", err)
			http.Error(w, "无效的认证令牌", http.StatusUnauthorized)
			return
		}

		// 将用户信息存储到请求上下文中，供后续 handler 使用
		r = r.WithContext(util.WithClaims(r.Context(), claims))

		// 继续处理请求
		next(w, r)
	}
}
