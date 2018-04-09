package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	"time"
	"github.com/chnzrb/myadmin/crons"
)

type NoticeController struct {
	BaseController
}


func (c *NoticeController) SendNotice() {
	//type platformServer struct {
	//	PlatformId   int
	//	ServerIdList [] string
	//}
	var params struct {
		Id           int
		PlatformId   int
		ServerIdList [] string
		//platformServerList [] *platformServer
		//PlayerNameList     string
		//Title              string
		Content    string
		NoticeType int
		NoticeTime int
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("发送公告:%+v", params)
	//platformId, err := c.GetInt("platformId")
	//c.CheckError(err)
	//serverId:= c.GetString("serverId")
	//playerId, err := c.GetInt("playerId")
	//c.CheckError
	serverIdList, err := json.Marshal(params.ServerIdList)
	c.CheckError(err)
	noticeLog := &models.NoticeLog{
		Id:           params.Id,
		PlatformId:   params.PlatformId,
		ServerIdList: string(serverIdList),
		Content:      params.Content,
		Time:         time.Now().Unix(),
		UserId:       c.curUser.Id,
		NoticeType:   params.NoticeType,
		NoticeTime:   params.NoticeTime,
		Status:       0,
	}
	err = models.Db.Save(&noticeLog).Error
	//logs.Debug("noticeLog:%+v", noticeLog.Id)
	c.CheckError(err, "写公告日志失败")
	// 异步处理日志
	go crons.DealNoticeLog(noticeLog.Id)
	c.Result(enums.CodeSuccess, "发送公告成功", 0)
}


func (c *NoticeController) NoticeLogList() {
	var params models.NoticeLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询公告日志:%+v", params)
	data, total := models.GetNoticeLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取公告日志", result)
}

func (c *NoticeController) DelNoticeLog() {
	var idList []int
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &idList)
	utils.CheckError(err)
	logs.Info("删除公告列表:%+v", idList)
	err = models.DeleteNoticeLog(idList)
	c.CheckError(err, "删除公告失败")
	c.Result(enums.CodeSuccess, "成功删除公告", idList)
}

