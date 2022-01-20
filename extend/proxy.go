package extend

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
	"time"
)

var pList = cache.New(10*time.Minute, 2*time.Minute)

type Proxy struct {
	rateLimit rate.Limit
	password  string
}

func Auth(username, password, srcIp, cid string) bool {
	itm, found := pList.Get(username)
	if found && password == itm.(Proxy).password {
		// 记录上线日志
		onlineLogChan <- fmt.Sprintf("%s:%s", username, srcIp)

		// 关联连接ID
		setCid(cid, username)

		return true
	}

	return false
}

func IsExistAccount(username string) bool {
	_, found := pList.Get(username)

	return found
}

func getProxy(username string) (Proxy, bool) {
	itm, found := pList.Get(username)
	if !found {

		return Proxy{}, found
	}

	return itm.(Proxy), found
}
