package extend

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
	"time"
)

var pList = cache.New(time.Minute, time.Minute)
var recent = cache.New(time.Minute*3, time.Second)

type Proxy struct {
	rateLimit rate.Limit
	password  string
}

func Auth(username, password, srcIp string) bool {
	itm, found := pList.Get(username)
	if found && password == itm.(Proxy).password {
		recent.Set(username, true, time.Minute)
		onlineLogChan <- fmt.Sprintf("%s:%s", username, srcIp)

		return true
	}

	return false
}

func getProxy(username string) (Proxy, bool) {
	itm, found := pList.Get(username)
	if !found {

		return Proxy{}, found
	}

	return itm.(Proxy), found
}
