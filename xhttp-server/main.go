package main

import (
	"fmt"
	"time"

	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
	_ "github.com/zhiyunliu/glue/contrib/config/nacos"   //配置中心
	_ "github.com/zhiyunliu/glue/contrib/registry/nacos" //注册中心
	_ "github.com/zhiyunliu/glue/contrib/xrpc/grpc"      //grpc协议
	"github.com/zhiyunliu/glue/global"
	"golang.org/x/sync/errgroup"

	"github.com/zhiyunliu/glue/server/api"
	"github.com/zhiyunliu/golibs/xsse"
)

func main() {
	global.AppName = "xhttp-server"

	apiSrv := api.New("apiserver", api.WithServiceName(global.AppName))

	apiSrv.Handle("/normal", func(ctx context.Context) interface{} {
		param := streamParams{}

		if err := ctx.Bind(&param); err != nil {
			return err
		}
		return map[string]any{
			"param":  param,
			"result": "normal",
		}
	})

	apiSrv.Handle("/stream", func(ctx context.Context) interface{} {
		param := streamParams{}
		if err := ctx.Bind(&param); err != nil {
			return err
		}
		result := &SSEResult{
			List: param.List,
		}
		return result
	})

	apiSrv.Handle("/stream-chan", func(ctx context.Context) interface{} {
		param := streamParams{}
		if err := ctx.Bind(&param); err != nil {
			return err
		}
		result := &SSEResultChan{
			List: make(chan *Item),
		}

		group := errgroup.Group{}
		group.Go(func() error {
			for _, item := range param.List {
				time.Sleep(time.Second)
				result.List <- item
			}
			close(result.List)
			return nil
		})

		return result
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

type SSEResult struct {
	List []*Item `json:"list"`
	idx  int
}

func (s *SSEResult) GetEvent() (evt *xsse.Event, ok bool) {
	if s.idx >= len(s.List) {
		return nil, false
	}
	item := s.List[s.idx]
	s.idx++
	return &xsse.Event{
		Id:   fmt.Sprint(s.idx),
		Data: item,
	}, true
}

type SSEResultChan struct {
	List chan *Item
}

func (s *SSEResultChan) GetEvent() (evt *xsse.Event, ok bool) {
	item, ok := <-s.List
	if !ok {
		return nil, false
	}
	return &xsse.Event{
		Id:   fmt.Sprint(item.Id),
		Data: item,
	}, ok
}
