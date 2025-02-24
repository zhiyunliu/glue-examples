package demos

import (
	"mime/multipart"

	"github.com/zhiyunliu/glue/context"
)

type MultiDemo struct {
}

type DataItem struct {
	A        string                `json:"a" form:"a"`
	B        int                   `json:"b" form:"b"`
	TestFile *multipart.FileHeader `json:"testfile" form:"testfile"`
}

func (d *DataItem) MaxMemory() int64 {
	return 0
}

func (d *MultiDemo) Handle(ctx context.Context) interface{} {
	var item = &DataItem{}

	err := ctx.Bind(&item)
	if err != nil {
		return err
	}

	return item
}
