package models

import (
	"github.com/chnzrb/myadmin/utils"
	//"github.com/zaaksam/dproxy/go/db"
)

type PlayerPropLog struct {
	Id          int `json:"id"`
	PlayerId    int `json:"playerId"`
	PlayerName string `json:"playerName" gorm:"-"`
	PropType    int `json:"propType"`
	PropId      int `json:"propId"`
	OpType      int `json:"opType"`
	OpTime      int `json:"opTime"`
	ChangeValue int `json:"changeValue"`
	NewValue    int `json:"newValue"`
}

type PlayerPropLogQueryParam struct {
	BaseQueryParam
	PlatformId int
	ServerId   string
	Ip         string
	PlayerId   int
	PlayerName string
	StartTime  int
	EndTime    int
	PropType   int
	PropId     int
}

func GetPlayerPropLogList(params *PlayerPropLogQueryParam) ([]*PlayerPropLog, int64) {
	gameDb, err := GetGameDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	defer gameDb.Close()
	data := make([]*PlayerPropLog, 0)
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
	gameDb.Model(&PlayerPropLog{}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	for _,e := range data {
		e.PlayerName = GetPlayerName(gameDb, e.PlayerId)
	}
	return data, count
}
