package controllers

import (
	"github.com/chnzrb/myadmin/models"
	"encoding/json"
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"github.com/chnzrb/myadmin/enums"
	"strconv"
	"strings"
	"github.com/astaxie/beego/logs"
)

type ServerNodeController struct {
	BaseController
}

func (c *ServerNodeController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	c.checkLogin()
}
// DataGrid 角色管理首页 表格获取数据
func (c *ServerNodeController) List() {
	//直接反序化获取json格式的requestbody里的值
	var params models.ServerNodeQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	//获取数据列表和总数
	data, total := models.ServerNodePageList(&params)
	logs.Debug("DataGrid params:%#v", params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取节点列表成功",result)
}


// Edit 添加 编辑 页面
func (c *ServerNodeController) Edit() {
	m := models.ServerNode{}
	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		logs.Error("%s%v", err.Error(), m)
		c.Result(enums.CodeFail, "获取数据失败666", m.Node + err.Error())
	}
	//fmt.Println("Save:", m, c.Input())
	out, err := utils.Nodetool(
		"mod_server_mgr",
		"add_server_node",
		m.Node,
		strconv.Itoa(m.Port),
		strconv.Itoa(m.Type),
		strconv.Itoa(m.OpenTime),
		strconv.Itoa(m.PlatformId),
		"null",
	)

	if err != nil {
		fmt.Println("保存失败:"+ out)
		c.Result(enums.CodeFail, "保存失败:" + out, m.Node)
	}

	c.Result(enums.CodeSuccess, "保存成功", m.Node)
	//
	////如果是Post请求，则由Save处理
	//if c.Ctx.Request.Method == "POST" {
	//	c.Save()
	//}
	//Id := c.GetString(":node", "")
	////fmt.Println("Edit", c.Ctx.Request.Method, Id)
	//m := &models.ServerNode{}
	//var err error
	//if Id != "" {
	//	m, err = models.GetServerNodeById(Id)
	//
	//	fmt.Printf("%#v", m)
	//	if err != nil {
	//		fmt.Println("err", err)
	//		//c.pageError("数据无效，请刷新后重试")
	//	}
	//}
	//c.Data["m"] = m
	//models.ShowPlatformJson(c.Data)
	//models.ShowNodeTypeJsone(c.Data)
	////c.setTpl("servernode/edit.html", "shared/layout_pullbox.html")
	//c.LayoutSections = make(map[string]string)
	//c.LayoutSections["footerjs"] = "servernode/edit_footerjs.html"
}

//func (c *ServerNodeController) Save() {
//
//	m := models.ServerNode{}
//	//获取form里的值
//	if err := c.ParseForm(&m); err != nil {
//		logs.Error("%s%v", err.Error(), m)
//		c.Result(enums.CodeFail, "获取数据失败666", m.Node + err.Error())
//	}
//	//fmt.Println("Save:", m, c.Input())
//	out, err := utils.Nodetool(
//		"mod_server_mgr",
//		"add_server_node",
//		m.Node,
//		strconv.Itoa(m.Port),
//		strconv.Itoa(m.Type),
//		strconv.Itoa(m.OpenTime),
//		strconv.Itoa(m.PlatformId),
//		"null",
//		)
//
//	if err != nil {
//		fmt.Println("保存失败:"+ out)
//		c.Result(enums.CodeFail, "保存失败:" + out, m.Node)
//	}
//
//	c.Result(enums.CodeSuccess, "保存成功", m.Node)
//}

func (c *ServerNodeController) Delete() {
	strs := c.GetString("ids")
	ids := strings.Split(strs, ",")
	//fmt.Println("Delete:", strs, ids)
	for _, str := range ids {
		out, err := utils.Nodetool(
			"mod_server_mgr",
			"delete_server_node",
			str,
		)
		if err != nil {
			fmt.Println("删除失败:", strs, out, err)
			c.Result(enums.CodeFail, "删除失败:" + out, 0)
		}
	}
	c.Result(enums.CodeSuccess, fmt.Sprintf("成功删除 %d 项", len(ids)), 0)
}

