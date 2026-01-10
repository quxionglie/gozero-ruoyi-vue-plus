// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package monitor

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero-ruoyi-vue-plus/internal/logic/monitor"
	"gozero-ruoyi-vue-plus/internal/svc"
)

// 清理操作日志记录
func OperLogCleanHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := monitor.NewOperLogCleanLogic(r.Context(), svcCtx)
		resp, err := l.OperLogClean()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
