package config

import (
	"io/ioutil"

	"github.com/tidwall/gjson"
)

var configuration gjson.Result

func Init() {
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	configuration = gjson.ParseBytes(content).Get("backup")
}

func GetDBCreds() (string, string) {
	uname := configuration.Get("db.username").String()
	pwd := configuration.Get("db.password").String()

	return uname, pwd
}

func Get(key string) gjson.Result {
	return configuration.Get(key)
}
