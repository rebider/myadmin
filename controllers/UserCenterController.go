package controllers

import (
	//"strings"

	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	//"github.com/astaxie/beego/logs"
	"strconv"
	//"github.com/chnzrb/myadmin/utils"
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
	Id := c.curUser.Id
	m, err := models.UserOne(Id)
	c.CheckError(err, "用户不存在")
	//if m == nil || err != nil {
	//	logs.Error(err)
	//	c.pageError("数据无效，请刷新后重试")
	//}


	roleIds := make([] string, len(m.RoleIds))
	//logs.Debug("roleIds:%v", m.RoleIds)
	for i, e := range m.RoleIds{
		//t, _:= strconv.Atoi(e)
		roleIds[i] = strconv.Itoa(e)
	}
	c.Result(enums.Success, "获取用户信息成功",
		struct {
			Roles []string `json:"roles"`
			Name string `json:"name"`
			Resources []*models.Resource
		}{Roles: roleIds, Name:m.Name, Resources: models.ResourceTreeGrid()})
}
