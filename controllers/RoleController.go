package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

//RoleController 角色管理
type RoleController struct {
	BaseController
}

//func (c *RoleController) Prepare() {
//	//先执行
//	c.BaseController.Prepare()
//	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
//	c.checkAuthor()
//	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
//	//权限控制里会进行登录验证，因此这里不用再作登录验证
//	//c.checkLogin()
//}


// DataGrid 角色管理首页 表格获取数据
func (c *RoleController) List() {
	//直接反序化获取json格式的requestbody里的值
	var params models.RoleQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.RolePageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取角色列表成功", result)
}

////DataList 角色列表
//func (c *RoleController) DataList() {
//	var params = models.RoleQueryParam{}
//	//获取数据列表和总数
//	data := models.RoleDataList(&params)
//	//定义返回的数据结构
//	c.jsonResult(enums.JRCodeSucc, "", data)
//}

//Edit 添加、编辑角色界面
func (c *RoleController) Edit() {
	m := models.Role{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑角色")
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err = o.Insert(&m); err == nil {
			c.Result(enums.CodeSuccess, "添加角色成功", m.Id)
		} else {
			c.Result(enums.CodeFail, "添加角色失败", m.Id)
		}

	} else {
		if _, err = o.Update(&m); err == nil {
			c.Result(enums.CodeSuccess, "编辑角色成功", m.Id)
		} else {
			c.Result(enums.CodeFail, "编辑角色失败", m.Id)
		}
	}
}

//Delete 批量删除
func (c *RoleController) Delete() {
	var ids []int
	json.Unmarshal(c.Ctx.Input.RequestBody, &ids)
	logs.Info("删除角色:%+v", ids)
	if num, err := models.RoleBatchDelete(ids); err == nil {
		c.Result(enums.CodeSuccess, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.Result(enums.CodeFail, "删除失败", 0)
	}
}

//Allocate 给角色分配资源界面
func (c *RoleController) Allocate() {

	var params struct {
		Id          int    `json:"id"`
		ResourceIds [] int `json:"resourceIds"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Info("给角色分配资源界面:%+v", params)
	utils.CheckError(err)

	roleId := params.Id
	resourceIds := params.ResourceIds

	o := orm.NewOrm()
	m := models.Role{Id: roleId}
	if err := o.Read(&m); err != nil {
		c.Result(enums.CodeFail, "数据无效，请刷新后重试", "")
	}
	//删除已关联的历史数据
	if _, err := o.QueryTable(models.RoleResourceRelTBName()).Filter("role__id", m.Id).Delete(); err != nil {
		c.Result(enums.CodeFail, "删除历史关系失败", "")
	}
	var relations []models.RoleResourceRel
	for _, id := range resourceIds {
		r := models.Resource{Id: id}
		relation := models.RoleResourceRel{Role: &m, Resource: &r}
		relations = append(relations, relation)
		//}
	}
	if len(relations) > 0 {
		//批量添加
		if _, err := o.InsertMulti(len(relations), relations); err == nil {
			c.Result(enums.CodeSuccess, "保存成功", "")
		}
	}
	c.Result(0, "保存失败", "")
}
