package controllers

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
	"strconv"
	"encoding/json"
	"strings"
	"github.com/chnzrb/myadmin/merge"
	"github.com/linclin/gopub/src/github.com/pkg/errors"
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
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	//node := getNode(params.ServerId)
	//gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	//node := gameServer.Node
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	commandArgs := []string{
		"nodetool",
		"-name",
		gameServer.Node,
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
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	//node := getNode(params.ServerId)
	//gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	//node := gameServer.Node
	commandArgs := []string{
		"nodetool",
		"-name",
		gameServer.Node,
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
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	//node := getNode(params.ServerId)
	//gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	//node := gameServer.Node
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	commandArgs := []string{
		"nodetool",
		"-name",
		gameServer.Node,
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

func (c *ToolController) FinishBranchTask() {
	var params struct {
		PlayerId   int    `json:"playerId"`
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	commandArgs := []string{
		"nodetool",
		"-name",
		gameServer.Node,
		"-setcookie",
		"game",
		"rpc",
		"tool",
		"finish_branch_task",
		strconv.Itoa(params.PlayerId),
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
		PlatformId    string `json:"platformId"`
		ServerId      string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("Action:%+v", params)
	//node := getNode(params.ServerId)
	//gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	//node := gameServer.Node
	commandArgs := []string{
		"nodetool",
		"-name",
		gameServer.Node,
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
		Ip string `json:"Ip"`
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

func (c *ToolController) Merge() {
	var params struct {
		Nodes    [] string `json:"nodes"`
		ZoneNode string    `json:"zone_node"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	c.CheckError(err)
	logs.Info("合服:%+v", params)
	//logs.Debug("ip origin:%v %v", params.Ip, o)
	nodeList := make([]*models.ServerNode, 0)
	for _, e := range params.Nodes {
		serverNode, err := models.GetServerNode(e)
		c.CheckError(err)
		nodeList = append(nodeList, serverNode)
	}
	_, err = models.GetServerNode(params.ZoneNode)
	c.CheckError(err)
	err = merge.Merge(nodeList[1:], nodeList[0], params.ZoneNode)
	if err != nil {
		c.Result(enums.CodeFail2, "合服失败:", "")
	} else {
		c.Result(enums.CodeSuccess, "合服成功!", "")
	}
}

func (c *ToolController) GetWeixinArgs() {
	c.Result(enums.CodeSuccess, "成功!", models.GetWeiXinArgs())
}

func (c *ToolController) UpdateWeixinArgs() {
	params :=  &models.WeiXinArgs{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("UpdateWeixinArgs:%+v", params)
	err = models.UpdateWeixinArgs(params)
	c.CheckError(err)
	c.Result(enums.CodeSuccess, "成功!", 0)
}


func (c *ToolController) SqlQuery() {
	type col struct {
		Label      string    `json:"label"`
		Prop    string `json:"prop"`
	}
	var params struct {
		Sql      string    `json:"sql"`
		PlatformId    string `json:"platformId"`
		ServerId      string `json:"serverId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	c.CheckError(err)
	logs.Info("SqlQuery:%+v", params)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	gameDb, err := models.GetVisitorGameDbByNode(gameServer.Node)
	c.CheckError(err)
	defer gameDb.Close()
	rows, err := gameDb.DB().Query(params.Sql)
	c.CheckError(err)
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	records := make([]map[string]string, 0)
	cols := make([]*col, 0)
	for _, e := range columns {
		cols = append(cols, &col{
			Label:e,
			Prop:e,
		})
	}
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		record := make(map[string]string, 0)
		for i, col := range values {
			if col != nil {
				//record = append(record, string(col.([]byte)))
				record[columns[i]] = string(col.([]byte))
			}
		}
		records = append(records, record)
		//fmt.Println(record)
	}
	if len(records) > 200 {
		c.CheckError(errors.New("返回结果超过200行， 请加限制语句！"))
	}
	result := make(map[string]interface{})
	result["cols"] = cols
	result["rows"] = records
	c.Result(enums.CodeSuccess, "成功!", result)
}
