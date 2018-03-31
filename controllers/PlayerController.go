package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
)

type PlayerController struct {
	BaseController
}


func (c *PlayerController) List() {
	var params models.PlayerQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询用户列表:%+v", params)
	data, total := models.GetPlayerList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取玩家列表成功", result)
}

func (c *PlayerController) Detail() {
	var params struct {
		PlatformId int `json:"platformId"`
		ServerId string `json:"serverId"`
		PlayerId int `json:"playerId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询玩家详细信息:%+v", params)
	playerDetail, err := models.GetPlayerDetail(params.PlatformId, params.ServerId, params.PlayerId)
	c.CheckError(err, "查询玩家详细信息失败")
	c.Result(enums.CodeSuccess, "获取玩家详细信息成功", playerDetail)
}

func (c *PlayerController) PlayerLoinLogList() {
	var params models.PlayerLoginLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询玩家详细信息:%+v", params)
	data, total := models.GetPlayerLoginLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取玩家登录日志", result)
}

func (c *PlayerController) PlayerOnlineLogList() {
	var params models.PlayerOnlineLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询在线日志:%+v", params)
	data, total := models.GetPlayerOnlineLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取在线日志", result)
}

func (c *PlayerController) GetServerGeneralize() {
	var params models.ServerGeneralizeQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询服务器概况:%+v", params)
	data, err := models.GetServerGeneralize(params.PlatformId, params.ServerId)
	c.CheckError(err)
	c.Result(enums.CodeSuccess, "获取服务器概况", data)
}

