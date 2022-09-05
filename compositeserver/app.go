package main

import (
	"net/http"
	"time"

	gel "github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue-examples/compositeserver/handles"
	"github.com/zhiyunliu/glue/context"
	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/log"
	"github.com/zhiyunliu/glue/middleware/tracing"
	"github.com/zhiyunliu/glue/transport"
	"github.com/zhiyunliu/glue/xhttp"
	"github.com/zhiyunliu/glue/xrpc"

	"github.com/zhiyunliu/glue/server/api"
	"github.com/zhiyunliu/glue/server/cron"
	"github.com/zhiyunliu/glue/server/mqc"
	"github.com/zhiyunliu/glue/server/rpc"
	"github.com/zhiyunliu/golibs/xtypes"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var Name = "compositeserver"

func init() {
	global.AppName = Name
	srvOpt := gel.Server(
		apiserver(),
		mqcserver(),
		cronserver(),
		rpcserver(),
	)
	opts = append(opts, srvOpt, gel.LogConcurrency(1))
	setTracerProvider("http://127.0.0.1:14268/api/traces")
}

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

func apiserver() transport.Server {
	apiSrv := api.New("apiserver", api.WithServiceName("apiserver"), api.Log(log.WithRequest(), log.WithResponse()))
	apiSrv.Use(tracing.Server(tracing.WithPropagator(propagation.TraceContext{}), tracing.WithTracerProvider(otel.GetTracerProvider())))
	apiSrv.Handle("/log", handles.NewLogDemo())
	apiSrv.Handle("/xxx", func(ctx context.Context) interface{} {
		body, err := gel.Http().Swap(ctx, "http://192.168.1.155:8080/demoapi", xhttp.WithMethod(http.MethodPost))
		if err != nil {
			ctx.Log().Error("gel.Http().GetHttp().Swap:", err)
		}
		ctx.Log().Debug(string(body.GetResult()))
		ctx.Log().Debug(body.GetHeader())
		ctx.Log().Debug(body.GetStatus())
		return string(body.GetResult())
	})
	apiSrv.Handle("/yyy", func(ctx context.Context) interface{} {
		body, err := gel.Http().Swap(ctx, "xhttp://apiserver/demoapi", xhttp.WithMethod(http.MethodPost))
		if err != nil {
			ctx.Log().Error("gel.Http().GetHttp().Swap:", err)
		}
		ctx.Log().Debug(string(body.GetResult()))
		ctx.Log().Debug(body.GetHeader())
		ctx.Log().Debug(body.GetStatus())
		return string(body.GetResult())
	})
	apiSrv.Handle("/demoapi", func(ctx context.Context) interface{} {
		ctx.Log().Debug("api.demoapi")

		body, err := gel.RPC().Swap(ctx, "grpc://rpcserver/demorpc", xrpc.WithWaitForReady(false))
		if err != nil {
			ctx.Log().Error("gel.RPC().GetRPC().Swap:", err)
		}
		ctx.Log().Debug(string(body.GetResult()))
		ctx.Log().Debug(body.GetHeader())
		ctx.Log().Debug(body.GetStatus())
		//time.Sleep(time.Second)
		return xtypes.XMap{
			"a": 1,
			"b": 2,
		}
	})
	return apiSrv
}

func mqcserver() transport.Server {
	mqcSrv := mqc.New("mqcserver", mqc.Log(log.WithRequest(), log.WithResponse()))
	//mqcSrv.Use(tracing.Server(tracing.WithPropagator(propagation.TraceContext{}), tracing.WithTracerProvider(otel.GetTracerProvider())))

	mqcSrv.Handle("/demomqc", func(ctx context.Context) interface{} {
		ctx.Log().Debug("demomqc")
		body, err := gel.Http().Swap(ctx, "xhttp://apiserver/demoapi", xhttp.WithMethod(http.MethodPost))
		if err != nil {
			ctx.Log().Error("gel.Http().GetHttp().Swap:", err)
		}
		ctx.Log().Debug(string(body.GetResult()))
		ctx.Log().Debug(body.GetHeader())
		ctx.Log().Debug(body.GetStatus())
		//time.Sleep(time.Second * 2)
		return xtypes.XMap{
			"a": 1,
			"b": 2,
		}
	})

	return mqcSrv
}

func rpcserver() transport.Server {
	rpcSrv := rpc.New("rpcserver", rpc.WithServiceName("rpcserver"), rpc.Log(log.WithRequest(), log.WithResponse()))
	rpcSrv.Use(tracing.Server(tracing.WithPropagator(propagation.TraceContext{}), tracing.WithTracerProvider(otel.GetTracerProvider())))
	rpcSrv.Handle("/demorpc", func(ctx context.Context) interface{} {
		//time.Sleep(time.Second * 1)
		ctx.Log().Debug("demorpc")
		return xtypes.XMap{
			"a": 1,
			"b": 2,
		}
	})

	return rpcSrv
}

func cronserver() transport.Server {
	cronSrv := cron.New("cronserver", cron.Log(log.WithRequest(), log.WithResponse()))

	cronSrv.Handle("/democron", func(ctx context.Context) interface{} {
		ctx.Log().Debug("democron")

		gel.Queue("default").Send(ctx.Context(), "xx.xx.xx", map[string]interface{}{
			"a": time.Now().Unix(),
		})

		return xtypes.XMap{
			"a": 1,
			"b": 2,
		}
	})
	return cronSrv
}
