package controllers

import (
	"github.com/chnzrb/myadmin/models"
	"encoding/json"
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"github.com/chnzrb/myadmin/enums"
	"strconv"
	"github.com/astaxie/beego/logs"
)

type ServerNodeController struct {
	BaseController
}

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
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	logs.Debug("Edit:%v",m )
	utils.CheckError(err, "编辑节点")
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
}

func (c *ServerNodeController) Delete() {
	var ids []string
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ids)
	utils.CheckError(err)
	logs.Info("删除节点:%+v", ids)
	for _, str := range ids {
		out, err := utils.Nodetool(
			"mod_server_mgr",
			"delete_server_node",
			str,
		)
		if err != nil {
			fmt.Println("删除失败:", ids, out, err)
			c.Result(enums.CodeFail, "删除失败:" + out, 0)
		}
	}
	c.Result(enums.CodeSuccess, fmt.Sprintf("成功删除 %d 项", len(ids)), 0)
}

