package exprs

import (
	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
)

type Like struct {
}

func (c *Like) Handle(ctx context.Context) (res any) {
	const compare = `
	select * from ljy_test t where t.a=@{a} &{like b} &{like %t.c} |{t.d like %d%}
	`
	dbObj := glue.DB(GetDbName(ctx))
	results, err := dbObj.Query(ctx.Context(), compare, ctx.Request().Query().Values())
	if err != nil {
		return err
	}
	return results
}
