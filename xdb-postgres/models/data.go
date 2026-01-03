package models

type QueryParams struct {
	ID      int64  `json:"id" form:"id" `
	Name    string `json:"name" form:"name"`
	Address string `json:"address" form:"address"`
	Age     int    `json:"age" form:"age"`
	Ids     []int  `json:"ids" form:"ids"`
}

type QueryResultRow struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}
