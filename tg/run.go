package tg

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"sync"
	"time"
)

var accounts sync.Map

func Start(configFile string) {
	logs = make(chan string, 10240)
	err := loadConfig(configFile)
	if err != nil {
		Error(fmt.Sprintf("配置文件加载失败：[%s]%s", configFile, err.Error()))
		return
	}

	go run()
}

func Auth(account, password, srcIp string) bool {
	storedPassed, found := accounts.Load(account)
	if found && password == storedPassed {
		// srcIp := inbound.Source.Address.String()

		// 记录上线日志
		logs <- fmt.Sprintf("%s:%s", account, srcIp)

		return true
	}

	return false
}

func run() {
	for i := 0; ; i++ {
		getAccounts()
		uploadLog()

		time.Sleep(time.Second * 30)
	}
}

func getAccounts() {
	var url = fmt.Sprintf("%sapi.php?act=get_user&token=%s&sid=%s", getC("tg.api"), getC("tg.key"), getC("tg.sid"))

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

	// 清空
	accounts.Range(func(k, v interface{}) bool {
		accounts.Delete(k)
		return true
	})
	total := int(gjson.Get(text, "total").Num)
	for i := 0; i < total; i++ {
		user := gjson.Get(text, fmt.Sprintf("accounts.%d.user", i)).String()
		pass := gjson.Get(text, fmt.Sprintf("accounts.%d.pass", i)).String()
		accounts.Store(user, pass)
	}

	Warning(fmt.Sprintf("账号获取成功：%d", total))
}
