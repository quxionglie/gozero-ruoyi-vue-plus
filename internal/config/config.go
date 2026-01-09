package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Mysql   MysqlConfig
	Redis   redis.RedisConf
	Captcha CaptchaConfig
}

type MysqlConfig struct {
	DataSource string
}

type CaptchaConfig struct {
	Enable       bool   `json:",default=true"`   // 是否启用验证码校验
	Type         string `json:",default=CHAR"`   // 验证码类型: MATH 数组计算, CHAR 字符验证
	Category     string `json:",default=CIRCLE"` // line 线段干扰 circle 圆圈干扰 shear 扭曲干扰
	NumberLength int    `json:",default=1"`      // 数字验证码位数
	CharLength   int    `json:",default=4"`      // 字符验证码长度
}
