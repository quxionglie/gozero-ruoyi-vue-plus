// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 退出登录
func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout() (resp *types.BaseResp, err error) {
	// 1. 从 JWT token 中获取用户信息（可选，如果 token 无效也不影响退出）
	userId, err := util.GetUserIdFromContext(l.ctx)
	username, _ := util.GetUsernameFromContext(l.ctx)
	tenantId, _ := util.GetTenantIdFromContext(l.ctx)

	// 2. 从请求头获取 token，清除 Redis 中的在线 token 缓存
	// 注意：这里需要从 handler 传递 token，或者从 context 中获取
	// 由于 go-zero 的 JWT 中间件会将 token 信息存储到 context，我们需要从请求头获取原始 token
	token := l.getTokenFromContext()
	if token != "" {
		// 清除 Redis 中的在线 token 缓存
		onlineTokenKey := "online_tokens:" + token
		_, delErr := l.svcCtx.RedisConn.DelCtx(l.ctx, onlineTokenKey)
		if delErr != nil {
			l.Errorf("清除在线 token 缓存失败: %v", delErr)
			// 清除失败不影响退出流程，继续执行
		} else {
			l.Infof("已清除在线 token 缓存: %s", onlineTokenKey)
		}
	}

	// 3. 记录登出日志
	if err == nil {
		l.Infof("用户退出登录：userId=%d, username=%s, tenantId=%s", userId, username, tenantId)
	} else {
		l.Infof("退出登录：无法获取用户信息，可能 token 已过期")
	}

	// 4. 返回成功响应
	// 注意：JWT 是无状态的，服务端无法主动使 token 失效
	// 客户端需要删除本地存储的 token
	return &types.BaseResp{
		Code: 200,
		Msg:  "退出成功",
	}, nil
}

// getTokenFromContext 从 context 中获取 token
func (l *LogoutLogic) getTokenFromContext() string {
	// 从 context 中获取 token（由 handler 传递）
	tokenValue := l.ctx.Value("token")
	if tokenValue != nil {
		if token, ok := tokenValue.(string); ok {
			return token
		}
	}
	return ""
}
