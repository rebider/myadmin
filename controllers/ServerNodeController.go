package controllers

import (
	"github.com/chnzrb/myadmin/models"
	"encoding/json"
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"github.com/chnzrb/myadmin/enums"
	"strconv"
	"github.com/astaxie/beego/logs"
	"os"
)

type ServerNodeController struct {
	BaseController
}

//获取节点列表
func (c *ServerNodeController) List() {
	var params models.ServerNodeQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Debug("获取节点列表:%#v", params)
	utils.CheckError(err)
	data, total := models.ServerNodePageList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取节点列表成功",result)
}


//  添加 编辑 节点
func (c *ServerNodeController) Edit() {
	m := models.ServerNode{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	logs.Debug("编辑 节点:%v",m )
	utils.CheckError(err, "编辑节点")
	if m.IsAdd == 1 && models.IsServerNodeExists(m.Node){
		c.Result(enums.CodeFail, "节点已经存在", m.Node)
	}
	if m.IsAdd == 0 && models.IsServerNodeExists(m.Node) == false{
		c.Result(enums.CodeFail, "节点不存在", m.Node)
	}

	out, err := utils.NodeTool(
		"mod_server_mgr",
		"add_server_node",
		m.Node,
		m.Ip,
		strconv.Itoa(m.Port),
		strconv.Itoa(m.WebPort),
		strconv.Itoa(m.Type),
		//strconv.Itoa(m.OpenTime),
		m.PlatformId,
		//m.ZoneNode,
		//strconv.Itoa(m.State),
		m.DbHost,
		strconv.Itoa(m.DbPort),
		m.DbName,
	)
	c.CheckError(err, "保存节点失败:"+ out)
	c.Result(enums.CodeSuccess, "保存成功", m.Node)
}

// 删除节点
func (c *ServerNodeController) Delete() {
	var ids []string
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ids)
	utils.CheckError(err)
	logs.Info("删除节点:%+v", ids)
	for _, str := range ids {
		out, err := utils.NodeTool(
			"mod_server_mgr",
			"delete_server_node",
			str,
		)
		c.CheckError(err, "删除节点失败:"+ out)
	}
	c.Result(enums.CodeSuccess, fmt.Sprintf("成功删除 %d 项", len(ids)), 0)
}


////  启动 节点
//func (c *ServerNodeController) Start() {
//	var params struct {
//		Node     string
//	}
//	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
//	logs.Debug("启动 节点:%v",params )
//	c.CheckError(err, "启动节点")
//	curDir := utils.GetCurrentDirectory()
//	defer os.Chdir(curDir)
//	toolDir := utils.GetToolDir()
//	err = os.Chdir(toolDir)
//
//	c.CheckError(err, "启动节点")
//	commandArgs := []string{
//		"node_tool.sh",
//		params.Node,
//		"start",
//	}
//	out, err := utils.Cmd("sh", commandArgs)
//	c.CheckError(err, "启动节点失败:"+ out)
//	c.Result(enums.CodeSuccess, "启动成功", params.Node)
//}
//
////  停止 节点
//func (c *ServerNodeController) Stop() {
//	var params struct {
//		Node     string
//	}
//	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
//	logs.Debug("停止节点:%v",params )
//	c.CheckError(err, "停止节点")
//	curDir := utils.GetCurrentDirectory()
//	defer os.Chdir(curDir)
//	toolDir := utils.GetToolDir()
//	err = os.Chdir(toolDir)
//	c.CheckError(err, "停止节点")
//	commandArgs := []string{
//		"node_tool.sh",
//		params.Node,
//		"stop",
//	}
//	out, err := utils.Cmd("sh", commandArgs)
//	c.CheckError(err, "停止节点失败:"+ out)
//	c.Result(enums.CodeSuccess, "停止成功", params.Node)
//}

func (c *ServerNodeController) Action() {
	var params struct {
		Node     string
		Action string
	}
	logs.Debug("节点操作:%v",params )
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	c.CheckError(err)
	curDir := utils.GetCurrentDirectory()
	defer os.Chdir(curDir)
	toolDir := utils.GetToolDir()
	err = os.Chdir(toolDir)
	c.CheckError(err)
	commandArgs := []string{
		"node_tool.sh",
		params.Node,
		params.Action,
	}
	out, err := utils.Cmd("sh", commandArgs)
	c.CheckError(err, fmt.Sprintf("操作节点失败:%v %v", params, out))
	c.Result(enums.CodeSuccess, "操作节点成功", params.Node)
}
