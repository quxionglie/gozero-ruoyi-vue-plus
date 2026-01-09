package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Mysql MysqlConfig
	Redis redis.RedisConf
}

type MysqlConfig struct {
	DataSource string
}
