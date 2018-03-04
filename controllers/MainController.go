package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	logs.Debug("ddddddddddddddddddddd")
	c.TplName = "index.html"
}
