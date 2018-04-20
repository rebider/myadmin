package models

import (
	"github.com/chnzrb/myadmin/utils"
)

type PlayerChargeInfoRecord struct {
	PlayerId      int    `json:"playerId"`
	PlayerName    string `json:"playerName" gorm:"-"`
	PlatformId    int    `json:"platformId" gorm:"column:part_id"`
	ServerId      string `json:"serverId"`
	TotalMoney    int    `json:"totalMoney"`
	ChargeCount   int    `json:"chargeCount"`
	LastLoginTime int    `json:"lastLoginTime" gorm:"-"`
	RegisterTime  int    `json:"registerTime" gorm:"-"`
	LastChargeTime  int    `json:"lastChargeTime" gorm:"column:last_time"`
	FirstChargeTime  int    `json:"firstChargeTime" gorm:"column:first_time"`
}

type PlayerChargeDataQueryParam struct {
	BaseQueryParam
	PlatformId int
	Node   string `json:"serverId"`
}

func GetPlayerChargeDataList(params *PlayerChargeDataQueryParam) ([]*PlayerChargeInfoRecord, int64) {
	data := make([]*PlayerChargeInfoRecord, 0)
	var count int64
	sortOrder := "total_money desc"
	//if params.Order == "descending" {
	//	sortOrder = sortOrder + " desc"
	//}
	if params.Node == "" {
		DbCharge.Model(&PlayerChargeInfoRecord{}).Count(&count).Where(&PlayerChargeInfoRecord{
			PlatformId: params.PlatformId,
		}).Where("charge_count > 0 ").Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	} else {
		DbCharge.Model(&PlayerChargeInfoRecord{}).Count(&count).Where(&PlayerChargeInfoRecord{
			PlatformId: params.PlatformId,
		}).Where("charge_count > 0 ").Where("server_id in (?)", GetGameServerIdListByNode(params.Node)).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	}

	for _, e := range data {
		gameDb, err := GetGameDbByPlatformIdAndSid(e.PlatformId, e.ServerId)
		utils.CheckError(err)
		defer gameDb.Close()
		player, err := GetPlayerByDb(gameDb, e.PlayerId)
		e.PlayerName = player.Nickname
		e.LastLoginTime = player.LastLoginTime
		e.RegisterTime = player.RegTime
	}
	return data, count
}
