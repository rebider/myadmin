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
		PlatformId int    `json:"platformId"`
		ServerId       string `json:"serverId"`
		PlayerId   int    `json:"playerId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询玩家详细信息:%+v", params)
	playerDetail, err := models.GetPlayerDetail(params.PlatformId, params.ServerId, params.PlayerId)
	c.CheckError(err, "查询玩家详细信息失败")
	c.Result(enums.CodeSuccess, "获取玩家详细信息成功", playerDetail)
}

func (c *PlayerController) One() {
	platformId, err := c.GetInt("platformId")
	c.CheckError(err)
	//serverId := c.GetString("serverId")
	playerName := c.GetString("playerName")
	c.CheckError(err)
	player, err := models.GetPlayerByPlatformIdAndNickname(platformId, playerName)
	c.CheckError(err, "查询玩家失败")
	c.Result(enums.CodeSuccess, "获取玩家成功", player)
}

