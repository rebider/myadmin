package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	"github.com/golang/protobuf/proto"
	"github.com/chnzrb/myadmin/proto"
	"time"
)

type MailController struct {
	BaseController
}

func (c *MailController) MailLogList() {
	var params models.MailLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询邮件日志:%+v", params)
	data, total := models.GetMailLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取邮件日志", result)
	//v := [] string{
	//	"d",
	//	"1",
	//	"2",
	//	"d",
	//}
	//utils.RemoveDuplicateArray(v)
}

func (c *MailController) DelMailLog() {
	var idList []int
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &idList)
	utils.CheckError(err)
	logs.Info("删除邮件列表:%+v", idList)
	err = models.DeleteMailLog(idList)
	c.CheckError(err, "删除邮件失败")
	c.Result(enums.CodeSuccess, "成功删除邮件", idList)
}
func (c *MailController) SendMail() {
	//type platformServer struct {
	//	PlatformId   int
	//	ServerIdList [] string
	//}
	var params struct {
		PlatformId   int
		ServerIdList [] string
		//platformServerList [] *platformServer
		PlayerNameList string
		MailItemList   [] *gm.MSendMailTosProp
		Title          string
		Content        string
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("发送邮件:%+v", params)
	logs.Info("发送邮件:%+v", params.MailItemList)
	//platformId, err := c.GetInt("platformId")
	//c.CheckError(err)
	//serverId:= c.GetString("serverId")
	//playerId, err := c.GetInt("playerId")
	//c.CheckError
	serverIdList, err := json.Marshal(params.ServerIdList)
	c.CheckError(err)
	itemList, err := json.Marshal(params.MailItemList)
	c.CheckError(err)
	mailLog := &models.MailLog{
		PlatformId:     params.PlatformId,
		ServerIdList:   string(serverIdList),
		Title:          params.Title,
		Content:        params.Content,
		Time:           time.Now().Unix(),
		UserId:         c.curUser.Id,
		ItemList:       string(itemList),
		PlayerNameList: params.PlayerNameList,
		Status:         0,
	}
	err = models.Db.Save(&mailLog).Error
	c.CheckError(err, "写邮件日志失败")
	for _, serverId := range params.ServerIdList {
		conn, err := models.GetWsByPlatformIdAndSid(params.PlatformId, serverId)
		c.CheckError(err)
		defer conn.Close()
		request := gm.MSendMailTos{
			Token:          proto.String(""),
			Title:          proto.String(params.Title),
			Content:        proto.String(params.Content),
			PlayerNameList: proto.String(params.PlayerNameList),
			PropList:       params.MailItemList,
		}
		mRequest, err := proto.Marshal(&request)
		c.CheckError(err)

		_, err = conn.Write(utils.Packet(9903, mRequest))
		c.CheckError(err)
		var receive = make([]byte, 100, 100)
		n, err := conn.Read(receive)
		c.CheckError(err)
		respone := &gm.MSendMailToc{}
		data := receive[5:n]
		err = proto.Unmarshal(data, respone)
		c.CheckError(err)

		if *respone.Result == gm.MSendMailToc_success {
			logs.Info("发送邮件成功:%+v", request)
		} else {
			c.Result(enums.CodeFail, "发送邮件失败", 0)
		}
	}
	c.Result(enums.CodeSuccess, "发送邮件成功", 0)
	//conn.Read()
}

