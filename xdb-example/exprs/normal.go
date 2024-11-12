package exprs

import (
	"strings"

	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
	"github.com/zhiyunliu/golibs/xtypes"
)

type Normal struct {
}

func (c *Normal) NHandle(ctx context.Context) (res any) {
	const compare = `
	select * from ljy_test t where t.a=@{a} &{b} |{t.c} and t.d in (${t.d})
	`

	dataMap := xtypes.XMap{}

	for k, v := range ctx.Request().Query().Values() {
		dataMap[k] = v
		vals := strings.Split(v, ",")
		if len(vals) > 1 {
			dataMap[k] = vals
		}
	}

	dbObj := glue.DB(GetDbName(ctx))
	results, err := dbObj.Query(ctx.Context(), compare, dataMap)
	if err != nil {
		return err
	}
	return results
}

func (c *Normal) CachedHandle(ctx context.Context) (res any) {
	const compare = `
	select * from ljy_test t where t.a=@{a} and t.b=@{b}
	`
	dbObj := glue.DB(GetDbName(ctx))
	results, err := dbObj.Query(ctx.Context(), compare, ctx.Request().Query().Values())
	if err != nil {
		return err
	}
	return results
}
