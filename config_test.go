package main

import (
	"log"
	"os"
	"testing"
)

func Test_GetConfig(t *testing.T) {
	logfile, err := os.OpenFile("TimerGo.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer logfile.Close()
	if err != nil {
		log.Fatalf("open file err :%v\n", err)

	}
	LogInit(logfile)

	a, b, c := GetConfig("./config.ini")
	t.Logf("a: %s b :%s c :%s\n",a,b,c)

}
