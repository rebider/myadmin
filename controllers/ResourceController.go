package controllers

import (
	"fmt"
	"strings"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
	"encoding/json"
)

type ResourceController struct {
	BaseController
}

//获取资源列表
func (c *ResourceController) List() {
	//获取数据列表和总数
	data := models.TranResourceList2ResourceTree(models.GetResourceList())
	result := make(map[string]interface{})
	c.UrlFor2Link(data)
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取资源列表成功", result)
}

func (c *ResourceController) ResourceTree() {
	c.Result(enums.CodeSuccess, "获取资源树成功", models.TranResourceList2ResourceTree(models.GetResourceList()))
}

//获取可以成为某节点的父节点列表
func (c *ResourceController) GetParentResourceList() {
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	tree := models.ResourceTreeGrid4Parent(params.Id)
	c.Result(enums.CodeSuccess, "", tree)
}

// 将资源表里的UrlFor值转成LinkUrl
func (c *ResourceController) UrlFor2LinkOne(urlfor string) string {
	if len(urlfor) == 0 {
		return ""
	}
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

//编辑添加资源
func (c *ResourceController) Edit() {
	m := models.Resource{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑资源")
	logs.Info("编辑资源:%+v", m)
	parent := &models.Resource{}
	parentId := m.ParentId
	//获取父节点
	if parentId > 0 {
		parent, err = models.GetResourceOne(parentId)
		c.CheckError(err, "父节点无效")
		m.Parent = parent
	}
	if c.UrlFor2LinkOne(m.UrlFor) != "" || strings.Contains(m.UrlFor, ".*") {
	} else {
		c.Result(enums.CodeFail, "控制器解析失败: "+m.UrlFor, "")
	}
	if m.Id == 0 {
		err = models.Db.Save(&m).Error
		c.CheckError(err, "添加资源失败")
		c.Result(enums.CodeSuccess, "添加资源成功", m.Id)

	} else {
		if parentId > 0 {
			if models.CanParent(m.Id, parentId) == false {
				c.Result(enums.CodeFail, "请重新选择父节点", "")
			}
		}
		err = models.Db.Save(&m).Error
		c.CheckError(err, "编辑资源失败")
		c.Result(enums.CodeSuccess, "编辑资源成功", m.Id)
	}
}

//删除资源
func (c *ResourceController) Delete() {
	var m []int
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err)
	logs.Info("删除资源:%+v", m)
	_, err = models.DeleteResources(m)
	c.CheckError(err, "删除资源失败")
	c.Result(enums.CodeSuccess, fmt.Sprintf("删除资源成功"), 0)
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
