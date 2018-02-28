package controllers

import (
	"encoding/json"
	"fmt"
	//"strconv"
	"strings"

	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

type UserController struct {
	BaseController
}

func (c *UserController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}
//func (c *UserController) Index() {
//	//是否显示更多查询条件的按钮
//	c.Data["showMoreQuery"] = true
//	//将页面左边菜单的某项激活
//	//c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
//	//页面模板设置
//	c.setTpl()
//	c.LayoutSections = make(map[string]string)
//	c.LayoutSections["headcssjs"] = "user/index_headcssjs.html"
//	c.LayoutSections["footerjs"] = "user/index_footerjs.html"
//	//页面里按钮权限控制
//	c.Data["canEdit"] = c.checkActionAuthor("UserController", "Edit")
//	c.Data["canDelete"] = c.checkActionAuthor("UserController", "Delete")
//}
func (c *UserController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值（要求配置文件里 copyrequestbody=true）
	var params models.UserQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Info("查询用户列表:%v", params)
	//fmt.Println(params)
	//获取数据列表和总数
	data, total := models.UserPageList(&params)
	//定义返回的数据结构

	//for _, r := range data {
	//	logs.Debug("rows:%v, %v", r.Id, r.RoleUserRel)
	//}
	//a := make([] int, 0)
	//b:= make([] *models.RoleUserRel, 0)
	//for _, r := range data {
	//	r.RoleIds = a
	//	r.RoleUserRel = b
	//}
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	//c.Data["json"] = result
	c.Result(enums.Success, "获取用户列表成功", result)
}

// Edit 添加 编辑 页面
func (c *UserController) Edit() {
	m := models.User{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	fmt.Printf("编辑用户:%v", m)
	//c.Result(enums.Success, "编辑用户列表成功", 0)

	o := orm.NewOrm()
	//var err error
	//获取form里的值
	//if err = c.ParseForm(&m); err != nil {
	//	c.Result(enums.JRCodeFailed, "解析数据失败", m.Id)
	//}
	//删除已关联的历史数据
	if _, err := o.QueryTable(models.RoleUserRelTBName()).Filter("user__id", m.Id).Delete(); err != nil {
		c.Result(enums.JRCodeFailed, "删除历史关系失败", "")
	}
	if m.Id == 0 {
		//对密码进行加密
		m.Password = utils.String2md5(m.ModifyPassword)
		if _, err := o.Insert(&m); err != nil {
			c.Result(enums.JRCodeFailed, "添加失败", m.Id)
		}
	} else {
		if oM, err := models.UserOne(m.Id); err != nil {
			c.Result(enums.JRCodeFailed, "未找到该用户，请刷新后重试", m.Id)
		} else {
			m.Password = strings.TrimSpace(m.ModifyPassword)
			logs.Info("修改密码:%v", m.Password)
			if len(m.Password) == 0 {
				//如果密码为空则不修改
				m.Password = oM.Password
			} else {
				m.Password = utils.String2md5(m.Password)
			}
			//本页面不修改头像和密码，直接将值附给新m
			//m.Avatar = oM.Avatar
		}
		if _, err := o.Update(&m); err != nil {
			c.Result(enums.JRCodeFailed, "保存失败", m.Id)
		}
	}
	//添加关系
	var relations []models.RoleUserRel
	for _, roleId := range m.RoleIds {
		r := models.Role{Id: roleId}
		relation := models.RoleUserRel{User: &m, Role: &r}
		relations = append(relations, relation)
	}
	logs.Debug("RoleIds:%v", m.RoleIds)
	logs.Debug("relations:%v", relations)
	if len(relations) > 0 {
		//批量添加
		if _, err := o.InsertMulti(len(relations), relations); err == nil {
			c.Result(enums.JRCodeSucc, "保存成功", m.Id)
		} else {
			c.Result(enums.JRCodeFailed, "保存失败", m.Id)
		}
	} else {
		c.Result(enums.JRCodeSucc, "保存成功", m.Id)
	}


	////如果是Post请求，则由Save处理
	//if c.Ctx.Request.Method == "POST" {
	//	c.Save()
	//}
	//Id, _ := c.GetInt(":id", 0)
	//m := &models.User{}
	//var err error
	//if Id > 0 {
	//	m, err = models.UserOne(Id)
	//	if err != nil {
	//		c.pageError("数据无效，请刷新后重试")
	//	}
	//	o := orm.NewOrm()
	//	o.LoadRelated(m, "RoleUserRel")
	//} else {
	//	//添加用户时默认状态为启用
	//	m.Status = enums.Enabled
	//}
	//c.Data["m"] = m
	////获取关联的roleId列表
	//var roleIds []string
	//for _, item := range m.RoleUserRel {
	//	roleIds = append(roleIds, strconv.Itoa(item.Role.Id))
	//}
	//c.Data["roles"] = strings.Join(roleIds, ",")
	//c.setTpl("user/edit.html", "shared/layout_pullbox.html")
	//c.LayoutSections = make(map[string]string)
	//c.LayoutSections["footerjs"] = "user/edit_footerjs.html"
}
func (c *UserController) Save() {
	m := models.User{}
	o := orm.NewOrm()
	var err error
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	//删除已关联的历史数据
	if _, err := o.QueryTable(models.RoleUserRelTBName()).Filter("user__id", m.Id).Delete(); err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除历史关系失败", "")
	}
	if m.Id == 0 {
		//对密码进行加密
		m.Password = utils.String2md5(m.Password)
		if _, err := o.Insert(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}
	} else {
		if oM, err := models.UserOne(m.Id); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据无效，请刷新后重试", m.Id)
		} else {
			m.Password = strings.TrimSpace(m.Password)
			if len(m.Password) == 0 {
				//如果密码为空则不修改
				m.Password = oM.Password
			} else {
				m.Password = utils.String2md5(m.Password)
			}
			//本页面不修改头像和密码，直接将值附给新m
			//m.Avatar = oM.Avatar
		}
		if _, err := o.Update(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}
	//添加关系
	var relations []models.RoleUserRel
	for _, roleId := range m.RoleIds {
		r := models.Role{Id: roleId}
		relation := models.RoleUserRel{User: &m, Role: &r}
		relations = append(relations, relation)
	}
	if len(relations) > 0 {
		//批量添加
		if _, err := o.InsertMulti(len(relations), relations); err == nil {
			c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "保存失败", m.Id)
		}
	} else {
		c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
	}
}
func (c *UserController) Delete() {
	var m  []int
	json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	logs.Info("删除用户:%v",  m)
	query := orm.NewOrm().QueryTable(models.UserTBName())
	if num, err := query.Filter("id__in", m).Delete(); err == nil {
		c.Result(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), m)
	} else {
		c.Result(enums.JRCodeFailed, "删除失败", m)
	}
}
