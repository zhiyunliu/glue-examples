package exprs

import (
	"database/sql"

	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
	"github.com/zhiyunliu/glue/xdb"
	"github.com/zhiyunliu/golibs/xtypes"
)

type Output struct {
}

func (c *Output) Handle(ctx context.Context) (res any) {
	const compare = `
	declare @c varchar(50)
		select top 1 @c = t.c from ljy_test t where t.a=@{a}
		set @{outputvalue}=@c
	
	
		`

	var outputValue string
	dataParam := xtypes.XMap{
		"a":           ctx.Request().Query().Values().Get("a"),
		"outputvalue": sql.Named("outputvalue", sql.Out{Dest: &outputValue}),
	}

	dbObj := glue.DB(GetDbName(ctx))
	results, err := dbObj.Query(ctx.Context(), compare, dataParam, xdb.WithExprCache(false))
	if err != nil {
		return err
	}

	return map[string]any{
		"out":  outputValue,
		"data": results,
	}
}
