package main

import (
	"time"

	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue-examples/mqcserver/demos"
	"github.com/zhiyunliu/glue/context"
	_ "github.com/zhiyunliu/glue/contrib/config/nacos"
	_ "github.com/zhiyunliu/glue/contrib/queue/redis"
	_ "github.com/zhiyunliu/glue/contrib/queue/streamredis"
	_ "github.com/zhiyunliu/glue/contrib/registry/nacos"
	"github.com/zhiyunliu/glue/server/api"
	"github.com/zhiyunliu/glue/server/mqc"
)

func main() {

	apiSrv := api.New("apiserver")

	apiSrv.Handle("/delaysend", func(ctx context.Context) (res interface{}) {
		queueType := ctx.Request().Query().Get("t")

		return glue.Queue(queueType).DelaySend(ctx.Context(), queueType, map[string]interface{}{
			"delay": time.Now().Unix(),
		}, 10)
	})

	apiSrv.Handle("/send", func(ctx context.Context) (res interface{}) {
		queueType := ctx.Request().Query().Get("t")

		return glue.Queue(queueType).Send(ctx.Context(), queueType, map[string]interface{}{
			"send": time.Now().Unix(),
		})
	})

	mqcSrv1 := mqc.New("mqcserver_redis")
	mqcSrv1.Handle("redis", &demos.Orgdemo{})

	mqcSrv2 := mqc.New("mqcserver_streamredis")
	mqcSrv2.Handle("streamredis", &demos.Orgdemo{})

	app := glue.NewApp(glue.Server(apiSrv, mqcSrv1, mqcSrv2))

	app.Start()
}
