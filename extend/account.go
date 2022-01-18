package extend

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"
	"golang.org/x/time/rate"
	"io"
	"net/http"
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
		SetCid(cid, username)

		return true
	}

	return false
}

func IsExistAccount(username string) bool {
	_, found := pList.Get(username)

	return found
}

func getAccounts() {
	var url = fmt.Sprintf("%sapi.php?act=get_user&token=%s&sid=%s", getC("api"), getC("key"), getC("sid"))

	resp, err := http.Get(url)
	if err != nil {
		Warning("账号获取失败：" + err.Error())
		return
	}

	bytes, _ := io.ReadAll(resp.Body)
	text := string(bytes)

	if gjson.Get(text, "code").String() != "0" {

		return
	}

	// 清空旧列表
	pList.Flush()
	total := int(gjson.Get(text, "total").Num)
	for i := 0; i < total; i++ {
		user := gjson.Get(text, fmt.Sprintf("accounts.%d.user", i)).String()
		pass := gjson.Get(text, fmt.Sprintf("accounts.%d.pass", i)).String()
		rateLimit := gjson.Get(text, fmt.Sprintf("accounts.%d.limit", i)).Float()
		if rateLimit == 0 { // 不限速

			rateLimit = 999
		}

		pList.Set(user, Proxy{password: pass, rateLimit: rate.Limit(rateLimit)}, cache.NoExpiration)

		setProxyRate(user, rate.Limit(rateLimit))
	}

	// 删除失效连接
	deleteExpireCid()

	// 删除失效限速器
	deleteExpireLimiter()

	Warning(fmt.Sprintf("账号数量：%d 连接数量：%d", total, cacheCidOfUser.ItemCount()))
}

func getProxy(username string) (Proxy, bool) {
	itm, found := pList.Get(username)
	if !found {

		return Proxy{}, found
	}

	return itm.(Proxy), found
}
