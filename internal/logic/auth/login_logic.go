// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登录方法，支持密码登录等授权类型
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 1. 验证客户端ID和授权类型
	client, err := l.svcCtx.SysClientModel.FindOneByClientId(l.ctx, req.ClientId)
	if err != nil {
		if err == sys.ErrNotFound {
			l.Errorf("客户端id: %s 不存在", req.ClientId)
			return &types.LoginResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "客户端认证类型错误",
				},
			}, nil
		}
		l.Errorf("查询客户端失败: %v", err)
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询客户端失败",
			},
		}, err
	}

	// 检查客户端状态
	if client.Status != "0" {
		l.Errorf("客户端id: %s 已被停用", req.ClientId)
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "客户端已被停用",
			},
		}, nil
	}

	// 检查授权类型是否支持
	if !client.GrantType.Valid || !strings.Contains(client.GrantType.String, req.GrantType) {
		l.Errorf("客户端id: %s 不支持授权类型: %s", req.ClientId, req.GrantType)
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "客户端认证类型错误",
			},
		}, nil
	}

	// 2. 校验租户（如果启用多租户）
	if l.svcCtx.Config.Tenant.Enable {
		if err := l.checkTenant(req.TenantId); err != nil {
			return &types.LoginResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  err.Error(),
				},
			}, nil
		}
	}

	// 3. 验证验证码（如果启用）
	if err := l.validateCaptcha(req.TenantId, req.Username, req.Code, req.Uuid); err != nil {
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  err.Error(),
			},
		}, nil
	}

	// 4. 查询用户（如果租户未启用，使用默认租户ID）
	tenantId := req.TenantId
	if !l.svcCtx.Config.Tenant.Enable {
		tenantId = "000000" // 默认租户
	}
	if tenantId == "" {
		tenantId = "000000" // 默认租户
	}
	user, err := l.svcCtx.SysUserModel.FindOneByUserName(l.ctx, req.Username, tenantId)
	if err != nil {
		if err == sys.ErrNotFound {
			l.Errorf("登录用户：%s 不存在", req.Username)
			return &types.LoginResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  fmt.Sprintf("用户 %s 不存在", req.Username),
				},
			}, nil
		}
		l.Errorf("查询用户失败: %v", err)
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询用户失败",
			},
		}, err
	}

	// 检查用户状态
	if user.Status != "0" {
		l.Errorf("登录用户：%s 已被停用", req.Username)
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  fmt.Sprintf("用户 %s 已被停用", req.Username),
			},
		}, nil
	}

	// 5. 验证密码（BCrypt）
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		l.Errorf("用户 %s 密码错误", req.Username)
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "用户名或密码错误",
			},
		}, nil
	}

	// 6. 生成JWT token
	timeout := client.Timeout
	if timeout == 0 {
		timeout = 1800 // 默认30分钟（秒）
	}
	activeTimeout := client.ActiveTimeout
	if activeTimeout == 0 {
		activeTimeout = 1800 // 默认30分钟（秒）
	}

	accessToken, err := util.GenerateToken(l.svcCtx.Config.Auth.AccessSecret, user.UserId, user.UserName, user.TenantId, timeout)
	if err != nil {
		l.Errorf("生成token失败: %v", err)
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "生成token失败",
			},
		}, err
	}

	refreshToken, err := util.GenerateToken(l.svcCtx.Config.Auth.AccessSecret, user.UserId, user.UserName, user.TenantId, activeTimeout)
	if err != nil {
		l.Errorf("生成刷新token失败: %v", err)
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "生成刷新token失败",
			},
		}, err
	}

	// 7. 存储在线 token 缓存到 Redis
	l.saveOnlineToken(accessToken, user, req.ClientId, timeout)

	// 8. 延迟 5 秒后发送 SSE 欢迎消息
	go func() {
		time.Sleep(5 * time.Second)
		l.sendWelcomeMessage(user.UserId)
	}()

	// 9. 返回登录信息
	resp = &types.LoginResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "登录成功",
		},
		Data: types.LoginVo{
			AccessToken:     accessToken,
			RefreshToken:    refreshToken,
			ExpireIn:        timeout,
			RefreshExpireIn: activeTimeout,
			ClientId:        req.ClientId,
		},
	}

	return resp, nil
}

// saveOnlineToken 存储在线 token 缓存到 Redis
func (l *LoginLogic) saveOnlineToken(token string, user *sys.SysUser, clientId string, expireSeconds int64) {
	// 从 context 中获取请求信息
	clientIP := ""
	if ipValue := l.ctx.Value("clientIP"); ipValue != nil {
		if ip, ok := ipValue.(string); ok {
			clientIP = ip
		}
	}

	userAgent := ""
	if uaValue := l.ctx.Value("userAgent"); uaValue != nil {
		if ua, ok := uaValue.(string); ok {
			userAgent = ua
		}
	}

	// 解析 User-Agent 获取浏览器和操作系统信息（简化版）
	browser := parseBrowser(userAgent)
	os := parseOS(userAgent)

	// 获取部门名称
	deptName := ""
	if user.DeptId.Valid {
		dept, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, user.DeptId.Int64)
		if err == nil {
			deptName = dept.DeptName
		}
	}

	// 构建在线用户信息
	onlineInfo := map[string]interface{}{
		"tokenId":       token,
		"userName":      user.UserName,
		"tenantId":      user.TenantId,
		"clientKey":     clientId,
		"deviceType":    "web", // 默认设备类型
		"ipaddr":        clientIP,
		"loginLocation": "", // 登录地址（需要 IP 地址库，暂时留空）
		"browser":       browser,
		"os":            os,
		"loginTime":     time.Now().UnixMilli(),
		"deptName":      deptName,
	}

	// 序列化为 JSON
	onlineInfoJSON, err := json.Marshal(onlineInfo)
	if err != nil {
		l.Errorf("序列化在线用户信息失败: %v", err)
		return
	}

	// 存储到 Redis，key: online_tokens:{token}
	onlineTokenKey := "online_tokens:" + token
	if expireSeconds > 0 {
		// 设置过期时间（秒）
		err = l.svcCtx.RedisConn.SetexCtx(l.ctx, onlineTokenKey, string(onlineInfoJSON), int(expireSeconds))
	} else {
		// 永不过期
		err = l.svcCtx.RedisConn.SetCtx(l.ctx, onlineTokenKey, string(onlineInfoJSON))
	}

	if err != nil {
		l.Errorf("存储在线 token 缓存失败: %v", err)
	} else {
		l.Infof("已存储在线 token 缓存: %s", onlineTokenKey)
	}
}

// parseBrowser 从 User-Agent 解析浏览器类型（简化版）
func parseBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)
	if strings.Contains(ua, "chrome") && !strings.Contains(ua, "edg") {
		return "Chrome"
	}
	if strings.Contains(ua, "firefox") {
		return "Firefox"
	}
	if strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") {
		return "Safari"
	}
	if strings.Contains(ua, "edg") {
		return "Edge"
	}
	if strings.Contains(ua, "opera") || strings.Contains(ua, "opr") {
		return "Opera"
	}
	return "Unknown"
}

// parseOS 从 User-Agent 解析操作系统类型（简化版）
func parseOS(userAgent string) string {
	ua := strings.ToLower(userAgent)
	if strings.Contains(ua, "windows") {
		return "Windows"
	}
	if strings.Contains(ua, "mac") {
		return "macOS"
	}
	if strings.Contains(ua, "linux") {
		return "Linux"
	}
	if strings.Contains(ua, "android") {
		return "Android"
	}
	if strings.Contains(ua, "ios") || strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		return "iOS"
	}
	return "Unknown"
}

// checkTenant 校验租户
func (l *LoginLogic) checkTenant(tenantId string) error {
	// 如果租户ID为空，使用默认租户
	if tenantId == "" {
		tenantId = "000000"
	}

	// 默认租户直接通过
	if tenantId == "000000" {
		return nil
	}

	// 查询租户
	tenant, err := l.svcCtx.SysTenantModel.FindOneByTenantId(l.ctx, tenantId)
	if err != nil {
		if err == sys.ErrNotFound {
			l.Errorf("登录租户：%s 不存在", tenantId)
			return fmt.Errorf("租户不存在")
		}
		return err
	}

	// 检查租户状态
	if tenant.Status != "0" {
		l.Errorf("登录租户：%s 已被停用", tenantId)
		return fmt.Errorf("租户已被停用")
	}

	// 检查租户是否过期
	if tenant.ExpireTime.Valid && tenant.ExpireTime.Time.Before(time.Now()) {
		l.Errorf("登录租户：%s 已超过有效期", tenantId)
		return fmt.Errorf("租户已超过有效期")
	}

	return nil
}

// sendWelcomeMessage 发送登录欢迎消息
func (l *LoginLogic) sendWelcomeMessage(userId int64) {
	// 获取当前时间并生成问候语
	now := time.Now()
	hour := now.Hour()
	var timeGreeting string
	switch {
	case hour <= 6:
		timeGreeting = "凌晨"
	case hour <= 9:
		timeGreeting = "早上"
	case hour <= 12:
		timeGreeting = "上午"
	case hour <= 14:
		timeGreeting = "中午"
	case hour <= 18:
		timeGreeting = "下午"
	case hour <= 22:
		timeGreeting = "晚上"
	default:
		timeGreeting = "深夜"
	}

	message := fmt.Sprintf("%s好，欢迎登录 RuoYi-Vue-Plus 后台管理系统", timeGreeting)

	// 获取 SSE 管理器并发送消息
	sseManager := util.GetSseEmitterManager()
	sseManager.SendMessage(userId, message)
}

// validateCaptcha 验证验证码
func (l *LoginLogic) validateCaptcha(tenantId, username, code, uuid string) error {
	// 检查验证码是否启用（从配置文件读取）
	captchaEnabled := l.svcCtx.Config.Captcha.Enable
	if !captchaEnabled {
		return nil // 验证码未启用，直接通过
	}

	if uuid == "" || code == "" {
		return fmt.Errorf("验证码不能为空")
	}

	// 从Redis获取验证码
	verifyKey := fmt.Sprintf("captcha_code:%s", uuid)
	captcha, err := l.svcCtx.RedisConn.GetCtx(l.ctx, verifyKey)
	if err != nil {
		l.Errorf("获取验证码失败: %v", err)
		return fmt.Errorf("验证码已过期")
	}

	// 删除验证码（一次性使用）
	_, _ = l.svcCtx.RedisConn.DelCtx(l.ctx, verifyKey)

	// 验证验证码（不区分大小写）
	if !strings.EqualFold(code, captcha) {
		l.Errorf("用户 %s 验证码错误: 输入=%s, 正确=%s", username, code, captcha)
		return fmt.Errorf("验证码错误")
	}

	return nil
}
