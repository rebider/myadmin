package sysinit

import (
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego"
	_ "github.com/chnzrb/myadmin/models"
	_ "github.com/chnzrb/myadmin/crons"
)

func init() {
	//启用Session
	beego.BConfig.WebConfig.Session.SessionOn = true
	//初始化日志
	utils.InitLogs()
	//初始化缓存
	utils.InitCache()
}
