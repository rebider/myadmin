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

		c.Result(enums.CodeSuccess, "登录成功",
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
	c.Result(enums.CodeSuccess, "退出登录成功","")
	//c.pageLogin()
}
