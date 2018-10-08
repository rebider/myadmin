package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	"errors"
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
	if params.PlatformId == "" {
		c.CheckError(errors.New("平台ID不能为空"))
	}
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.Result(enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total, playerCount, moneyCount := models.GetChargeInfoRecordList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["playerCount"] = playerCount
	result["moneyCount"] = moneyCount
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

// 充值活动分布
func (c *ChargeController) ChargeActivityDistribution() {
	var params models.ChargeActivityDistributionQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("获取充值获得统计:%+v", params)
	data := models.GetChargeActivityDistribution(params)
	result := make(map[string]interface{})
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取充值活动分布", result)
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


// 每日LTV
func (c *ChargeController) GetDailyLTV() {
	var params models.DailyLTVQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("获取每日LTV:%+v", params)
	data := models.GetDailyLTVList(params)
	result := make(map[string]interface{})
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取每日LTV", result)
}
