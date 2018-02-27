package controllers

import (
	"fmt"
	"strings"
	//"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	//"github.com/chnzrb/myadmin/utils"
)

type BaseController struct {
	beego.Controller
	curUser        models.User //当前用户信息
}
//func (c *BaseController) Init() {
//	logs.Debug("Init")
//}

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
	//controllerName, actionName := c.GetControllerAndAction()
	//从Session里获取数据 设置用户信息
	c.adapterUserInfo()
	//userId := c.curUser.Id
	//logs.Debug("BaseController Prepare()")
	//if userId > 0 {
	//	c.Data["activeSidebarUrl"] = c.URLFor(controllerName + "." + actionName)
	//	c.ShowSidebar(userId)
	//}
}
//func (c *BaseController) ShowSidebar(userId int) {
//	tree := models.ResourceTreeGridByUserId(userId, 1)
//	for _, item := range tree {
//		item.LinkUrl = c.UrlFor2LinkOne(item.UrlFor)
//	}
//	out,err := json.Marshal(tree)
//	utils.CheckError(err)
//	c.Data["sidebarJson"] = string(out)
//}
//func (c *BaseController) UrlFor2Link(src []*models.Resource) {
//	for _, item := range src {
//		item.LinkUrl = c.UrlFor2LinkOne(item.UrlFor)
//	}
//}
//
//func (c *BaseController) UrlFor2LinkOne(urlfor string) string {
//	if len(urlfor) == 0 {
//		return ""
//	}
//	// ResourceController.Edit,:id,1
//	strs := strings.Split(urlfor, ",")
//	if len(strs) == 1 {
//		return c.URLFor(strs[0])
//	} else if len(strs) > 1 {
//		var values []interface{}
//		for _, val := range strs[1:] {
//			values = append(values, val)
//		}
//		return c.URLFor(strs[0], values...)
//	}
//	return ""
//}


// checkLogin判断用户是否登录，未登录则跳转至登录页面
// 一定要在BaseController.Prepare()后执行
func (c *BaseController) checkLogin() {
	if c.curUser.Id == 0 {
		logs.Debug("未登录")
		c.Result(enums.NoLogin, "未登录", "")
		////登录页面地址
		//urlstr := c.URLFor("HomeController.Login") + "?url="
		////登录成功后返回的址为当前
		//returnURL := c.Ctx.Request.URL.Path
		////如果ajax请求则返回相应的错码和跳转的地址
		//if c.Ctx.Input.IsAjax() {
		//	//由于是ajax请求，因此地址是header里的Referer
		//	returnURL := c.Ctx.Input.Refer()
		//	c.jsonResult(enums.JRCode302, "请登录", urlstr+returnURL)
		//}
		//c.Redirect(urlstr+returnURL, 302)
		//c.StopRun()
	}
}

// 判断某 Controller.Action 当前用户是否有权访问
func (c *BaseController) checkActionAuthor(ctrlName, ActName string) bool {
	if c.curUser.Id == 0 {
		return false
	}
	//从session获取用户信息
	user := c.GetSession("user")
	//类型断言
	v, ok := user.(models.User)
	if ok {
		//如果是超级管理员，则直接通过
		//if v.IsSuper == true {
		//	return true
		//}
		//遍历用户所负责的资源列表
		for i, _ := range v.ResourceUrlForList {
			urlfor := strings.TrimSpace(v.ResourceUrlForList[i])
			if len(urlfor) == 0 {
				continue
			}
			// TestController.Get,:last,xie,:first,asta
			strs := strings.Split(urlfor, ",")
			if len(strs) > 0 && strs[0] == (ctrlName+"."+ActName) {
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
	//先判断是否登录
	c.checkLogin()
	//如果Action在忽略列表里，则直接通用
	controllerName, actionName := c.GetControllerAndAction()
	for _, ignore := range ignores {
		if ignore == actionName {
			return
		}
	}

	hasAuthor := true
	//hasAuthor := c.checkActionAuthor(controllerName, actionName)
	if !hasAuthor {
		logs.Debug(fmt.Sprintf("author control: path=%s.%s userid=%v  无权访问", controllerName, actionName, c.curUser.Id))
		//如果没有权限
		c.jsonResult(enums.JRCode401, "无权访问", "")
		//if !hasAuthor {
		//	if c.Ctx.Input.IsAjax() {
		//		c.jsonResult(enums.JRCode401, "无权访问", "")
		//	} else {
		//		c.pageError("无权访问")
		//	}
		//}
	}
}

//从session里取用户信息
func (c *BaseController) adapterUserInfo() {
	a := c.GetSession("user")
	logs.Debug("adapterUserInfo:%v", a)
	if a != nil {
		c.curUser = a.(models.User)
		c.Data["user"] = a
	}
}

//SetUser2Session 获取用户信息（包括资源UrlFor）保存至Session
func (c *BaseController) setUser2Session(userId int) error {
	m, err := models.UserOne(userId)
	if err != nil {
		return err
	}
	//获取这个用户能获取到的所有资源列表
	resourceList := models.ResourceTreeGridByUserId(userId, 1000)
	for _, item := range resourceList {
		m.ResourceUrlForList = append(m.ResourceUrlForList, strings.TrimSpace(item.UrlFor))
	}
	c.SetSession("user", *m)
	logs.Info("CruSession0000000:%v",   c.CruSession)
	return nil
}

// 设置模板
// 第一个参数模板，第二个参数为layout
func (c *BaseController) setTpl(template ...string) {
	var tplName string
	layout := "shared/layout_page.html"
	switch {
	case len(template) == 1:
		tplName = template[0]
	case len(template) == 2:
		tplName = template[0]
		layout = template[1]
	default:
		//不要Controller这个10个字母
		controllerName, actionName := c.GetControllerAndAction()
		ctrlName := strings.ToLower(controllerName[0 : len(controllerName)-10])
		actionName = strings.ToLower(actionName)
		tplName = ctrlName + "/" + actionName + ".html"
	}
	//fmt.Println("666666666666666:",layout, tplName)
	c.Layout = layout
	c.TplName = tplName
}
func (c *BaseController) jsonResult(code enums.JsonResultCode, msg string, obj interface{}) {
	data   := [3]string {"Jerry", "Tom", "Jerry & Tom"}
	r := &models.Result{Code:20000, Data:data}
	c.Data["json"] = r
	c.AllowCross()
	c.ServeJSON()
	c.StopRun()
}
func (c *BaseController) Result(code enums.JsonResultCode, msg string, data interface{}) {
	//c.AllowCross()
	r := &models.Result{Code:code, Data:data, Msg:msg}
	c.Data["json"] = r
	//c.AllowCross()
	c.ServeJSON()
	c.StopRun()
}
//// 重定向
//func (c *BaseController) redirect(url string) {
//	c.Redirect(url, 302)
//	c.StopRun()
//}

// 重定向 去错误页
func (c *BaseController) pageError(msg string) {
	errorurl := c.URLFor("HomeController.Error") + "/" + msg
	c.Redirect(errorurl, 302)
	c.StopRun()
}
//
//// 重定向 去登录页
//func (c *BaseController) pageLogin() {
//	url := c.URLFor("HomeController.Login")
//	c.Redirect(url, 302)
//	c.StopRun()
//}

// 是否POST提交
func (this *BaseController) isPost() bool {
	return this.Ctx.Request.Method == "POST"
}
