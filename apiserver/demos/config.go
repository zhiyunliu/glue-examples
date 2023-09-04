package demos

import (
	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
)

type Cfgdemo struct{}

func NewCfg() *Cfgdemo {
	return &Cfgdemo{}
}

func (d *Cfgdemo) Q1Handle(ctx context.Context) interface{} {
	dbobj := glue.DB("microsql")
	idval := ctx.Request().Query().Get("id")
	sql := `select 1 aaa`

	rows, err := dbobj.Query(ctx.Context(), sql, map[string]interface{}{
		"id": idval,
	})
	if err != nil {
		ctx.Log().Error(err)
	}
	return rows
}

func (d *Cfgdemo) Q2Handle(ctx context.Context) interface{} {
	dbobj := glue.DB("xunyou0725")
	idval := ctx.Request().Query().Get("id")
	sql := `select 2 bbb`

	rows, err := dbobj.Query(ctx.Context(), sql, map[string]interface{}{
		"id": idval,
	})
	if err != nil {
		ctx.Log().Error(err)
	}
	return rows
}
