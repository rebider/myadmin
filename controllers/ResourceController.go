package controllers

import (
	"fmt"
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

func (c *ResourceController) List() {
	//获取数据列表和总数
	data := models.TranResourceList2ResourceTree(models.ResourceList())
	result := make(map[string]interface{})
	c.UrlFor2Link(data)
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取资源列表成功", result)
}

func (c *ResourceController) ResourceTree() {
	c.Result(enums.CodeSuccess, "获取资源树成功", models.TranResourceList2ResourceTree(models.ResourceList()))
}

//ParentTreeGrid 获取可以成为某节点的父节点列表
func (c *ResourceController) GetParentResourceList() {
	//Id, _ := c.GetInt("id", 0)
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Debug("获取可以成为某节点的父节点列表:%+v", params)
	tree := models.ResourceTreeGrid4Parent(params.Id)
	////转换UrlFor 2 LinkUrl
	//c.UrlFor2Link(tree)
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
		item.Url = c.UrlFor2LinkOne(item.UrlFor)
		c.UrlFor2Link(item.Children)
	}
}

//Edit 资源编辑页面
func (c *ResourceController) Edit() {
	m := models.Resource{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑资源")
	logs.Info("编辑资源:%+v", m)
	//var err error
	o := orm.NewOrm()
	parent := &models.Resource{}
	//m := models.Resource{}
	parentId := m.ParentId
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
	if m.Type == 1 {
		if c.UrlFor2LinkOne(m.UrlFor) != "" || strings.Contains(m.UrlFor, ".*"){
		} else {
			c.Result(enums.CodeFail, "控制器解析失败: " + m.UrlFor, "")
		}
	}
	if m.Id == 0 {
		if _, err = o.Insert(&m); err == nil {
			c.Result(enums.CodeSuccess, "添加成功", m.Id)
		} else {
			c.Result(enums.CodeFail, "添加失败", m.Id)
		}

	} else {
		if parentId > 0 {
			if models.CanParent(m.Id, parentId) == false {
				c.Result(enums.CodeFail, "请重新选择父节点", "")
			}
		}
		if _, err = o.Update(&m); err == nil {
			c.Result(enums.CodeSuccess, "编辑成功", m.Id)
		} else {
			c.Result(enums.CodeFail, "编辑失败", m.Id)
		}
	}
}

// Delete 删除
func (c *ResourceController) Delete() {
	var m []int
	json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	logs.Info("删除资源:%+v", m)
	query := orm.NewOrm().QueryTable(models.ResourceTBName())
	if _, err := query.Filter("id", m[0]).Delete(); err == nil {
		c.Result(enums.CodeSuccess, fmt.Sprintf("删除成功"), 0)
	} else {
		c.Result(enums.CodeFail, "删除失败", 0)
	}
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
