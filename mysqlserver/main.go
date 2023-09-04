package main

import (
	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
	_ "github.com/zhiyunliu/glue/contrib/cache/redis"
	_ "github.com/zhiyunliu/glue/contrib/config/consul"
	_ "github.com/zhiyunliu/glue/contrib/config/nacos"
	_ "github.com/zhiyunliu/glue/contrib/queue/redis"
	_ "github.com/zhiyunliu/glue/contrib/registry/nacos"
	_ "github.com/zhiyunliu/glue/contrib/xdb/mysql"
	"github.com/zhiyunliu/glue/log"
	"github.com/zhiyunliu/glue/xdb"

	_ "github.com/zhiyunliu/glue/contrib/xhttp/http"

	//_ "github.com/zhiyunliu/glue/contrib/xdb/oracle"
	_ "github.com/zhiyunliu/glue/contrib/xdb/postgres"
	_ "github.com/zhiyunliu/glue/contrib/xdb/sqlite"
	_ "github.com/zhiyunliu/glue/contrib/xdb/sqlserver"

	_ "github.com/zhiyunliu/glue/contrib/dlocker/redis"

	"github.com/zhiyunliu/glue/server/api"
)

var Name = "apiserver"

func main() {

	apiSrv := api.New("apiserver", api.Log(log.WithResponse()))

	apiSrv.Handle("/mysql", func(ctx context.Context) interface{} {
		a := ctx.Request().Query().Get("a")
		execSql := `select a,b,c,d,e,f,g,h from tmp_glue_test where a = @{a} &{d} or h in(${h}) or h = @{a} |{d} and c in (${c}) and g=@{a} `
		//execSql := `select a,b,c,d  from tmp_glue_test where a = @p_a`

		dbObj := getDbObj(ctx)

		rows, err := dbObj.Query(ctx.Context(), execSql, map[string]interface{}{"a": a, "d": "d", "h": []string{"h2", "h3"}, "c": []int{1, 2, 3}})
		if err != nil {
			if dberr, ok := err.(xdb.DbError); ok {
				ctx.Log().Error(dberr.SQL(), dberr.Args())
			}
			return err
		}
		// 遍历查询到的数据

		return rows

	})

	apiSrv.Handle("/cache", func(ctx context.Context) interface{} {
		a := ctx.Request().Query().Get("a")
		execSql := `select a,b,c,d,e,f,g,h from tmp_glue_test where a = @{a}  and g=@{g}  `
		//execSql := `select a,b,c,d  from tmp_glue_test where a = @p_a`

		dbObj := getDbObj(ctx)

		rows, err := dbObj.Query(ctx.Context(), execSql, map[string]interface{}{
			"a": a,
			"d": "d",
			"g": ctx.Request().Query().Get("g"),
			"h": []string{"h2", "h3"},
			"c": []int{1, 2, 3}})
		if err != nil {
			if dberr, ok := err.(xdb.DbError); ok {
				ctx.Log().Error(dberr.SQL(), dberr.Args())
			}
			return err
		}
		// 遍历查询到的数据

		return rows

	})

	app := glue.NewApp(glue.Server(apiSrv), glue.LogConcurrency(1))
	app.Start()
}

func getDbObj(ctx context.Context) xdb.IDB {
	dbtype := ctx.Request().Query().Get("dbtype")
	if dbtype == "" {
		dbtype = "localhost"
	}
	return glue.DB(dbtype)
}
