package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Mysql      MysqlConfig
	Redis      redis.RedisConf
	Captcha    CaptchaConfig
	Tenant     TenantConfig
	ApiDecrypt ApiDecryptConfig
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

type TenantConfig struct {
	Enable bool `json:",default=false"` // 是否开启多租户
}

type ApiDecryptConfig struct {
	Enabled    bool   `json:",default=false"`       // 是否开启全局接口加密
	HeaderFlag string `json:",default=encrypt-key"` // AES 加密头标识
	PublicKey  string // 响应加密公钥
	PrivateKey string // 请求解密私钥
}
