package extend

import (
	"github.com/tidwall/gjson"
	"os"
)

var config string

func loadConfig(configFile string) error {
	bytes, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	config = string(bytes)

	return nil
}

func getC(path string) gjson.Result {

	return gjson.Get(config, path)
}
