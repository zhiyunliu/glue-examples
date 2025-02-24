package demos

import (
	sctx "context"
	"strconv"
	"time"

	gel "github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
	"github.com/zhiyunliu/glue/dlocker"
)

type Dlockdemo struct{}

func NewDLock() *Dlockdemo {
	return &Dlockdemo{}
}

func (d *Dlockdemo) CreateHandle(ctx context.Context) interface{} {
	lockKey := ctx.Request().Query().Get("key")
	if lockKey == "" {
		lockKey = "lockkey"
	}
	opts := []dlocker.Option{}
	data := ctx.Request().Query().Get("data")
	if data != "" {
		opts = append(opts, dlocker.WithData(data))
	}
	renewal := ctx.Request().Query().Get("renewal")
	if renewal != "" {
		opts = append(opts, dlocker.WithAutoRenewal())
	}
	reentrant := ctx.Request().Query().Get("reentrant")
	if reentrant != "" {
		ret, _ := strconv.ParseBool(reentrant)
		opts = append(opts, dlocker.WithReentrant(ret))
	}

	locker := gel.DLocker(lockKey, opts...)
	isok, err := locker.Acquire(ctx.Context(), 2)
	if err != nil {
		ctx.Log().Errorf("Acquire err:%+v", err)
		return err
	}
	if isok {
		sleeptime := 8
		ctx.Log().Debugf("sleep %d seconds", sleeptime)
		time.Sleep(time.Duration(sleeptime) * time.Second)
		err = locker.Renewal(ctx.Context(), 10)
		ctx.Log().Debugf("Renewal 10 seconds %+v", err)

		go func() {
			time.Sleep(time.Second * 20)
			locker.Release(sctx.Background())
		}()
		return "success"
	}
	return "lock failure"
}

func (d *Dlockdemo) ReentrantHandle(ctx context.Context) interface{} {
	lockKey := ctx.Request().Query().Get("key")
	if lockKey == "" {
		lockKey = "lockkey"
	}
	locker := gel.DLocker(lockKey, dlocker.WithAutoRenewal(), dlocker.WithReentrant(false))
	isok, err := locker.Acquire(ctx.Context(), 2)
	if err != nil {
		ctx.Log().Errorf("Acquire err:%+v", err)
		return err
	}
	if isok {
		sleeptime := 8
		ctx.Log().Debugf("sleep %d seconds", sleeptime)
		time.Sleep(time.Duration(sleeptime) * time.Second)
		err = locker.Renewal(ctx.Context(), 10)
		ctx.Log().Debugf("Renewal 10 seconds %+v", err)

		go func() {
			time.Sleep(time.Second * 20)
			locker.Release(sctx.Background())
		}()
		return "success"
	}
	return "lock failure"
}
