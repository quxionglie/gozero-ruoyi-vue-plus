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

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Type         string // 验证码类型: MATH 数组计算, CHAR 字符验证
	Category     string // line 线段干扰 circle 圆圈干扰 shear 扭曲干扰
	NumberLength int    // 数字验证码位数
	CharLength   int    // 字符验证码长度
}

// GenerateCaptcha 根据配置生成验证码
func GenerateCaptcha(cfg CaptchaConfig) (code string, imgBase64 string, captchaId string, err error) {
	var driver base64Captcha.Driver

	// 根据类型创建不同的驱动器
	// 优化参数：增大图片尺寸，完全移除干扰，提高清晰度
	// height: 高度（增大到80，让字符更大更清晰）
	// width: 宽度（根据字符长度动态调整，每个字符约50像素，给字体更大空间）
	// noiseCount: 干扰点数量（设为0，完全移除干扰）
	// showLineOptions: 干扰线类型（设为0，完全移除干扰线）

	switch strings.ToUpper(cfg.Type) {
	case "MATH":
		// 数学计算验证码
		// NewDriverMath(height, width, noiseCount, showLineOptions, bgColor, fontsStorage, fonts)
		// 数学验证码宽度固定，因为格式固定（如 "3+5=?"）
		// 增大尺寸：高度80，宽度200，让字体更大
		driver = base64Captcha.NewDriverMath(80, 200, 0, 0, nil, nil, nil)
	case "CHAR":
		fallthrough
	default:
		// 字符验证码
		// NewDriverString(height, width, noiseCount, showLineOptions, length, source, bgColor, fontsStorage, fonts)
		// source: 字符集，排除容易混淆的字符（0,o,O,l,1,i,I）
		charSource := "23456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ"
		if cfg.CharLength <= 0 {
			cfg.CharLength = 4 // 默认长度
		}
		// 根据字符长度动态调整宽度，每个字符约50像素，加上边距，让字体更大
		width := cfg.CharLength*50 + 40
		if width < 200 {
			width = 200 // 最小宽度
		}
		// 完全移除干扰：noiseCount设为0，showLineOptions设为0
		driver = base64Captcha.NewDriverString(80, width, 0, 0, cfg.CharLength, charSource, nil, nil, nil)
	}

	captcha := base64Captcha.NewCaptcha(driver, store)

	// 生成验证码
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return "", "", "", err
	}

	// 去掉 base64 字符串的前缀（如 data:image/png;base64, 或 data:image/gif;base64,）
	imgBase64 = removeBase64Prefix(b64s)

	// 对于MATH类型，直接返回PNG base64
	if strings.ToUpper(cfg.Type) == "MATH" {
		return answer, imgBase64, id, nil
	}

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

	// 创建优化的调色板：白色背景 + 黑色数字
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// 创建简单的二色调色板
	palette := make(color.Palette, 2)
	palette[0] = white // 索引0：白色背景
	palette[1] = black // 索引1：黑色数字

	gifImg := image.NewPaletted(bounds, palette)

	// 先分析图片，确定主要背景色和文字色
	// 统计亮色和暗色像素的数量
	var lightPixels, darkPixels int
	var totalBrightness uint64
	var pixelCount int

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := pngImg.At(x, y)
			r, g, b, a := c.RGBA()

			// 跳过透明像素
			if a < 0x8000 {
				continue
			}

			// 将 RGBA 值从 16 位转换为 8 位
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 计算亮度（使用加权平均）
			brightness := float64(r8)*0.299 + float64(g8)*0.587 + float64(b8)*0.114
			totalBrightness += uint64(brightness)
			pixelCount++

			if brightness > 128 {
				lightPixels++
			} else {
				darkPixels++
			}
		}
	}

	// 判断图片的主要颜色：如果亮色像素多，说明是白色背景黑色文字；否则可能是黑色背景白色文字
	// 计算平均亮度
	var avgBrightness float64
	if pixelCount > 0 {
		avgBrightness = float64(totalBrightness) / float64(pixelCount)
	}

	// 如果平均亮度较高（>128），说明主要是亮色，可能是白色背景；否则可能是深色背景
	// 为了确保文字清晰，我们统一转换为：白色背景 + 黑色文字
	// 如果原图是深色背景白色文字，需要反转

	// 将 PNG 图片转换为 GIF，确保白色背景和黑色文字
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := pngImg.At(x, y)
			r, g, b, a := c.RGBA()

			// 如果透明度很低，设置为白色背景
			if a < 0x8000 {
				gifImg.Set(x, y, white)
				continue
			}

			// 将 RGBA 值从 16 位转换为 8 位
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 计算亮度（使用加权平均）
			brightness := float64(r8)*0.299 + float64(g8)*0.587 + float64(b8)*0.114

			// 如果平均亮度较低（<100），说明原图可能是深色背景白色文字，需要反转
			// 否则原图是浅色背景深色文字，直接使用
			if avgBrightness < 100 {
				// 反转：原图深色背景白色文字 -> 白色背景黑色文字
				if brightness > 100 {
					// 原图的亮色（文字）-> 黑色
					gifImg.Set(x, y, black)
				} else {
					// 原图的暗色（背景）-> 白色
					gifImg.Set(x, y, white)
				}
			} else {
				// 不反转：原图浅色背景深色文字 -> 白色背景黑色文字
				if brightness < 150 {
					// 暗色区域 -> 黑色文字
					gifImg.Set(x, y, black)
				} else {
					// 亮色区域 -> 白色背景
					gifImg.Set(x, y, white)
				}
			}
		}
	}

	// 编码为 GIF
	var buf bytes.Buffer
	err = gif.Encode(&buf, gifImg, &gif.Options{
		NumColors: 2, // 只使用黑白两色
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
