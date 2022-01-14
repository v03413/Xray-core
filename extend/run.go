package extend

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"time"
)

var userList = cache.New(10*time.Minute, 2*time.Minute)
var TrafficLogChan = make(chan string, 1000000)
var CacheUuidOfUser = cache.New(3*time.Minute, 10*time.Minute)

func Start(configFile string) {
	logs = make(chan string, 10240)
	err := loadConfig(configFile)
	if err != nil {
		Error(fmt.Sprintf("配置文件加载失败：[%s]%s", configFile, err.Error()))
		return
	}

	go run()
}

// Auth 账号授权验证
func Auth(account, password, srcIp, cid string) bool {
	storedPassed, found := userList.Get(account)
	if found && password == storedPassed {
		// 记录上线日志
		logs <- fmt.Sprintf("%s:%s", account, srcIp)

		// 关联连接ID
		CacheUuidOfUser.Add(cid, account, time.Minute*3)

		return true
	}

	return false
}

// IsExistAccount 账号是否存在
func IsExistAccount(account string) bool {
	_, found := userList.Get(account)

	return found
}

func run() {
	var interval = getC("extend.interval").Int()
	if interval == 0 {

		interval = 30 // 默认三十秒
	}

	for i := 0; ; i++ {
		getAccounts()
		uploadLog()
		CacheUuidOfUser.DeleteExpired()

		time.Sleep(time.Second * time.Duration(interval))
	}
}

func getAccounts() {
	var url = fmt.Sprintf("%sapi.php?act=get_user&token=%s&sid=%s", getC("extend.api"), getC("extend.key"), getC("extend.sid"))

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
	userList.Flush()
	total := int(gjson.Get(text, "total").Num)
	for i := 0; i < total; i++ {
		user := gjson.Get(text, fmt.Sprintf("accounts.%d.user", i)).String()
		pass := gjson.Get(text, fmt.Sprintf("accounts.%d.pass", i)).String()

		userList.Set(user, pass, cache.NoExpiration)
	}

	Warning(fmt.Sprintf("账号获取成功：%d 当前连接数：%d", total, CacheUuidOfUser.ItemCount()))
}
