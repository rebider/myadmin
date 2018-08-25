package routers

import (
	"github.com/chnzrb/myadmin/controllers"
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego/context"
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
	if beego.BConfig.RunMode == "dev" {
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
		beego.SetStaticPath("/swagger", "swagger")
	}

	var FilterNoCache = func(ctx *context.Context) {
		ctx.Output.Header("Cache-Control", "no-cache, no-store")
		ctx.Output.Header("Pragma", "no-cache")
		ctx.Output.Header("Expires", "0")
	}
	beego.InsertFilter("*", beego.BeforeStatic, FilterNoCache)


	beego.SetStaticPath("/static", "views/static")
	beego.SetStaticPath("/", "views")
	//ns := beego.NewNamespace("/v2",
	//	beego.NSNamespace("/resource",
	//		beego.NSInclude(
	//			&controllers.ResourceController{},
	//		),
	//	),
	//)
	//beego.AddNamespace(ns)

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
	beego.Router("/tool/get_ip_origin", &controllers.ToolController{}, "*:GetIpOrigin")

	//游戏服
	beego.Router("/game_server/list", &controllers.GameServerController{}, "*:List")
	beego.Router("/game_server/edit/?:id", &controllers.GameServerController{}, "*:Edit")
	beego.Router("/game_server/delete", &controllers.GameServerController{}, "*:Delete")
	beego.Router("/game_server/refresh", &controllers.GameServerController{}, "*:Refresh")
	beego.Router("/game_server/batch_update_state", &controllers.GameServerController{}, "*:BatchUpdateState")
	beego.Router("/game_server/open_server", &controllers.GameServerController{}, "*:OpenServer")


	//节点
	beego.Router("/server_node/list", &controllers.ServerNodeController{}, "*:List")
	beego.Router("/server_node/edit/?:node", &controllers.ServerNodeController{}, "*:Edit")
	beego.Router("/server_node/delete", &controllers.ServerNodeController{}, "*:Delete")
	//beego.Router("/server_node/start", &controllers.ServerNodeController{}, "*:Start")
	//beego.Router("/server_node/stop", &controllers.ServerNodeController{}, "*:Stop")
	beego.Router("/server_node/action", &controllers.ServerNodeController{}, "*:Action")
	beego.Router("/server_node/install", &controllers.ServerNodeController{}, "*:Install")


	//日志
	beego.Router("/log/login_log/", &controllers.LogController{}, "*:PlayerLoinLogList")
	beego.Router("/log/online_log/", &controllers.LogController{}, "*:PlayerOnlineLogList")
	beego.Router("/log/challenge_mission_log/", &controllers.LogController{}, "*:PlayerChallengeMissionLogList")
	beego.Router("/log/prop_log/", &controllers.LogController{}, "*:PlayerPropLogList")


	//玩家
	beego.Router("/player/list", &controllers.PlayerController{}, "*:List")
	beego.Router("/player/one", &controllers.PlayerController{}, "*:One")
	beego.Router("/player/detail/", &controllers.PlayerController{}, "*:Detail")
	beego.Router("/player/set_account_type/", &controllers.PlayerController{}, "*:SetAccountType")


	//统计
	beego.Router("/statistics/online_statistics/", &controllers.StatisticsController{}, "*:OnlineStatisticsList")
	beego.Router("/statistics/register_statistics/", &controllers.StatisticsController{}, "*:RegisterStatisticsList")
	beego.Router("/statistics/consume_analysis/", &controllers.StatisticsController{}, "*:ConsumeAnalysis")
	beego.Router("/statistics/get_server_generalize/", &controllers.StatisticsController{}, "*:GetServerGeneralize")
	beego.Router("/statistics/real_time_online/", &controllers.StatisticsController{}, "*:GetRealTimeOnline")
	beego.Router("/statistics/get_active_statistics/", &controllers.StatisticsController{}, "*:ActiveStatisticsList")

	//封禁
	beego.Router("/forbid/set_forbid/", &controllers.ForbidController{}, "*:SetForbid")
	beego.Router("/forbid/forbid_log/", &controllers.ForbidController{}, "*:ForbidLogList")

	//公告
	beego.Router("/notice/send_notice/", &controllers.NoticeController{}, "*:SendNotice")
	beego.Router("/notice/notice_log/", &controllers.NoticeController{}, "*:NoticeLogList")
	beego.Router("/notice/del_notice_log/", &controllers.NoticeController{}, "*:DelNoticeLog")

	//登录公告
	beego.Router("/login_notice/set/", &controllers.LoginNoticeController{}, "*:SetNotice")
	beego.Router("/login_notice/batch_set/", &controllers.LoginNoticeController{}, "*:BatchSetNotice")
	beego.Router("/login_notice/list/", &controllers.LoginNoticeController{}, "*:LoginNoticeList")
	beego.Router("/login_notice/del/", &controllers.LoginNoticeController{}, "*:DelLoginNotice")

	//邮件
	beego.Router("/mail/send_mail/", &controllers.MailController{}, "*:SendMail")
	beego.Router("/mail/mail_log/", &controllers.MailController{}, "*:MailLogList")
	beego.Router("/mail/del_mail_log/", &controllers.MailController{}, "*:DelMailLog")

	//平台
	beego.Router("/platform/list/", &controllers.PlatformController{}, "*:List")
	beego.Router("/platform/edit/", &controllers.PlatformController{}, "*:Edit")
	beego.Router("/platform/del/", &controllers.PlatformController{}, "*:Del")

	// 充值
	beego.Router("/charge/charge_list/", &controllers.ChargeController{}, "*:ChargeList")
	beego.Router("/charge/charge_rank/", &controllers.ChargeController{}, "*:ChargeRankList")
	beego.Router("/charge/charge_statistics/", &controllers.ChargeController{}, "*:ChargeStatisticsList")
	beego.Router("/charge/charge_task_distribution/", &controllers.ChargeController{}, "*:ChargeTaskDistribution")
	beego.Router("/charge/charge_activity_distribution/", &controllers.ChargeController{}, "*:ChargeActivityDistribution")
	beego.Router("/charge/charge_money_distribution/", &controllers.ChargeController{}, "*:ChargeMoneyDistribution")
	beego.Router("/charge/charge_level_distribution/", &controllers.ChargeController{}, "*:ChargeLevelDistribution")

	//后台充值
	beego.Router("/charge/background_charge/", &controllers.BackgroundController{}, "*:Charge")
	beego.Router("/charge/background_charge_list/", &controllers.BackgroundController{}, "*:List")

	//留存
	beego.Router("/remain/total_remain/", &controllers.RemainController{}, "*:GetTotalRemain")
	beego.Router("/remain/active_remain/", &controllers.RemainController{}, "*:GetActiveRemain")
	beego.Router("/remain/task_remain/", &controllers.RemainController{}, "*:GetTaskRemain")
	beego.Router("/remain/level_remain/", &controllers.RemainController{}, "*:GetLevelRemain")
	beego.Router("/remain/time_remain/", &controllers.RemainController{}, "*:GetTimeRemain")

	//资产
	beego.Router("/inventory/server_list", &controllers.InventoryServerController{}, "*:ServerList")
	beego.Router("/inventory/edit_server/?:id", &controllers.InventoryServerController{}, "*:EditServer")
	beego.Router("/inventory/delete_server", &controllers.InventoryServerController{}, "*:DeleteServer")
	beego.Router("/inventory/create_ansible_inventory", &controllers.InventoryServerController{}, "*:CreateAnsibleInventory")


	beego.Router("/inventory/database_list", &controllers.InventoryDatabaseController{}, "*:DatabaseList")
	beego.Router("/inventory/edit_database/?:id", &controllers.InventoryDatabaseController{}, "*:EditDatabase")
	beego.Router("/inventory/delete_database", &controllers.InventoryDatabaseController{}, "*:DeleteDatabase")

	//主页
	beego.Router("/", &controllers.MainController{})
}
