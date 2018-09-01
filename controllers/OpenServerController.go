package controllers
//
//import (
//	"github.com/chnzrb/myadmin/enums"
//	"github.com/chnzrb/myadmin/models"
//	"github.com/chnzrb/myadmin/utils"
//	"github.com/astaxie/beego/logs"
//	"encoding/json"
//)
//
//type OpenServerController struct {
//	BaseController
//}
//
//func (c *OpenServerController) List() {
//	list := models.GetOpenServerList()
//	logs.Info("OpenServerList:%v", list)
//	//c.CheckError(err)
//	result := make(map[string]interface{})
//	//result["total"] = total
//	result["rows"] = list
//	c.Result(enums.CodeSuccess, "获取开服策略列表成功", result)
//}
//
//// 编辑 添加开服策略
//func (c *OpenServerController) Edit() {
//	m := models.OpenServer{}
//	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
//	utils.CheckError(err, "编辑开服策略")
//	logs.Info("编辑开服策略:%+v", m)
//	//删除旧的用户角色关系
//	_, err = models.DeletePlatformInventorySeverRelByPlatformIdList([] string{m.PlatformId})
//	c.CheckError(err, "删除旧的用户角色关系失败")
//	for _, inventorySeverId := range m.InventorySeverIds {
//		relation := models.PlatformInventorySeverRel{OpenServerId: m.PlatformId, InventoryServerId: inventorySeverId}
//		m.OpenServerInventorySeverRel = append(m.OpenServerInventorySeverRel, &relation)
//	}
//	now := utils.GetTimestamp()
//	m.Time = now
//	logs.Info("%+++v",m.OpenServerInventorySeverRel)
//	err = models.Db.Debug().Save(&m).Error
//	c.CheckError(err, "编辑开服策略失败")
//	c.Result(enums.CodeSuccess, "保存成功", "")
//}
//
//// 删除开服策略
//func (c *OpenServerController) Del() {
//	var idList []string
//	err := json.Unmarshal(c.Ctx.Input.RequestBody, &idList)
//	utils.CheckError(err)
//	logs.Info("删除开服策略:%+v", idList)
//	err = models.DeleteOpenServer(idList)
//	c.CheckError(err, "删除开服策略失败")
//	c.Result(enums.CodeSuccess, "成功删除开服策略", idList)
//}
