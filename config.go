package main

import (
	"github.com/Unknwon/goconfig"
)

func GetConfig(path string) (count string, freq string,cal, cmd string) {
	config, err := goconfig.LoadConfigFile(path)
	if err != nil {
		LogErr("can't open config file :%v\n", err)
		return
	}
	ret , err := config.GetSection("Timer")
	if err != nil {
		LogErr("config get section faild :%v \n",err)
	}
	count = ret["CountdownTime"]
	freq = ret["FreqTime"]
	cal = ret["Calendar"]
	cmd = ret["Exec"]

	return
}
