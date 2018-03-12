package controllers

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
	//"github.com/astaxie/beego/orm"
	"fmt"
	"strings"
	"github.com/astaxie/beego/orm"
	//"os/exec"
	//"bytes"
	//"os"
)

type ToolController struct {
	BaseController
}

//func (c *ToolController) Prepare() {
//	//先执行
//	c.BaseController.Prepare()
//	logs.Debug("ToolController Prepare()")
//	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
//	c.checkLogin()
//}
func (c *ToolController) Build() {
	//获取数据列表和总数
	//gameServerList, _ := models.GetAllGameServer()
	platformId :=  c.Ctx.GetCookie(enums.ChosePlatformId)
	if platformId == "" {
		platformId = "0"
	}
	serverId :=  c.Ctx.GetCookie(enums.ChoseServerId)
	var params models.GameServerQueryParam
	params.PlatformId = platformId
	//获取数据列表和总数
	gameServerList, _ := models.GetGameServerList(&params)

	//if platformId == "" {
	//	//c.pageError("数据无效，请刷新后重试")
	//}
	//if serverId != "" {
	//	c.Ctx.SetCookie("serverId", serverId)
	//} else {
	//	serverId = c.Ctx.GetCookie("serverId")
	//}
	//c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.Data["pageTitle"] = "调试工具"
	c.Data["game_server_list"] = gameServerList
	models.ShowPlatformList(c.Data)
	logs.Debug("serverId:%v, platformId:%v", serverId, platformId)
	//logs.Debug(platformId)
	c.Data[enums.ChoseServerId] = serverId
	c.Data[enums.ChosePlatformId] = platformId
	//c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "tool/tool_headcssjs.html"
	c.LayoutSections["footerjs"] = "tool/tool_footerjs.html"
}

func (this *ToolController) Action() {
	//if this.isPost() {

	action := this.GetString("action")
	desc := this.GetString("desc")
	serverId := this.GetString("serverId")
	fmt.Println("Action:", action, desc)
	this.Data["redirect"] = "/tool/build"
	//fmt.Println("serverId:", serverId)
	node := getNode(serverId)
	commandArgs := []string{
		"nodetool",
		"-name",
		node.(string),
		"-setcookie",
		"game",
		"rpc",
		"tool",
		"project_helper",
		action,
	}
	out, err := utils.Cmd("escript", commandArgs)

	if err != nil {
		out = strings.Replace(out, " ", "&nbsp", -1)
		out = strings.Replace(out, "\n", "<br>", -1)
		out = strings.Replace(out, "\\n", "<br>", -1)
		this.Result(enums.MSG_ERR, desc+"失败:<br>"+out+"<br>"+err.Error(), 0)
	} else {
		this.Result(enums.MSG_OK, desc+"成功!", 0)
	}
}

func getNode(serverId string) interface{} {
	var o = orm.NewOrm()
	var lists []orm.ParamsList
	o.Using("center")
	o.Raw("select `node` from c_game_server where `sid` = ?", serverId).ValuesList(&lists)
	node := lists[0][0]
	return node
}
