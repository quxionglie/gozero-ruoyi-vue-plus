package util

import (
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

// GenerateCaptcha 生成4位数字验证码
func GenerateCaptcha() (code string, imgBase64 string, captchaId string, err error) {
	// 创建4位数字验证码驱动器
	driver := base64Captcha.NewDriverDigit(40, 120, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, store)

	// 生成验证码
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return "", "", "", err
	}

	return answer, b64s, id, nil
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id string, code string) bool {
	return store.Verify(id, code, true)
}
