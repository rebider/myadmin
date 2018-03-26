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
	list, err := models.GetPlatformList("data/json/Platform.json")
	utils.CheckError(err)
	logs.Info("platformList:%v", list)
	c.Result(enums.CodeSuccess, "获取平台列表成功", list)
}


func (c *ToolController) Action() {
	var params struct {
		Action string `json:"action"`
		PlatformId int `json:"platformId"`
		ServerId string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	//node := getNode(params.ServerId)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	node := gameServer.Node
	commandArgs := []string{
		"nodetool",
		"-name",
		node,
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
		c.Result(enums.CodeSuccess, "成功!", 0)
	}
}

func (c *ToolController) SendProp() {
	var params struct {
		PlayerId int `json:"playerId"`
		PropType int `json:"propType"`
		PropId int `json:"propId"`
		PropNum int `json:"propNum"`
		PlatformId int `json:"platformId"`
		ServerId string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	//node := getNode(params.ServerId)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	node := gameServer.Node
	commandArgs := []string{
		"nodetool",
		"-name",
		node,
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

func (c *ToolController)SetTask() {
	var params struct {
		PlayerId int `json:"playerId"`
		TaskId int `json:"taskId"`
		PlatformId int `json:"platformId"`
		ServerId string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	//node := getNode(params.ServerId)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	node := gameServer.Node
	commandArgs := []string{
		"nodetool",
		"-name",
		node,
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

func (c *ToolController)ActiveFunction() {
	var params struct {
		PlayerId int `json:"playerId"`
		FunctionId int `json:"functionId"`
		FunctionParam int `json:"functionParam"`
		FunctionValue int `json:"functionValue"`
		PlatformId int `json:"platformId"`
		ServerId string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	//node := getNode(params.ServerId)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	node := gameServer.Node
	commandArgs := []string{
		"nodetool",
		"-name",
		node,
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
