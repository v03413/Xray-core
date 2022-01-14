package extend

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var cacheCidOfUser = cache.New(3*time.Minute, 10*time.Minute)

func SetCid(cid string, v interface{}) {

	cacheCidOfUser.Set(cid, v, cache.NoExpiration)
}

func DelCid(cid string) {

	cacheCidOfUser.Delete(cid)
}

func GetAccountByCid(cid string) (interface{}, bool) {

	return cacheCidOfUser.Get(cid)
}
