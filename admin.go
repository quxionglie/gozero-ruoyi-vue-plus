package main

import (
	"flag"
	"fmt"
	"gozero-ruoyi-vue-plus/internal/config"
	"gozero-ruoyi-vue-plus/internal/handler"
	"gozero-ruoyi-vue-plus/internal/middleware"
	"gozero-ruoyi-vue-plus/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/admin-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 添加 CORS 中间件
	server.Use(middleware.NewCorsMiddleware().Handle)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
