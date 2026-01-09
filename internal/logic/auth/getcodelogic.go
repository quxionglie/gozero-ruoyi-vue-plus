// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"fmt"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 生成验证码，返回4位数字验证码图片和唯一标识
func NewGetCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCodeLogic {
	return &GetCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCodeLogic) GetCode() (resp *types.CaptchaResp, err error) {
	// 检查验证码是否启用（从配置文件读取）
	captchaEnabled := l.svcCtx.Config.Captcha.Enable

	resp = &types.CaptchaResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: types.CaptchaVo{
			CaptchaEnabled: captchaEnabled,
		},
	}

	// 如果验证码未启用，直接返回
	if !captchaEnabled {
		return resp, nil
	}

	// 生成4位数字验证码
	code, imgBase64, _, err := util.GenerateCaptcha()
	if err != nil {
		l.Errorf("生成验证码失败: %v", err)
		return &types.CaptchaResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "生成验证码失败",
			},
		}, err
	}

	// 生成唯一标识（使用 UUID）
	uid := uuid.New().String()

	// 将验证码存储到 Redis，过期时间 2 分钟
	verifyKey := fmt.Sprintf("captcha_code:%s", uid)
	err = l.svcCtx.RedisConn.SetexCtx(l.ctx, verifyKey, code, 120) // 120秒 = 2分钟
	if err != nil {
		l.Errorf("存储验证码到Redis失败: %v", err)
		return &types.CaptchaResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "存储验证码失败",
			},
		}, err
	}

	resp.Data.Uuid = uid
	resp.Data.Img = imgBase64

	return resp, nil
}
