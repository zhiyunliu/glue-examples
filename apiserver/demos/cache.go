package demos

import (
	"strconv"
	"time"

	gel "github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/cache"
	"github.com/zhiyunliu/glue/context"
)

type Cachedemo struct{}

func NewCache() *Cachedemo {
	return &Cachedemo{}
}

func (d *Cachedemo) GetHandle(ctx context.Context) interface{} {
	tmpIdx := ctx.Request().Query().Get("idx")
	idx, _ := strconv.ParseInt(tmpIdx, 10, 32)

	cacheObj := gel.Cache("default", cache.WithDBIndex(int(idx)))

	sctx := ctx.Context()
	val, err := cacheObj.Get(sctx, "key")
	return map[string]interface{}{
		"val": val,
		"err": err,
	}
}

func (d *Cachedemo) SetHandle(ctx context.Context) interface{} {

	tmpIdx := ctx.Request().Query().Get("idx")
	idx, _ := strconv.ParseInt(tmpIdx, 10, 32)

	cacheObj := gel.Cache("default", cache.WithDBIndex(int(idx)))
	err := cacheObj.Set(ctx.Context(), "key", time.Now().Nanosecond(), 10)

	return map[string]interface{}{
		"err": err,
	}
}

func (d *Cachedemo) DelHandle(ctx context.Context) interface{} {
	cacheObj := gel.Cache("default")
	err := cacheObj.Del(ctx.Context(), "key")
	return map[string]interface{}{
		"err": err,
	}
}

func (d *Cachedemo) HgetHandle(ctx context.Context) interface{} {
	cacheObj := gel.Cache("default")
	val, err := cacheObj.HashGet(ctx.Context(), "hash", "key")
	return map[string]interface{}{
		"val": val,
		"err": err,
	}
}

func (d *Cachedemo) HSetHandle(ctx context.Context) interface{} {
	cacheObj := gel.Cache("default")
	val, err := cacheObj.HashSet(ctx.Context(), "hash", "key", time.Now().GoString())
	return map[string]interface{}{
		"val": val,
		"err": err,
	}
}

func (d *Cachedemo) IncreaseHandle(ctx context.Context) interface{} {
	cacheObj := gel.Cache("default")
	val, err := cacheObj.Increase(ctx.Context(), "increase")
	return map[string]interface{}{
		"val": val,
		"err": err,
	}
}

func (d *Cachedemo) DecreaseHandle(ctx context.Context) interface{} {
	cacheObj := gel.Cache("default")
	val, err := cacheObj.Decrease(ctx.Context(), "increase")
	return map[string]interface{}{
		"val": val,
		"err": err,
	}
}

func (d *Cachedemo) ExpireHandle(ctx context.Context) interface{} {
	cacheObj := gel.Cache("default")
	err := cacheObj.Set(ctx.Context(), "expire", 10, -1)
	err = cacheObj.Expire(ctx.Context(), "expire", 10)
	return map[string]interface{}{
		"err": err,
	}
}

// Get(key string) (string, error)
// Set(key string, val interface{}, expire int) error
// Del(key string) error
// HashGet(hk, key string) (string, error)
// HashDel(hk, key string) error
// Increase(key string) error
// Decrease(key string) error
// Expire(key string, expire int) error
// GetImpl() interface{}
