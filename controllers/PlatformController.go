package controllers

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
	"encoding/json"
)

type PlatformController struct {
	BaseController
}

func (c *PlatformController) List() {
	list := models.GetPlatformList()
	logs.Info("platformList:%v", list)
	//c.CheckError(err)
	result := make(map[string]interface{})
	//result["total"] = total
	result["rows"] = list
	c.Result(enums.CodeSuccess, "获取平台列表成功", result)
}

// 编辑 添加平台
func (c *PlatformController) Edit() {
	m := models.Platform{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑平台")
	logs.Info("编辑平台:%+v", m)
	//删除旧的用平台服务器关系
	_, err = models.DeletePlatformInventorySeverRelByPlatformIdList([] string{m.Id})
	c.CheckError(err, "删除旧的用户角色关系失败")
	for _, inventorySeverId := range m.InventorySeverIds {
		relation := models.PlatformInventorySeverRel{PlatformId: m.Id, InventoryServerId: inventorySeverId}
		m.PlatformInventorySeverRel = append(m.PlatformInventorySeverRel, &relation)
	}
	now := utils.GetTimestamp()
	m.Time = now
	err = models.Db.Save(&m).Error
	c.CheckError(err, "编辑平台失败")
	c.Result(enums.CodeSuccess, "保存成功", m.Id)
}

// 删除平台
func (c *PlatformController) Del() {
	var idList []string
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &idList)
	utils.CheckError(err)
	logs.Info("删除平台:%+v", idList)
	err = models.DeletePlatform(idList)
	c.CheckError(err, "删除平台失败")
	c.Result(enums.CodeSuccess, "成功删除平台", idList)
}
