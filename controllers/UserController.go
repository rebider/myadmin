package controllers

import (
	"encoding/json"
	"strings"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
)

type UserController struct {
	BaseController
}

func (c *UserController) Info() {
	m := c.curUser
	platformList, err := models.GetPlatformList("data/json/Platform.json")
	utils.CheckError(err)
	gameServerList, _ := models.GetAllGameServer()
	c.Result(enums.CodeSuccess, "获取用户信息成功",
		struct {
			Name           string         `json:"name"`
			ResourceTree   []*models.Menu `json:"menuTree"`
			PlatformList   []*models.Platform
			GameServerList []*models.GameServer
		}{
			Name:           m.Name,
			ResourceTree:   models.TranMenuList2MenuTree(models.GetMenuListByUserId(m.Id)),
			PlatformList:   platformList,
			GameServerList: gameServerList,
		})
}

// 获取用户列表
func (c *UserController) List() {
	var params models.UserQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询用户列表:%+v", params)
	data, total := models.GetUserList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取用户列表成功", result)
}

// 编辑 添加用户
func (c *UserController) Edit() {
	m := models.User{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	utils.CheckError(err, "编辑用户")
	logs.Info("编辑用户:%+v", m)

	//删除旧的用户角色关系
	_, err = models.DeleteRoleUserRelByUserIdList([] int{m.Id})
	c.CheckError(err, "删除旧的用户角色关系失败")
	for _, roleId := range m.RoleIds {
		relation := models.RoleUserRel{UserId: m.Id, RoleId: roleId}
		m.RoleUserRel = append(m.RoleUserRel, &relation)
	}
	if m.Id == 0 {
		m.Password = utils.String2md5(m.ModifyPassword)
		err = models.Db.Save(&m).Error
		c.CheckError(err, "添加用户失败")
	} else {
		oM, err := models.GetUserOne(m.Id)
		c.CheckError(err, "未找到该用户")
		m.Password = strings.TrimSpace(m.ModifyPassword)
		if len(m.Password) == 0 {
			//密码为空则不修改
			m.Password = oM.Password
		} else {
			m.Password = utils.String2md5(m.Password)
		}
		err = models.Db.Save(&m).Error
		c.CheckError(err, "保存用户失败")
	}
	c.Result(enums.CodeSuccess, "保存成功", m.Id)
}

// 删除用户
func (c *UserController) Delete() {
	var userIdList []int
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userIdList)
	utils.CheckError(err)
	logs.Info("删除用户:%+v", userIdList)
	err = models.DeleteUsers(userIdList)
	c.CheckError(err, "删除用户失败")
	c.Result(enums.CodeSuccess, "成功删除用户", userIdList)
}

//修改密码
func (c *UserController) ChangePassword() {
	var params struct {
		OldPwd string `json:"oldPwd"`
		NewPwd string `json:"newPwd"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("修改密码:%+v", params)
	Id := c.curUser.Id
	user, err := models.GetUserOne(Id)
	c.CheckError(err, "未找到该用户")
	md5str := utils.String2md5(params.OldPwd)
	if user.Password != md5str {
		c.Result(enums.CodeFail, "原密码错误", "")
	}
	if len(params.NewPwd) == 0 {
		c.Result(enums.CodeFail, "请输入新密码", "")
	}
	user.Password = utils.String2md5(params.NewPwd)
	err = models.Db.Model(&user).Updates(models.User{Password: user.Password}).Error
	c.CheckError(err, "保存失败")
	c.setUser2Session(Id)
	c.Result(enums.CodeSuccess, "保存成功", user.Id)
}
