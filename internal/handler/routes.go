package handler

import (
	"net/http"

	"gozero-ruoyi-vue-plus/internal/svc"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/ping",
				Handler: PingHandler(serverCtx),
			},
		},
	)
}

func PingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"message": "pong",
			"status":  "ok",
		})
	}
}
