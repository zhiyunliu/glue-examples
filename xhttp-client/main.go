package main

import (
	"bufio"
	"encoding/json"
	"net/http"

	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/constants"
	"github.com/zhiyunliu/glue/context"
	_ "github.com/zhiyunliu/glue/contrib/config/nacos"   //配置中心
	_ "github.com/zhiyunliu/glue/contrib/registry/nacos" //注册中心
	_ "github.com/zhiyunliu/glue/contrib/xhttp/http"     //http协议
	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/xhttp"

	"github.com/zhiyunliu/glue/server/api"
	"github.com/zhiyunliu/golibs/xsse"
)

func main() {
	global.AppName = "xhttp-client"

	apiSrv := api.New("apiserver", api.WithServiceName(global.AppName))

	apiSrv.Handle("/normal", func(ctx context.Context) interface{} {
		param := streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}
		body, err := glue.Http("").Request(ctx.Context(), "xhttp://xhttp-server/normal", param,
			xhttp.WithXRequestID(ctx.Log().SessionID()),
			xhttp.WithContentType(constants.ContentTypeApplicationJSON),
			xhttp.WithMethod(http.MethodPost),
		)

		return map[string]any{
			"err":    err,
			"result": string(body.GetResult()),
		}
	})

	apiSrv.Handle("/stream", func(ctx context.Context) interface{} {
		param := streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}
		result := []*Item{}

		//调用grpc服务

		var respHandler = func(resp *http.Response) (body xhttp.Body, err error) {
			ctx.Log().Info("stream response start")
			body = xhttp.NewEmptyBody()

			reader := bufio.NewReader(resp.Body) // 使用 bufio.NewReader 来读取数据
			for {
				eventData, err := reader.ReadString('\n') // 读取直到换行符，假设每个事件是一个独立的行
				if err != nil {
					if err.Error() == "EOF" { // 如果是 EOF，说明数据流结束
						break
					}
					return nil, err // 发生其他错误
				}
				//do something
				ctx.Log().Info(eventData)
			}
			ctx.Log().Info("stream response end")
			return
		}

		body, err := glue.Http("").Request(ctx.Context(), "xhttp://xhttp-server/stream", param,
			xhttp.WithXRequestID(ctx.Log().SessionID()),
			xhttp.WithContentType(constants.ContentTypeApplicationJSON),
			xhttp.WithMethod(http.MethodPost),
			xhttp.WithRespHandler(respHandler),
		)

		ctx.Log().Infof("stream.body:%s,ERR:%v", string(body.GetResult()), err)

		return map[string]any{
			"err":    err,
			"result": result,
		}
	})

	apiSrv.Handle("/stream2", func(ctx context.Context) interface{} {
		param := streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}
		result := []*Item{}

		//调用grpc服务

		var sseHandler xhttp.SSEHandler = func(evt *xsse.Event) (err error) {
			ctx.Log().Info("stream response", evt)
			result = append(result, evt.Data.(*Item))
			return
		}

		body, err := glue.Http("").Request(ctx.Context(), "xhttp://xhttp-server/stream", param,
			xhttp.WithXRequestID(ctx.Log().SessionID()),
			xhttp.WithContentType(constants.ContentTypeApplicationJSON),
			xhttp.WithMethod(http.MethodPost),
			xhttp.WithSSEHandler(sseHandler, xsse.WithDecoderUnmarshal(func(bytes []byte) (data any, err error) {
				var item *Item = &Item{}
				err = json.Unmarshal(bytes, item)
				return item, err
			})),
		)

		ctx.Log().Infof("stream.body:%s,ERR:%v", string(body.GetResult()), err)

		return map[string]any{
			"err":    err,
			"result": result,
		}
	})

	apiSrv.Handle("/stream-chan", func(ctx context.Context) interface{} {
		param := streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}
		result := []*Item{}

		//调用grpc服务

		var sseHandler xhttp.SSEHandler = func(evt *xsse.Event) (err error) {
			ctx.Log().Info("stream response", evt)
			result = append(result, evt.Data.(*Item))
			return
		}

		body, err := glue.Http("").Request(ctx.Context(), "xhttp://xhttp-server/stream-chan", param,
			xhttp.WithXRequestID(ctx.Log().SessionID()),
			xhttp.WithContentType(constants.ContentTypeApplicationJSON),
			xhttp.WithMethod(http.MethodPost),
			xhttp.WithSSEHandler(sseHandler, xsse.WithDecoderUnmarshal(func(bytes []byte) (data any, err error) {
				var item *Item = &Item{}
				err = json.Unmarshal(bytes, item)
				return item, err
			})),
		)

		ctx.Log().Infof("stream.body:%s,ERR:%v", string(body.GetResult()), err)

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
