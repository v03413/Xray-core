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

func setCid(cid string, v interface{}) {

	cacheCidOfUser.Set(cid, v, cache.NoExpiration)
}

func DelCid(cid string) {

	cacheCidOfUser.Delete(cid)
}

func GetUsernameByCid(cid string) (interface{}, bool) {

	return cacheCidOfUser.Get(cid)
}

func getOnlineNum() int {
	var arr = make([]string, 100)
	var total = 0

	for _, v := range cacheCidOfUser.Items() {

		arr = append(arr, v.Object.(string))
	}

	arr = arrUnique(arr)
	for _, v := range arr {
		if v != "" {
			total++
		}

	}

	return total
}
