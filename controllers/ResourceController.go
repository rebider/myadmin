package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"encoding/json"
)

type ResourceController struct {
	BaseController
}

func (c *ResourceController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的少数Action需要权限控制，则将验证放到需要控制的Action里
	//c.checkAuthor("TreeGrid", "UserMenuTree", "ParentTreeGrid", "Select")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//这里注释了权限控制，因此这里需要登录验证
	c.checkLogin()
}
//func (c *ResourceController) Index() {
//	//需要权限控制
//	c.checkAuthor()
//	//将页面左边菜单的某项激活
//	//c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
//	c.setTpl()
//	c.LayoutSections = make(map[string]string)
//	c.LayoutSections["headcssjs"] = "resource/index_headcssjs.html"
//	c.LayoutSections["footerjs"] = "resource/index_footerjs.html"
//	//页面里按钮权限控制
//	c.Data["canEdit"] = c.checkActionAuthor("ResourceController", "Edit")
//	c.Data["canDelete"] = c.checkActionAuthor("ResourceController", "Delete")
//}

// @Title 获取资源列表1
// @Description 获取资源列表
// @Success 200 {object} models.GameServer
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /resource/List [get]
//TreeGrid 获取所有资源的列表
//func (c *ResourceController) List() {
//	tree := models.ResourceTreeGrid()
//	//转换UrlFor 2 LinkUrl
//	//c.UrlFor2Link(tree)
//	logs.Debug("tree:%v", tree)
//	c.Result(enums.JRCodeSucc, "", tree)
//}

func (c *ResourceController) List() {
	logs.Info("查询资源列表")
	//fmt.Println(params)
	//获取数据列表和总数
	data := models.ResourceList()
	result := make(map[string]interface{})
	//result["total"] = total
	result["rows"] = data
	//c.Data["json"] = result
	c.Result(enums.CodeSuccess, "获取资源列表成功", result)
}

func (c *ResourceController) ResourceTree() {
	c.Result(enums.CodeSuccess, "获取资源树成功", models.TranResourceList2ResourceTree(models.ResourceList()))
}

////UserMenuTree 获取用户有权管理的菜单、区域列表
//func (c *ResourceController) UserMenuTree() {
//	userid := c.curUser.Id
//	//获取用户有权管理的菜单列表（包括区域）
//	tree := models.ResourceTreeGridByUserId(userid, 1)
//	//转换UrlFor 2 LinkUrl
//
//	c.UrlFor2Link(tree)
//	c.Result(enums.CodeSuccess, "", tree)
//}


//ParentTreeGrid 获取可以成为某节点的父节点列表
func (c *ResourceController) ParentTreeGrid() {
	Id, _ := c.GetInt("id", 0)
	tree := models.ResourceTreeGrid4Parent(Id)
	//转换UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.Result(enums.CodeSuccess, "", tree)
}

// UrlFor2LinkOne 使用URLFor方法，将资源表里的UrlFor值转成LinkUrl
func (c *ResourceController) UrlFor2LinkOne(urlfor string) string {
	if len(urlfor) == 0 {
		return ""
	}
	// ResourceController.Edit,:id,1
	strs := strings.Split(urlfor, ",")
	if len(strs) == 1 {
		return c.URLFor(strs[0])
	} else if len(strs) > 1 {
		var values []interface{}
		for _, val := range strs[1:] {
			values = append(values, val)
		}
		return c.URLFor(strs[0], values...)
	}
	return ""
}

//UrlFor2Link 使用URLFor方法，批量将资源表里的UrlFor值转成LinkUrl
func (c *ResourceController) UrlFor2Link(src []*models.Resource) {
	for _, item := range src {
		item.LinkUrl = c.UrlFor2LinkOne(item.UrlFor)
	}
}

//Edit 资源编辑页面
func (c *ResourceController) Edit() {
	m := models.Resource{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑资源")
	logs.Info("编辑资源:%+v", m)
	logs.Info("编辑资源:%+v", m.Parent)
	//var err error
	o := orm.NewOrm()
	parent := &models.Resource{}
	//m := models.Resource{}
	parentId := m.Parent.Id
	//parentId, _ := c.GetInt("Parent", 0)
	//获取form里的值
	//if err = c.ParseForm(&m); err != nil {
	//	c.Result(enums.JRCodeFailed, "获取数据失败", m.Id)
	//}
	//获取父节点
	if parentId > 0 {
		parent, err = models.ResourceOne(parentId)
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


// Delete 删除
func (c *ResourceController) Delete() {
	var m  []int
	json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	logs.Info("删除资源:%+v",  m)
	query := orm.NewOrm().QueryTable(models.ResourceTBName())
	if _, err := query.Filter("id", m[0]).Delete(); err == nil {
		c.Result(enums.CodeSuccess, fmt.Sprintf("删除成功"), 0)
	} else {
		c.Result(enums.CodeFail, "删除失败", 0)
	}
}

// Select 通用选择面板
func (c *ResourceController) Select() {
	//获取调用者的类别 1表示 角色
	desttype, _ := c.GetInt("desttype", 0)
	//获取调用者的值
	destval, _ := c.GetInt("destval", 0)
	//返回的资源列表
	var selectedIds []string
	o := orm.NewOrm()
	if desttype > 0 && destval > 0 {
		//如果都大于0,则获取已选择的值，例如：角色，就是获取某个角色已关联的资源列表
		switch desttype {
		case 1:
			{
				role := models.Role{Id: destval}
				o.LoadRelated(&role, "RoleResourceRel")
				for _, item := range role.RoleResourceRel {
					selectedIds = append(selectedIds, strconv.Itoa(item.Resource.Id))
				}
			}
		}
	}
	c.Data["selectedIds"] = strings.Join(selectedIds, ",")
	//c.setTpl("resource/select.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "resource/select_headcssjs.html"
	c.LayoutSections["footerjs"] = "resource/select_footerjs.html"
}

//CheckUrlFor 填写UrlFor时进行验证
func (c *ResourceController) CheckUrlFor() {
	urlfor := c.GetString("urlfor")
	link := c.UrlFor2LinkOne(urlfor)
	if len(link) > 0 {
		c.Result(enums.CodeSuccess, "解析成功", link)
	} else {
		c.Result(enums.CodeFail, "解析失败", link)
	}
}
