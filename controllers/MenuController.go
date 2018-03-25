package controllers

import (
	"fmt"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	//"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	//"github.com/astaxie/beego/orm"
)

type MenuController struct {
	BaseController
}

func (c *MenuController) List() {
	//获取数据列表和总数
	data := models.TranMenuList2MenuTree(models.GetMenuList())
	result := make(map[string]interface{})
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取菜单列表成功", result)
}

func (c *MenuController) MenuTree() {
	c.Result(enums.CodeSuccess, "获取菜单树成功", models.TranMenuList2MenuTree(models.GetMenuList()))
}

//ParentTreeGrid 获取可以成为某节点的父节点列表
func (c *MenuController) GetParentMenuList() {
	//Id, _ := c.GetInt("id", 0)
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Debug("获取可以成为某节点的父节点列表:%+v", params)
	tree := models.MenuTreeGrid4Parent(params.Id)
	////转换UrlFor 2 LinkUrl
	//c.UrlFor2Link(tree)
	c.Result(enums.CodeSuccess, "", tree)
}

//Edit 资源编辑页面
func (c *MenuController) Edit() {
	m := models.Menu{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑资源")
	logs.Info("编辑资源:%+v", m)
	//var err error
	//o := orm.NewOrm()
	//parent := &models.Menu{}
	parentId := m.ParentId
	//parentId, _ := c.GetInt("Parent", 0)
	//获取form里的值
	//if err = c.ParseForm(&m); err != nil {
	//	c.Result(enums.JRCodeFailed, "获取数据失败", m.Id)
	//}
	//获取父节点
	if parentId > 0 {
		parent, err := models.GetMenuOne(parentId)
		if err == nil && parent != nil {
			m.Parent = parent
		} else {
			c.Result(enums.CodeFail, "父节点无效", "")
		}
	}
	if m.Id == 0 {
		err = models.Db.Save(&m).Error
		if  err == nil {
			c.Result(enums.CodeSuccess, "添加成功", m.Id)
		} else {
			c.Result(enums.CodeFail, "添加失败", m.Id)
		}

	} else {
		if parentId > 0 {
			if models.CanParentMenu(m.Id, parentId) == false {
				c.Result(enums.CodeFail, "请重新选择父节点", "")
			}
		}
		err = models.Db.Save(&m).Error
		if  err == nil {
			c.Result(enums.CodeSuccess, "编辑成功", m.Id)
		} else {
			c.Result(enums.CodeFail, "编辑失败", m.Id)
		}
	}
}

// Delete 删除
func (c *MenuController) Delete() {
	var m []int
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err)
	logs.Info("删除菜单:%+v", m)
	//query := orm.NewOrm().QueryTable(models.MenuTBName())
	_, err = models.DeleteMenus(m)
	if  err == nil {
		c.Result(enums.CodeSuccess, fmt.Sprintf("删除成功"), 0)
	} else {
		c.Result(enums.CodeFail, "删除失败", 0)
	}
}
