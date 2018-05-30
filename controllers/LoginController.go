package controllers

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"time"
	"fmt"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Login() {
	var params struct {
		Account  string
		Password string
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	account := params.Account
	password := params.Password
	if len(account) == 0 || len(password) == 0 {
		c.Result(enums.CodeFail, "请输入用户名和密码", "")
	}
	password = utils.String2md5(enums.PasswordSalt + password)
	user, err := models.GetUserOneByAccount(account)
	c.CheckError(err)
	//c.CheckError(err)
	now := utils.GetTimestamp()
	if user.CanLoginTime >= now {
		c.Result(enums.CodeFail, fmt.Sprintf("您已经连续输错5次密码, 请%d秒后重试", user.CanLoginTime-now), user.CanLoginTime-now)
	}

	if user.Password != password {
		user.ContinueLoginErrorTimes += 1
		if user.ContinueLoginErrorTimes >= 5 { // map[string]interface{}{"continue_login_error_times": "0", "can_login_time": now + 180}
			err = models.Db.Model(&user).Updates(map[string]interface{}{"continue_login_error_times": 0, "can_login_time": now + 180}).Error
			c.CheckError(err)
		} else {
			err = models.Db.Model(&user).Updates(map[string]interface{}{"continue_login_error_times": user.ContinueLoginErrorTimes}).Error
			c.CheckError(err)
		}
		c.Result(enums.CodeFail, "密码错误", "")
		//c.CheckError(err, "用户名或者密码错误")
	}
	if user.Status == enums.Disabled {
		c.Result(enums.CodeFail, "用户被禁用，请联系管理员", "")
	}


	//更新用户登录时间
	err = models.Db.Model(&user).Updates(&models.User{LastLoginIp: c.Ctx.Input.IP(), LastLoginTime: int(time.Now().Unix())}).Error
	c.CheckError(err)

	//保存用户信息到session
	c.setUser2Session(user.Id)
	logs.Info("登录成功:%v, %v, %v", user.Id, c.GetSession("user"), c.curUser.Id)

	c.Result(enums.CodeSuccess, "登录成功",
		struct {
			Token string `json:"token"`
		}{Token: c.CruSession.SessionID()})

	//if user != nil && err == nil {
	//	if user.CanLoginTime >= now {
	//		c.Result(enums.CodeFail, "您已经连续输错5次密码, 请稍后重试", user.CanLoginTime-now)
	//	}
	//	if user.Status == enums.Disabled {
	//		c.Result(enums.CodeFail, "用户被禁用，请联系管理员", "")
	//	}
	//	//保存用户信息到session
	//	c.setUser2Session(user.Id)
	//	logs.Info("登录成功:%v, %v, %v", user.Id, c.GetSession("user"), c.curUser.Id)
	//
	//	//更新用户登录时间
	//	err = models.Db.Model(&user).Updates(&models.User{LastLoginIp: c.Ctx.Input.IP(), LastLoginTime: int(time.Now().Unix())}).Error
	//	c.CheckError(err)
	//
	//	c.Result(enums.CodeSuccess, "登录成功",
	//		struct {
	//			Token string `json:"token"`
	//		}{Token: c.CruSession.SessionID()})
	//} else {
	//	userOne, err := models.GetUserOneByAccount(account)
	//	logs.Debug("%+v%+v", account, userOne)
	//	if err == nil {
	//		userOne.ContinueLoginErrorTimes += 1
	//		if userOne.ContinueLoginErrorTimes > 5 {
	//			err = models.Db.Model(&userOne).Updates(&models.User{ContinueLoginErrorTimes: 0, CanLoginTime: now + 180}).Error
	//			c.CheckError(err)
	//		} else {
	//			err = models.Db.Model(&userOne).Updates(&models.User{ContinueLoginErrorTimes: userOne.ContinueLoginErrorTimes}).Error
	//			c.CheckError(err)
	//		}
	//	}
	//	c.Result(enums.CodeFail, "用户名或者密码错误", "")
	//}
}
func (c *LoginController) Logout() {
	//user := models.User{}
	//c.SetSession("user", user)
	c.DelSession("user")
	c.Result(enums.CodeSuccess, "退出登录成功", "")
}
