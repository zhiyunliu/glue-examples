package main

import (
	"context"

	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue-examples/xdb-postgres/demos"
	_ "github.com/zhiyunliu/glue/contrib/cache/redis"
	_ "github.com/zhiyunliu/glue/contrib/metrics/prometheus"
	_ "github.com/zhiyunliu/glue/contrib/registry/nacos"

	_ "github.com/zhiyunliu/queue-redis"
	_ "github.com/zhiyunliu/xdb-postgres"

	"github.com/zhiyunliu/glue/server/api"
)

func main() {

	apiSrv := api.New("apiserver")
	//mqcSrv := mqc.New("bb")

	apiSrv.Handle("/xdb-postgres/query", demos.QueryData)
	apiSrv.Handle("/xdb-postgres/like", demos.LikeData)
	apiSrv.Handle("/xdb-postgres/like2", demos.Like2Data)
	apiSrv.Handle("/xdb-postgres/in", demos.InData)
	apiSrv.Handle("/xdb-postgres/notin", demos.NotInData)
	apiSrv.Handle("/xdb-postgres/notin2", demos.NotIn2Data)
	apiSrv.Handle("/xdb-postgres/compare", demos.CompareData)
	apiSrv.Handle("/xdb-postgres/compare2", demos.Compare2Data)
	apiSrv.Handle("/xdb-postgres/cte", demos.CteData)

	app := glue.NewApp(glue.Server(apiSrv),
		glue.StartingHook(func(ctx context.Context) error {
			return nil
		}))
	app.Start()
}
