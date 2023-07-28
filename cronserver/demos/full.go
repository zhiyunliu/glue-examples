package demos

import (
	"strconv"
	"time"

	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
)

type Fulldemo struct{}

func (d *Fulldemo) Handle(ctx context.Context) interface{} {
	ctx.Log().Infof("cron.demo:%s", time.Now().Format("2006-01-02 15:04:05"))

	ctx.Log().Infof("body-1:%s", ctx.Request().Body().Bytes())
	time.Sleep(time.Millisecond * 200)

	mapData := map[string]interface{}{}
	ctx.Request().Body().Scan(&mapData)
	ctx.Log().Infof("body-2:%+v", mapData)

	return "success"
}

func (d *Fulldemo) MqHandle(ctx context.Context) interface{} {

	mapData := map[string]string{}
	ctx.Request().Body().Scan(&mapData)

	qt := mapData["qt"]
	qn := mapData["qn"]
	tmpcnt := mapData["cnt"]
	cnt, _ := strconv.ParseInt(tmpcnt, 10, 32)
	if cnt == 0 {
		cnt = 10
	}

	ctx.Log().Infof("qt:%s,qn:%s", qt, qn)
	for c := int64(0); c < cnt; c++ {
		glue.Queue(qt).Send(ctx.Context(), qn, map[string]interface{}{
			"send": time.Now().Unix(),
		})
	}
	return nil
}

func (d *Fulldemo) NoneBodyHandle(ctx context.Context) interface{} {
	ctx.Log().Infof("cron.NoneBody:%s", time.Now().Format("2006-01-02 15:04:05"))

	ctx.Log().Infof("NoneBody-1:%s", ctx.Request().Body().Bytes())
	time.Sleep(time.Millisecond * 200)

	mapData := map[string]interface{}{}
	ctx.Request().Body().Scan(&mapData)
	ctx.Log().Infof("NoneBody-2:%+v", mapData)

	return "success"
}

func (d *Fulldemo) NoneBody2Handle(ctx context.Context) interface{} {
	ctx.Log().Infof("cron. NoneBody2:%s", time.Now().Format("2006-01-02 15:04:05"))

	ctx.Log().Infof(" NoneBody2-1:%s", ctx.Request().Body().Bytes())
	time.Sleep(time.Millisecond * 200)

	mapData := map[string]interface{}{}
	ctx.Request().Body().Scan(&mapData)
	ctx.Log().Infof(" NoneBody2-2:%+v", mapData)

	return "success"
}

func (d *Fulldemo) NotRunHandle(ctx context.Context) interface{} {
	ctx.Log().Infof("cron.NotRun:%s", time.Now().Format("2006-01-02 15:04:05"))

	ctx.Log().Infof("NoneBody-1:%s", ctx.Request().Body().Bytes())
	time.Sleep(time.Millisecond * 5000)

	mapData := map[string]interface{}{}
	ctx.Request().Body().Scan(&mapData)
	ctx.Log().Infof("NoneBody-2:%+v", mapData)

	return "success"
}
