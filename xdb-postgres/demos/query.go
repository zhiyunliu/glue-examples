package demos

import (
	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue-examples/xdb-postgres/demos/sqls"
	"github.com/zhiyunliu/glue-examples/xdb-postgres/models"
	"github.com/zhiyunliu/glue/context"
)

func QueryData(ctx context.Context) any {
	params := &models.QueryParams{}

	if err := ctx.Bind(params); err != nil {
		return err
	}
	result := make([]models.QueryResultRow, 0)

	dbObj := glue.DB("xdb-postgres")
	err := dbObj.QueryAs(ctx.Context(), sqls.QueryData, params, &result)
	if err != nil {
		return err
	}
	return result
}
func LikeData(ctx context.Context) any {
	params := &models.QueryParams{}

	if err := ctx.Bind(params); err != nil {
		return err
	}
	result := make([]models.QueryResultRow, 0)

	dbObj := glue.DB("xdb-postgres")
	err := dbObj.QueryAs(ctx.Context(), sqls.QueryLikeData, params, &result)
	if err != nil {
		return err
	}
	return result
}
func Like2Data(ctx context.Context) any {
	params := &models.QueryParams{}

	if err := ctx.Bind(params); err != nil {
		return err
	}
	result := make([]models.QueryResultRow, 0)

	dbObj := glue.DB("xdb-postgres")
	err := dbObj.QueryAs(ctx.Context(), sqls.QueryLike2Data, params, &result)
	if err != nil {
		return err
	}
	return result
}
func InData(ctx context.Context) any {
	params := &models.QueryParams{}

	if err := ctx.Bind(params); err != nil {
		return err
	}
	result := make([]models.QueryResultRow, 0)

	dbObj := glue.DB("xdb-postgres")
	err := dbObj.QueryAs(ctx.Context(), sqls.QueryInData, params, &result)
	if err != nil {
		return err
	}
	return result
}

func NotInData(ctx context.Context) any {
	params := &models.QueryParams{}

	if err := ctx.Bind(params); err != nil {
		return err
	}
	result := make([]models.QueryResultRow, 0)

	dbObj := glue.DB("xdb-postgres")
	err := dbObj.QueryAs(ctx.Context(), sqls.QueryNotInData, params, &result)
	if err != nil {
		return err
	}
	return result
}

func NotIn2Data(ctx context.Context) any {
	params := &models.QueryParams{}

	if err := ctx.Bind(params); err != nil {
		return err
	}
	result := make([]models.QueryResultRow, 0)

	dbObj := glue.DB("xdb-postgres")
	err := dbObj.QueryAs(ctx.Context(), sqls.QueryNotIn2Data, params, &result)
	if err != nil {
		return err
	}
	return result
}

func CompareData(ctx context.Context) any {
	params := &models.QueryParams{}

	if err := ctx.Bind(params); err != nil {
		return err
	}
	result := make([]models.QueryResultRow, 0)

	dbObj := glue.DB("xdb-postgres")
	err := dbObj.QueryAs(ctx.Context(), sqls.QueryCompareData, params, &result)
	if err != nil {
		return err
	}
	return result
}

func Compare2Data(ctx context.Context) any {
	params := &models.QueryParams{}

	if err := ctx.Bind(params); err != nil {
		return err
	}
	result := make([]models.QueryResultRow, 0)

	dbObj := glue.DB("xdb-postgres")
	err := dbObj.QueryAs(ctx.Context(), sqls.QueryCompare2Data, params, &result)
	if err != nil {
		return err
	}
	return result
}

func CteData(ctx context.Context) any {
	params := &models.QueryParams{}

	if err := ctx.Bind(params); err != nil {
		return err
	}
	result := make([]models.QueryResultRow, 0)

	dbObj := glue.DB("xdb-postgres")
	err := dbObj.QueryAs(ctx.Context(), sqls.QueryCTE, params, &result)
	if err != nil {
		return err
	}
	return result
}
