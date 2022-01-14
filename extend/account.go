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

// Auth 账号授权验证
func Auth(account, password, srcIp, cid string) bool {
	storedPassed, found := userList.Get(account)
	if found && password == storedPassed {
		// 记录上线日志
		onlineLogChan <- fmt.Sprintf("%s:%s", account, srcIp)

		// 关联连接ID
		SetCid(cid, account)

		return true
	}

	return false
}

// IsExistAccount 账号是否存在
func IsExistAccount(account string) bool {
	_, found := userList.Get(account)

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
	userList.Flush()
	total := int(gjson.Get(text, "total").Num)
	for i := 0; i < total; i++ {
		user := gjson.Get(text, fmt.Sprintf("accounts.%d.user", i)).String()
		pass := gjson.Get(text, fmt.Sprintf("accounts.%d.pass", i)).String()

		userList.Set(user, pass, cache.NoExpiration)
	}

	Warning(fmt.Sprintf("账号数量：%d 连接数量：%d", total, cacheCidOfUser.ItemCount()))
}
