package models

import (
)

type ForbidLog struct {
	PlatformId int    `json:"platformId" gorm:"primary_key"`
	ServerId   string `json:"serverId" gorm:"primary_key"`
	PlayerName string `json:"playerName" gorm:"primary_key"`
	ForbidType int32  `json:"forbidType"`
	ForbidTime int32  `json:"forbidTime"`
	Time       int64  `json:"time"`
	UserId     int    `json:"userId"`
	UserName   string `json:"userName" gorm:"-"`
}

type ForbidLogQueryParam struct {
	BaseQueryParam
	PlatformId int
	ServerId   string
	PlayerName string
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
	Db.Model(&ForbidLog{}).Where(&ForbidLog{
		PlatformId: params.PlatformId,
		ServerId:   params.ServerId,
		PlayerName: params.PlayerName,
	}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	for _, e := range data {
		u, err := GetUserOne(e.UserId)
		if err == nil {
			e.UserName = u.Name
		}
	}
	return data, count
}
