package controllers

import (
	//"strings"

	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"strconv"
)

type UserCenterController struct {
	BaseController
}

func (c *UserCenterController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	c.checkLogin()
}
func (c *UserCenterController) Info() {
	logs.Debug("66666666666666")
	Id := c.curUser.Id
	m, err := models.UserOne(Id)
	if m == nil || err != nil {
		logs.Error(err)
		c.pageError("数据无效，请刷新后重试")
	}


	roleIds := make([] string, len(m.RoleIds))
	logs.Debug("roleIds:%v", m.RoleIds)
	for i, e := range m.RoleIds{
		//t, _:= strconv.Atoi(e)
		roleIds[i] = strconv.Itoa(e)
	}
	c.Result(enums.Success, "获取用户信息成功",
		struct {
			Roles []string `json:"roles"`
			Name string `json:"name"`
		}{Roles: roleIds, Name:m.Name})
}

//func (c *UserCenterController) Profile() {
//	Id := c.curUser.Id
//	m, err := models.UserOne(Id)
//	if m == nil || err != nil {
//		c.pageError("数据无效，请刷新后重试")
//	}
//	c.Data["hasAvatar"] = len(m.Avatar) > 0
//	//utils.Debug(m.Avatar)
//	c.Data["m"] = m
//	c.setTpl()
//	c.LayoutSections = make(map[string]string)
//	c.LayoutSections["headcssjs"] = "usercenter/profile_headcssjs.html"
//	c.LayoutSections["footerjs"] = "usercenter/profile_footerjs.html"
//}
//func (c *UserCenterController) BasicInfoSave() {
//	Id := c.curUser.Id
//	oM, err := models.UserOne(Id)
//	if oM == nil || err != nil {
//		c.jsonResult(enums.JRCodeFailed, "数据无效，请刷新后重试", "")
//	}
//	m := models.User{}
//	//获取form里的值
//	if err = c.ParseForm(&m); err != nil {
//		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
//	}
//	oM.RealName = m.RealName
//	oM.Mobile = m.Mobile
//	oM.Email = m.Email
//	oM.Avatar = c.GetString("ImageUrl")
//	o := orm.NewOrm()
//	if _, err := o.Update(oM); err != nil {
//		c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
//	} else {
//		c.setUser2Session(Id)
//		c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
//	}
//}
func (c *UserCenterController) ChangePasswd() {
	var params struct {
		OldPwd string `json:"oldPwd"`
		NewPwd string `json:"newPwd"`
		NewPwd2 string `json:"newPwd2"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Debug("RequestBody:%v", c.Ctx.Input.RequestBody)
	logs.Debug("params:%v", params)
	Id := c.curUser.Id
	oM, err := models.UserOne(Id)
	if oM == nil || err != nil {
		c.pageError("数据无效，请刷新后重试")
	}

	//oldPwd := strings.TrimSpace(c.GetString("UserPwd", ""))
	//newPwd := strings.TrimSpace(c.GetString("NewUserPwd", ""))
	//confirmPwd := strings.TrimSpace(c.GetString("ConfirmPwd", ""))
	md5str := utils.String2md5(params.OldPwd)
	if oM.Password != md5str {
		c.Result(enums.JRCodeFailed, "原密码错误", "")
	}
	if len(params.NewPwd2) == 0 {
		c.Result(enums.JRCodeFailed, "请输入新密码", "")
	}
	if params.NewPwd != params.NewPwd2 {
		c.Result(enums.JRCodeFailed, "两次输入的新密码不一致", "")
	}
	oM.Password = utils.String2md5(params.NewPwd)
	o := orm.NewOrm()
	if _, err := o.Update(oM); err != nil {
		c.Result(enums.JRCodeFailed, "保存失败", oM.Id)
	} else {
		c.setUser2Session(Id)
		c.Result(enums.Success, "保存成功", oM.Id)
	}
}
//func (c *UserCenterController) UploadImage() {
//	//这里type没有用，只是为了演示传值
//	stype, _ := c.GetInt32("type", 0)
//	if stype > 0 {
//		f, h, err := c.GetFile("fileImageUrl")
//		if err != nil {
//			c.jsonResult(enums.JRCodeFailed, "上传失败", "")
//		}
//		defer f.Close()
//		filePath := "static/upload/" + h.Filename
//		// 保存位置在 static/upload, 没有文件夹要先创建
//		c.SaveToFile("fileImageUrl", filePath)
//		c.jsonResult(enums.JRCodeSucc, "上传成功", "/"+filePath)
//	} else {
//		c.jsonResult(enums.JRCodeFailed, "上传失败", "")
//	}
//}
