package extend

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const Version = "1.2"

var apiUrl, listenPort string

func Start(configFile string) {
	err := loadConfig(configFile)
	if err != nil {
		Error(fmt.Sprintf("配置文件加载失败：[%s]%s", configFile, err.Error()))
		return
	}

	apiUrl = getC("extend.api").String()
	listenPort = getC("inbounds.0.port").String()

	log.Println("[监听端口] ==> ", listenPort)

	go run()
}

func run() {
	for range time.Tick(time.Second * 5) {

		heartbeat()
	}
}

func heartbeat() {
	var post = getStates()

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

	// 删除失效限速器
	deleteExpireLimiter()

	Warning(fmt.Sprintf("{%s}代理数:%s 「%s」", Version, total, msg))
}
