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

func (c *Output) OutParamHandle(ctx context.Context) (res any) {
	type param struct {
		A         string `json:"a"`
		OutValue1 string `xdb:"outputvalue1,dbtype:output"`
		OutValue2 string `json:"outputvalue2,dbtype:output"`
	}

	const compare = `
		declare @c varchar(50) = newid()
		set @{outputvalue1}=@c
		set @{outputvalue2}=@c
		`

	dataParam := param{
		A: ctx.Request().Query().Values().Get("a"),
	}

	dbObj := glue.DB(GetDbName(ctx))
	results, err := dbObj.Query(ctx.Context(), compare, &dataParam, xdb.WithExprCache(false))
	if err != nil {
		return err
	}

	return map[string]any{
		"out1": dataParam.OutValue1,
		"out2": dataParam.OutValue2,
		"data": results,
	}
}
