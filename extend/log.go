package extend

import (
	"fmt"
	"github.com/xtls/xray-core/common/errors"
	"strings"
)

type trafficLog struct {
	username string
	total    int32
}

var onlineLogChan = make(chan string, 100000)
var trafficLogChan = make(chan trafficLog, 100000)

func getStates() string {
	var online []string
	var traffic []string

	var serverId = getC("id").Int()
	var onlineNum = getOnlineNum()

	// 各用户流量统计
	trafficMap := make(map[interface{}]int32)
	for len(trafficLogChan) > 0 {
		var log = <-trafficLogChan
		if log.total == 0 {

			continue
		}

		if v, ok := trafficMap[log.username]; ok {
			trafficMap[log.username] = log.total + v
		} else {
			trafficMap[log.username] = log.total
		}
	}

	// 汇总流量数据
	for username, v := range trafficMap {

		traffic = append(traffic, fmt.Sprintf("%s:%d", username, v))
	}

	// 账号上线IP汇总
	for len(onlineLogChan) != 0 {

		online = append(online, <-onlineLogChan)
	}

	return fmt.Sprintf(`{"id":%d,"total":%d,"log":{"online":"%s","traffic":"%s"}}`, serverId, onlineNum, strings.Join(arrUnique(online), ","), strings.Join(traffic, ","))
}

func PushTrafficLog(cid string, total int32) {
	if username, ok := GetUsernameByCid(cid); ok {
		trafficLogChan <- trafficLog{
			username: username.(string),
			total:    total,
		}
	}
}

func Warning(values ...interface{}) {

	errors.New(values).AtWarning().WriteToLog()
}

func Error(values ...interface{}) {

	errors.New(values).AtError().WriteToLog()
}

func Info(values ...interface{}) {

	errors.New(values).AtInfo().WriteToLog()
}

func arrUnique(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {

			newArr = append(newArr, arr[i])
		}
	}

	return
}
