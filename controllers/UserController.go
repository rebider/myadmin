package controllers

import (
	"encoding/json"
	"fmt"
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

func (c *UserController) Info() {
	m := c.curUser
	c.Result(enums.CodeSuccess, "获取用户信息成功",
		struct {
			Name         string             `json:"name"`
			ResourceTree []*models.Menu `json:"menuTree"`
		}{Name: m.Name, ResourceTree: models.TranMenuList2MenuTree(models.GetMenuListByUserId(m.Id))})
}


func (c *UserController) List() {
	var params models.UserQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Info("查询用户列表:%+v", params)
	//获取数据列表和总数
	data, total := models.UserPageList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取用户列表成功", result)
}

// Edit 添加 编辑 页面
func (c *UserController) Edit() {
	m := models.User{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑用户")
	logs.Info("编辑用户:%+v", m)

	o := orm.NewOrm()

	//删除旧的用户角色关系
	_, err = o.QueryTable(models.RoleUserRelTBName()).Filter("user__id", m.Id).Delete()
	utils.CheckError(err, "删除旧的用户角色关系")

	if m.Id == 0 {
		//对密码进行加密
		m.Password = utils.String2md5(m.ModifyPassword)
		if _, err := o.Insert(&m); err != nil {
			c.Result(enums.CodeFail, "添加用户失败", m.Id)
		}
	} else {
		if oM, err := models.UserOne(m.Id); err != nil {
			c.Result(enums.CodeFail, "未找到该用户", m.Id)
		} else {
			m.Password = strings.TrimSpace(m.ModifyPassword)
			if len(m.Password) == 0 {
				//如果密码为空则不修改
				m.Password = oM.Password
			} else {
				logs.Info("修改密码:%+v", m.Password)
				m.Password = utils.String2md5(m.Password)
			}
		}
		if _, err := o.Update(&m); err != nil {
			c.Result(enums.CodeFail, "保存用户失败", m.Id)
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
			c.Result(enums.CodeSuccess, "保存用户角色关系成功", m.Id)
		} else {
			c.Result(enums.CodeFail, "保存用户角色关系失败", m.Id)
		}
	} else {
		c.Result(enums.CodeSuccess, "保存成功", m.Id)
	}
}

func (c *UserController) Delete() {
	logs.Debug(c.GetControllerAndAction())
	var m  []int
	json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	logs.Info("删除用户:%+v",  m)
	query := orm.NewOrm().QueryTable(models.UserTBName())
	if num, err := query.Filter("id__in", m).Delete(); err == nil {
		c.Result(enums.CodeSuccess, fmt.Sprintf("成功删除 %d 项", num), m)
	} else {
		c.Result(enums.CodeFail, "删除失败", m)
	}
}


func (c *UserController) ChangePassword() {
	var params struct {
		OldPwd string `json:"oldPwd"`
		NewPwd string `json:"newPwd"`
		NewPwd2 string `json:"newPwd2"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("修改密码:%+v", params)
	Id := c.curUser.Id
	oM, err := models.UserOne(Id)
	utils.CheckError(err)
	md5str := utils.String2md5(params.OldPwd)
	if oM.Password != md5str {
		c.Result(enums.CodeFail, "原密码错误", "")
	}
	if len(params.NewPwd2) == 0 {
		c.Result(enums.CodeFail, "请输入新密码", "")
	}
	if params.NewPwd != params.NewPwd2 {
		c.Result(enums.CodeFail, "两次输入的新密码不一致", "")
	}
	oM.Password = utils.String2md5(params.NewPwd)
	o := orm.NewOrm()
	if _, err := o.Update(oM); err != nil {
		c.Result(enums.CodeFail, "保存失败", oM.Id)
	} else {
		c.setUser2Session(Id)
		c.Result(enums.CodeSuccess, "保存成功", oM.Id)
	}
}
