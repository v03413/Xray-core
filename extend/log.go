package extend

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/xtls/xray-core/common/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var logs chan string

func uploadLog() {
	var result []string
	var traffic []string

	// 各IP流量统计
	trafficMap := make(map[string]int)
	for len(TrafficLogChan) > 0 {
		tmp := strings.Split(<-TrafficLogChan, "|")
		if tmp[1] == "0" {

			continue
		}

		total, _ := strconv.Atoi(tmp[1])
		if v, ok := trafficMap[tmp[0]]; ok {
			trafficMap[tmp[0]] = total + v
		} else {
			trafficMap[tmp[0]] = total
		}
	}

	// 汇总流量数据
	for ip, v := range trafficMap {

		traffic = append(traffic, fmt.Sprintf("%s:%d", ip, v))
	}

	// 账号上线IP汇总
	for len(logs) != 0 {

		result = append(result, <-logs)
	}

	var unique = elementUnique(result) // 去重
	var post = fmt.Sprintf(`{"online":"%s","traffic":"%s"}`, strings.Join(unique, ","), strings.Join(traffic, ","))
	var url = fmt.Sprintf("%sapi.php?act=upload_log&v=2", getC("extend.api"))

	resp, err := http.Post(url, "application/json", strings.NewReader(post))
	if err != nil {
		Error("日志上报错误：" + err.Error())

		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		Warning("日志上报错误：", err.Error())
	} else {

		Warning("日志上报成功：", gjson.Get(string(data), "msg"))
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

// 元素去重
func elementUnique(arr []string) (newArr []string) {
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
