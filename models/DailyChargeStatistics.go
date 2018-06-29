package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
)

type DailyChargeStatistics struct {
	Node                 string  `json:"node"`
	Time                 int     `json:"time"`
	ChargeMoney          int     `json:"chargeMoney"`
	ChargePlayerCount    int     `json:"chargePlayerCount"`
	ARPU                 float32 `json:"arpu" gorm:"-"`
	Active_ARPU          float32 `json:"activeArpu" gorm:"-"`
	AcriveChargeRate     int     `json:"activeChargeRate" gorm:"-"`
	NewChargePlayerCount int     `json:"newChargePlayerCount"`
}

type DailyChargeStatisticsQueryParam struct {
	BaseQueryParam
	PlatformId string
	Node       string `json:"serverId"`
	StartTime  int
	EndTime    int
}

func GetDailyChargeStatisticsList(params *DailyChargeStatisticsQueryParam) []*DailyChargeStatistics {
	//data := make([]*DailyChargeStatistics, 0)
	//var count int64
	//f := func(db *gorm.DB) *gorm.DB {
	//	if params.StartTime > 0 {
	//		return db.Where("time between ? and ?", params.StartTime, params.EndTime)
	//	}
	//	return db
	//}
	//f(Db.Model(&DailyChargeStatistics{}).Where(&DailyChargeStatistics{Node: params.Node})).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data)
	//for _, e := range data {
	//	e.ARPU = CaclARPU(e.ChargeMoney, e.ChargePlayerCount)
	//}
	//return data, count
	data := make([]*DailyChargeStatistics, 0)
	for i := params.StartTime; i <= params.EndTime; i = i + 86400 {
		tmpData := make([]*DailyChargeStatistics, 0)
		err := Db.Model(&DailyChargeStatistics{}).Where(&DailyChargeStatistics{Node: params.Node, Time: i}).Find(&tmpData).Error
		utils.CheckError(err)
		if len(tmpData) > 0 {
			tmpE := &DailyChargeStatistics{
				Node: params.Node,
				Time: i,
			}
			for _, e := range tmpData {
				tmpE.ChargeMoney += e.ChargeMoney
				tmpE.ChargePlayerCount += e.ChargePlayerCount
				tmpE.NewChargePlayerCount += e.NewChargePlayerCount
				//tmpE.CreateRoleCount += e.CreateRoleCount
			}
			data = append(data, tmpE)
		}
	}
	for _, e := range data {
		e.ARPU = CaclARPU(e.ChargeMoney, e.ChargePlayerCount)
	}
	return data
}

func UpdateDailyChargeStatistics(node string, timestamp int) error {
	logs.Info("更新DailyChargeStatistics:%v, %v", node, timestamp)
	m := &DailyChargeStatistics{
		Node:                 node,
		Time:                 timestamp,
		ChargeMoney:          GetThatDayServerTotalChargeMoney(node, timestamp),
		ChargePlayerCount:    GetThatDayServerChargePlayerCount(node, timestamp),
		NewChargePlayerCount: GetThadDayServerFirstChargePlayerCount(node, timestamp),
	}
	err := Db.Save(&m).Error
	return err
}
