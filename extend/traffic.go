package extend

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var cacheCidOfUser = cache.New(3*time.Minute, 10*time.Minute)

func deleteExpireCid() {
	cacheCidOfUser.DeleteExpired()

	for cid, v := range cacheCidOfUser.Items() {
		if !IsExistAccount(v.Object.(string)) {

			DelCid(cid)
		}
	}
}

func SetCid(cid string, v interface{}) {

	cacheCidOfUser.Set(cid, v, cache.NoExpiration)
}

func DelCid(cid string) {

	cacheCidOfUser.Delete(cid)
}

func GetUsernameByCid(cid string) (interface{}, bool) {

	return cacheCidOfUser.Get(cid)
}
