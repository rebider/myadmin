package models

import (
	//"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
)

type TenMinuteStatistics struct {
	//Node          string `json:"node" gorm:"primary_key"`
	PlatformId        string `json:"platformId" gorm:"primary_key"`
	ServerId          string `json:"serverId" gorm:"primary_key"`
	Channel           string `json:"channel" gorm:"primary_key"`
	Time              int    `json:"time" gorm:"primary_key"`
	OnlineCount       int    `json:"onlineNum"`
	RegisterCount     int    `json:"registerCount"`
	CreateRoleCount   int    `json:"createRoleCount"`
	ChargeCount       int    `json:"chargeCount"`
	ChargePlayerCount int    `json:"chargePlayerCount"`
}

func UpdateTenMinuteStatistics(platformId string, serverId string, channel string, timestamp int) error {
	//logs.Info("UpdateTenMinuteStatistics:%v, %v, %v, %v", platformId, serverId, channel, timestamp)
	serverNode, err := GetGameServerOne(platformId, serverId)
	if err != nil {
		return err
	}
	//node := serverNode.Node
	gameDb, err := GetGameDbByNode(serverNode.Node)
	if err != nil {
		return err
	}
	defer gameDb.Close()

	zeroTime := utils.GetThatZeroTimestamp(int64(timestamp))

	if timestamp == zeroTime {
		zeroTime = zeroTime - 86400
	}

	m := &TenMinuteStatistics{
		//Node:          serverNode.Node,
		PlatformId:      platformId,
		ServerId:        serverId,
		Channel:         channel,
		Time:            timestamp,
		OnlineCount:     GetNowOnlineCount2(gameDb, serverId, [] string{channel}),
		RegisterCount:   GetRegisterRoleCount(gameDb, serverId, channel, timestamp-600, timestamp-1),
		CreateRoleCount: GetCreateRoleCount(gameDb, serverId, channel, timestamp-600, timestamp-1),
		ChargeCount:     GetTotalChargeMoney(platformId, serverId, channel, timestamp-600, timestamp-1),
		ChargePlayerCount:GetTotalChargePlayerCount(platformId, serverId, channel, zeroTime, timestamp),
	}
	err = Db.Save(&m).Error
	return err
}

func RepireTenMinuteStatistics() {
	logs.Info("开始修复每10分钟统计")
	gameServerList, _ := GetAllGameServerDirty()
	for i := 1538928000; i <= 1539156000; i += 600 {
		for _, gameServer := range gameServerList {
			//err := models.UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
			//utils.CheckError(err)
			platformId := gameServer.PlatformId
			//if platformId == "af" || platformId == "djs" {
			serverId := gameServer.Sid
			channelList := GetChannelListByPlatformId(platformId)
			if len(channelList) == 0 {
				logs.Error("渠道未配置:%v %+v", platformId, channelList)
			}

			if i < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
				continue
			}
			for _, channel := range channelList {
				err := DoRepireTenMinuteStatistics(platformId, serverId, channel.Channel, i)
				utils.CheckError(err)
			}
			//}
		}
	}
	logs.Info("修复每10分钟统计完毕")
}
func DoRepireTenMinuteStatistics(platformId string, serverId string, channel string, timestamp int) error {
	m := &TenMinuteStatistics{
		PlatformId: platformId,
		ServerId:   serverId,
		Channel:    channel,
		Time:       timestamp,
	}
	zeroTime := utils.GetThatZeroTimestamp(int64(timestamp))
	if timestamp == zeroTime {
		zeroTime = zeroTime - 86400
	}
	c := GetTotalChargePlayerCount(platformId, serverId, channel, zeroTime, timestamp)
	err := Db.Model(&m).Update("charge_player_count", c).Error
	return err
}


//func RepireTenMinuteStatistics() {
//	logs.Info("开始修复每10分钟统计")
//	gameServerList, _ := GetAllGameServerDirty()
//	for i := 1538755200; i <= 1538978400; i += 600 {
//		for _, gameServer := range gameServerList {
//			//err := models.UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
//			//utils.CheckError(err)
//			platformId := gameServer.PlatformId
//			//if platformId == "af" || platformId == "djs" {
//			serverId := gameServer.Sid
//			channelList := GetChannelListByPlatformId(platformId)
//			if len(channelList) == 0 {
//				logs.Error("渠道未配置:%v %+v", platformId, channelList)
//			}
//
//			if i < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
//				continue
//			}
//			for _, channel := range channelList {
//				err := DoRepireTenMinuteStatistics(platformId, serverId, channel.Channel, i)
//				utils.CheckError(err)
//			}
//			//}
//		}
//	}
//	logs.Info("修复每10分钟统计完毕")
//}
//func DoRepireTenMinuteStatistics(platformId string, serverId string, channel string, timestamp int) error {
//	m := &TenMinuteStatistics{
//		PlatformId: platformId,
//		ServerId:   serverId,
//		Channel:    channel,
//		Time:       timestamp,
//	}
//	c := GetTotalChargeMoney(platformId, serverId, channel, timestamp-600, timestamp-1)
//	err := Db.Model(&m).Update("charge_count", c).Error
//	return err
//}
