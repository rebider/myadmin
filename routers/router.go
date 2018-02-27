// @APIVersion 1.0.0
// @Title mobile API
// @Description mobile has every tool to get any job done, so codename for the new mobile APIs.
// @Contact astaxie@gmail.com
package routers

import (
	"github.com/chnzrb/myadmin/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//beego.SetStaticPath("/public", "public")
	//beego.SetStaticPath("/", "static")
	beego.SetStaticPath("/static", "views/static")
	beego.SetStaticPath("/swagger", "swagger")

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/gameserver",
			beego.NSInclude(
				&controllers.GameServerController{},
			),
		),
	)
	beego.AddNamespace(ns)
	//角色路由
	//beego.Router("/role/index", &controllers.RoleController{}, "*:Index")
	beego.Router("/role/roleList", &controllers.RoleController{}, "*:RoleList")
	beego.Router("/role/edit/?:id", &controllers.RoleController{}, "*:Edit")
	beego.Router("/role/delete", &controllers.RoleController{}, "*:Delete")
	beego.Router("/role/datalist", &controllers.RoleController{}, "Post:DataList")
	beego.Router("/role/allocate", &controllers.RoleController{}, "Post:Allocate")
	beego.Router("/role/updateseq", &controllers.RoleController{}, "Post:UpdateSeq")

	//资源路由
	beego.Router("/resource/index", &controllers.ResourceController{}, "*:Index")
	beego.Router("/resource/treegrid", &controllers.ResourceController{}, "POST:TreeGrid")
	beego.Router("/resource/edit/?:id", &controllers.ResourceController{}, "Get,Post:Edit")
	beego.Router("/resource/parent", &controllers.ResourceController{}, "Post:ParentTreeGrid")
	beego.Router("/resource/delete", &controllers.ResourceController{}, "Post:Delete")
	//快速修改顺序
	beego.Router("/resource/updateseq", &controllers.ResourceController{}, "Post:UpdateSeq")

	//通用选择面板
	beego.Router("/resource/select", &controllers.ResourceController{}, "Get:Select")
	//用户有权管理的菜单列表（包括区域）
	beego.Router("/resource/usermenutree", &controllers.ResourceController{}, "POST:UserMenuTree")
	beego.Router("/resource/checkurlfor", &controllers.ResourceController{}, "POST:CheckUrlFor")

	//后台用户路由
	//beego.Router("/user/index", &controllers.UserController{}, "*:Index")
	beego.Router("/user/userList", &controllers.UserController{}, "*:DataGrid")
	beego.Router("/user/edit/?:id", &controllers.UserController{}, "*,Post:Edit")
	beego.Router("/user/delete", &controllers.UserController{}, "*:Delete")


	//后台用户中心
	beego.Router("/user/info", &controllers.UserCenterController{}, "*:Info")
	beego.Router("/user/login", &controllers.HomeController{}, "*:DoLogin")
	beego.Router("/user/logout", &controllers.HomeController{}, "*:Logout")

	//beego.Router("/usercenter/profile", &controllers.UserCenterController{}, "Get:Profile")
	//beego.Router("/usercenter/basicinfosave", &controllers.UserCenterController{}, "Post:BasicInfoSave")
	//beego.Router("/usercenter/uploadimage", &controllers.UserCenterController{}, "Post:UploadImage")
	beego.Router("/user/changePasswd", &controllers.UserCenterController{}, "Post:ChangePasswd")




	//beego.Router("/home/index", &controllers.HomeController{}, "*:Index")
	//beego.Router("/home/login", &controllers.HomeController{}, "*:Login")
	beego.Router("/home/dologin", &controllers.HomeController{}, "*:DoLogin")
	beego.Router("/home/logout", &controllers.HomeController{}, "*:Logout")

	//beego.Router("/home/404", &controllers.HomeController{}, "*:Page404")
	//beego.Router("/home/error/?:error", &controllers.HomeController{}, "*:Error")

	beego.Router("/home/ChangePlatformId", &controllers.HomeController{}, "POST:ChangePlatformId")
	beego.Router("/home/ChangeServerId", &controllers.HomeController{}, "POST:ChangeServerId")

	//工具
	beego.Router("/tool/build", &controllers.ToolController{}, "*:Build")
	beego.Router("/tool/action", &controllers.ToolController{}, "*:Action")

	//游戏服
	beego.Router("/game_server/list", &controllers.GameServerController{}, "*:List")
	beego.Router("/game_server/datagrid", &controllers.GameServerController{}, "POST:DataGrid")
	beego.Router("/game_server/edit/?:id", &controllers.GameServerController{}, "*:Edit")
	beego.Router("/game_server/delete", &controllers.GameServerController{}, "Post:Delete")

	//节点
	beego.Router("/server_node/list", &controllers.ServerNodeController{}, "*:List")
	beego.Router("/server_node/datagrid", &controllers.ServerNodeController{}, "POST:DataGrid")
	beego.Router("/server_node/edit/?:node", &controllers.ServerNodeController{}, "*:Edit")
	beego.Router("/server_node/delete", &controllers.ServerNodeController{}, "Post:Delete")
	//主页
	beego.Router("/", &controllers.MainController{})
	//beego.Router("*", &controllers.MainController{})
}
