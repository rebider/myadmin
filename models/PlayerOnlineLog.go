package models

import (
	"github.com/chnzrb/myadmin/utils"
	"github.com/jinzhu/gorm"
)

type PlayerOnlineLog struct {
	Id          int    `json:"id"`
	PlayerId    int    `json:"playerId"`
	PlayerName  string `json:"playerName" gorm:"-"`
	LoginTime   int    `json:"loginTime"`
	OfflineTime int    `json:"offlineTime"`
	OnlineTime  int    `json:"onlineTime"`
}

type PlayerOnlineLogQueryParam struct {
	BaseQueryParam
	PlatformId int
	Node        string `json:"serverId"`
	PlayerId   int
	PlayerName string
	StartTime  int
	EndTime    int
}

func GetPlayerOnlineLogList(params *PlayerOnlineLogQueryParam) ([]*PlayerOnlineLog, int64) {
	gameDb, err := GetGameDbByNode(params.Node)
	utils.CheckError(err)
	defer gameDb.Close()
	data := make([]*PlayerOnlineLog, 0)
	var count int64
	sortOrder := "id"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	//if params.PlayerId != 0 {
	//	gameDb = gameDb.Where("player_id = ?", params.PlayerId)
	//}
	//if params.StartTime != 0 {
	//	gameDb = gameDb.Where("offline_time >= ?", params.StartTime)
	//}
	//if params.EndTime != 0 {
	//	gameDb = gameDb.Where("offline_time <= ?", params.EndTime)
	//}
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("offline_time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}
	f(gameDb.Model(&PlayerOnlineLog{}).Where(&PlayerOnlineLog{
		PlayerId: params.PlayerId,
	})).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	for _, e := range data {
		e.PlayerName = GetPlayerName(gameDb, e.PlayerId)
	}
	return data, count
}
