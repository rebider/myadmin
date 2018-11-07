package models

import (
	"github.com/astaxie/beego/logs"
)

type ForbidLog struct {
	PlatformId string    `json:"platformId" gorm:"primary_key"`
	ServerId   string `json:"serverId" gorm:"primary_key"`
	PlayerId int `json:"playerId" gorm:"primary_key"`
	ForbidType int32  `json:"forbidType"`
	ForbidTime int32  `json:"forbidTime"`
	Time       int64  `json:"time"`
	UserId     int    `json:"userId"`
	UserName   string `json:"userName" gorm:"-"`
	PlayerName   string `json:"playerName" gorm:"-"`
}

type ForbidLogQueryParam struct {
	BaseQueryParam
	PlatformId string
	ServerId       string `json:"serverId"`
	PlayerName string
	PlayerId  int
	StartTime  int
	EndTime    int
	UserId     int
}

func GetForbidLogList(params *ForbidLogQueryParam) ([]*ForbidLog, int64) {
	data := make([]*ForbidLog, 0)
	var count int64
	sortOrder := "time"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	if params.ServerId == "" {
		Db.Model(&ForbidLog{}).Where(&ForbidLog{
			PlatformId: params.PlatformId,
			PlayerId: params.PlayerId,
		}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	} else {
		//serverIdList := GetGameServerIdListByNode(params.ServerId)
		Db.Model(&ForbidLog{}).Where(&ForbidLog{
			PlatformId: params.PlatformId,
			PlayerId: params.PlayerId,
			ServerId:params.ServerId,
		}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	}

	for _, e := range data {
		u, err := GetUserOne(e.UserId)
		if err == nil {
			e.UserName = u.Name
			e.PlayerName = GetPlayerName_2(params.PlatformId, e.ServerId, e.PlayerId)
		}
	}
	logs.Info(count)
	return data, count
}
