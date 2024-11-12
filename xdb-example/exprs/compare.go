package exprs

import (
	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
)

type Compare struct {
}

func (c *Compare) Handle(ctx context.Context) (res any) {

	const compare = `
	select * from ljy_test t where t.a=@{a} &{b>b} &{>t.b} &{t.c>=c} |{t.d=d} &{t.a>b}
	`

	dbObj := glue.DB(GetDbName(ctx))
	results, err := dbObj.Query(ctx.Context(), compare, ctx.Request().Query().Values())
	if err != nil {
		return err
	}
	return results
}
