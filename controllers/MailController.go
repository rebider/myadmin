package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	"github.com/chnzrb/myadmin/proto"
	"time"
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/base64"
)

type MailController struct {
	BaseController
}

// 获取邮件列表
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
}

// 删除邮件
func (c *MailController) DelMailLog() {
	var idList []int
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &idList)
	utils.CheckError(err)
	logs.Info("删除邮件列表:%+v", idList)
	err = models.DeleteMailLog(idList)
	c.CheckError(err, "删除邮件失败")
	c.Result(enums.CodeSuccess, "成功删除邮件", idList)
}



//发送邮件
func (c *MailController) SendMail() {
	var params struct {
		PlatformId     int `json:"platformId"`
		NodeList       [] string `json:"serverIdList"`
		PlayerNameList string `json:"playerNameList"`
		MailItemList   [] *gm.MSendMailTosProp `json:"mailItemList"`
		Title          string `json:"title"`
		Content        string `json:"content"`
		PlayerIdList [] int  `json:"playerIdList"`
	}
	var result struct {
		ErrorCode int
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("发送邮件:%+v", params)
	if params.PlayerNameList != "" && len(params.NodeList) != 1 {
		logs.Error("参数错误!!!!")
		c.Result(enums.CodeFail, "参数错误!!!", 0)
	}

	if params.PlayerNameList != "" {
		params.PlayerIdList, err = models.TranPlayerNameSting2PlayerIdList(params.PlatformId, params.PlayerNameList)
		c.CheckError(err, "解析玩家失败")
	} else {
		params.PlayerIdList = make([] int, 0)
	}


	request, err := json.Marshal(params)
	c.CheckError(err)

	nodeList, err := json.Marshal(params.NodeList)
	c.CheckError(err)
	itemList, err := json.Marshal(params.MailItemList)
	c.CheckError(err)

	mailLog := &models.MailLog{
		PlatformId:     params.PlatformId,
		NodeList:       string(nodeList),
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


	for _, node :=range params.NodeList {
		url := models.GetGameURLByNode(node) + "/send_mail"
		data := string(request)
		sign := utils.String2md5(data + enums.GmSalt)
		base64Data := base64.URLEncoding.EncodeToString([]byte(data))
		requestBody := "data=" + base64Data+ "&sign=" + sign
		resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(requestBody))
		utils.CheckError(err)
		if err != nil {
			continue
		}

		defer resp.Body.Close()
		responseBody, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			logs.Info("发送邮件失败 区服:%s, 标题:%s, 内容:%s, 玩家:%v", node, params.Title, params.Content, params.PlayerNameList)
			continue
		}

		//logs.Debug("result:%v", string(responseBody))

		err = json.Unmarshal(responseBody, &result)

		//c.CheckError(err)
		if result.ErrorCode != 0 || err != nil{
			logs.Info("发送邮件失败 区服:%s, 标题:%s, 内容:%s, 玩家:%v", node, params.Title, params.Content, params.PlayerNameList)
		} else {
			logs.Info("发送邮件成功 区服:%s, 标题:%s, 内容:%s, 玩家:%v", node, params.Title, params.Content, params.PlayerNameList)
		}

	}
	//args := fmt.Sprintf("platform_id=%d&server_id=%s&player_id=%d&type=%d&sec=%d", params.PlatformId, params.ServerId, params.PlayerId, params.Type, params.Sec)
	//sign := utils.String2md5(args + enums.GmSalt)

	//if result.ErrorCode != 0 {
	//	c.Result(enums.CodeFail, "发送邮件失败", 0)
	//}

	//for _, node := range params.NodeList {
	//	conn, err := models.GetWsByNode(node)
	//	c.CheckError(err)
	//	defer conn.Close()
	//	request := gm.MSendMailTos{
	//		Token:          proto.String(""),
	//		Title:          proto.String(params.Title),
	//		Content:        proto.String(params.Content),
	//		PlayerNameList: proto.String(params.PlayerNameList),
	//		PropList:       params.MailItemList,
	//	}
	//	mRequest, err := proto.Marshal(&request)
	//	c.CheckError(err)
	//
	//	_, err = conn.Write(utils.Packet(9903, mRequest))
	//	c.CheckError(err)
	//	var receive = make([]byte, 100, 100)
	//	n, err := conn.Read(receive)
	//	c.CheckError(err)
	//	respone := &gm.MSendMailToc{}
	//	data := receive[5:n]
	//	err = proto.Unmarshal(data, respone)
	//	c.CheckError(err)
	//
	//	if *respone.Result == gm.MSendMailToc_success {
	//		logs.Info("发送邮件成功:%+v, %+v", node, request)
	//	} else {
	//		c.Result(enums.CodeFail, fmt.Sprintf("发送邮件失败:%+v, %+v", node, request), 0)
	//	}
	//}
	c.Result(enums.CodeSuccess, "发送邮件成功", 0)
}
