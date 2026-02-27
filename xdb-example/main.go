package main

import (
	"context"
	"fmt"
	"time"

	"github.com/urfave/cli"
	"github.com/zhiyunliu/glue"
	_ "github.com/zhiyunliu/glue/contrib/cache/redis"
	_ "github.com/zhiyunliu/glue/contrib/config/nacos"

	_ "github.com/zhiyunliu/glue/contrib/registry/nacos"
	"github.com/zhiyunliu/glue/xdb"

	_ "github.com/zhiyunliu/glue/contrib/xhttp/http"

	_ "github.com/zhiyunliu/glue/contrib/metrics/prometheus"

	_ "github.com/zhiyunliu/queue-redis"
	_ "github.com/zhiyunliu/xdb-mssql"

	"github.com/zhiyunliu/glue-examples/xdb-example/exprs"
	"github.com/zhiyunliu/glue-examples/xdb-example/matchers"
	"github.com/zhiyunliu/glue/server/api"
)

func main() {

	apiSrv := api.New("apiserver", api.WithServiceName("xxxx"))
	//mqcSrv := mqc.New("bb")

	apiSrv.Handle("/tests/compare", &exprs.Compare{})
	apiSrv.Handle("/tests/in", &exprs.In{})
	apiSrv.Handle("/tests/like", &exprs.Like{})
	apiSrv.Handle("/tests/normal", &exprs.Normal{})
	apiSrv.Handle("/tests/output", &exprs.Output{})

	app := glue.NewApp(glue.Server(apiSrv),
		glue.Command(appendCli()),
		glue.StartingHook(func(ctx context.Context) error {
			return nil
		}))
	app.Start()
}

func appendCli() *cli.Command {
	var matcher string

	return &cli.Command{
		Name: "matcher",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Destination: &matcher,
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println(matcher)
			xdb.RegistExpressionMatcher(matcher, &matchers.CustomMatcher{})
			time.Sleep(time.Second)
			return nil
		},
	}
}
