package models

import (
	"github.com/chnzrb/myadmin/utils"
	//"github.com/zaaksam/dproxy/go/db"
)

type PlayerChallengeMissionLog struct {
	Id          int `json:"id"`
	PlayerId    int `json:"playerId"`
	PlayerName string `json:"playerName" gorm:"-"`
	MissionType int `json:"missionType"`
	MissionId   int `json:"missionId"`
	Result      int `json:"result"`
	Time        int `json:"time"`
	UsedTime    int `json:"usedTime"`
}

type PlayerChallengeMissionLogQueryParam struct {
	BaseQueryParam
	PlatformId  int
	ServerId    string
	Ip          string
	PlayerId    int
	PlayerName  string
	StartTime   int
	EndTime     int
	MissionType int
}

func GetPlayerChallengeMissionLogList(params *PlayerChallengeMissionLogQueryParam) ([]*PlayerChallengeMissionLog, int64) {
	gameDb, err := GetGameDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	defer gameDb.Close()
	data := make([]*PlayerChallengeMissionLog, 0)
	var count int64
	sortOrder := "id"
	//if params.Order == "descending" {
	//	sortOrder = sortOrder + " desc"
	//}
	//if params.Ip != "" {
	//	gameDb = gameDb.Where("ip = ?", params.Ip)
	//}
	//if params.PlayerId != 0 {
	//	gameDb = gameDb.Where("player_id = ?", params.PlayerId)
	//}
	//if params.StartTime != 0 {
	//	gameDb = gameDb.Where("timestamp >= ?", params.StartTime)
	//}
	//if params.EndTime != 0 {
	//	gameDb = gameDb.Where("timestamp <= ?", params.EndTime)
	//}
	gameDb.Model(&PlayerChallengeMissionLog{}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	for _,e := range data {
		e.PlayerName = GetPlayerName(gameDb, e.PlayerId)
	}
	return data, count
}
