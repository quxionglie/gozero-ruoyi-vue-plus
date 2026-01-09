package svc

import (
	"gozero-ruoyi-vue-plus/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	RedisConn *redis.Redis
	DB        sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化MySQL连接
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	// 初始化Redis连接
	rds := redis.MustNewRedis(c.Redis)

	return &ServiceContext{
		Config:    c,
		RedisConn: rds,
		DB:        conn,
	}
}
