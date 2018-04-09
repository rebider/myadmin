package models

import (
	"github.com/chnzrb/myadmin/utils"
)

type PlayerLoginLog struct {
	Id        int    `json:"id"`
	PlayerId  int    `json:"playerId"`
	Ip        string `json:"ip"`
	Timestamp int    `json:"time"`
}

type PlayerLoginLogQueryParam struct {
	BaseQueryParam
	PlatformId int
	ServerId   string
	Ip         string
	PlayerId   int
	PlayerName string
	StartTime  int
	EndTime    int
}

func GetPlayerLoginLogList(params *PlayerLoginLogQueryParam) ([]*PlayerLoginLog, int64) {
	db, err := GetDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	defer db.Close()
	data := make([]*PlayerLoginLog, 0)
	var count int64
	sortOrder := "id"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	if params.Ip != "" {
		db = db.Where("ip = ?", params.Ip)
	}
	if params.PlayerId != 0 {
		db = db.Where("player_id = ?", params.PlayerId)
	}
	if params.StartTime != 0 {
		db = db.Where("timestamp >= ?", params.StartTime)
	}
	if params.EndTime != 0 {
		db = db.Where("timestamp <= ?", params.EndTime)
	}
	db.Model(&PlayerLoginLog{}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	return data, count
}
