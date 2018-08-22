package controllers

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
	"strconv"
	"encoding/json"
	"strings"
)

type ToolController struct {
	BaseController
}

func (c *ToolController) GetJson() {
	list := models.GetPlatformList()
	//utils.CheckError(err)
	logs.Info("platformList:%v", list)
	c.Result(enums.CodeSuccess, "获取平台列表成功 ", list)
}

func (c *ToolController) Action() {
	var params struct {
		Action     string `json:"action"`
		PlatformId string    `json:"platformId"`
		Node       string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	//node := getNode(params.ServerId)
	//gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	//node := gameServer.Node
	commandArgs := []string{
		"nodetool",
		"-name",
		params.Node,
		"-setcookie",
		"game",
		"rpc",
		"tool",
		"project_helper",
		params.Action,
	}
	out, err := utils.Cmd("escript", commandArgs)

	if err != nil {
		out = strings.Replace(out, " ", "&nbsp", -1)
		out = strings.Replace(out, "\n", "<br>", -1)
		out = strings.Replace(out, "\\n", "<br>", -1)
		c.Result(enums.CodeFail2, "失败:"+out+err.Error(), 0)
	} else {
		//if params.Action == "build_table" {
		//	commandArgs = []string{
		//		"ci",
		//		"/opt/h5/trunk/client/client_enum",
		//		"-m",
		//		"web_submit",
		//	}
		//	_, err = utils.Cmd("svn", commandArgs)
		//	c.CheckError(err, "提交客户端枚举")
		//}
		c.Result(enums.CodeSuccess, "成功!", 0)
	}
}

func (c *ToolController) SendProp() {
	var params struct {
		PlayerId   int    `json:"playerId"`
		PropType   int    `json:"propType"`
		PropId     int    `json:"propId"`
		PropNum    int    `json:"propNum"`
		PlatformId string    `json:"platformId"`
		Node       string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	//node := getNode(params.ServerId)
	//gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	//node := gameServer.Node
	commandArgs := []string{
		"nodetool",
		"-name",
		params.Node,
		"-setcookie",
		"game",
		"rpc",
		"tool",
		"give_prop",
		strconv.Itoa(params.PlayerId),
		strconv.Itoa(params.PropType),
		strconv.Itoa(params.PropId),
		strconv.Itoa(params.PropNum),
	}
	out, err := utils.Cmd("escript", commandArgs)

	if err != nil {
		out = strings.Replace(out, " ", "&nbsp", -1)
		out = strings.Replace(out, "\n", "<br>", -1)
		out = strings.Replace(out, "\\n", "<br>", -1)
		c.Result(enums.CodeFail2, "失败:"+out+err.Error(), 0)
	} else {
		c.Result(enums.CodeSuccess, "成功!", 0)
	}
}

func (c *ToolController) SetTask() {
	var params struct {
		PlayerId   int    `json:"playerId"`
		TaskId     int    `json:"taskId"`
		PlatformId string    `json:"platformId"`
		Node       string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	//node := getNode(params.ServerId)
	//gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	//node := gameServer.Node
	commandArgs := []string{
		"nodetool",
		"-name",
		params.Node,
		"-setcookie",
		"game",
		"rpc",
		"tool",
		"debug_set_task",
		strconv.Itoa(params.PlayerId),
		strconv.Itoa(params.TaskId),
	}
	out, err := utils.Cmd("escript", commandArgs)

	if err != nil {
		out = strings.Replace(out, " ", "&nbsp", -1)
		out = strings.Replace(out, "\n", "<br>", -1)
		out = strings.Replace(out, "\\n", "<br>", -1)
		c.Result(enums.CodeFail2, "失败:"+out+err.Error(), 0)
	} else {
		c.Result(enums.CodeSuccess, "成功!", 0)
	}
}

func (c *ToolController) ActiveFunction() {
	var params struct {
		PlayerId      int    `json:"playerId"`
		FunctionId    int    `json:"functionId"`
		FunctionParam int    `json:"functionParam"`
		FunctionValue int    `json:"functionValue"`
		PlatformId    string    `json:"platformId"`
		Node          string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	//node := getNode(params.ServerId)
	//gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	//node := gameServer.Node
	commandArgs := []string{
		"nodetool",
		"-name",
		params.Node,
		"-setcookie",
		"game",
		"rpc",
		"tool",
		"active_function",
		strconv.Itoa(params.PlayerId),
		strconv.Itoa(params.FunctionId),
		strconv.Itoa(params.FunctionParam),
		strconv.Itoa(params.FunctionValue),
	}
	out, err := utils.Cmd("escript", commandArgs)

	if err != nil {
		out = strings.Replace(out, " ", "&nbsp", -1)
		out = strings.Replace(out, "\n", "<br>", -1)
		out = strings.Replace(out, "\\n", "<br>", -1)
		c.Result(enums.CodeFail2, "失败:"+out+err.Error(), 0)
	} else {
		c.Result(enums.CodeSuccess, "成功!", 0)
	}
}


func (c *ToolController) GetIpOrigin() {
	var params struct {
		Ip string    `json:"Ip"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	c.CheckError(err)
	o := utils.GetIpLocation(params.Ip)
	//logs.Debug("ip origin:%v %v", params.Ip, o)
	if err != nil {
		c.Result(enums.CodeFail2, "失败:", "")
	} else {
		c.Result(enums.CodeSuccess, "成功!", o)
	}
}
