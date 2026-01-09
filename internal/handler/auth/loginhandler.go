// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"gozero-ruoyi-vue-plus/internal/logic/auth"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 登录方法，支持密码登录等授权类型
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq

		// 检查是否启用 API 解密
		if svcCtx.Config.ApiDecrypt.Enabled {
			// 检查请求头是否有加密标识
			encryptFlag := r.Header.Get(svcCtx.Config.ApiDecrypt.HeaderFlag)
			if encryptFlag != "" {
				// 步骤1: 使用 RSA 私钥解密请求头，得到 Base64 编码的 AES 密钥
				decryptAesBase64, err := util.RSADecrypt(svcCtx.Config.ApiDecrypt.PrivateKey, encryptFlag)
				if err != nil {
					logx.Errorf("RSA解密请求头失败: %v", err)
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}

				// 步骤2: Base64 解码，得到 AES 密钥
				aesKeyBytes, err := base64.StdEncoding.DecodeString(decryptAesBase64)
				if err != nil {
					logx.Errorf("Base64解码AES密钥失败: %v", err)
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
				aesKey := string(aesKeyBytes)

				// 步骤3: 读取请求体（Base64 编码的加密数据）
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					logx.Errorf("读取请求体失败: %v", err)
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
				defer r.Body.Close()

				// 步骤4: 使用 AES 密钥解密请求体（内部会自动去除前后双引号）
				// 去除前后双引号（如果存在）
				aesBody := strings.Trim(string(bodyBytes), "\"")
				decryptedBody, err := util.AESDecrypt(aesBody, aesKey)
				if err != nil {
					logx.Errorf("AES解密请求体失败: %v", err)
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}

				// 步骤5: 解析解密后的 JSON 数据
				if err := json.Unmarshal([]byte(decryptedBody), &req); err != nil {
					logx.Errorf("解析解密后的数据失败: %v", err)
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			} else {
				// 没有加密标识，正常解析
				if err := httpx.Parse(r, &req); err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			}
		} else {
			// 未启用加密，正常解析
			if err := httpx.Parse(r, &req); err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
		}

		l := auth.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
