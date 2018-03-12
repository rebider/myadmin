package controllers

import (
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"encoding/json"
	"fmt"
	"github.com/chnzrb/myadmin/enums"
	"github.com/astaxie/beego/logs"
	"strconv"
	//"strings"
)

type GameServerController struct {
	BaseController
}
//
//func (c *GameServerController) Prepare() {
//	//先执行
//	c.BaseController.Prepare()
//	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
//	c.checkLogin()
//}


func (c *GameServerController) List() {
	//直接反序化获取json格式的requestbody里的值
	var params models.GameServerQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	//获取数据列表和总数
	data, total := models.GetGameServerList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取游戏服列表成功",result)
}


// Edit 添加 编辑 页面
func (c *GameServerController) Edit() {
	m := models.GameServer{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑游戏服")
	logs.Info("编辑游戏服:%+v", m)
	//fmt.Println("Save:", m, c.Input())
	out, err := utils.Nodetool(
		"mod_server_mgr",
		"add_game_server",
		strconv.Itoa(m.PlatformId),
		m.Sid,
		m.Desc,
		m.Node)

	if err != nil {
		fmt.Println("保存失败:"+ out)
		c.Result(enums.CodeFail, "保存失败:" + out, m.Sid)
	}

	logs.Info("修改游戏服:%s", m)
	c.Result(enums.CodeSuccess, "保存成功", m.Sid)
	////如果是Post请求，则由Save处理
	//if c.Ctx.Request.Method == "POST" {
	//	c.Save()
	//}
	//Id := c.GetString(":id", "")
	////fmt.Println("Edit", c.Ctx.Request.Method, Id)
	//m := &models.GameServer{}
	//var err error
	//if Id != "0" {
	//	m, err = models.GetGameServer(1, Id)
	//	fmt.Printf("%#v", m)
	//	if err != nil {
	//		fmt.Println("err", err)
	//		logs.Error(err)
	//		//c.pageError("数据无效，请刷新后重试")
	//	}
	//}
	//c.Data["m"] = m
	//models.ShowPlatformJson(c.Data)
	//models.ShowGameNodeJson(c.Data)
	////c.setTpl("gameserver/edit.html", "shared/layout_pullbox.html")
	//c.LayoutSections = make(map[string]string)
	//c.LayoutSections["footerjs"] = "gameserver/edit_footerjs.html"
}
//func (c *GameServerController) Save() {
//
//	m := models.GameServer{}
//	//获取form里的值
//	if err := c.ParseForm(&m); err != nil {
//		c.Result(enums.CodeFail, "获取数据失败", m.Sid)
//	}
//	//fmt.Println("Save:", m, c.Input())
//	out, err := utils.Nodetool(
//		"mod_server_mgr",
//		"add_game_server",
//		strconv.Itoa(m.PlatformId),
//		m.Sid,
//		m.Desc,
//		m.Node)
//
//	if err != nil {
//		fmt.Println("保存失败:"+ out)
//		c.Result(enums.CodeFail, "保存失败:" + out, m.Sid)
//	}
//
//	logs.Info("修改游戏服:%s", m)
//	c.Result(enums.CodeSuccess, "保存成功", m.Sid)
//}

func (c *GameServerController) Delete() {
	var ids []int
	json.Unmarshal(c.Ctx.Input.RequestBody, &ids)
	logs.Info("删除游戏服:%+v", ids)

	for _, id := range ids {

		out, err := utils.Nodetool(
			"mod_server_mgr",
			"delete_game_server",
			"1",
			strconv.Itoa(id),
		)
		if err != nil {
			fmt.Println("删除失败:", ids, out, err)
			c.Result(enums.CodeFail, "删除失败:" + out, 0)
		}
		logs.Info("删除游戏服:%s", id)
	}

	c.Result(enums.CodeSuccess, fmt.Sprintf("成功删除 %d 项", len(ids)), 0)
}
