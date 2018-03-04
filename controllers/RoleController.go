package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	//"strconv"
	//"strings"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

//RoleController 角色管理
type RoleController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *RoleController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DataList", "UpdateSeq")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

//Index 角色管理首页
func (c *RoleController) Index() {
	//是否显示更多查询条件的按钮
	c.Data["showMoreQuery"] = true
	//将页面左边菜单的某项激活
	//c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "role/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "role/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("RoleController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("RoleController", "Delete")
	c.Data["canAllocate"] = c.checkActionAuthor("RoleController", "Allocate")
}

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
	c.Result(enums.Success, "获取角色列表成功", result)
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
			c.Result(enums.JRCodeSucc, "添加角色成功", m.Id)
		} else {
			c.Result(enums.JRCodeFailed, "添加角色失败", m.Id)
		}

	} else {
		if _, err = o.Update(&m); err == nil {
			c.Result(enums.JRCodeSucc, "编辑角色成功", m.Id)
		} else {
			c.Result(enums.JRCodeFailed, "编辑角色失败", m.Id)
		}
	}
}

//
////Save 添加、编辑页面 保存
//func (c *RoleController) Save() {
//	var err error
//	m := models.Role{}
//	//获取form里的值
//	if err = c.ParseForm(&m); err != nil {
//		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
//	}
//	o := orm.NewOrm()
//	if m.Id == 0 {
//		if _, err = o.Insert(&m); err == nil {
//			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
//		} else {
//			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
//		}
//
//	} else {
//		if _, err = o.Update(&m); err == nil {
//			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
//		} else {
//			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
//		}
//	}
//
//}

//Delete 批量删除
func (c *RoleController) Delete() {
	//strs := c.GetString("ids")
	//ids := make([]int, 0, len(strs))
	//for _, str := range strings.Split(strs, ",") {
	//	if id, err := strconv.Atoi(str); err == nil {
	//		ids = append(ids, id)
	//	}
	//}
	var ids []int
	json.Unmarshal(c.Ctx.Input.RequestBody, &ids)
	logs.Info("删除角色:%+v", ids)
	if num, err := models.RoleBatchDelete(ids); err == nil {
		c.Result(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.Result(enums.JRCodeFailed, "删除失败", 0)
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
	//strs := c.GetString("ids")

	o := orm.NewOrm()
	m := models.Role{Id: roleId}
	if err := o.Read(&m); err != nil {
		c.Result(enums.JRCodeFailed, "数据无效，请刷新后重试", "")
	}
	//删除已关联的历史数据
	if _, err := o.QueryTable(models.RoleResourceRelTBName()).Filter("role__id", m.Id).Delete(); err != nil {
		c.Result(enums.JRCodeFailed, "删除历史关系失败", "")
	}
	var relations []models.RoleResourceRel
	for _, id := range resourceIds {
		//if id, err := strconv.Atoi(str); err == nil {
		r := models.Resource{Id: id}
		relation := models.RoleResourceRel{Role: &m, Resource: &r}
		relations = append(relations, relation)
		//}
	}
	//for _, str := range strings.Split(strs, ",") {
	//	if id, err := strconv.Atoi(str); err == nil {
	//		r := models.Resource{Id: id}
	//		relation := models.RoleResourceRel{Role: &m, Resource: &r}
	//		relations = append(relations, relation)
	//	}
	//}
	if len(relations) > 0 {
		//批量添加
		if _, err := o.InsertMulti(len(relations), relations); err == nil {
			c.Result(enums.JRCodeSucc, "保存成功", "")
		}
	}
	c.Result(0, "保存失败", "")
}

//func (c *RoleController) UpdateSeq() {
//	Id, _ := c.GetInt("pk", 0)
//	oM, err := models.RoleOne(Id)
//	if err != nil || oM == nil {
//		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
//	}
//	value, _ := c.GetInt("value", 0)
//	oM.Seq = value
//	o := orm.NewOrm()
//	if _, err := o.Update(oM); err == nil {
//		c.jsonResult(enums.JRCodeSucc, "修改成功", oM.Id)
//	} else {
//		c.jsonResult(enums.JRCodeFailed, "修改失败", oM.Id)
//	}
//}
