package extend

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/xtls/xray-core/common/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var logs chan string

func uploadLog() {
	var result = make(map[string]bool)
	for {
		if len(logs) == 0 {
			break
		}

		result[<-logs] = true
	}

	if len(result) == 0 {

		return
	}

	var post string
	for k, _ := range result {
		post += k + ","
	}

	var url = fmt.Sprintf("%sapi.php?act=upload_log", getC("extend.api"))

	resp, err := http.Post(url, "application/json", strings.NewReader(post[:len(post)-1]))
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
