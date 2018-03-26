package controllers

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
	"encoding/json"
)

type MenuController struct {
	BaseController
}
//获取菜单列表
func (c *MenuController) List() {
	data := models.TranMenuList2MenuTree(models.GetMenuList())
	result := make(map[string]interface{})
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取菜单列表成功", result)
}

func (c *MenuController) MenuTree() {
	c.Result(enums.CodeSuccess, "获取菜单树成功", models.TranMenuList2MenuTree(models.GetMenuList()))
}

//获取可以成为某节点的父节点列表
func (c *MenuController) GetParentMenuList() {
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Debug("获取可以成为某节点的父节点列表:%+v", params)
	tree := models.MenuTreeGrid4Parent(params.Id)
	c.Result(enums.CodeSuccess, "", tree)
}

//编辑添加菜单
func (c *MenuController) Edit() {
	m := models.Menu{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑菜单")
	logs.Info("编辑菜单:%+v", m)
	parentId := m.ParentId
	//获取父节点
	if parentId > 0 {
		parent, err := models.GetMenuOne(parentId)
		c.CheckError(err, "父节点无效")
		m.Parent = parent
	}
	if m.Id == 0 {
		err = models.Db.Save(&m).Error
		c.CheckError(err, "添加菜单失败")
		c.Result(enums.CodeSuccess, "添加菜单成功", m.Id)
	} else {
		if parentId > 0 {
			if models.CanParentMenu(m.Id, parentId) == false {
				c.Result(enums.CodeFail, "请重新选择父节点", "")
			}
		}
		err = models.Db.Save(&m).Error
		c.CheckError(err, "编辑菜单失败")
		c.Result(enums.CodeSuccess, "编辑菜单成功", m.Id)
	}
}

// 删除菜单
func (c *MenuController) Delete() {
	var m []int
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err)
	logs.Info("删除菜单:%+v", m)
	_, err = models.DeleteMenus(m)
	c.CheckError(err, "删除菜单失败")
	c.Result(enums.CodeSuccess, "删除菜单成功", "")
}
