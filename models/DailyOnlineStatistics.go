package models

import (
	"github.com/astaxie/beego/logs"
	//"github.com/jinzhu/gorm"
	//"sort"
	"github.com/chnzrb/myadmin/utils"
)

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
	Node      string `json:"serverId"`
	StartTime int
	EndTime   int
}

func GetDailyOnlineStatisticsList(params *DailyOnlineStatisticsQueryParam) []*DailyOnlineStatistics {
	data := make([]*DailyOnlineStatistics, 0)
	//mapData := make(map[int]*DailyOnlineStatistics, 0)
	//result := make([] *DailyOnlineStatistics, 0)

	//count := 0

	//serverNode, err := GetGameServerOne(params.PlatformId, params.ServerId)
	//if err != nil {
	//	return nil, 0
	//}
	//params.Node = serverNode.Node

	//for i
	//logs.Info("aaaaaaaaaaaaaaaaaa")
	for i := params.StartTime; i <= params.EndTime; i = i + 86400 {
		tmpData := make([]*DailyOnlineStatistics, 0)
		//logs.Info("gggggggggggggggg:%+v", i)
		//f := func(db *gorm.DB) *gorm.DB {
		//	if params.StartTime > 0 {
		//		return db.Where("time between ? and ?", i, i + 86400)
		//	}
		//	return db
		//}
		err := Db.Model(&DailyOnlineStatistics{}).Where(&DailyOnlineStatistics{Node: params.Node, Time:i}).Find(&tmpData).Error
		utils.CheckError(err)
		//logs.Info("isNotFound:%+v", isNotFound)
		//logs.Info("len:%+v", len(tmpData))
		if len(tmpData) > 0 {
			tmpE := &DailyOnlineStatistics{
				Node:params.Node,
				Time:i,
			}
			for _, e := range tmpData {
				tmpE.MaxOnlineCount += e.MaxOnlineCount
				tmpE.MinOnlineCount += e.MinOnlineCount
				tmpE.AvgOnlineCount += e.AvgOnlineCount
				tmpE.AvgOnlineTime += e.AvgOnlineTime
			}
			data = append(data, tmpE)
		}
	}
	//for _, e := range data {
	//	//logs.Info("%+v", e)
	//	if r, ok := mapData[e.Time]; ok == true {
	//		r.MaxOnlineCount += e.MaxOnlineCount
	//		r.MinOnlineCount += e.MinOnlineCount
	//		r.AvgOnlineCount += e.AvgOnlineCount
	//		r.AvgOnlineTime += e.AvgOnlineTime
	//	} else {
	//			mapData[e.Time] = e
	//	}
	//}
	//var keys [] int
	//for key, _ := range mapData {
	//	keys = append(keys, key)
	//}
	//sort.Ints(keys)
	//
	//for _, id := range keys {
	//	e := mapData[id]
	//	//logs.Info("666:%+v", e)
	//	result = append(result, e)
	//	count ++
	//}

	return data
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
