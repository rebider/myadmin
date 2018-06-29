package controllers

import (
	//"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	"encoding/json"
)

type RemainController struct {
	BaseController
}

// 总体留存
func (c *RemainController) GetTotalRemain() {
	var params models.TotalRemainQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询总体留存:+", params)
	data, total := models.GetRemainTotalList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "查询总体留存成功", result)
}

// 活跃留存
func (c *RemainController) GetActiveRemain() {
	var params models.ActiveRemainQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询活跃留存:+", params)
	data, total := models.GetRemainActiveList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "查询活跃留存成功", result)
}

// 任务留存
func (c *RemainController) GetTaskRemain() {
	var params struct {
		PlatformId string    `json:"platformId"`
		Node       string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Info("查询任务留存:%+v", params)
	utils.CheckError(err)
	data := models.GetRemainTask(params.Node)
	result := make(map[string]interface{})
	result["rows"] = data
	c.Result(enums.CodeSuccess, "查询任务留存成功", result)
}

// 等级留存
func (c *RemainController) GetLevelRemain() {
	var params struct {
		PlatformId string    `json:"platformId"`
		Node       string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询等级留存:%+v", params)
	data := models.GetRemainLevel(params.Node)
	result := make(map[string]interface{})
	result["rows"] = data
	c.Result(enums.CodeSuccess, "查询等级留存成功", result)
}

//时长留存
func (c *RemainController) GetTimeRemain() {
	var params struct {
		PlatformId string    `json:"platformId"`
		Node       string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询时长留存:%+v", params)
	data := models.GetRemainTime(params.Node)
	result := make(map[string]interface{})
	result["rows"] = data
	c.Result(enums.CodeSuccess, "查询时长留存成功", result)
}
