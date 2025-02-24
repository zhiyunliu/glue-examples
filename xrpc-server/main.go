package main

import (
	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
	_ "github.com/zhiyunliu/glue/contrib/config/nacos"   //配置中心
	_ "github.com/zhiyunliu/glue/contrib/registry/nacos" //注册中心
	_ "github.com/zhiyunliu/glue/contrib/xrpc/grpc"      //grpc协议
	"github.com/zhiyunliu/glue/log"
	"golang.org/x/sync/errgroup"

	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/xrpc"

	"github.com/zhiyunliu/glue/server/rpc"
)

func main() {
	global.AppName = "xrpc-server"

	rcpSrv := rpc.New("rpcserver", rpc.WithServiceName(global.AppName), rpc.Log(log.WithRequest(), log.WithResponse()))

	rcpSrv.Handle("/normal", func(ctx context.Context) interface{} {
		return "normal"
	})

	rcpSrv.Handle("/bidirectional", func(ctx context.Context) interface{} {
		ctx.Log().Info("stream bidirectional start")
		req := ctx.Request().GetImpl()
		streamReq, ok := req.(xrpc.BidirectionalStreamRequest)
		if !ok {
			return "invalid request"
		}

		recvChan := make(chan *Item, 10)
		errGroup := errgroup.Group{}

		errGroup.Go(func() error {
			for {
				//异步接收数据流信息
				item := &Item{}
				if closed, err := streamReq.Recv(item); err != nil || closed {
					close(recvChan)
					return err
				}
				ctx.Log().Info("server.recv item", item)
				recvChan <- item
			}
		})

		errGroup.Go(func() error {
			idx := 0
			for item := range recvChan {
				//特殊场景判定不需要返回
				idx++
				if idx%2 == 0 {
					continue
				}
				//业务处理
				item.Name = "server_" + item.Name
				//异步写回相应内容
				err := streamReq.Send(item)
				if err != nil {
					return err
				}
				ctx.Log().Info("server.send item", item)
			}
			return nil
		})

		err := errGroup.Wait()
		ctx.Log().Info("stream bidirectional end", err)
		return nil
	})

	rcpSrv.Handle("/client", func(ctx context.Context) interface{} {
		ctx.Log().Info("stream client start")
		req := ctx.Request().GetImpl()
		streamReq, ok := req.(xrpc.ClientStreamRequest)
		if !ok {
			return "invalid request"
		}

		errGroup := errgroup.Group{}

		errGroup.Go(func() error {
			for {
				//异步接收数据流信息
				item := &Item{}
				if closed, err := streamReq.Recv(item); err != nil || closed {
					return err
				}
				ctx.Log().Info("server.recv item", item)
			}
		})

		err := errGroup.Wait()
		ctx.Log().Info("stream client end", err)
		return map[string]any{
			"client": "ok",
		}
	})

	rcpSrv.Handle("/server", func(ctx context.Context) interface{} {
		ctx.Log().Info("stream server start")

		param := &streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}
		ctx.Log().Info("stream server param", param)

		req := ctx.Request().GetImpl()
		streamReq, ok := req.(xrpc.ServerStreamRequest)
		if !ok {
			return "invalid request"
		}

		recvChan := make(chan *Item, 10)
		recvChan <- &Item{Id: 1, Name: "1 stream server"}
		recvChan <- &Item{Id: 2, Name: "2 stream server"}
		recvChan <- &Item{Id: 3, Name: "3 stream server"}
		recvChan <- &Item{Id: 4, Name: "4 stream server"}
		recvChan <- &Item{Id: 5, Name: "5 stream server"}
		close(recvChan)
		errGroup := errgroup.Group{}

		errGroup.Go(func() error {
			for item := range recvChan {
				//异步接收数据流信息
				if err := streamReq.Send(item); err != nil {
					return err
				}
				ctx.Log().Info("server.Send item", item)
			}
			return nil
		})

		err := errGroup.Wait()
		ctx.Log().Info("stream server end", err)
		return nil
	})

	app := glue.NewApp(glue.Server(rcpSrv))
	app.Start()
}

type streamParams struct {
	List []*Item `json:"list"`
}

type Item struct {
	Id   int
	Name string
}
