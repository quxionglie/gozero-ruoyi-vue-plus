// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"net/http"
	"strings"

	"gozero-ruoyi-vue-plus/internal/logic/auth"
	"gozero-ruoyi-vue-plus/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 退出登录
func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取 token
		authHeader := r.Header.Get("Authorization")
		token := ""
		if authHeader != "" {
			// 支持 Bearer token 格式
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				token = authHeader
			}
		}

		// 将 token 存储到 context 中
		ctx := r.Context()
		if token != "" {
			ctx = context.WithValue(ctx, "token", token)
		}

		l := auth.NewLogoutLogic(ctx, svcCtx)
		resp, err := l.Logout()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
