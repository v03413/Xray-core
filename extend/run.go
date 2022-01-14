package extend

import (
	"fmt"
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
	var interval = getC("extend.interval").Int()
	if interval == 0 {

		interval = 30 // 默认三十秒
	}

	for i := 0; ; i++ {
		getAccounts()
		uploadLog()
		cacheCidOfUser.DeleteExpired()

		time.Sleep(time.Second * time.Duration(interval))
	}
}
