package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/chnzrb/myadmin/enums"
)

type StatisticsController struct {
	BaseController
}


func (c *StatisticsController) ChargeStatisticsList() {
	var params models.DailyChargeStatisticsQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("获取充值统计:%+v", params)
	data, total := models.GetDailyChargeStatisticsList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取充值统计", result)
}

func (c *StatisticsController) OnlineStatisticsList() {
	var params models.DailyOnlineStatisticsQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("获取在线统计:%+v", params)
	data, total := models.GetDailyOnlineStatisticsList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取在线统计", result)
}

func (c *StatisticsController) RegisterStatisticsList() {
	var params models.DailyRegisterStatisticsQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("获取注册统计:%+v", params)
	data, total := models.GetDailyRegisterStatisticsList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取注册统计", result)
}

func (c *StatisticsController) ConsumeStatistics() {
	var params models.PropConsumeStatisticsQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("获取消费分析统计:%+v", params)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.Result(enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, err := models.GetPropConsumeStatistics(&params)
	c.CheckError(err)
	result := make(map[string]interface{})
	result["rows"] = data
	c.Result(enums.CodeSuccess, "消费分析统计", result)
}


func (c *StatisticsController) ChargeRankList() {
	var params models.PlayerChargeDataQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询充值排行:%+v", params)
	data, total := models.GetPlayerChargeDataList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "查询充值排行", result)
}

