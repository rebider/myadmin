package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"github.com/chnzrb/myadmin/utils"
)

type DailyRegisterStatistics struct {
	Node            string `json:"node" gorm:"primary_key"`
	Time            int    `json:"time" gorm:"primary_key"`
	RegisterCount   int    `json:"registerCount"`
	CreateRoleCount int    `json:"createRoleCount"`
	ValidRoleCount  int    `json:"validRoleCount"`
}

type DailyRegisterStatisticsQueryParam struct {
	BaseQueryParam
	PlatformId int
	//ServerId   string
	Node       string `json:"serverId"`
	StartTime  int
	EndTime    int
}

func GetDailyRegisterStatisticsOne(node string, timestamp int) (*DailyRegisterStatistics, error) {
	data := &DailyRegisterStatistics{
		Node: node,
		Time: timestamp,
	}
	err := Db.Model(&DailyRegisterStatistics{}).First(&data).Error
	return data, err
}

func GetDailyRegisterStatisticsList(params *DailyRegisterStatisticsQueryParam) ([]*DailyRegisterStatistics, int64) {
	data := make([]*DailyRegisterStatistics, 0)
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
	f(Db.Model(&DailyRegisterStatistics{}).Where(&DailyRegisterStatistics{Node: params.Node})).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data)
	return data, count
}

func UpdateDailyRegisterStatistics(node string, timestamp int) error {
	logs.Info("UpdateDailyRegisterStatistics:%v, %v", node, timestamp)
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	defer gameDb.Close()

	m := &DailyRegisterStatistics{
		Node:            node,
		Time:            timestamp,
		RegisterCount:   GetThatDayRegister(gameDb, timestamp),
		CreateRoleCount: GetThatDayCreateRoleCount(gameDb, timestamp),
		ValidRoleCount:  GetThatDayValidCreateRoleCount(gameDb, timestamp),
	}
	err = Db.Save(&m).Error
	return err
}
