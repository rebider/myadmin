package controllers

import (
	"fmt"
	"strings"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	curUser        models.User //当前用户信息
}

func (c *BaseController) AllowCross() {
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:9528")       //允许访问源
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "*")    //允许post访问
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Token") //header的类型
	//c.Ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	//c.Ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
}

//

func (c *BaseController) Prepare() {
	controllerName, actionName := c.GetControllerAndAction()

	//c.AllowCross()
	user := c.GetSession("user")
	if user != nil {
		c.curUser = user.(models.User)
	}
	logs.Info("[%v:%v][%v] 请求:%v.%v", c.curUser.Id, c.curUser.Account, c.Ctx.Input.IP(), controllerName, actionName)

	if controllerName == "LoginController" {
		//登录控制器不判断
	} else {
		//判断是否登录
		c.checkLogin()
		//判断是否有权限
		c.checkAuthor()
	}
}

func (c *BaseController) CheckError(err error, msg... string) {
	if err != nil {
		errMsg := fmt.Sprintf("%s %v", msg, err)
		logs.GetBeeLogger().Error(errMsg)
		c.Result(enums.CodeFail, "[ERROR] " + errMsg, "")
	}
}

//检查是否登录
func (c *BaseController) checkLogin() {
	if c.IsLogin() == false {
		c.Result(enums.CodeNoLogin, "未登录", "")
	}
}

//是否登录
func (c *BaseController) IsLogin() bool{
	return c.curUser.Id > 0
}

//是否帐号有效
func (c *BaseController) IsAccountEnable() bool{
	return c.curUser.Status == enums.Enabled
}

//权限验证
func (c *BaseController) checkAuthor() {
	ignoreAuthorMap := map[string] [] string{
		"LoginController": {"*"},
		"UserController":  {"ChangePassword", "Info"},
	}

	controllerName, actionName := c.GetControllerAndAction()
	ignoreActionList, ok := ignoreAuthorMap[controllerName]
	if ok {
		for _, ignoreAction := range ignoreActionList {
			if ignoreAction == actionName || ignoreAction == "*" {
				return
			}
		}
	}
	isHasAuthor := c.checkActionAuthor(controllerName, actionName)
	if !isHasAuthor {
		//没有权限
		logs.Error(fmt.Sprintf("无权访问!!! 路径: %s.%s, 用户: %v.", controllerName, actionName, c.curUser.Id))
		c.Result(enums.CodeUnauthorized, fmt.Sprintf("无权访问: %v.%v", controllerName, actionName), "")
	}
}

// 判断某 Controller.Action 当前用户是否有权访问
func (c *BaseController) checkActionAuthor(ctrlName, ActName string) bool {
	if c.IsLogin() == false || c.IsAccountEnable() == false{
		return false
	}
	user := c.GetSession("user")
	v, ok := user.(models.User)
	if ok {
		//如果是超级管理员，则直接通过
		if v.IsSuper == 1 {
			return true
		}
		//遍历用户有权限的资源列表
		for _, v := range v.ResourceUrlForList {
			if v == ctrlName+"."+ActName || v == ctrlName+".*"{
				return true
			}
		}
	}
	return false
}

func (c *BaseController) setUser2Session(userId int) error {
	m, err := models.UserOne(userId)
	if err != nil {
		return err
	}
	resourceList := models.GetResourceListByUserId(userId)
	for _, item := range resourceList {
		m.ResourceUrlForList = append(m.ResourceUrlForList, strings.TrimSpace(item.UrlFor))
	}
	c.SetSession("user", *m)
	return nil
}

//返回json
func (c *BaseController) Result(code enums.JsonResultCode, msg string, data interface{}) {
	r := &models.Result{Code:code, Data:data, Msg:msg}
	c.Data["json"] = r
	c.ServeJSON()
	c.StopRun()
}
