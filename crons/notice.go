package crons

import (
	"github.com/astaxie/beego/logs"
	"time"
	"encoding/json"
	"github.com/chnzrb/myadmin/proto"
	"github.com/chnzrb/myadmin/utils"
	"github.com/chnzrb/myadmin/models"
	"github.com/golang/protobuf/proto"
	"fmt"
)



func ClockDealAllNoticeLog() {
	logs.Info("定时处理公告")
	noticeLogList := models.GetAllNoticeLog()
	for _, noticeLog := range noticeLogList {
		if noticeLog.Status == 0 {
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
	if noticeLog.Status == 0 {
		if noticeLog.NoticeType == 0 {
			// 立即发送
			isGo = true
		} else if noticeLog.NoticeType == 1 {
			//定时发送
			if now > noticeLog.NoticeTime {
				isGo = true
			}
		} else if noticeLog.NoticeType == 2 {
			// 循环发送
			if (now - noticeLog.LastSendTime) / 60 > noticeLog.NoticeTime {
				isGo = true
			}
		}
	}
	if isGo == true {
		logs.Info("处理公告:%+v", noticeLog)
		serverIdList := make([] string, 0)
		err := json.Unmarshal([]byte(noticeLog.ServerIdList), &serverIdList)
		if err != nil {
			logs.Error("解析公告(%+v)区服列表失败%+v", id, err)
			return
		}
		for _, serverId := range serverIdList {
			logs.Info("发送公告 PlatformId: %v, serverId: %v", noticeLog.PlatformId, serverId)
			conn, err := models.GetWsByPlatformIdAndSid(noticeLog.PlatformId, serverId)
			utils.CheckError(err, fmt.Sprintf("连接游戏服失败 PlatformId: %v, serverId: %v", noticeLog.PlatformId, serverId))
			if err != nil {
				continue
			}
			defer conn.Close()
			request := gm.MSendNoticeTos{
				Token:   proto.String(""),
				Content: proto.String(noticeLog.Content),
			}
			mRequest, err := proto.Marshal(&request)
			utils.CheckError(err)
			if err != nil {
				continue
			}
			_, err = conn.Write(utils.Packet(9905, mRequest))
			utils.CheckError(err)
			if err != nil {
				continue
			}
			var receive = make([]byte, 100, 100)
			n, err := conn.Read(receive)
			utils.CheckError(err)
			if err != nil {
				continue
			}
			respone := &gm.MSendNoticeToc{}
			data := receive[5:n]
			err = proto.Unmarshal(data, respone)
			utils.CheckError(err)
			if err != nil {
				continue
			}
			if *respone.Result == gm.MSendNoticeToc_success {
				//logs.Info("发送公告 PlatformId: %v, serverId: %v", noticeLog.PlatformId, serverId)
			} else {
				logs.Info("发送公告失败 PlatformId: %v, serverId: %v", noticeLog.PlatformId, serverId)
				//logs.Error("发送公告失败:%+v", request)
			}
		}
		noticeLog.Status = 1
		noticeLog.LastSendTime = now
		err = models.Db.Save(&noticeLog).Error
		utils.CheckError(err, "保存公告日志失败")
	}
}
