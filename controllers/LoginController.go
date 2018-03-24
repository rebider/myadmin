package controllers

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	//"github.com/astaxie/beego/orm"
	"time"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Login() {
	var params struct {
		Account string
		Password string
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	account := params.Account
	password := params.Password
	if len(account) == 0 || len(password) == 0 {
		c.Result(enums.CodeFail, "请输入用户名和密码", "")
	}
	password = utils.String2md5(password)
	user, err := models.UserOneByAccount(account, password)
	utils.CheckError(err)
	if user != nil && err == nil {
		if user.Status == enums.Disabled {
			c.Result(enums.CodeFail, "用户被禁用，请联系管理员", "")
		}
		//保存用户信息到session
		c.setUser2Session(user.Id)
		logs.Info("登录成功:%v, %v, %v", user.Id, c.GetSession("user"), c.curUser.Id)

		//更新用户登录时间
		err = models.Db.Model(&user).Updates(&models.User{LastLoginIp:c.Ctx.Input.IP(), LastLoginTime:int(time.Now().Unix())}).Error
		utils.CheckError(err)


		c.Result(enums.CodeSuccess, "登录成功",
			struct {
				Token string `json:"token"`
			}{Token: c.CruSession.SessionID()})
	} else {
		c.Result(enums.CodeFail, "用户名或者密码错误", "")
	}
}
func (c *LoginController) Logout() {
	user := models.User{}
	c.SetSession("user", user)
	c.Result(enums.CodeSuccess, "退出登录成功","")
}
