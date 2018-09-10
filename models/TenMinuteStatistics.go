package models

import (
	//"github.com/astaxie/beego/logs"
)

type TenMinuteStatistics struct {
	//Node          string `json:"node" gorm:"primary_key"`
	PlatformId    string `json:"platformId" gorm:"primary_key"`
	ServerId      string `json:"serverId" gorm:"primary_key"`
	Channel       string `json:"channel" gorm:"primary_key"`
	Time          int    `json:"time" gorm:"primary_key"`
	OnlineCount   int    `json:"onlineNum"`
	RegisterCount int    `json:"registerCount"`
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

	m := &TenMinuteStatistics{
		//Node:          serverNode.Node,
		PlatformId:    platformId,
		ServerId:      serverId,
		Channel:       channel,
		Time:          timestamp,
		OnlineCount:   GetNowOnlineCount2(gameDb, serverId, [] string{channel}),
		RegisterCount: GetRegisterRoleCount(gameDb, serverId, channel, timestamp-600, timestamp),
	}
	err = Db.Save(&m).Error
	return err
}
