
package routers

import (
	"github.com/chnzrb/myadmin/controllers"

	"github.com/astaxie/beego"
)
// @APIVersion 1.0.0
// @Title 后台管理系统
// @Description documents of server API powered by swagger, you can also generate client code by swagger. refer : https://github.com/swagger-api/swagger-codegen
// @Contact ming.zhao@hobot.cc
// @TermsOfServiceUrl http://www.horizon.ai/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
func init() {
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
	//角色路由
	beego.Router("/role/list", &controllers.RoleController{}, "*:List")
	beego.Router("/role/edit/?:id", &controllers.RoleController{}, "*:Edit")
	beego.Router("/role/delete", &controllers.RoleController{}, "*:Delete")
	beego.Router("/role/allocate", &controllers.RoleController{}, "*:Allocate")


	//资源路由
	beego.Router("/resource/list", &controllers.ResourceController{}, "*:List")
	beego.Router("/resource/edit/?:id", &controllers.ResourceController{}, "*:Edit")
	beego.Router("/resource/getParentResourceList", &controllers.ResourceController{}, "*:GetParentResourceList")
	beego.Router("/resource/delete", &controllers.ResourceController{}, "*:Delete")
	beego.Router("/resource/resourceTree", &controllers.ResourceController{}, "*:ResourceTree")

	//通用选择面板
	//beego.Router("/resource/usermenutree", &controllers.ResourceController{}, "POST:UserMenuTree")
	beego.Router("/resource/checkurlfor", &controllers.ResourceController{}, "POST:CheckUrlFor")

	//后台用户路由
	beego.Router("/user/list", &controllers.UserController{}, "*:List")
	beego.Router("/user/edit/?:id", &controllers.UserController{}, "*:Edit")
	beego.Router("/user/delete", &controllers.UserController{}, "*:Delete")

	//登录
	beego.Router("/login", &controllers.LoginController{}, "*:Login")
	beego.Router("/logout", &controllers.LoginController{}, "*:Logout")

	//后台用户中心
	beego.Router("/user/info", &controllers.UserCenterController{}, "*:Info")
	beego.Router("/user/changePassword", &controllers.UserController{}, "*:ChangePassword")

	//工具
	beego.Router("/tool/build", &controllers.ToolController{}, "*:Build")
	beego.Router("/tool/action", &controllers.ToolController{}, "*:Action")

	//游戏服
	beego.Router("/game_server/list", &controllers.GameServerController{}, "*:List")
	beego.Router("/game_server/edit/?:id", &controllers.GameServerController{}, "*:Edit")
	beego.Router("/game_server/delete", &controllers.GameServerController{}, "*:Delete")

	//节点
	beego.Router("/server_node/list", &controllers.ServerNodeController{}, "*:List")
	beego.Router("/server_node/edit/?:node", &controllers.ServerNodeController{}, "*:Edit")
	beego.Router("/server_node/delete", &controllers.ServerNodeController{}, "Post:Delete")

	//主页
	beego.Router("/", &controllers.MainController{})
}
