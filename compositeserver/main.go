package main

import (
	"github.com/zhiyunliu/glue"
	_ "github.com/zhiyunliu/glue/contrib/cache/redis"
	_ "github.com/zhiyunliu/glue/contrib/config/consul"
	_ "github.com/zhiyunliu/glue/contrib/config/nacos"
	_ "github.com/zhiyunliu/glue/contrib/queue/redis"
	_ "github.com/zhiyunliu/glue/contrib/queue/streamredis"
	_ "github.com/zhiyunliu/glue/contrib/registry/nacos"
	_ "github.com/zhiyunliu/glue/contrib/xdb/mysql"
	_ "github.com/zhiyunliu/glue/contrib/xdb/sqlite"
	_ "github.com/zhiyunliu/glue/contrib/xdb/sqlserver"

	_ "github.com/zhiyunliu/glue/contrib/dlocker/redis"

	_ "github.com/zhiyunliu/glue/contrib/xhttp/http"
	_ "github.com/zhiyunliu/glue/contrib/xrpc/grpc"
)

var (
	opts = []glue.Option{glue.LogConcurrency(1)}
)

func main() {

	app := glue.NewApp(opts...)
	app.Start()
}
