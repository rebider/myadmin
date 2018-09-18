package controllers

import (
	"github.com/chnzrb/myadmin/models"
	//"github.com/astaxie/beego"
	"github.com/chnzrb/myadmin/utils"
	"encoding/json"
	"fmt"
	"github.com/chnzrb/myadmin/enums"
	"github.com/astaxie/beego/logs"
	//"strconv"
	"strconv"
	//"encoding/base64"
	//"net/http"
	//"io/ioutil"
)

type GameServerController struct {
	BaseController
}

// 获取游戏服列表
func (c *GameServerController) List() {
	var params models.GameServerQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Debug("查询游戏服列表:%+v", params)
	utils.CheckError(err)
	data, total := models.GetGameServerList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取游戏服列表成功", result)
}

// 添加 编辑 游戏服
func (c *GameServerController) Edit() {
	m := models.GameServer{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑游戏服")
	logs.Info("编辑游戏服:%+v", m)
	if m.OpenTime == 0 {
		c.Result(enums.CodeFail, "开服时间不能为0", 0)
	}
	if m.IsAdd == 1 && models.IsGameServerExists(m.PlatformId, m.Sid) {
		c.Result(enums.CodeFail, "游戏服已经存在", m.Node)
	}
	if m.IsAdd == 0 && models.IsGameServerExists(m.PlatformId, m.Sid) == false {
		c.Result(enums.CodeFail, "游戏服不存在", m.Node)
	}

	out, err := models.AddGameServer(m.PlatformId, m.Sid, m.Desc, m.Node, m.ZoneNode, m.State, m.OpenTime, m.IsShow)
	//out, err := utils.NodeTool(
	//	"mod_server_mgr",
	//	"add_game_server",
	//	m.PlatformId,
	//	m.Sid,
	//	m.Desc,
	//	m.Node,
	//	m.ZoneNode,
	//	strconv.Itoa(m.State),
	//	strconv.Itoa(m.OpenTime),
	//	strconv.Itoa(m.IsShow),
	//)

	c.CheckError(err, "保存游戏服失败:"+out)
	c.Result(enums.CodeSuccess, "保存成功", m.Sid)
}

// 删除游戏服
func (c *GameServerController) Delete() {
	var ids [] struct {
		PlatformId string `json:platformId`
		ServerId   string `json:serverId`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ids)
	utils.CheckError(err)
	logs.Info("删除游戏服:%+v", ids)

	for _, id := range ids {

		out, err := utils.NodeTool(
			"mod_server_mgr",
			"delete_game_server",
			id.PlatformId,
			id.ServerId,
		)
		c.CheckError(err, "删除游戏服失败:"+out)
		logs.Info("删除游戏服:%s", id)
	}

	c.Result(enums.CodeSuccess, fmt.Sprintf("成功删除 %d 项", len(ids)), 0)
}

// 批量修改区服状态
func (c *GameServerController) BatchUpdateState() {
	var param struct {
		PlatformId string `json:platformId`
		Nodes []string `json:node`
		State int
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &param)
	utils.CheckError(err)
	logs.Info("批量修改区服状态:%+v", param)

	if len(param.Nodes) == 0 {
		if param.PlatformId == "" {
			logs.Error("平台id不能为空")
			c.Result(enums.CodeFail, "平台id不能为空", 0)
		}
		out, err := utils.NodeTool(
			"mod_server_mgr",
			"update_all_game_server_state",
			param.PlatformId,
			strconv.Itoa(param.State),
		)
		c.CheckError(err, "修改所有区服状态:"+out)
	} else {
		for _, node := range param.Nodes {

			out, err := utils.NodeTool(
				"mod_server_mgr",
				"update_node_state",
				node,
				strconv.Itoa(param.State),
			)
			c.CheckError(err, "批量修改区服状态:"+out)
		}
	}

	c.Result(enums.CodeSuccess, fmt.Sprintf("批量修改区服状态 %d 项", len(param.Nodes)), 0)
}

// 刷新区服入口
func (c *GameServerController) Refresh() {
	var params struct {
	}
	//var result struct {
	//	ErrorCode int
	//}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	//logs.Info("刷新区服入口:%+v", params)
	//
	//data := fmt.Sprintf("time=%d", utils.GetTimestamp())
	//sign := utils.String2md5(data + enums.GmSalt)
	//base64Data := base64.URLEncoding.EncodeToString([]byte(data))
	//
	//baseUrl := beego.AppConfig.String("login_server" + "::url")
	//url := fmt.Sprintf("%s?data=%s&sign=%s", baseUrl, base64Data, sign)
	////url := "http://192.168.31.100:16667/refresh?" + "data=" + base64Data+ "&sign=" + sign
	//
	//logs.Info("url:%s", url)
	//resp, err := http.Get(url)
	//c.CheckError(err)
	//
	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//c.CheckError(err)
	//
	//logs.Info("刷新区服入口 result:%v", string(body))
	//
	//err = json.Unmarshal(body, &result)
	//
	//c.CheckError(err)
	err = models.AfterAddGameServer()
	c.CheckError(err)
	//if result.ErrorCode != 0 {
	//	c.Result(enums.CodeFail, "刷新区服入口失败", 0)
	//}
	c.Result(enums.CodeSuccess, "刷新区服入口成功", 0)

}


// 开服
func (c *GameServerController) OpenServer() {
	var params struct {
		PlatformId string `json:platformId`
		Time int `json:time`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	c.CheckError(err)
	logs.Info("开服:%+v", params)

	openServerTime := 0
	if params.Time == 0 {
		// 立即开服
		openServerTime = utils.GetTimestamp()
	} else {
		//定时开服
		openServerTime = params.Time
	}
	err = models.AutoCreateAndOpenServer(params.PlatformId,false, openServerTime)
	c.CheckError(err)
	c.Result(enums.CodeSuccess, "开服成功", 0)
}

