package main

import (
	"flag"
	"fmt"
	"gozero-ruoyi-vue-plus/internal/config"
	"gozero-ruoyi-vue-plus/internal/handler"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/util"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/admin-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 初始化雪花算法（使用默认机器ID和数据中心ID：0, 0）
	// 如果需要多实例部署，可以从配置文件中读取不同的 workerID 和 datacenterID
	if err := util.InitDefaultSnowflake(0, 0); err != nil {
		logx.Errorf("初始化雪花算法失败: %v", err)
		panic(fmt.Sprintf("初始化雪花算法失败: %v", err))
	}
	logx.Infof("雪花算法初始化成功")

	//server := rest.MustNewServer(c.RestConf)
	// 跨域
	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(func(header http.Header) {
		header.Set("Access-Control-Allow-Headers", "*")
		header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		header.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		header.Set("Access-Control-Allow-Credentials", "true")
	}, nil, "*"))

	defer server.Stop()

	// 添加 CORS 中间件
	//server.Use(middleware.NewCorsMiddleware().Handle)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
