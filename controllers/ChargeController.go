package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
)

type ChargeController struct {
	BaseController
}

// 获取充值列表
func (c *ChargeController) ChargeList() {
	var params models.ChargeInfoRecordQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询充值记录日志:%+v", params)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.Result(enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetChargeInfoRecordList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取充值记录日志", result)
}



//获取充值排行榜
func (c *ChargeController) ChargeRankList() {
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



// 充值统计
func (c *ChargeController) ChargeStatisticsList() {
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

// 充值任务分布
func (c *ChargeController) ChargeTaskDistribution() {
	var params models.ChargeTaskDistributionQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("获取充值统计:%+v", params)
	data := models.GetChargeTaskDistribution(params)
	result := make(map[string]interface{})
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取充充值任务分布", result)
}
// 充值金额分布
func (c *ChargeController) ChargeMoneyDistribution() {
	var params models.ChargeMoneyDistributionQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("获取金额分布:%+v", params)
	data := models.GetChargeMoneyDistribution(params)
	result := make(map[string]interface{})
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取金额分布", result)
}

// 充值等级分布
func (c *ChargeController) ChargeLevelDistribution() {
	var params models.ChargeLevelDistributionQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("获取等级分布:%+v", params)
	data := models.GetChargeLevelDistribution(params)
	result := make(map[string]interface{})
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取等级分布", result)
}
