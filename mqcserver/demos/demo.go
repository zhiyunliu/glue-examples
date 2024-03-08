package demos

import (
	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
	"github.com/zhiyunliu/golibs/xtypes"
)

type Orgdemo struct{}

func (d *Orgdemo) Handle(ctx context.Context) interface{} {

	param := xtypes.XMap{}

	if err := ctx.Request().Body().Scan(&param); err != nil {
		ctx.Log().Error("scan", err)
		return nil
	}
	cnt, _ := param.GetInt("cnt")
	param["cnt"] = cnt + 1
	if cnt > 5 {
		return nil
	}

	err := glue.Queue("streamredis").DelaySend(ctx.Context(), "streamredis3", param, 100)
	if err != nil {
		ctx.Log().Error("Orgdemo.err", err)
	}
	return nil
}
