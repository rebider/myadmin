package controllers

import (
	"github.com/chnzrb/myadmin/models"
	"encoding/json"
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"github.com/chnzrb/myadmin/enums"
	//"strconv"
	"github.com/astaxie/beego/logs"
	//"os"
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
	c.Result(enums.CodeSuccess, "获取节点列表成功", result)
}

//  添加 编辑 节点
func (c *ServerNodeController) Edit() {
	m := models.ServerNode{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	logs.Debug("编辑 节点:%v", m)
	utils.CheckError(err, "编辑节点")
	if m.IsAdd == 1 && models.IsServerNodeExists(m.Node) {
		c.Result(enums.CodeFail, "节点已经存在", m.Node)
	}
	if m.IsAdd == 0 && models.IsServerNodeExists(m.Node) == false {
		c.Result(enums.CodeFail, "节点不存在", m.Node)
	}

	out, err := models.AddServerNode(m.Node, m.Ip, m.Port, m.WebPort, m.Type, m.PlatformId, m.DbHost, m.DbPort, m.DbName)

	//out, err := utils.NodeTool(
	//	"mod_server_mgr",
	//	"add_server_node",
	//	m.Node,
	//	m.Ip,
	//	strconv.Itoa(m.Port),
	//	strconv.Itoa(m.WebPort),
	//	strconv.Itoa(m.Type),
	//	m.PlatformId,
	//	m.DbHost,
	//	strconv.Itoa(m.DbPort),
	//	m.DbName,
	//)
	c.CheckError(err, "保存节点失败:"+out)
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
		c.CheckError(err, "删除节点失败:"+out)
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
		Nodes  [] string `json:"nodes"`
		Action string
	}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Debug("节点操作666:%v", params)

	c.CheckError(err)

	err = models.NodeAction(params.Nodes, params.Action)
	c.CheckError(err)
	//curDir := utils
	// .GetCurrentDirectory()
	//defer os.Chdir(curDir)
	//toolDir := utils.GetToolDir()
	//err = os.Chdir(toolDir)
	//c.CheckError(err)
	//var commandArgs []string
	//for _, node := range params.Nodes {
	//	switch params.Action {
	//	case "start":
	//		commandArgs = []string{"node_tool.sh", node, params.Action,}
	//	case "stop":
	//		commandArgs = []string{"node_tool.sh", node, params.Action,}
	//	case "hot_reload":
	//		commandArgs = []string{"node_hot_reload.sh", node, "server",}
	//	case "cold_reload":
	//		commandArgs = []string{"node_cold_reload.sh", node, "server",}
	//	}
	//	out, err := utils.Cmd("sh", commandArgs)
	//	c.CheckError(err, fmt.Sprintf("操作节点失败:%v %v", params, out))
	//}

	c.Result(enums.CodeSuccess, "操作节点成功", "")
}

func (c *ServerNodeController) Install() {
	var params struct {
		Node string `json:"node"`
	}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Info("部署节点:%v", params)
	c.CheckError(err)
	//curDir := utils.GetCurrentDirectory()
	//defer os.Chdir(curDir)
	//toolDir := utils.GetToolDir()
	//err = os.Chdir(toolDir)
	//c.CheckError(err)
	//var commandArgs []string
	//serverNode, err := models.GetServerNode(params.Node)
	//c.CheckError(err, "节点不存在:"+params.Node)
	//app := ""
	//switch serverNode.Type {
	//case 0:
	//	app = "center"
	//case 1:
	//	app = "game"
	//case 2:
	//	app = "zone"
	//case 4:
	//	app = "login_server"
	//case 5:
	//	app = "unique_id"
	//case 6:
	//	app = "charge"
	//}
	//
	//commandArgs = []string{"/data/tool/ansible/do-install.sh", serverNode.Node, app,serverNode.DbName, serverNode.DbHost, strconv.Itoa(serverNode.DbPort), "root"}
	//out, err := utils.Cmd("sh", commandArgs)
	//c.CheckError(err, fmt.Sprintf("操作节点失败:%v %v", params, out))
	//logs.Info("部署节点成功:%v", params)
	err = models.InstallNode(params.Node)
	c.CheckError(err)
	c.Result(enums.CodeSuccess, "部署节点成功", params.Node)
}
