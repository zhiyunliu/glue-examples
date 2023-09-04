package main

import (
	sctx "context"
	"os"
	"reflect"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/config/file"
	"github.com/zhiyunliu/glue/context"
	_ "github.com/zhiyunliu/glue/contrib/cache/redis"
	_ "github.com/zhiyunliu/glue/contrib/config/consul"
	_ "github.com/zhiyunliu/glue/contrib/config/nacos"
	_ "github.com/zhiyunliu/glue/contrib/queue/redis"
	_ "github.com/zhiyunliu/glue/contrib/registry/nacos"
	_ "github.com/zhiyunliu/glue/contrib/xdb/mysql"
	"github.com/zhiyunliu/glue/xdb"

	_ "github.com/zhiyunliu/glue/contrib/xhttp/http"

	_ "github.com/zhiyunliu/glue/contrib/metrics/prometheus"
	_ "github.com/zhiyunliu/glue/contrib/xdb/postgres"
	_ "github.com/zhiyunliu/glue/contrib/xdb/sqlite"
	_ "github.com/zhiyunliu/glue/contrib/xdb/sqlserver"
	_ "github.com/zhiyunliu/glue/contrib/xdb/xgorm"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"

	_ "github.com/zhiyunliu/glue/contrib/dlocker/redis"

	"github.com/zhiyunliu/glue-examples/apiserver/dbconn"
	"github.com/zhiyunliu/glue-examples/apiserver/demos"
	"github.com/zhiyunliu/glue/errors"
	"github.com/zhiyunliu/glue/server/api"
	"github.com/zhiyunliu/golibs/xtypes"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var Name = "apiserver"

// Set global trace provider
func setTracerProvider(url string) error {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.AlwaysSample())),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
			attribute.String("env", "dev"),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

type redirect struct {
	code int
	url  string
}

func (r redirect) Render(ctx context.Context) error {
	gctx := ctx.GetImpl().(*gin.Context)
	gctx.Redirect(r.code, r.url)
	return nil
}

func main() {
	//setTracerProvider("http://127.0.0.1:14268/api/traces")

	apiSrv := api.New("apiserver", api.WithServiceName("xxxx"))
	//mqcSrv := mqc.New("bb")
	apiSrv.Handle("/redirect", func(ctx context.Context) interface{} {
		//ctx.Response().Redirect(301, "https://www.baidu.com")
		return redirect{code: 301, url: "https://www.qq.com"}
	})
	apiSrv.Handle("/demo", func(ctx context.Context) interface{} {
		ctx.Log().Debug("demo")

		body, err := glue.Http().Request(ctx.Context(), "http://www.baidu.com", nil)
		if err != nil {
			ctx.Log().Error("request", err)
		}
		ctx.Log().Info("body:", string(body.GetResult()))

		return xtypes.XMap{
			"a": 1,
			"b": 2,
		}
	})

	apiSrv.Handle("/error", func(ctx context.Context) interface{} {
		ctx.Log().Debug("error")
		return errors.New(300, "xxx")
	})

	apiSrv.Handle("/panic", func(ctx context.Context) interface{} {
		ctx.Log().Debug("panic")
		panic(fmt.Errorf("xx i am panic"))
	})

	apiSrv.Handle("/cfg", demos.NewCfg())
	apiSrv.Handle("/db", demos.NewDb())
	apiSrv.Handle("/cache", demos.NewCache())
	apiSrv.Handle("/queue", demos.NewQueue())
	apiSrv.Handle("/log", demos.NewLogDemo())
	apiSrv.Handle("/rpc", demos.NewGrpcDemo())
	apiSrv.Handle("/dlock", demos.NewDLock())
	apiSrv.Handle("/ppp/:aaa/:bbb", func(ctx context.Context) interface{} {
		return ctx.Request().Path().Params()
	})

	//apiSrv.Use(jwt.Server(jwt.WithSecret("123456")))
	//apiSrv.Use(ratelimit.Server())
	//apiSrv.Use(tracing.Server(tracing.WithTracerProvider(provider)))
	//apiSrv.Use(tracing.Server(tracing.WithPropagator(propagation.TraceContext{}), tracing.WithTracerProvider(otel.GetTracerProvider())))

	var binaryType = reflect.TypeOf((*Binary)(nil)).Elem()
	xdb.RegisterDbType("BINARY", binaryType)

	defaultConfigFile := os.Getenv("MS_DEFAULT_CONFIG_FILE")
	if defaultConfigFile == "" {
		defaultConfigFile = `E:\projects\golang\work\config.json`
	}

	app := glue.NewApp(
		glue.Server(apiSrv),
		glue.LogConcurrency(1),
		glue.WithConfigSource(file.NewSource(defaultConfigFile)),
		glue.StartingHook(func(ctx sctx.Context) error {
			return dbconn.Refactor(
				dbconn.WithSlowsql("xxx", "yyy"),
				dbconn.WithConfig("microsql", xdb.WithMaxOpen(10), xdb.WithShowQueryLog(true), xdb.WithLongQueryTime(2000)),
			)
		}))
	app.Start()
}

type Binary struct {
	data []uint8
}

func (b *Binary) Scan(src any) error {
	b.data = src.([]uint8)
	return nil
}

func (b Binary) GetBit(offset int) uint8 {
	return b.data[offset]
}

// func (b Binary) BitwiseAND(in Binary) uint8 {

// }

// func (b Binary) BitwiseOR(in Binary) uint8 {

// }
