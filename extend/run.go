package extend

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"strconv"
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

	for range time.Tick(time.Second * time.Duration(interval)) {
		// 开始心跳

		heartbeat()
	}
}

func heartbeat() {
	var post = getStates()
	var apiUrl = getC("api").String()

	resp, err := http.Post(apiUrl, "application/json", strings.NewReader(post))
	if err != nil {

		Warning("心跳错误：" + err.Error())
		return
	}

	bytes, _ := io.ReadAll(resp.Body)
	text := string(bytes)

	var data = strings.Split(text, "\n")
	var code = data[0]
	var msg = data[1]
	var total = data[2]
	var list = strings.Split(data[3], "|")

	if code != "0" {

		Warning(fmt.Sprintf("心跳错误：状态码code必须为0（%s）", text))
		return
	}

	// 清空旧列表
	pList.Flush()
	for _, v := range list {
		var tmp = strings.Split(v, ":")

		user := tmp[0]
		pass := tmp[1]
		rateLimit, _ := strconv.ParseFloat(tmp[2], 64)
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

	Warning(fmt.Sprintf("代理数:%s 连接数:%d 「%s」", total, cacheCidOfUser.ItemCount(), msg))
}
