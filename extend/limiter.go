package extend

import (
	"context"
	"golang.org/x/time/rate"
	"sync"
)

type limit struct {
	size    rate.Limit
	limiter *rate.Limiter
}

var limiters = &sync.Map{}

func WaitForCid(cid string) {
	if username, ok := GetUsernameByCid(cid); ok {
		itm, _ := limiters.Load(username)

		itm.(*limit).limiter.Wait(context.Background())
	}
}

func setProxyRate(account string, size rate.Limit) {
	itm, found := limiters.Load(account)
	if !found || itm.(*limit).size != size {
		limiters.Store(account, &limit{
			size:    size,
			limiter: rate.NewLimiter(size, int(size)),
		})

		return
	}

	// 速率未发生变化
}

func deleteExpireLimiter() {
	limiters.Range(func(key, value interface{}) bool {
		_, ok := getProxy(key.(string))
		if !ok {
			// 删除失效数据

			limiters.Delete(key)
		}

		return true
	})
}
