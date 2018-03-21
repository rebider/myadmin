package controllers

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
)

type PlatformController struct {
	BaseController
}

func (c *PlatformController) List() {
	list, err := models.GetPlatformList("data/json/Platform.json")
	utils.CheckError(err)
	logs.Info("platformList:%v", list)
	c.Result(enums.CodeSuccess, "获取平台列表成功", list)
}
