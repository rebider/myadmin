package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)


func InitLogs() {
	beego.BConfig.Log.AccessLogsFormat = ""
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
