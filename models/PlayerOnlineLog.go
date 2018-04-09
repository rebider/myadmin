package models

import (
	"github.com/chnzrb/myadmin/utils"
)

type PlayerOnlineLog struct {
	Id          int `json:"id"`
	PlayerId    int `json:"playerId"`
	LoginTime   int `json:"loginTime"`
	OfflineTime int `json:"offlineTime"`
	OnlineTime  int `json:"onlineTime"`
}

type PlayerOnlineLogQueryParam struct {
	BaseQueryParam
	PlatformId int
	ServerId   string
	PlayerId   int
	PlayerName string
	StartTime  int
	EndTime    int
}

func GetPlayerOnlineLogList(params *PlayerOnlineLogQueryParam) ([]*PlayerOnlineLog, int64) {
	db, err := GetDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	defer db.Close()
	data := make([]*PlayerOnlineLog, 0)
	var count int64
	sortOrder := "id"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	if params.PlayerId != 0 {
		db = db.Where("player_id = ?", params.PlayerId)
	}
	if params.StartTime != 0 {
		db = db.Where("offline_time >= ?", params.StartTime)
	}
	if params.EndTime != 0 {
		db = db.Where("offline_time <= ?", params.EndTime)
	}
	db.Model(&PlayerOnlineLog{}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	return data, count
}
