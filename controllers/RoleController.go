package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"github.com/astaxie/beego/logs"
)

//角色管理
type RoleController struct {
	BaseController
}

// 获取角色列表
func (c *RoleController) List() {
	var params models.RoleQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	data, total := models.GetRoleList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取角色列表成功", result)
}

//添加,编辑角色
func (c *RoleController) Edit() {
	m := models.Role{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑角色")
	if m.Id == 0 {
		err = models.Db.Save(&m).Error
		c.CheckError(err, "添加角色失败")
		c.Result(enums.CodeSuccess, "添加角色成功", m.Id)
	} else {
		err = models.Db.Save(&m).Error
		c.CheckError(err, "编辑角色失败")
		c.Result(enums.CodeSuccess, "编辑角色成功", m.Id)
	}
}

//删除角色
func (c *RoleController) Delete() {
	var ids []int
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ids)
	utils.CheckError(err)
	logs.Info("删除角色:%+v", ids)
	err = models.DeleteRoles(ids)
	c.CheckError(err, "删除角色失败")
	c.Result(enums.CodeSuccess, "成功删除角色", 0)
}

//角色分配资源
func (c *RoleController) AllocateResource() {
	var params struct {
		Id          int    `json:"id"`
		ResourceIds [] int `json:"resourceIds"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Info("角色分配资源:%+v", params)
	utils.CheckError(err)

	roleId := params.Id
	resourceIds := params.ResourceIds

	m, err := models.GetRoleOne(roleId)
	c.CheckError(err, fmt.Sprintf("角色不存在:%d", roleId))

	_, err = models.DeleteRoleResourceRelByRoleIdList([] int {roleId})
	c.CheckError(err, "删除旧的角色资源关系失败")

	for _, resourceId := range resourceIds {
		_, err := models.GetResourceOne(resourceId)
		c.CheckError(err, fmt.Sprintf("资源不存在:%d", resourceId))

		relation := models.RoleResourceRel{RoleId: roleId, ResourceId: resourceId}
		m.RoleResourceRel = append(m.RoleResourceRel, &relation)
	}
	err = models.Db.Save(&m).Error
	c.CheckError(err, "保存失败")
	c.Result(enums.CodeSuccess, "保存成功", "")
}

//角色分配渠道
func (c *RoleController) AllocateChannel() {
	var params struct {
		Id          int    `json:"id"`
		ChannelIds [] int `json:"channelIds"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Info("角色分配渠道:%+v", params)
	utils.CheckError(err)

	roleId := params.Id
	channelIds := params.ChannelIds

	m, err := models.GetRoleOne(roleId)
	c.CheckError(err, fmt.Sprintf("角色不存在:%d", roleId))

	_, err = models.DeleteRoleChannelRelByRoleIdList([] int {roleId})
	c.CheckError(err, "删除旧的角色平台关系失败")

	for _, channelId := range channelIds {
		_, err := models.GetChannelOne(channelId)
		c.CheckError(err, fmt.Sprintf("渠道不存在:%s", channelId))

		relation := models.RoleChannelRel{RoleId: roleId, ChannelId: channelId}
		m.RoleChannelRel = append(m.RoleChannelRel, &relation)
	}
	err = models.Db.Save(&m).Error
	c.CheckError(err, "保存失败")
	c.Result(enums.CodeSuccess, "保存成功", "")
}

//角色分配菜单
func (c *RoleController) AllocateMenu() {

	var params struct {
		Id          int    `json:"id"`
		MenuIds [] int `json:"menuIds"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Info("角色分配菜单:%+v", params)
	utils.CheckError(err)

	roleId := params.Id
	menuIds := params.MenuIds

	m, err := models.GetRoleOne(roleId)
	c.CheckError(err, fmt.Sprintf("角色不存在:%d", roleId))

	_, err = models.DeleteRoleMenuRelByRoleIdList([] int{roleId})
	c.CheckError(err, "删除旧的角色菜单关系失败")

	for _, menuId := range menuIds {
		_, err := models.GetMenuOne(menuId)
		c.CheckError(err, fmt.Sprintf("菜单不存在:%d", menuId))

		relation := models.RoleMenuRel{RoleId: roleId, MenuId: menuId}
		m.RoleMenuRel = append(m.RoleMenuRel, &relation)
	}
	err = models.Db.Save(&m).Error
	c.CheckError(err, "保存失败")
	c.Result(enums.CodeSuccess, "保存成功", "")
}
