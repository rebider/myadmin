package utils

import (
	//"strings"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	//"fmt"
	//"github.com/hzwy23/panda/logger"
)

// consoleLogs开发模式下日志
//var consoleLogs *logs.BeeLogger

// fileLogs 生产环境下日志
//var fileLogs *logs.BeeLogger
//var logger *logs.BeeLogger
//运行方式
//var runmode string

func InitLogs() {
	beego.BConfig.Log.AccessLogsFormat = ""
	//consoleLogs = logs.NewLogger(1)
	//consoleLogs.SetLogger(logs.AdapterConsole)
	//consoleLogs.Async() //异步
	//
	//fileLogs = logs.NewLogger(10000)
	//level := beego.AppConfig.String("logs::level")
	//fileLogs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/admin.log",
	//	"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],
	//	"level":`+level+`,
	//	"daily":true,
	//	"maxdays":10}`)
	//fileLogs.Async() //异步
	//runmode := strings.TrimSpace(strings.ToLower(beego.AppConfig.DefaultString("runmode", "dev")))
	//if runmode == "dev" {
	//	//level := beego.AppConfig.String("logs::level")
	//	////"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],
	//	//logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/admin.log",
	//	//"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],
	//	//"level":`+ level+ `,
	//	//"daily":true,
	//	//"maxdays":10}`)
	//	//fmt.Println("")
	//	//logger = logs.NewLogger(1)
	//	//logger.SetLogger(logs.AdapterConsole)
	//	//logger.Async() //异步
	//} else {
	//	//logger = logs.NewLogger(10000)
	//	level := beego.AppConfig.String("logs::level")
	//	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/admin.log",
	//	"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],
	//	"level":`+ level+ `,
	//	"daily":true,
	//	"maxdays":10}`)
	//}
	level := beego.AppConfig.String("logs::level")
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/admin.log",
		"separate":["critical", "error", "warning", "info", "debug"],
		"level":`+ level+ `,
		"daily":true,
		"maxdays":10}`)
	logs.Async() //异步
	//输出文件名和行号
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
}

//func Debug(format string, v ...interface{}) {
//	logger.Debug(format, v)
//}
//func Debug(v ...interface{}) {
//	format := "%s"
//	logger.Debug(format, v)
//}
//func InfoFormat(format string, v ...interface{}) {
//	logger.Info(format, v)
//}
//func Info(v ...interface{}) {
//	format := "%s"
//	logger.Info(format, v)
//}
//
//func Error(v ...interface{}) {
//	format := "%s"
//	logger.Error(format, v)
//}
