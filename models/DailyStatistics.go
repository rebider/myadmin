package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
)

//import "github.com/chnzrb/myadmin/utils"

type DailyStatistics struct {
	Node          string `json:"node"`
	Time          int    `json:"time"`
	RegisterNum   int    `json:"registerNum"`
	CreateRoleNum int    `json:"createRoleNum"`
	MaxOnline     int    `json:"maxOnline"`
	MinOnline     int    `json:"minOnline"`
	AvgOnline     int    `json:"avgOnline"`
}

type DailyStatisticsQueryParam struct {
	BaseQueryParam
	PlatformId int
	ServerId string
	Node      string
	StartTime int
	EndTime   int
}

func GetDailyStatisticsOne(node string, timestamp int) (*DailyStatistics, error) {
	data := &DailyStatistics{
		Node: node,
		Time: timestamp,
	}
	err := Db.Model(&DailyStatistics{}).First(&data).Error
	return data, err
}

func GetDailyStatisticsList(params *DailyStatisticsQueryParam) ([]*DailyStatistics, int64) {
	data := make([]*DailyStatistics, 0)
	var count int64
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}
	f(Db.Model(&DailyStatistics{}).Where(&DailyStatistics{Node:params.Node})).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data)
	return data, count
}

func UpdateDailyStatistics(node string, timestamp int) error {
	logs.Info("更新DailyStatistics:%v, %v", node, timestamp)
	serverNode, err := GetServerNode(node)
	if err != nil {
		return err
	}
	gameDb, err := GetGameDbByServerNode(serverNode)
	if err != nil {
		return err
	}
	defer gameDb.Close()
	m := &DailyStatistics{
		Node:          node,
		Time:          timestamp,
		RegisterNum:   GetThatDayCreateRole(gameDb, timestamp),
		CreateRoleNum: GetThatDayCreateRole(gameDb, timestamp),
		MaxOnline:GetThatDayMaxOnlineCount(node, timestamp),
		MinOnline:GetThatDayMinOnlineCount(node, timestamp),
		AvgOnline:int(GetThatDayAverageOnlineCount(node, timestamp)),
	}
	err = Db.Save(&m).Error
	return err
}
