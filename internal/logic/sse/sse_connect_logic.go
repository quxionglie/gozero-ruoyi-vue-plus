// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sse

import (
	"context"
	"fmt"
	"net/http"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SseConnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 建立 SSE 连接
func NewSseConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SseConnectLogic {
	return &SseConnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SseConnectLogic) SseConnect(w http.ResponseWriter, r *http.Request) error {
	// 获取 Flusher
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("ResponseWriter 不支持 Flush")
	}

	// 从 form 参数中获取 Authorization
	authorization := r.FormValue("Authorization")
	if authorization == "" {
		// 如果 form 参数中没有，尝试从 header 获取
		authorization = r.Header.Get("Authorization")
	}
	if authorization == "" {
		l.Errorf("未找到有效的 token")
		return fmt.Errorf("未找到有效的 token")
	}

	// 提取 token（去掉 "Bearer " 前缀）
	token := authorization
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// 验证 JWT token 是否有效
	claims, err := util.ParseToken(token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		l.Errorf("JWT token 验证失败: %v", err)
		return fmt.Errorf("JWT token 验证失败: %v", err)
	}

	// 从 Claims 中获取用户信息
	userId := claims.UserId

	// 获取 SSE 管理器并建立连接
	sseManager := util.GetSseEmitterManager()
	conn := sseManager.Connect(userId, token, w, flusher)

	// 保持连接打开，直到客户端断开
	<-conn.Done()

	return nil
}
