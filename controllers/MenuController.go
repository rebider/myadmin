package controllers

import (
	"fmt"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"encoding/json"
)

type MenuController struct {
	BaseController
}

func (c *MenuController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的少数Action需要权限控制，则将验证放到需要控制的Action里
	//c.checkAuthor("TreeGrid", "UserMenuTree", "ParentTreeGrid", "Select")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//这里注释了权限控制，因此这里需要登录验证
	c.checkLogin()
}
func (c *MenuController) List() {
	logs.Info("查询菜单列表")
	//fmt.Println(params)
	//获取数据列表和总数
	menuList := models.MenuList()
	c.Result(enums.CodeSuccess, "获取资源列表成功", menuList)
}

func (c *MenuController) MenuTree() {
	c.Result(enums.CodeSuccess, "获取资源树成功", models.GetMenuTree())
}

//编辑菜单
func (c *MenuController) Edit() {
	m := models.Menu{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑菜单")
	logs.Info("编辑菜单:%+v", m)
	logs.Info("编辑菜单:%+v", m.Parent)

	o := orm.NewOrm()
	parent := &models.Menu{}
	parentId := m.Parent.Id

	//获取父节点
	if parentId > 0 {
		parent, err = models.MenuOne(parentId)
		if err == nil && parent != nil {
			m.Parent = parent
		} else {
			c.Result(enums.CodeFail, "父节点无效", "")
		}
	}
	if m.Id == 0 {
		if _, err = o.Insert(&m); err == nil {
			c.Result(enums.CodeSuccess, "添加成功", m.Id)
		} else {
			c.Result(enums.CodeFail, "添加失败", m.Id)
		}

	} else {
		if _, err = o.Update(&m); err == nil {
			c.Result(enums.CodeSuccess, "编辑成功", m.Id)
		} else {
			c.Result(enums.CodeFail, "编辑失败", m.Id)
		}
	}
}

// 删除菜单
func (c *MenuController) Delete() {
	var m  []int
	json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	logs.Info("删除菜单:%+v",  m)
	query := orm.NewOrm().QueryTable(models.MenuTBName())
	if _, err := query.Filter("id", m[0]).Delete(); err == nil {
		c.Result(enums.CodeSuccess, fmt.Sprintf("删除成功"), 0)
	} else {
		c.Result(enums.CodeFail, "删除失败", 0)
	}
}
