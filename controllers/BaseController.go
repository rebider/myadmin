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

func (c *BaseController) Prepare() {
	c.AllowCross()
	//从Session里获取数据 设置用户信息
	c.adapterUserInfo()
}

func (c *BaseController) CheckError(err error, msg... string) {
	if err != nil {
		errMsg := fmt.Sprintf("%s %v", msg, err)
		logs.GetBeeLogger().Error(errMsg)
		c.Result(enums.CodeFail, "[ERROR] " + errMsg, "")
	}
}

// checkLogin判断用户是否登录，未登录则跳转至登录页面
// 一定要在BaseController.Prepare()后执行
func (c *BaseController) checkLogin() {
	if c.IsLogin() == false {
		c.Result(enums.CodeNoLogin, "未登录", "")
	}
}

//是否登录
func (c *BaseController) IsLogin() bool{
	return c.curUser.Id > 0
}

// 判断某 Controller.Action 当前用户是否有权访问
func (c *BaseController) checkActionAuthor(ctrlName, ActName string) bool {
	if c.IsLogin() == false {
		return false
	}
	//从session获取用户信息
	user := c.GetSession("user")
	//类型断言
	v, ok := user.(models.User)
	if ok {
		//如果是超级管理员，则直接通过
		if v.IsSuper == 1 {
			return true
		}
		//logs.Debug("checkActionAuthor: %+v%+v%+v", ctrlName, ActName, v.ResourceUrlForList)
		//遍历用户所负责的资源列表
		for _, v := range v.ResourceUrlForList {
			if v == ctrlName+"."+ActName {
				return true
			}
		}
	}
	return false
}

// checkLogin判断用户是否有权访问某地址，无权则会跳转到错误页面
//一定要在BaseController.Prepare()后执行
// 会调用checkLogin
// 传入的参数为忽略权限控制的Action
func (c *BaseController) checkAuthor(ignores ...string) {
	logs.Debug("权限验证")
	//先判断是否登录
	c.checkLogin()
	//如果Action在忽略列表里，则直接通用
	controllerName, actionName := c.GetControllerAndAction()
	for _, ignore := range ignores {
		if ignore == actionName {
			return
		}
	}

	hasAuthor := c.checkActionAuthor(controllerName, actionName)
	if !hasAuthor {
		logs.Error(fmt.Sprintf("无权访问!!! 路径: %s.%s, 用户: %v.", controllerName, actionName, c.curUser.Id))
		//如果没有权限
		c.Result(enums.CodeUnauthorized, "无权访问", "")
	}
}

//从session里取用户信息
func (c *BaseController) adapterUserInfo() {
	a := c.GetSession("user")
	if a != nil {
		c.curUser = a.(models.User)
	}
}

//SetUser2Session 获取用户信息（包括资源UrlFor）保存至Session
func (c *BaseController) setUser2Session(userId int) error {
	m, err := models.UserOne(userId)
	if err != nil {
		return err
	}
	//获取这个用户能获取到的所有资源列表
	resourceList := models.GetResourceListByUserId(userId, 1)
	for _, item := range resourceList {
		m.ResourceUrlForList = append(m.ResourceUrlForList, strings.TrimSpace(item.UrlFor))
	}
	c.SetSession("user", *m)
	return nil
}

func (c *BaseController) Result(code enums.JsonResultCode, msg string, data interface{}) {
	r := &models.Result{Code:code, Data:data, Msg:msg}
	c.Data["json"] = r
	c.ServeJSON()
	c.StopRun()
}
