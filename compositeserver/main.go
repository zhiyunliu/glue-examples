package main

import (
	"github.com/zhiyunliu/glue"
	_ "github.com/zhiyunliu/glue/contrib/cache/redis"
	_ "github.com/zhiyunliu/glue/contrib/config/nacos"
	_ "github.com/zhiyunliu/glue/contrib/metrics/prometheus"
	_ "github.com/zhiyunliu/glue/contrib/registry/nacos"
	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/log"
	_ "github.com/zhiyunliu/queue-rabbitmq"
	_ "github.com/zhiyunliu/queue-redis"
	_ "github.com/zhiyunliu/xdb-gorm"
	_ "github.com/zhiyunliu/xdb-mssql"
	_ "github.com/zhiyunliu/xdb-mysql"
	_ "github.com/zhiyunliu/xdb-sqlite"

	_ "github.com/zhiyunliu/glue/contrib/dlocker/redis"
	_ "github.com/zhiyunliu/glue/contrib/xhttp/http"
	_ "github.com/zhiyunliu/glue/contrib/xrpc/grpc"
)

var (
	opts = []glue.Option{glue.LogParams(log.WithConcurrency(1))}
)

func main() {
	global.AppName = "compositeserver"
	app := glue.NewApp(opts...)
	app.Start()
}
