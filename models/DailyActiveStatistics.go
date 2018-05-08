package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	//"github.com/jinzhu/gorm"
)

type DailyActiveStatistics struct {
	Node             string `json:"node"`
	Time             int    `json:"time"`
	LoginTimes       int    `json:"loginTimes"`
	LoginPlayerCount int    `json:"loginPlayerCount"`
	//AvgLoginCount     int    `json:"avgLoginCount" gorm:"-"`
	ActivePlayerCount    int `json:"activePlayerCount"`
	CreateRoleCount      int `json:"createRoleCount"`
	TotalCreateRoleCount int `json:"totalCreateRoleCount"`
	//ActiveRate        int    `json:"activeRate" gorm:"-"`
}

type DailyActiveStatisticsQueryParam struct {
	BaseQueryParam
	PlatformId int
	//ServerId   string
	Node      string `json:"serverId"`
	StartTime int
	EndTime   int
}

func GetDailyActiveStatisticsList(params *DailyActiveStatisticsQueryParam) []*DailyActiveStatistics {
	//data := make([]*DailyActiveStatistics, 0)
	//var count int64
	//f := func(db *gorm.DB) *gorm.DB {
	//	if params.StartTime > 0 {
	//		return db.Where("time between ? and ?", params.StartTime, params.EndTime)
	//	}
	//	return db
	//}
	//f(Db.Model(&DailyActiveStatistics{}).Where(&DailyActiveStatistics{Node: params.Node})).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data)
	//return data, count
	data := make([]*DailyActiveStatistics, 0)
	for i := params.StartTime; i <= params.EndTime; i = i + 86400 {
		tmpData := make([]*DailyActiveStatistics, 0)
		err := Db.Model(&DailyActiveStatistics{}).Where(&DailyActiveStatistics{Node: params.Node, Time: i}).Find(&tmpData).Error
		utils.CheckError(err)
		if len(tmpData) > 0 {
			tmpE := &DailyActiveStatistics{
				Node: params.Node,
				Time: i,
			}
			for _, e := range tmpData {
				tmpE.LoginTimes += e.LoginTimes
				tmpE.LoginPlayerCount += e.LoginPlayerCount
				tmpE.ActivePlayerCount += e.ActivePlayerCount
				tmpE.CreateRoleCount += e.CreateRoleCount
				tmpE.TotalCreateRoleCount += e.TotalCreateRoleCount
			}
			data = append(data, tmpE)
		}
	}
	return data
}

func UpdateDailyActiveStatistics(node string, timestamp int) error {
	logs.Info("UpdateDailyActiveStatistics:%v, %v", node, timestamp)
	gameDb, err := GetGameDbByNode(node)
	if err != nil {
		return err
	}
	defer gameDb.Close()
	m := &DailyActiveStatistics{
		Node:              node,
		Time:              timestamp,
		LoginTimes:        GetThatDayLoginTimes(gameDb, timestamp),
		LoginPlayerCount:  GetThatDayLoginPlayerCount(gameDb, timestamp),
		ActivePlayerCount: GetThatDayActivePlayerCount(gameDb, timestamp),
		CreateRoleCount:   GetThatDayCreateRoleCount(gameDb, timestamp),
		TotalCreateRoleCount: GetTotalCreateRoleCount(gameDb),
	}
	err = Db.Save(&m).Error
	return err
}
