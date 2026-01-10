// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero-ruoyi-vue-plus/internal/logic/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
)

// 查询字典数据列表
func DictDataListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sys.NewDictDataListLogic(r.Context(), svcCtx)
		resp, err := l.DictDataList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
