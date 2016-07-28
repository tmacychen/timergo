package main

import (
	"log"
	"os"
)

var logger *log.Logger

func LogInit(logfile *os.File) {
	logger = log.New(logfile, "[unset]", log.LstdFlags|log.Lshortfile)
}

func LogDebug(msg string,v ...interface{}) {
	logger.SetPrefix("[Debug]")
	LogFmt(msg,v...)
}
func LogInfo(msg string,v ...interface{}) {
	logger.SetPrefix("[Info]")
	LogFmt(msg,v...)
}
func LogWarn(msg string,v ...interface{}) {
	logger.SetPrefix("[Warn]")
	LogFmt(msg,v...)
}

func LogErr(msg string,v ...interface{}) {
	logger.SetPrefix("[Error]")
	LogFmt(msg,v...)
}
func LogFmt(format string, v ...interface{}) {
	logger.Printf(format, v...)
}
