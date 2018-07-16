package crons

import (
	"github.com/astaxie/beego/logs"
	"time"
	"encoding/json"
	//"github.com/chnzrb/myadmin/proto"
	"github.com/chnzrb/myadmin/utils"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/enums"
	//"github.com/golang/protobuf/proto"
	//"fmt"
	"encoding/base64"
	"net/http"
	"strings"
	"io/ioutil"
)

func ClockDealAllNoticeLog() {
	//logs.Info("定时处理公告")
	noticeLogList := models.GetAllNoticeLog()
	for _, noticeLog := range noticeLogList {
		if noticeLog.Status == 0 {
			// 处理未完成的公告
			go DealNoticeLog(noticeLog.Id)
		}
	}
}

func DealNoticeLog(id int) {
	noticeLog, err := models.GetNoticeLogOne(id)
	if err != nil {
		logs.Error("获取公告(%v)失败：%+v", id, err)
	}
	isGo := false
	now := int(time.Now().Unix())
	logs.Debug(id, now, noticeLog.Status, noticeLog.NoticeTime, noticeLog.LastSendTime)
	if noticeLog.Status == 0 {
		if noticeLog.NoticeType == enums.NoticeTypeMoment {
			// 立即发送
			isGo = true
		} else if noticeLog.NoticeType == enums.NoticeTypeClock {
			//定时发送
			if now > noticeLog.NoticeTime {
				isGo = true
			}
		} else if noticeLog.NoticeType == enums.NoticeTypeLoop {
			// 循环发送
			if (now-noticeLog.LastSendTime)/60 >= noticeLog.NoticeTime {
				isGo = true
			}
		}
	}
	if isGo == true {
		logs.Info("处理公告:%+v", noticeLog)
		nodeList := make([] string, 0)
		if noticeLog.IsAllServer == 0 {
			err := json.Unmarshal([]byte(noticeLog.NodeList), &nodeList)
			if err != nil {
				logs.Error("解析公告(%+v)区服列表失败%+v", id, err)
				return
			}
		} else {
			// 全服
			nodeList = models.GetAllGameNodeByPlatformId(noticeLog.PlatformId)
		}

		var result struct {
			ErrorCode int
		}

		var request struct {
			NodeList [] string `json:"nodeList"`
			Content  string    `json:"content"`
		}
		request.NodeList = nodeList
		request.Content = noticeLog.Content

		for _, node :=range nodeList {
			url := models.GetGameURLByNode(node) + "/send_notice"
			data, err := json.Marshal(request)
			utils.CheckError(err)
			sign := utils.String2md5(string(data) + enums.GmSalt)
			base64Data := base64.URLEncoding.EncodeToString([]byte(data))
			requestBody := "data=" + base64Data + "&sign=" + sign
			resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(requestBody))
			utils.CheckError(err)
			if err != nil {
				continue
			}
			defer resp.Body.Close()
			responseBody, err := ioutil.ReadAll(resp.Body)
			utils.CheckError(err)
			if err != nil {
				continue
			}
			logs.Info("result:%v", string(responseBody))

			err = json.Unmarshal(responseBody, &result)

			utils.CheckError(err)
			if err != nil {
				continue
			}
			if result.ErrorCode == 0 {
				logs.Info("发送公告成功 PlatformId: %v, id: %v", node, noticeLog.Id)
			} else {
				logs.Info("发送公告失败 PlatformId: %v, id: %v", node, noticeLog.Id)
			}
		}


		//for _, node := range nodeList {
		//	logs.Info("开始发送公告 PlatformId: %v, node: %v", noticeLog.PlatformId, node)
		//	conn, err := models.GetWsByNode(node)
		//	utils.CheckError(err, fmt.Sprintf("连接游戏服websocket失败 PlatformId: %v, node: %v", noticeLog.PlatformId, node))
		//	if err != nil {
		//		continue
		//	}
		//	defer conn.Close()
		//	request := gm.MSendNoticeTos{
		//		Token:   proto.String(""),
		//		Content: proto.String(noticeLog.Content),
		//	}
		//	mRequest, err := proto.Marshal(&request)
		//	utils.CheckError(err, "发送公告协议失败")
		//	if err != nil {
		//		continue
		//	}
		//	_, err = conn.Write(utils.Packet(9905, mRequest))
		//	utils.CheckError(err)
		//	if err != nil {
		//		continue
		//	}
		//	var receive = make([]byte, 100, 100)
		//	n, err := conn.Read(receive)
		//	utils.CheckError(err)
		//	if err != nil {
		//		continue
		//	}
		//	response := &gm.MSendNoticeToc{}
		//	data := receive[5:n]
		//	err = proto.Unmarshal(data, response)
		//	utils.CheckError(err)
		//	if err != nil {
		//		continue
		//	}
		//	if *response.Result == gm.MSendNoticeToc_success {
		//		logs.Info("发送公告成功 PlatformId: %v, node: %v", noticeLog.PlatformId, node)
		//	} else {
		//		logs.Info("发送公告失败 PlatformId: %v, node: %v", noticeLog.PlatformId, node)
		//	}
		if noticeLog.NoticeType != enums.NoticeTypeLoop {
			noticeLog.Status = 1
		}
		noticeLog.LastSendTime = now
		err = models.Db.Save(&noticeLog).Error
		utils.CheckError(err, "保存公告日志失败")
	}
}
