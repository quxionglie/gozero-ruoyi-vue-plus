// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero-ruoyi-vue-plus/internal/logic/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
)

// 头像上传
func UserProfileAvatarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 文件上传使用 multipart/form-data
		// 解析文件
		err := r.ParseMultipartForm(32 << 20) // 32MB max memory
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		file, fileHeader, err := r.FormFile("avatarfile")
		if err != nil {
			// 尝试其他可能的字段名
			file, fileHeader, err = r.FormFile("file")
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
		}
		defer file.Close()

		l := sys.NewUserProfileAvatarLogic(r.Context(), svcCtx)
		resp, err := l.UserProfileAvatar(file, fileHeader)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
