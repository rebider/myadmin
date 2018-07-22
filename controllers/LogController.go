package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
)

type LogController struct {
	BaseController
}

func (c *LogController) PlayerLoinLogList() {
	var params models.PlayerLoginLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询玩家详细信息:%+v", params)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.Result(enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetPlayerLoginLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取玩家登录日志", result)
}

func (c *LogController) PlayerOnlineLogList() {
	var params models.PlayerOnlineLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询在线日志:%+v", params)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.Result(enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetPlayerOnlineLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取在线日志", result)
}

func (c *LogController) PlayerChallengeMissionLogList() {
	var params models.PlayerChallengeMissionLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询副本挑战日志:%+v", params)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.Result(enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetPlayerChallengeMissionLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取副本挑战日志", result)
}

func (c *LogController) PlayerPropLogList() {
	var params models.PlayerPropLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询道具日志:%+v", params)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.Result(enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetPlayerPropLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取道具日志", result)
}

//func (c *LogController) ChargeList() {
//	var params models.ChargeInfoRecordQueryParam
//	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
//	utils.CheckError(err)
//	logs.Info("查询充值记录日志:%+v", params)
//	if params.PlayerName != "" {
//		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
//		if player == nil || err != nil {
//			c.Result(enums.CodeFail, "玩家不存在", 0)
//		}
//		params.PlayerId = player.Id
//	}
//	data, total := models.GetChargeInfoRecordList(&params)
//	result := make(map[string]interface{})
//	result["total"] = total
//	result["rows"] = data
//	c.Result(enums.CodeSuccess, "获取充值记录日志", result)
//}
//
//

