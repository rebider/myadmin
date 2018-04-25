package models

import (
	"github.com/chnzrb/myadmin/utils"
	//"github.com/zaaksam/dproxy/go/db"
	"github.com/jinzhu/gorm"
)

type PlayerPropLog struct {
	Id          int    `json:"id"`
	PlayerId    int    `json:"playerId"`
	PlayerName  string `json:"playerName" gorm:"-"`
	PropType    int    `json:"propType"`
	PropId      int    `json:"propId"`
	OpType      int    `json:"opType"`
	OpTime      int    `json:"opTime"`
	ChangeValue int    `json:"changeValue"`
	NewValue    int    `json:"newValue"`
}

type PlayerPropLogQueryParam struct {
	BaseQueryParam
	PlatformId int
	Node   string `json:"serverId"`
	Ip         string
	PlayerId   int
	PlayerName string
	StartTime  int
	EndTime    int
	PropType   int
	PropId     int
	OpType    int
	Type       int //1：获得 2：消耗
}

func GetPlayerPropLogList(params *PlayerPropLogQueryParam) ([]*PlayerPropLog, int64) {
	gameDb, err := GetGameDbByNode(params.Node)
	utils.CheckError(err)
	if err != nil {
		return nil, 0
	}
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
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("op_time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}
	f1 := func(db *gorm.DB) *gorm.DB {
		if params.Type == 1 {
			return db.Where("change_value > 0")
		}
		if params.Type == 2 {
			return db.Where("change_value < 0")
		}
		return db
	}
	f1(f(gameDb.Model(&PlayerPropLog{}).Where(&PlayerPropLog{
		PlayerId: params.PlayerId,
		PropType: params.PropType,
		PropId:   params.PropId,
		OpType:params.OpType,
	}))).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	for _, e := range data {
		e.PlayerName = GetPlayerName(gameDb, e.PlayerId)
	}
	return data, count
}
