package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
)

//import "github.com/chnzrb/myadmin/utils"

type DailyOnlineStatistics struct {
	Node           string `json:"node"`
	Time           int    `json:"time"`
	MaxOnlineCount int    `json:"maxOnline"`
	MinOnlineCount int    `json:"minOnline"`
	AvgOnlineCount int    `json:"avgOnline"`
	AvgOnlineTime  int    `json:"avgOnlineTime"`
}

type DailyOnlineStatisticsQueryParam struct {
	BaseQueryParam
	PlatformId int
	//ServerId   string
	Node       string `json:"serverId"`
	StartTime  int
	EndTime    int
}

func GetDailyOnlineStatisticsList(params *DailyOnlineStatisticsQueryParam) ([]*DailyOnlineStatistics, int64) {
	data := make([]*DailyOnlineStatistics, 0)
	var count int64
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}
	//serverNode, err := GetGameServerOne(params.PlatformId, params.ServerId)
	//if err != nil {
	//	return nil, 0
	//}
	//params.Node = serverNode.Node
	f(Db.Model(&DailyOnlineStatistics{}).Where(&DailyOnlineStatistics{Node: params.Node})).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data)
	return data, count
}

func UpdateDailyOnlineStatistics(node string, timestamp int) error {
	logs.Info("UpdateDailyOnlineStatistics:%v, %v", node, timestamp)
	m := &DailyOnlineStatistics{
		Node:           node,
		Time:           timestamp,
		MaxOnlineCount: GetThatDayMaxOnlineCount(node, timestamp),
		MinOnlineCount: GetThatDayMinOnlineCount(node, timestamp),
		AvgOnlineCount: int(GetThatDayAverageOnlineCount(node, timestamp)),
		AvgOnlineTime:  GetAvgOnlineTime(node, timestamp),
	}
	err := Db.Save(&m).Error
	return err
}
