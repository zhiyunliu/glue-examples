package main

import (
	"net/http"

	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/constants"
	"github.com/zhiyunliu/glue/context"
	_ "github.com/zhiyunliu/glue/contrib/config/nacos"   //配置中心
	_ "github.com/zhiyunliu/glue/contrib/registry/nacos" //注册中心
	_ "github.com/zhiyunliu/glue/contrib/xrpc/grpc"      //grpc协议
	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/xrpc"
	"golang.org/x/sync/errgroup"

	"github.com/zhiyunliu/glue/server/api"
)

func main() {
	global.AppName = "xrpc-client"

	apiSrv := api.New("apiserver", api.WithServiceName(global.AppName))

	apiSrv.Handle("/normal", func(ctx context.Context) interface{} {
		param := streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}
		body, err := glue.RPC("").Request(ctx.Context(), "grpc://xrpc-server/normal", param,
			xrpc.WithXRequestID(ctx.Log().SessionID()),
			xrpc.WithContentType(constants.ContentTypeApplicationJSON),
			xrpc.WithMethod(http.MethodPost),
		)

		return map[string]any{
			"err":    err,
			"result": string(body.GetResult()),
		}
	})

	apiSrv.Handle("/stream/bidirectional", func(ctx context.Context) interface{} {
		param := streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}

		//调用grpc服务
		result := []*Item{}
		var processor xrpc.BidirectionalStreamProcessor = func(sc xrpc.BidirectionalStreamClient) error {
			errGroup := errgroup.Group{}

			ctx.Log().Info("grpc client start")
			errGroup.Go(func() error {
				for _, item := range param.List {
					if err := sc.Send(item); err != nil {
						return err
					}
				}
				return sc.CloseSend()
			})

			errGroup.Go(func() error {
				for {
					item := &Item{}
					closed, err := sc.Recv(item)
					if err != nil || closed {
						return err
					} else {
						ctx.Log().Info(item)
					}
					result = append(result, item)
				}
			})

			err := errGroup.Wait()
			ctx.Log().Info("grpc client end")
			return err
		}

		body, err := glue.RPC("").Request(ctx.Context(), "grpc://xrpc-server/bidirectional", nil,
			xrpc.WithXRequestID(ctx.Log().SessionID()),
			xrpc.WithContentType(constants.ContentTypeApplicationJSON),
			xrpc.WithMethod(http.MethodPost),
			xrpc.WithStreamProcessor(processor),
		)
		ctx.Log().Info("stream/bidirectional", string(body.GetResult()), err)
		return map[string]any{
			"err":    err,
			"result": result,
		}
	})

	apiSrv.Handle("/stream/client-default-array", func(ctx context.Context) interface{} {
		param := streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}

		body, err := glue.RPC("").Request(ctx.Context(), "grpc://xrpc-server/client", param.List,
			xrpc.WithXRequestID(ctx.Log().SessionID()),
			xrpc.WithContentType(constants.ContentTypeApplicationJSON),
			xrpc.WithMethod(http.MethodPost),
			xrpc.WithStreamDefaultProcessor())

		ctx.Log().Info("stream/client-default-array", string(body.GetResult()), err)

		return body
	})

	apiSrv.Handle("/stream/client-default-implement", func(ctx context.Context) interface{} {
		param := streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}

		body, err := glue.RPC("").Request(ctx.Context(), "grpc://xrpc-server/client", param,
			xrpc.WithXRequestID(ctx.Log().SessionID()),
			xrpc.WithContentType(constants.ContentTypeApplicationJSON),
			xrpc.WithMethod(http.MethodPost),
			xrpc.WithStreamDefaultProcessor())

		ctx.Log().Info("stream/client-default-implement", string(body.GetResult()), err)

		return body
	})

	apiSrv.Handle("/stream/client-default-chan", func(ctx context.Context) interface{} {
		param := streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}

		body, err := glue.RPC("").Request(ctx.Context(), "grpc://xrpc-server/client", streamParamsChan(param),
			xrpc.WithXRequestID(ctx.Log().SessionID()),
			xrpc.WithContentType(constants.ContentTypeApplicationJSON),
			xrpc.WithMethod(http.MethodPost),
			xrpc.WithStreamDefaultProcessor())

		ctx.Log().Info("stream/client-default-chan", string(body.GetResult()), err)

		return body
	})

	apiSrv.Handle("/stream/client", func(ctx context.Context) interface{} {
		param := streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}

		errGroup := errgroup.Group{}
		//调用grpc服务

		var processor xrpc.ClientStreamProcessor = func(sc xrpc.ClientStreamClient) error {
			ctx.Log().Info("grpc client start")

			errGroup.Go(func() error {
				for i := 0; i < len(param.List); i++ {
					err := sc.Send(param.List[i])
					if err != nil {
						return err
					}
				}
				return nil
			})

			err := errGroup.Wait()
			ctx.Log().Info("grpc client end")
			return err
		}

		body, err := glue.RPC("").Request(ctx.Context(), "grpc://xrpc-server/client", nil,
			xrpc.WithXRequestID(ctx.Log().SessionID()),
			xrpc.WithContentType(constants.ContentTypeApplicationJSON),
			xrpc.WithMethod(http.MethodPost),
			xrpc.WithStreamProcessor(processor),
		)

		ctx.Log().Info("stream/client", string(body.GetResult()), err)

		return body
	})

	apiSrv.Handle("/stream/server", func(ctx context.Context) interface{} {
		param := streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}

		result := []*Item{}

		errGroup := errgroup.Group{}
		//调用grpc服务

		var processor = func(sc xrpc.ServerStreamClient) error {
			ctx.Log().Info("grpc server start")

			errGroup.Go(func() error {
				for {
					item := &Item{}
					closed, err := sc.Recv(item)
					if err != nil || closed {
						return err
					}
					ctx.Log().Info(item)
					result = append(result, item)
				}
			})

			err := errGroup.Wait()
			ctx.Log().Info("grpc server end")
			return err
		}

		body, err := glue.RPC("").Request(ctx.Context(), "grpc://xrpc-server/server", param,
			xrpc.WithXRequestID(ctx.Log().SessionID()),
			xrpc.WithContentType(constants.ContentTypeApplicationJSON),
			xrpc.WithMethod(http.MethodPost),
			xrpc.WithStreamProcessor(processor),
		)

		ctx.Log().Info("stream/server", string(body.GetResult()), err)

		return map[string]any{
			"err":    err,
			"result": result,
		}
	})

	app := glue.NewApp(glue.Server(apiSrv))
	_ = app.Start()
}

type streamParams struct {
	List []*Item `json:"list"`
}

type Item struct {
	Id   int
	Name string
}

func (p streamParams) GetObjects() []any {
	objs := make([]any, 0, len(p.List))
	for _, item := range p.List {
		objs = append(objs, item)
	}
	return objs
}

type streamParamsChan streamParams

func (p streamParamsChan) GetObject() <-chan any {
	ch := make(chan any, len(p.List))
	go func() {
		for _, item := range p.List {
			ch <- item
		}
		close(ch)
	}()
	return ch
}
