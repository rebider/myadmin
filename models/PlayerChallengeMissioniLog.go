package models

import (
	"github.com/chnzrb/myadmin/utils"
	//"github.com/zaaksam/dproxy/go/db"
	"github.com/jinzhu/gorm"
)

type PlayerChallengeMissionLog struct {
	Id          int    `json:"id"`
	PlayerId    int    `json:"playerId"`
	PlayerName  string `json:"playerName" gorm:"-"`
	MissionType int    `json:"missionType"`
	MissionId   int    `json:"missionId"`
	Result      int    `json:"result"`
	Time        int    `json:"time"`
	UsedTime    int    `json:"usedTime"`
}

type PlayerChallengeMissionLogQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId        string `json:"serverId"`
	Ip          string
	PlayerId    int
	PlayerName  string
	StartTime   int
	EndTime     int
	MissionType int
}

func GetPlayerChallengeMissionLogList(params *PlayerChallengeMissionLogQueryParam) ([]*PlayerChallengeMissionLog, int64) {
	gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	if err != nil {
		return nil, 0
	}
	node := gameServer.Node
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, 0
	}
	defer gameDb.Close()
	data := make([]*PlayerChallengeMissionLog, 0)
	var count int64
	sortOrder := "id"
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}
	f(gameDb.Model(&PlayerChallengeMissionLog{}).Where(&PlayerChallengeMissionLog{
		PlayerId:    params.PlayerId,
		MissionType: params.MissionType,
	})).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	for _, e := range data {
		e.PlayerName = GetPlayerName(gameDb, e.PlayerId)
	}
	return data, count
}
