package util

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/gif"
	"strings"

	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

// GenerateCaptcha 生成4位数字验证码（GIF格式，白色背景，清晰数字）
func GenerateCaptcha() (code string, imgBase64 string, captchaId string, err error) {
	// 创建4位数字验证码驱动器
	// 参数：高度, 宽度, 数字长度, 最大倾斜度, 点数量
	// 减少干扰点数量（从80降到20），降低倾斜度（从0.7降到0.3），使数字更清晰
	driver := base64Captcha.NewDriverDigit(50, 140, 4, 0.3, 20)
	captcha := base64Captcha.NewCaptcha(driver, store)

	// 生成验证码
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return "", "", "", err
	}

	// 去掉 base64 字符串的前缀（如 data:image/png;base64, 或 data:image/gif;base64,）
	imgBase64 = removeBase64Prefix(b64s)

	// 将 PNG 转换为 GIF，并优化为白色背景和清晰数字
	gifBase64, err := convertPNGToGIFWithWhiteBackground(imgBase64)
	if err != nil {
		// 如果转换失败，返回原 PNG base64
		return answer, imgBase64, id, nil
	}

	return answer, gifBase64, id, nil
}

// removeBase64Prefix 去掉 base64 字符串的前缀
func removeBase64Prefix(base64Str string) string {
	// 去掉常见的前缀
	prefixes := []string{
		"data:image/gif;base64,",
		"data:image/png;base64,",
		"data:image/jpeg;base64,",
		"data:image/jpg;base64,",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(base64Str, prefix) {
			return strings.TrimPrefix(base64Str, prefix)
		}
	}

	return base64Str
}

// convertPNGToGIFWithWhiteBackground 将 PNG base64 转换为 GIF base64，并优化为白色背景和清晰数字
func convertPNGToGIFWithWhiteBackground(pngBase64 string) (string, error) {
	// 解码 PNG base64
	pngData, err := base64.StdEncoding.DecodeString(pngBase64)
	if err != nil {
		return "", err
	}

	// 解码 PNG 图片
	pngImg, _, err := image.Decode(bytes.NewReader(pngData))
	if err != nil {
		return "", err
	}

	// 创建 GIF 图片
	bounds := pngImg.Bounds()

	// 创建优化的调色板：白色背景 + 黑色数字 + 灰度过渡
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// 创建调色板，包含白色、黑色和灰度
	palette := make(color.Palette, 256)
	palette[0] = white // 索引0：白色背景
	palette[1] = black // 索引1：黑色数字

	// 填充灰度色阶
	for i := 2; i < 256; i++ {
		gray := uint8(i)
		palette[i] = color.RGBA{R: gray, G: gray, B: gray, A: 255}
	}

	gifImg := image.NewPaletted(bounds, palette)

	// 将 PNG 图片转换为 GIF，增强对比度
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := pngImg.At(x, y)
			r, g, b, a := c.RGBA()

			// 将 RGBA 值从 16 位转换为 8 位
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 计算亮度
			brightness := (uint32(r8) + uint32(g8) + uint32(b8)) / 3

			// 如果透明度很低，设置为白色背景
			if a < 0x8000 {
				gifImg.Set(x, y, white)
			} else if brightness > 200 {
				// 亮色区域（接近白色）-> 白色背景
				gifImg.Set(x, y, white)
			} else if brightness < 100 {
				// 暗色区域（接近黑色）-> 黑色数字，增强对比度
				gifImg.Set(x, y, black)
			} else {
				// 中间灰度，转换为更接近黑色或白色
				if brightness < 128 {
					gifImg.Set(x, y, black)
				} else {
					gifImg.Set(x, y, white)
				}
			}
		}
	}

	// 编码为 GIF
	var buf bytes.Buffer
	err = gif.Encode(&buf, gifImg, &gif.Options{
		NumColors: 256,
	})
	if err != nil {
		return "", err
	}

	// 转换为 base64
	gifBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return gifBase64, nil
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id string, code string) bool {
	return store.Verify(id, code, true)
}
