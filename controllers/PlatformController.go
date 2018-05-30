package controllers

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
)

type PlatformController struct {
	BaseController
}

func (c *PlatformController) List() {
	list, err := models.GetPlatformList()
	logs.Info("platformList:%v", list)
	c.CheckError(err)
	c.Result(enums.CodeSuccess, "获取平台列表成功", list)
}
