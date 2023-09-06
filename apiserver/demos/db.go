package demos

import (
	sctx "context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
	"github.com/zhiyunliu/glue/log"
	"github.com/zhiyunliu/glue/xdb"
)

type DBdemo struct{}

func NewDb() *DBdemo {
	return &DBdemo{}
}

func init() {
	xdb.RegistryLogger(&dblogger{})
}

type dblogger struct{}

func (l *dblogger) Name() string {
	return "default"
}

func (l *dblogger) Log(ctx sctx.Context, elapsed int64, sql string, args ...interface{}) {
	log.DefaultLogger.Warn(args...)
}

func (d *DBdemo) And1Handle(ctx context.Context) interface{} {
	dbobj := glue.DB("microsql")
	idval := ctx.Request().Query().Get("id")
	sql := `select * from ljy_test where ctime > '2022-1-1' &{status}`

	rows, err := dbobj.Query(ctx.Context(), sql, map[string]interface{}{
		"id":     idval,
		"status": ctx.Request().Query().Get("status"),
	})
	if err != nil {
		ctx.Log().Error(err)
	}
	return rows
}

func (d *DBdemo) And2Handle(ctx context.Context) interface{} {
	dbobj := glue.DB("microsql")
	idval := ctx.Request().Query().Get("id")
	sql := `select * from ljy_test t where ctime > '2022-1-1' &{t.status}`

	rows, err := dbobj.Query(ctx.Context(), sql, map[string]interface{}{
		"id":     idval,
		"status": ctx.Request().Query().Get("status"),
	})
	if err != nil {
		ctx.Log().Error(err)
	}
	return rows
}

func (d *DBdemo) Or1Handle(ctx context.Context) interface{} {
	dbobj := glue.DB("microsql")
	idval := ctx.Request().Query().Get("id")
	sql := `select * from ljy_test where id=1 |{id}`

	rows, err := dbobj.Query(ctx.Context(), sql, map[string]interface{}{
		"id": idval,
	})
	if err != nil {
		ctx.Log().Error(err)
	}

	return rows
}
func (d *DBdemo) Or2Handle(ctx context.Context) interface{} {
	dbobj := glue.DB("microsql")
	idval := ctx.Request().Query().Get("id")
	sql := `select * from ljy_test  t where id=1 |{t.id}`

	rows, err := dbobj.Query(ctx.Context(), sql, map[string]interface{}{
		"id": idval,
	})
	if err != nil {
		ctx.Log().Error(err)
	}

	return rows
}

func (d *DBdemo) QueryHandle(ctx context.Context) interface{} {
	dbobj := glue.DB("microsql")
	idval := ctx.Request().Query().Get("id")
	sql := `select * from ljy_test`

	rows, err := dbobj.Query(ctx.Context(), sql, map[string]interface{}{
		"id": idval,
	})
	if err != nil {
		ctx.Log().Error(err)
	}

	return rows

}

func (d *DBdemo) FirstHandle(ctx context.Context) interface{} {
	dbobj := glue.DB("localhost")
	row, err := dbobj.First(ctx.Context(), "select id,name from new_table t where t.id=@id", map[string]interface{}{
		"id": ctx.Request().Query().Get("id"),
	})
	if err != nil {
		ctx.Log().Error(err)
	}

	return row

}

func (d *DBdemo) InsertHandle(ctx context.Context) interface{} {
	dbobj := glue.DB("localhost")
	result, err := dbobj.Exec(ctx.Context(), "insert into new_table(name) values(@name) ", map[string]interface{}{
		"name": fmt.Sprintf("insert:%s:%s", ctx.Request().Query().Get("name"), time.Now().Format("2006-01-02 15:04:05")),
	})
	if err != nil {
		ctx.Log().Error(err)
	}

	lastId, err1 := result.LastInsertId()
	effCnt, err2 := result.RowsAffected()
	return map[string]interface{}{
		"LastInsertId": lastId,
		"RowsAffected": effCnt,
		"Error1":       err1,
		"Error2":       err2,
	}
}

func (d *DBdemo) UpdateHandle(ctx context.Context) interface{} {
	dbobj := glue.DB("localhost")
	result, err := dbobj.Exec(ctx.Context(), "update new_table set name=@name where id=@id ", map[string]interface{}{
		"id":   ctx.Request().Query().Get("id"),
		"name": fmt.Sprintf("update:%s", time.Now().Format("2006-01-02 15:04:05")),
	})
	if err != nil {
		ctx.Log().Error(err)
	}
	lastId, err1 := result.LastInsertId()
	effCnt, err2 := result.RowsAffected()
	return map[string]interface{}{
		"LastInsertId": lastId,
		"RowsAffected": effCnt,
		"Error1":       err1,
		"Error2":       err2,
	}
}

func (d *DBdemo) TransHandle(ctx context.Context) interface{} {
	dbobj := glue.DB("localhost")

	trans, err := dbobj.Begin()
	if err != nil {
		return err
	}

	istResult, err := trans.Exec(ctx.Context(), "insert into new_table(name) values(@name) ", map[string]interface{}{
		"name": "trans insert",
	})

	if err != nil {
		trans.Rollback()
		return err
	}

	lastId, err := istResult.LastInsertId()
	if err != nil {
		trans.Rollback()
		return err
	}

	result, err := trans.Exec(ctx.Context(), "update new_table set name=@name where id=@id ", map[string]interface{}{
		"id":   lastId,
		"name": fmt.Sprintf("update-trans:%s", time.Now().Format("2006-01-02 15:04:05")),
	})
	if err != nil {
		trans.Rollback()
		return err
	}
	trans.Commit()

	uefcnt, err := result.RowsAffected()

	return map[string]interface{}{
		"insertid": lastId,
		"uefcnt":   uefcnt,
		"err":      err,
	}

}

func (d *DBdemo) MultiHandle(ctx context.Context) interface{} {
	dbobj := glue.DB("microsql")

	var outArg string
	result, err := dbobj.Multi(ctx.Context(), `
DECLARE	@return_value int

EXEC @return_value = [dbo].[test_aaa]
	@id = @{id},
	@name = @{name} OUTPUT

	`, map[string]interface{}{
		"id":   ctx.Request().Query().Get("id"),
		"name": sql.Named("name", sql.Out{Dest: &outArg}),
	})
	if err != nil {
		ctx.Log().Error(err)
	}

	ctx.Log().Debug("outArg:", outArg)

	return result
}

func (d *DBdemo) SpHandle(ctx context.Context) interface{} {
	dbobj := glue.DB("localhost")
	result, err := dbobj.Exec(ctx.Context(), "update new_table set name=@name where id=@id ", map[string]interface{}{
		"id":   ctx.Request().Query().Get("id"),
		"name": fmt.Sprintf("update:%s", time.Now().Format("2006-01-02 15:04:05")),
	})
	if err != nil {
		ctx.Log().Error(err)
	}
	lastId, err1 := result.LastInsertId()
	effCnt, err2 := result.RowsAffected()
	return map[string]interface{}{
		"LastInsertId": lastId,
		"RowsAffected": effCnt,
		"Error1":       err1,
		"Error2":       err2,
	}
}

func (d *DBdemo) StructHandle(ctx context.Context) interface{} {
	dbobj := glue.DB("microsql")
	p := struct {
		Name   string     `json:"name" form:"name"`
		Sleep  int        `json:"sleep" form:"sleep"`
		Status *int       `json:"status" form:"status"`
		Ctime  time.Time  `json:"time" form:"time" time_format:"2006-01-02 15:04:05"`
		PCtime *time.Time `json:"ptime" form:"ptime" time_format:"2006-01-02 15:04:05"`
	}{}

	ctx.Bind(&p)

	result, err := dbobj.Query(ctx.Context(), `
	if @{sleep} > 0 
	begin 
	  waitfor delay '00:00:${sleep}:00'
	end  

	SELECT   [id]
	,[name]
	,[status]
	,[x]
FROM [dbo].[ljy_test] t
where  name=@{name}  &{t.status} &{t.ctime}	
	`, map[string]interface{}{
		"sleep":  p.Sleep,
		"status": p.Status,
		"name":   p.Name,
		"ctime":  p.PCtime,
	})
	if err != nil {
		if dberr, ok := err.(xdb.DbError); ok {
			ctx.Log().Error(err.Error(), dberr.SQL(), dberr.Args())
			return result
		}
		ctx.Log().Error(err)
	}
	return result
}
