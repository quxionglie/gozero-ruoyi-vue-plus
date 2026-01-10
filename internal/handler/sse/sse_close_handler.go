// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sse

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero-ruoyi-vue-plus/internal/logic/sse"
	"gozero-ruoyi-vue-plus/internal/svc"
)

// 关闭 SSE 连接
func SseCloseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sse.NewSseCloseLogic(r.Context(), svcCtx)
		resp, err := l.SseClose(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
