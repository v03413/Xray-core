package extend

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/xtls/xray-core/common/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type trafficLog struct {
	username string
	total    int32
}

var onlineLogChan = make(chan string, 100000)
var trafficLogChan = make(chan trafficLog, 100000)

func uploadLog() {
	var online []string
	var traffic []string

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

	var unique = arrUnique(online)
	var post = fmt.Sprintf(`{"online":"%s","traffic":"%s"}`, strings.Join(unique, ","), strings.Join(traffic, ","))
	var url = fmt.Sprintf("%sapi.php?act=upload_log&v=2", getC("api"))

	resp, err := http.Post(url, "application/json", strings.NewReader(post))
	if err != nil {
		Error("日志上报：" + err.Error())

		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		Warning("日志上报：", err.Error())
	} else {

		Warning("日志上报：", gjson.Get(string(data), "msg"))
	}
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
