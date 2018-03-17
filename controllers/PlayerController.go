package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
)

type PlayerController struct {
	BaseController
}


func (c *PlayerController) List() {
	var params models.PlayerQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Info("查询用户列表:%+v", params)
	//获取数据列表和总数
	data, total := models.PlayerPageList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取玩家列表成功", result)
}
