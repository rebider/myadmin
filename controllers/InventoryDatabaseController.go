package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
)

type InventoryDatabaseController struct {
	BaseController
}



func (c *InventoryDatabaseController) DatabaseList() {
	var params models.InventoryDatabaseParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询资产列表:%+v", params)
	data, total := models.GetInventoryDatabaseList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取资产列表成功", result)
}

// 编辑 添加用户
func (c *InventoryDatabaseController) EditDatabase() {
	m := models.InventoryDatabase{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑资产")
	logs.Info("编辑资产:%+v", m)
	now := utils.GetTimestamp()
	if m.Id == 0 {
		m.AddTime = now
		m.UpdateTime = now
		err = models.Db.Save(&m).Error
		c.CheckError(err, "添加资产失败")
	} else {
		om, err := models.GetInventoryDatabaseOne(m.Id)
		c.CheckError(err, "未找到该资产")
		m.UpdateTime = now
		m.AddTime = om.AddTime
		err = models.Db.Save(&m).Error
		c.CheckError(err, "保存资产失败")
	}
	c.Result(enums.CodeSuccess, "保存成功", m.Id)
}

// 删除用户
func (c *InventoryDatabaseController) DeleteDatabase() {
	var idList []int
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &idList)
	utils.CheckError(err)
	logs.Info("删除资产:%+v", idList)
	err = models.DeleteInventoryDatabases(idList)
	c.CheckError(err, "删除资产失败")
	c.Result(enums.CodeSuccess, "成功删除资产", idList)
}
