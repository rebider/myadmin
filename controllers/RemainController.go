package controllers

import (
	//"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	//"github.com/chnzrb/myadmin/utils"
	"encoding/json"
)

type RemainController struct {
	BaseController
}

func (c *RemainController) GetTotalRemain() {
	logs.Info("查询总体留存")
	var params models.TotalRemainQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	c.CheckError(err)
	data, total := models.GetRemainTotalList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "查询在线统计成功", result)
}


func (c *RemainController) GetLevelRemain() {
	logs.Info("查询等级留存")
	var params models.TotalRemainQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	c.CheckError(err)
	data, total := models.GetRemainTotalList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "查询在线统计成功", result)
}
