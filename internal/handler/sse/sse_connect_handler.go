// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sse

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero-ruoyi-vue-plus/internal/logic/sse"
	"gozero-ruoyi-vue-plus/internal/svc"
)

// 建立 SSE 连接
func SseConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sse.NewSseConnectLogic(r.Context(), svcCtx)
		if err := l.SseConnect(w, r); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
	}
}
