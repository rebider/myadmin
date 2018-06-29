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

//在线统计
func (c *StatisticsController) OnlineStatisticsList() {
	var params models.DailyOnlineStatisticsQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("获取在线统计:%+v", params)
	data := models.GetDailyOnlineStatisticsList(&params)
	result := make(map[string]interface{})
	//result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取在线统计", result)
}
//注册统计
func (c *StatisticsController) RegisterStatisticsList() {
	var params models.DailyRegisterStatisticsQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("获取注册统计:%+v", params)
	data:= models.GetDailyRegisterStatisticsList(&params)
	result := make(map[string]interface{})
	//result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取注册统计", result)
}

//注册统计
func (c *StatisticsController) ActiveStatisticsList() {
	var params models.DailyActiveStatisticsQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("获取活跃统计:%+v", params)
	data := models.GetDailyActiveStatisticsList(&params)
	result := make(map[string]interface{})
	//result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取活跃统计", result)
}

//消费分析
func (c *StatisticsController) ConsumeAnalysis() {
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

//服务器概况
func (c *StatisticsController) GetServerGeneralize() {
	var params models.ServerGeneralizeQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询服务器概况:%+v", params)
	data, err := models.GetServerGeneralize(params.PlatformId, params.Node)
	c.CheckError(err)
	c.Result(enums.CodeSuccess, "获取服务器概况", data)
}

func (c *StatisticsController) GetRealTimeOnline() {
	platformId := c.GetString("platformId")
	serverId := c.GetString("serverId")
	serverOnlineStatistics, err := models.GetServerOnlineStatistics(platformId, serverId)
	c.CheckError(err, "查询实时在线统计")
	c.Result(enums.CodeSuccess, "查询实时在线统计成功", serverOnlineStatistics)
}
