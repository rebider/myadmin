
package routers

import (
	"github.com/chnzrb/myadmin/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)
// @APIVersion 1.0.0
// @Title 后台管理系统
// @Description documents of server API powered by swagger, you can also generate client code by swagger. refer : https://github.com/swagger-api/swagger-codegen
// @Contact ming.zhao@hobot.cc
// @TermsOfServiceUrl http://www.horizon.ai/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//AllowAllOrigins:  true,
		AllowOrigins:     []string{"http://localhost:9528"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*", "Origin", "Authorization", "Cookie", "Host", "Referer", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "x-token"},
		ExposeHeaders:    []string{"*", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.InsertFilter("*", beego.BeforeStatic, cors.Allow(&cors.Options{
		//AllowAllOrigins:  true,
		AllowOrigins:     []string{"http://localhost:9528"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*", "Origin", "Authorization", "Cookie", "Host", "Referer", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "x-token"},
		ExposeHeaders:    []string{"*", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.SetStaticPath("/static", "views/static")
	beego.SetStaticPath("/swagger", "swagger")

	ns := beego.NewNamespace("/v2",
		beego.NSNamespace("/resource",
			beego.NSInclude(
				&controllers.ResourceController{},
			),
		),
	)
	beego.AddNamespace(ns)


	//角色
	beego.Router("/role/list", &controllers.RoleController{}, "*:List")
	beego.Router("/role/edit/?:id", &controllers.RoleController{}, "*:Edit")
	beego.Router("/role/delete", &controllers.RoleController{}, "*:Delete")
	beego.Router("/role/allocateResource", &controllers.RoleController{}, "*:AllocateResource")
	beego.Router("/role/allocateMenu", &controllers.RoleController{}, "*:AllocateMenu")

	//资源
	beego.Router("/resource/list", &controllers.ResourceController{}, "*:List")
	beego.Router("/resource/edit/?:id", &controllers.ResourceController{}, "*:Edit")
	beego.Router("/resource/getParentResourceList", &controllers.ResourceController{}, "*:GetParentResourceList")
	beego.Router("/resource/delete", &controllers.ResourceController{}, "*:Delete")
	beego.Router("/resource/resourceTree", &controllers.ResourceController{}, "*:ResourceTree")
	beego.Router("/resource/checkurlfor", &controllers.ResourceController{}, "POST:CheckUrlFor")

	//菜单
	beego.Router("/menu/list", &controllers.MenuController{}, "*:List")
	beego.Router("/menu/edit/?:id", &controllers.MenuController{}, "*:Edit")
	beego.Router("/menu/getParentMenuList", &controllers.MenuController{}, "*:GetParentMenuList")
	beego.Router("/menu/delete", &controllers.MenuController{}, "*:Delete")
	beego.Router("/menu/menuTree", &controllers.MenuController{}, "*:MenuTree")

	//用户
	beego.Router("/user/list", &controllers.UserController{}, "*:List")
	beego.Router("/user/edit/?:id", &controllers.UserController{}, "*:Edit")
	beego.Router("/user/delete", &controllers.UserController{}, "*:Delete")
	beego.Router("/user/info", &controllers.UserController{}, "*:Info")
	beego.Router("/user/changePassword", &controllers.UserController{}, "*:ChangePassword")

	//登录
	beego.Router("/login", &controllers.LoginController{}, "*:Login")
	beego.Router("/logout", &controllers.LoginController{}, "*:Logout")

	//工具
	beego.Router("/tool/action", &controllers.ToolController{}, "*:Action")
	beego.Router("/tool/send_prop", &controllers.ToolController{}, "*:SendProp")
	beego.Router("/tool/set_task", &controllers.ToolController{}, "*:SetTask")
	beego.Router("/tool/active_function", &controllers.ToolController{}, "*:ActiveFunction")
	//游戏服
	beego.Router("/game_server/list", &controllers.GameServerController{}, "*:List")
	beego.Router("/game_server/edit/?:id", &controllers.GameServerController{}, "*:Edit")
	beego.Router("/game_server/delete", &controllers.GameServerController{}, "*:Delete")

	//节点
	beego.Router("/server_node/list", &controllers.ServerNodeController{}, "*:List")
	beego.Router("/server_node/edit/?:node", &controllers.ServerNodeController{}, "*:Edit")
	beego.Router("/server_node/delete", &controllers.ServerNodeController{}, "*:Delete")

	//玩家
	beego.Router("/player/list", &controllers.PlayerController{}, "*:List")
	beego.Router("/player/detail/", &controllers.PlayerController{}, "*:Detail")
	beego.Router("/player/login_log/", &controllers.PlayerController{}, "*:PlayerLoinLogList")

	//主页
	beego.Router("/", &controllers.MainController{})
}
