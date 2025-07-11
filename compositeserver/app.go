package main

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue-examples/compositeserver/handles"
	"github.com/zhiyunliu/glue/context"
	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/log"
	"github.com/zhiyunliu/glue/queue"
	"github.com/zhiyunliu/glue/transport"
	"github.com/zhiyunliu/glue/xhttp"
	"github.com/zhiyunliu/glue/xrpc"

	"github.com/zhiyunliu/glue/server/api"
	"github.com/zhiyunliu/glue/server/cron"
	"github.com/zhiyunliu/glue/server/mqc"
	"github.com/zhiyunliu/glue/server/rpc"
	"github.com/zhiyunliu/golibs/xtypes"
)

var Name = "compositeserver"

func init() {
	global.AppName = Name
	srvOpt := glue.Server(
		apiserver(),
		mqcserver(),
		cronserver(),
		rpcserver(),
	)
	opts = append(opts, srvOpt)
	//setTracerProvider("127.0.0.1:14268")
}

func apiserver() transport.Server {
	apiSrv := api.New("apiserver", api.WithServiceName("apiserver"), api.Log(log.WithRequest(), log.WithResponse()))

	apiSrv.Handle("/manual", func(ctx context.Context) interface{} {
		body, err := glue.RPC("").Swap(ctx, "grpc://payment-rpc/rpc/paywithdraw/manual", xrpc.WithWaitForReady(false))
		if err != nil {
			ctx.Log().Error("glue.RPC().GetRPC().Swap:", err)
		}
		ctx.Log().Debug(string(body.GetResult()))
		ctx.Log().Debug(body.GetHeader())
		ctx.Log().Debug(body.GetStatus())
		return nil
	})

	apiSrv.Handle("/log", handles.NewLogDemo())
	apiSrv.Handle("/xxx", func(ctx context.Context) interface{} {
		body, err := glue.Http("").Swap(ctx, "http://192.168.1.155:8080/demoapi", xhttp.WithMethod(http.MethodPost))
		if err != nil {
			ctx.Log().Error("glue.Http().GetHttp().xxx:", err)
		}
		ctx.Log().Debug(string(body.GetResult()))
		ctx.Log().Debug(body.GetHeader())
		ctx.Log().Debug(body.GetStatus())
		return string(body.GetResult())
	})
	apiSrv.Handle("/yyy", func(ctx context.Context) interface{} {
		body, err := glue.Http("").Swap(ctx, "xhttp://apiserver/demoapi", xhttp.WithMethod(http.MethodPost))
		if err != nil {
			ctx.Log().Error("glue.Http().GetHttp().yyy:", err)
		}
		ctx.Log().Debug(string(body.GetResult()))
		ctx.Log().Debug(body.GetHeader())
		ctx.Log().Debug(body.GetStatus())
		return string(body.GetResult())
	})
	apiSrv.Handle("/demoapi", func(ctx context.Context) interface{} {
		ctx.Log().Debug("api.demoapi")

		msg := queue.NewMsg(map[string]interface{}{
			"a": time.Now().Unix(),
		}, queue.WithXRequestID(ctx.Log().SessionID()))

		bodyMap := xtypes.XMap{}
		if err := ctx.Request().Body().ScanTo(&bodyMap); err != nil {
			return err
		}
		queueName := bodyMap.GetString("queue_name")
		if queueName == "" {
			queueName = "default"
		}

		err := glue.Queue(queueName).Send(ctx.Context(), "ayy.xx.xx", msg)
		if err != nil {
			ctx.Log().Errorf("send:%+v", err)
		}

		return xtypes.XMap{
			"a": 1,
			"b": 2,
		}

	})

	apiSrv.Handle("/demofile", func(ctx context.Context) interface{} {
		ctx.Log().Debug("api.demofile")

		body, err := glue.RPC("").Swap(ctx, "grpc://rpcserver/demorpcfile", xrpc.WithWaitForReady(false))
		if err != nil {
			ctx.Log().Error("glue.RPC().GetRPC().Swap:", err)
		}
		//time.Sleep(time.Second)
		return json.RawMessage(body.GetResult())
	})

	return apiSrv
}

func mqcserver() transport.Server {
	mqcSrv := mqc.New("mqcserver", mqc.Log(log.WithRequest(), log.WithResponse()))
	//mqcSrv.Use(tracing.Server(tracing.WithPropagator(propagation.TraceContext{}), tracing.WithTracerProvider(otel.GetTracerProvider())))

	mqcSrv.Handle("/demomqc", func(ctx context.Context) interface{} {

		body, err := glue.RPC("").Swap(ctx, "grpc://rpcserver/demorpc", xrpc.WithWaitForReady(false))
		if err != nil {
			ctx.Log().Error("glue.RPC().GetRPC().Swap:", err)
		}
		ctx.Log().Info(string(body.GetResult()))
		ctx.Log().Info(body.GetHeader())
		ctx.Log().Info(body.GetStatus())
		//time.Sleep(time.Second)
		return xtypes.XMap{
			"a": 1,
			"b": 2,
		}

	})

	return mqcSrv
}

func rpcserver() transport.Server {
	rpcSrv := rpc.New("rpcserver", rpc.WithServiceName("rpcserver"), rpc.Log(log.WithRequest(), log.WithResponse()))
	//rpcSrv.Use(tracing.Server(tracing.WithPropagator(propagation.TraceContext{}), tracing.WithTracerProvider(otel.GetTracerProvider())))
	rpcSrv.Handle("/demorpc", func(ctx context.Context) interface{} {

		return map[string]any{
			"out":    "rpc",
			"affect": "rpc.resp",
		}

	})

	rpcSrv.Handle("/demorpcfile", func(ctx context.Context) interface{} {
		var item = &DataItem{}
		err := ctx.Bind(item)
		if err != nil {
			return err
		}
		return item
	})

	return rpcSrv
}

func cronserver() transport.Server {
	cronSrv := cron.New("cronserver", cron.Log(log.WithRequest(), log.WithResponse()))

	cronSrv.Handle("/democron", func(ctx context.Context) interface{} {
		ctx.Log().Debug("democron")

		ctx.Log().Info(string(ctx.Request().Body().Bytes()))
		body, err := glue.Http("").Swap(ctx, "xhttp://apiserver/demoapi", xhttp.WithMethod(http.MethodPost),
			xhttp.WithRespHandler(func(resp *http.Response) (xhttp.Body, error) {
				return xhttp.NewEmptyBody(), nil
			}))
		if err != nil {
			ctx.Log().Error("glue.Http().GetHttp().xhttp:", err)
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
	return cronSrv
}
func GetDbName(ctx context.Context) string {
	dbName := ctx.Request().Query().Get("db_name")
	if dbName == "" {
		dbName = "xdb-mssql"
	}
	return dbName
}

type DataItem struct {
	A        string                `json:"a" form:"a" xml:"a"`
	B        int                   `json:"b" form:"b" xml:"b"`
	TestFile *multipart.FileHeader `json:"testfile" form:"testfile"`
}
