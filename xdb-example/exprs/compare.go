package exprs

import (
	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
)

type Compare struct {
}

type param struct {
	A string `form:"a" json:"a,dbtype:varchar"`
	B string `form:"b" json:"b"`
	C string `form:"c" json:"c"`
	D string `form:"d" json:"d" xdb:"d,dbtype:varchar"`
}

func (c *Compare) Handle(ctx context.Context) (res any) {

	p := &param{}
	if err := ctx.Bind(p); err != nil {
		return err
	}

	const compare = `
	select * from ljy_test t where t.a=@{a} &{b>b} &{>t.b} &{t.c>=c} |{t.d=d} &{t.a>b}
	`

	dbObj := glue.DB(GetDbName(ctx))
	results, err := dbObj.Query(ctx.Context(), compare, p)
	if err != nil {
		return err
	}
	return []any{p, results}
}


