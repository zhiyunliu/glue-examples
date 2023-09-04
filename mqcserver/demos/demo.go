package demos

import (
	"database/sql"

	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
	"github.com/zhiyunliu/glue/xdb"
)

type Orgdemo struct{}

func (d *Orgdemo) Handle(ctx context.Context) interface{} {

	// time.Sleep(time.Second * 5)

	// ctx.Log().Infof("mqc.demo:%s", time.Now().Format("2006-01-02 15:04:05"))

	// ctx.Log().Infof("header.a:%+v", ctx.Request().GetHeader("a"))
	// time.Sleep(time.Millisecond * 200)
	// ctx.Log().Infof("header.b:%+v", ctx.Request().GetHeader("b"))
	// time.Sleep(time.Millisecond * 200)

	// ctx.Log().Infof("header.c:%+v", ctx.Request().GetHeader("c"))
	// time.Sleep(time.Millisecond * 200)

	// ctx.Log().Infof("body-1:%s", ctx.Request().Body().Bytes())

	// mapData := map[string]string{}
	// ctx.Request().Body().Scan(&mapData)
	// ctx.Log().Infof("body-2:%+v", mapData)

	dbObj := glue.DB("microsql")

	rowData, err := dbObj.First(ctx.Context(), "select 2", nil)
	if err != nil {
		ctx.Log().Error(1, err)
	}
	ctx.Log().Info(rowData)
	err = dbObj.Transaction(func(dbObj xdb.Executer) error {
		var outaaa int
		args := map[string]interface{}{
			"a":   1001,
			"aaa": sql.Named("aaa", sql.Out{Dest: &outaaa}),
		}
		trowData, err := dbObj.First(ctx.Context(), `
		select  1 
			
		`, args)
		if err != nil {
			if dberr, ok := err.(xdb.DbError); ok {
				ctx.Log().Error(dberr.SQL())
				ctx.Log().Error(dberr.Args()...)
			}

			return err
		}
		ctx.Log().Info(3, trowData)
		return nil
	})
	if err != nil {
		ctx.Log().Error(2, err)
	}
	return "success"
}
