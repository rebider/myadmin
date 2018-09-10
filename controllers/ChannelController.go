package controllers

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
	"encoding/json"
)

type ChannelController struct {
	BaseController
}

// 获取渠道列表
func (c *ChannelController) List() {
	list := models.GetChannelList()
	logs.Info("channelList:%v", list)
	//c.CheckError(err)
	result := make(map[string]interface{})
	//result["total"] = total
	result["rows"] = list
	c.Result(enums.CodeSuccess, "获取渠道列表成功", result)
}

// 编辑 添加渠道
func (c *ChannelController) Edit() {
	m := models.Channel{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑渠道")
	logs.Info("编辑平台:%+v", m)
	err = models.Db.Save(&m).Error
	c.CheckError(err, "编辑渠道失败")
	c.Result(enums.CodeSuccess, "保存成功", m.Id)
}

// 删除渠道
func (c *ChannelController) Del() {
	var idList []int
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &idList)
	utils.CheckError(err)
	logs.Info("删除渠道:%+v", idList)
	err = models.DeleteChannel(idList)
	c.CheckError(err, "删除渠道失败")
	c.Result(enums.CodeSuccess, "成功删除渠道", idList)
}
