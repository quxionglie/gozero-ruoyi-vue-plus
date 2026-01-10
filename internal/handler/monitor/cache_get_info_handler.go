// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package monitor

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero-ruoyi-vue-plus/internal/logic/monitor"
	"gozero-ruoyi-vue-plus/internal/svc"
)

// 获取缓存监控列表
func CacheGetInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := monitor.NewCacheGetInfoLogic(r.Context(), svcCtx)
		resp, err := l.CacheGetInfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
