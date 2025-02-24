package exprs

import (
	"strings"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
	"github.com/zhiyunliu/golibs/xtypes"
)

type In struct {
}

func (c *In) Handle(ctx context.Context) (res any) {
	const compare = `
	select * from ljy_test t where t.a=@{a} &{in b} &{in t.c} &{t.d in d}
	`

	dataMap := xtypes.XMap{}

	for k, v := range ctx.Request().Query().Values() {
		dataMap[k] = v
		vals := strings.Split(v, ",")
		if len(vals) > 1 {
			dataMap[k] = vals
		}
	}

	dataMap["a"] = mssql.VarChar(dataMap.GetString("a"))

	dbObj := glue.DB(GetDbName(ctx))
	results, err := dbObj.Query(ctx.Context(), compare, dataMap)
	if err != nil {
		return err
	}
	return results
}
