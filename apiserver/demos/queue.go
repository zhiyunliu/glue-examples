package demos

import (
	"strconv"

	gel "github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
	"github.com/zhiyunliu/glue/queue"
)

type Queuedemo struct{}

func NewQueue() *Queuedemo {
	return &Queuedemo{}
}

func (d *Queuedemo) GetHandle(ctx context.Context) interface{} {
	ctx.Log().Debug("Queuedemo.get")
	queueObj := gel.Queue("default")

	err := queueObj.Send(ctx.Context(), "key", map[string]interface{}{
		"a": "1",
		"b": "2",
	})
	return map[string]interface{}{
		"err": err,
	}
}

func (d *Queuedemo) WithOptHandle(ctx context.Context) interface{} {
	ctx.Log().Debug("Queuedemo.get")

	tmpIdx := ctx.Request().Query().Get("idx")
	idx, _ := strconv.ParseInt(tmpIdx, 10, 32)

	queueObj := gel.Queue("default", queue.WithDBIndex(int(idx)))

	err := queueObj.Send(ctx.Context(), "key", map[string]interface{}{
		"a": "1",
		"b": "2",
	})
	return map[string]interface{}{
		"err": err,
	}
}
