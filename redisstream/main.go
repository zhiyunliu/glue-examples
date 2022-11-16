package main

import (
	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue-examples/redisstream/demos"
	_ "github.com/zhiyunliu/glue/contrib/queue/redis"
	_ "github.com/zhiyunliu/glue/contrib/queue/streamredis"
	"github.com/zhiyunliu/glue/server/api"
	"github.com/zhiyunliu/glue/server/mqc"
)

func main() {

	apiSrv := api.New("apiserver")
	mqcSrv := mqc.New("mqcserver")

	apiSrv.Handle("/queue", demos.NewMQ())

	mqcSrv.Handle("queue1", demos.NewMQC())
	mqcSrv.Handle("yy.xx.xx", demos.NewMQC())

	app := glue.NewApp(glue.Server(mqcSrv), glue.LogConcurrency(1))
	app.Start()
}
