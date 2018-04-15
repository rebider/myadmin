package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
)

//import "github.com/chnzrb/myadmin/utils"

type DailyChargeStatistics struct {
	Node                 string `json:"node"`
	Time                 int    `json:"time"`
	ChargeMoney          int    `json:"chargeMoney"`
	ChargePlayerCount    int    `json:"chargePlayerCount"`
	ARPU                 float32    `json:"arpu"`
	NewChargePlayerCount int    `json:"newChargePlayerCount"`
}

type DailyChargeStatisticsQueryParam struct {
	BaseQueryParam
	PlatformId int
	ServerId   string
	Node       string
	StartTime  int
	EndTime    int
}

func GetDailyChargeStatisticsList(params *DailyChargeStatisticsQueryParam) ([]*DailyChargeStatistics, int64) {
	data := make([]*DailyChargeStatistics, 0)
	var count int64
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}
	serverNode, err := GetGameServerOne(params.PlatformId, params.ServerId)
	if err != nil {
		return nil, 0
	}
	params.Node = serverNode.Node
	f(Db.Model(&DailyChargeStatistics{}).Where(&DailyChargeStatistics{Node: params.Node})).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data)
	for _, e := range data {
		e.ARPU = CaclARPU(e.ChargeMoney, e.ChargePlayerCount)
	}
	return data, count
}

func UpdateDailyChargeStatistics(node string, timestamp int) error {
	logs.Info("更新DailyChargeStatistics:%v, %v", node, timestamp)
	m := &DailyChargeStatistics{
		Node:                 node,
		Time:                 timestamp,
		ChargeMoney:          GetThatDayServerTotalChargeMoney(node, timestamp),
		ChargePlayerCount:    GetThatDayServerChargePlayerCount(node, timestamp),
		NewChargePlayerCount: GetThadDayServerSecondChargePlayerCount(node, timestamp),
	}
	err := Db.Save(&m).Error
	return err
}
