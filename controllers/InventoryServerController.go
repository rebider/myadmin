package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
)

type InventoryServerController struct {
	BaseController
}



func (c *InventoryServerController) ServerList() {
	var params models.InventoryServerParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询资产列表:%+v", params)
	data, total := models.GetInventoryServerList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取资产列表成功", result)
}
func (c *InventoryServerController) AllServerList() {
	logs.Info("查询所有资产列表")
	data:= models.GetAllServerList()
	result := make(map[string]interface{})
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取所有资产列表成功", result)
}
// 编辑 添加服务器
func (c *InventoryServerController) EditServer() {
	m := models.InventoryServer{}
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
		om, err := models.GetInventoryServerOne(m.Id)
		c.CheckError(err, "未找到该资产")
		m.UpdateTime = now
		m.AddTime = om.AddTime
		err = models.Db.Save(&m).Error
		c.CheckError(err, "保存资产失败")
	}
	c.Result(enums.CodeSuccess, "保存成功", m.Id)
}

// 删除服务器
func (c *InventoryServerController) DeleteServer() {
	var idList []int
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &idList)
	utils.CheckError(err)
	logs.Info("删除资产:%+v", idList)
	err = models.DeleteInventoryServers(idList)
	c.CheckError(err, "删除资产失败")
	c.Result(enums.CodeSuccess, "成功删除资产", idList)
}


// 创建ansible
func (c *InventoryServerController) CreateAnsibleInventory() {
	err := models.CreateAnsibleInventory()
	c.CheckError(err, "生成ansible inventory 失败")
	c.Result(enums.CodeSuccess, "生成ansible inventory成功", "")
}
