package main

import (
	gel "github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue-examples/cronserver/demos"
	_ "github.com/zhiyunliu/glue/contrib/config/nacos"
	_ "github.com/zhiyunliu/glue/contrib/dlocker/redis"
	_ "github.com/zhiyunliu/glue/contrib/queue/redis"
	_ "github.com/zhiyunliu/glue/contrib/queue/streamredis"
	_ "github.com/zhiyunliu/glue/contrib/registry/nacos"
	"github.com/zhiyunliu/glue/server/cron"
)

func main() {
	cronSrv := cron.New("cronserver")

	cronSrv.Handle("/demo", &demos.Fulldemo{})

	app := gel.NewApp(gel.Server(cronSrv))

	app.Start()
}
