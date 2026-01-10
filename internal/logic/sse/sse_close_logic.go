// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sse

import (
	"context"
	"net/http"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SseCloseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 关闭 SSE 连接
func NewSseCloseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SseCloseLogic {
	return &SseCloseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SseCloseLogic) SseClose(r *http.Request) (resp *types.BaseResp, err error) {
	// 从 form 参数中获取 Authorization
	authorization := r.FormValue("Authorization")
	if authorization == "" {
		// 如果 form 参数中没有，尝试从 header 获取
		authorization = r.Header.Get("Authorization")
	}
	if authorization == "" {
		l.Errorf("未找到有效的 token")
		return &types.BaseResp{
			Code: 500,
			Msg:  "未找到有效的 token",
		}, nil
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
		return &types.BaseResp{
			Code: 401,
			Msg:  "JWT token 验证失败",
		}, nil
	}

	// 从 Claims 中获取用户信息
	userId := claims.UserId

	// 获取 SSE 管理器并断开连接
	sseManager := util.GetSseEmitterManager()
	sseManager.Disconnect(userId, token)

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
