package exprs

import "github.com/zhiyunliu/glue/context"

func GetDbName(ctx context.Context) string {
	dbName := ctx.Request().Query().Get("db_name")
	if dbName == "" {
		dbName = "xdb-mssql"
	}
	return dbName
}
