package controllers

import (
	//"strings"

	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	//"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type HomeController struct {
	BaseController
}

//func (c *HomeController) Index() {
//	//判断是否登录
//	c.checkLogin()
//	c.setTpl()
//}
//func (c *HomeController) Page404() {
//	c.setTpl()
//}
//func (c *HomeController) Error() {
//	c.Data["error"] = c.GetString(":error")
//	c.setTpl("home/error.html", "shared/layout_pullbox.html")
//}
//func (c *HomeController) Login() {
//
//	c.LayoutSections = make(map[string]string)
//	c.LayoutSections["headcssjs"] = "home/login_headcssjs.html"
//	c.LayoutSections["footerjs"] = "home/login_footerjs.html"
//	c.setTpl("home/login.html", "shared/layout_base.html")
//}
func (c *HomeController) ChangePlatformId() {
	platformId := c.GetString(enums.ChosePlatformId)
	logs.Debug("chose_platform_id:%v", platformId)
	c.Ctx.SetCookie(enums.ChosePlatformId, platformId)
	c.Ctx.SetCookie(enums.ChoseServerId, "0")
	c.ServeJSON()
	//c.StopRun()
}
func (c *HomeController) ChangeServerId() {
	serverId := c.GetString(enums.ChoseServerId)
	c.Ctx.SetCookie(enums.ChoseServerId, serverId)
	c.ServeJSON()
	//c.StopRun()
}

func (c *HomeController) DoLogin() {
	var params struct {
		Account string
		Password string
	}
	//c.AllowCross()
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Debug("登录:%v", params)
	account := params.Account
	password := params.Password
	//username := "admin"
	//userpwd := "123456"
	if len(account) == 0 || len(password) == 0 {
		c.Result(enums.JRCodeFailed, "用户名和密码不正确", "")
	}
	password = utils.String2md5(password)
	user, err := models.UserOneByAccount(account, password)
	if user != nil && err == nil {
		if user.Status == enums.Disabled {
			c.Result(enums.JRCodeFailed, "用户被禁用，请联系管理员", "")
		}
		//保存用户信息到session
		c.setUser2Session(user.Id)
		//获取用户信息
		//type Data struct {
		//	token string
		//}
		//logs.Info("CruSession:%v",   c.GetSession())
		logs.Info("CruSession:%v",   c.CruSession.SessionID())
		logs.Info("登录成功:%v, %v, %v", user.Id, c.GetSession("user"), c.curUser.Id)
		c.Ctx.SetCookie("myadminsessionid", c.CruSession.SessionID())

		o := orm.NewOrm()
		user.LastLoginTime = int(time.Now().Unix())
		user.LastLoginIp = c.Ctx.Input.IP()

		_, err := o.Update(user)
		utils.CheckError(err)

		c.Result(enums.Success, "登录成功",
			struct {
				Token string `json:"token"`
			}{Token: c.CruSession.SessionID()})
	} else {
		c.Result(enums.JRCodeFailed, "用户名或者密码错误", "")
	}
}
func (c *HomeController) Logout() {
	user := models.User{}
	c.SetSession("user", user)
	c.Result(enums.Success, "退出登录成功","")
	//c.pageLogin()
}
