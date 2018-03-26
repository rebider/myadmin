package controllers

import (
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"encoding/json"
	"fmt"
	"github.com/chnzrb/myadmin/enums"
	"github.com/astaxie/beego/logs"
	"strconv"
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
	c.Result(enums.CodeSuccess, "获取游戏服列表成功",result)
}


// 添加 编辑 游戏服
func (c *GameServerController) Edit() {
	m := models.GameServer{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑游戏服")
	logs.Info("编辑游戏服:%+v", m)
	out, err := utils.Nodetool(
		"mod_server_mgr",
		"add_game_server",
		strconv.Itoa(m.PlatformId),
		m.Sid,
		m.Desc,
		m.Node)

	c.CheckError(err, "保存游戏服失败:" + out)
	c.Result(enums.CodeSuccess, "保存成功", m.Sid)
}

// 删除游戏服
func (c *GameServerController) Delete() {
	var ids []string
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ids)
	utils.CheckError(err)
	logs.Info("删除游戏服:%+v", ids)

	for _, id := range ids {

		out, err := utils.Nodetool(
			"mod_server_mgr",
			"delete_game_server",
			"1",
			id,
		)
		c.CheckError(err, "删除游戏服失败:" + out)
		logs.Info("删除游戏服:%s", id)
	}

	c.Result(enums.CodeSuccess, fmt.Sprintf("成功删除 %d 项", len(ids)), 0)
}
