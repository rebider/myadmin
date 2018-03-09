package controllers

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
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

	c.Result(enums.CodeSuccess, "获取用户信息成功",
		struct {
			Name         string             `json:"name"`
			ResourceTree []*models.Resource `json:"menuTree"`
		}{Name: m.Name, ResourceTree: models.TranResourceList2ResourceTree(models.GetResourceListByUserId(Id, 0))})
}
