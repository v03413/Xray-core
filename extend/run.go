package extend

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"strings"
	"time"
)

func Start(configFile string) {
	err := loadConfig(configFile)
	if err != nil {
		Error(fmt.Sprintf("配置文件加载失败：[%s]%s", configFile, err.Error()))
		return
	}

	go run()
}

func run() {
	var interval = getC("interval").Int()
	if interval == 0 {

		interval = 30 // 默认三十秒
	}

	for i := 0; ; i++ {
		// 先传日志，再获取账号，顺序不能错
		heartbeat()

		time.Sleep(time.Second * time.Duration(interval))
	}
}

func heartbeat() {
	var post = getPostLog()
	var apiUrl = getC("api").String()

	resp, err := http.Post(apiUrl, "application/json", strings.NewReader(post))
	if err != nil {

		Warning("心跳错误：" + err.Error())
		return
	}

	bytes, _ := io.ReadAll(resp.Body)
	text := string(bytes)

	if gjson.Get(text, "code").String() != "0" {

		Warning(fmt.Sprintf("心跳错误：状态码code必须为0（%s）", text))
		return
	}

	// 清空旧列表
	pList.Flush()
	total := int(gjson.Get(text, "total").Num)
	for i := 0; i < total; i++ {
		user := gjson.Get(text, fmt.Sprintf("list.%d.user", i)).String()
		pass := gjson.Get(text, fmt.Sprintf("list.%d.pass", i)).String()
		rateLimit := gjson.Get(text, fmt.Sprintf("list.%d.rate", i)).Float()
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

	Warning(fmt.Sprintf("代理数:%d 连接数:%d 「%s」", total, cacheCidOfUser.ItemCount(), gjson.Get(text, "msg")))
}
