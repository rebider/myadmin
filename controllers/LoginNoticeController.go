// 登录公告管理
package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	//"github.com/chnzrb/myadmin/utils"
)

type LoginNoticeController struct {
	BaseController
}

// 设置登录公告
func (c *LoginNoticeController) SetNotice() {
	var params struct {
		PlatformId string `json:"platformId"`
		Notice    string	`json:"notice"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	c.CheckError(err)
	logs.Info("设置登录公告:%+v", params)
	//c.CheckError(err)
	//
	//data, err := json.Marshal(params)
	//utils.CheckError(err)
	//
	//url := utils.GetCenterURL() + "/set_login_notice"
	//err = utils.HttpRequest(url, string(data))
	////out, err := utils.NodeTool(
	////	"mod_login_notice",
	////	"update_login_notice",
	////	params.PlatformId,
	////	params.Notice,
	////)
	//c.CheckError(err, "设置中心服登录公告:" + url)
	//noticeLog := &models.LoginNotice{
	//	PlatformId:  params.PlatformId,
	//	Notice:     params.Notice,
	//	Time:        time.Now().Unix(),
	//	UserId:      c.curUser.Id,
	//}
	//err = models.Db.Save(&noticeLog).Error
	err = models.UpdateAndPushLoginNotice(c.curUser.Id, params.PlatformId, params.Notice)
	c.CheckError(err, "写登录公告日志失败")
	c.Result(enums.CodeSuccess, "设置登录公告成功", 0)
}

// 批量设置登录公告
func (c *LoginNoticeController) BatchSetNotice() {
	var params struct {
		PlatformIdList [] string
		Notice    string
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	c.CheckError(err)
	logs.Info("批量设置登录公告:%+v", params)
	for _, PlatformId := range params.PlatformIdList {
		//var request struct {
		//	PlatformId string `json:"platformId"`
		//	Notice    string	`json:"notice"`
		//}
		//request.PlatformId = PlatformId
		//request.Notice = params.Notice
		//data, err := json.Marshal(request)
		//utils.CheckError(err)
		//url := utils.GetCenterURL() + "/set_login_notice"
		//err = utils.HttpRequest(url, string(data))
		////out, err := utils.NodeTool(
		////	"mod_login_notice",
		////	"update_login_notice",
		////	PlatformId,
		////	params.Notice,
		////)
		//c.CheckError(err, "批量中心服登录公告")
		//noticeLog := &models.LoginNotice{
		//	PlatformId:  PlatformId,
		//	Notice:     params.Notice,
		//	Time:        time.Now().Unix(),
		//	UserId:      c.curUser.Id,
		//}
		//err = models.Db.Save(&noticeLog).Error

		err = models.UpdateAndPushLoginNotice(c.curUser.Id, PlatformId, params.Notice)
		c.CheckError(err, "写登录公告日志失败")
	}

	c.Result(enums.CodeSuccess, "批量设置登录公告成功", 0)
}

// 获取登录公告列表
func (c *LoginNoticeController) LoginNoticeList() {
	//logs.Info("查询公告日志:%+v", params)
	var params struct {
		PlatformIdList [] string
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	c.CheckError(err)
	logs.Info("获取登录公告列表:%+v", params)
	data := models.GetLoginNoticeListByPlatformIdList(params.PlatformIdList)
	result := make(map[string]interface{})
	//result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取登录公告列表", result)
}

//删除登录公告
func (c *LoginNoticeController) DelLoginNotice() {
	var idList []string
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &idList)
	c.CheckError(err)
	logs.Info("删除登录公告列表:%+v", idList)
	for _, PlatformId := range idList {
		err = models.UpdateAndPushLoginNotice(c.curUser.Id, PlatformId, "")
		c.CheckError(err, "写登录公告失败")
	}
	err = models.DeleteLoginNotice(idList)
	c.CheckError(err, "删除登录公告失败")
	c.Result(enums.CodeSuccess, "成功删除登录公告", idList)
}
