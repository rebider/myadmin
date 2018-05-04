package models

import (
	"github.com/chnzrb/myadmin/utils"
)

type PlayerChargeInfoRecord struct {
	PlayerId        int    `json:"playerId" gorm:"primary_key"`
	PlayerName      string `json:"playerName" gorm:"-"`
	Account         string `json:"account" gorm:"-"`
	PlatformId      int    `json:"platformId" gorm:"column:part_id"`
	ServerId        string `json:"serverId"`
	TotalMoney      int    `json:"totalMoney"`
	MaxMoney        int    `json:"maxMoney"`
	MinMoney        int    `json:"minMoney"`
	ChargeCount     int    `json:"chargeCount"`
	LastLoginTime   int    `json:"lastLoginTime" gorm:"-"`
	RegisterTime    int    `json:"registerTime" gorm:"-"`
	LastChargeTime  int    `json:"lastChargeTime" gorm:"column:last_time"`
	FirstChargeTime int    `json:"firstChargeTime" gorm:"column:first_time"`
}

type PlayerChargeDataQueryParam struct {
	BaseQueryParam
	PlatformId int
	Node       string `json:"serverId"`
}

func GetPlayerChargeDataOne(playerId int) (*PlayerChargeInfoRecord, error) {
	playerChargeInfo := &PlayerChargeInfoRecord{
		PlayerId: playerId,
	}
	err := DbCharge.FirstOrInit(&playerChargeInfo).Error
	return playerChargeInfo, err
}


func GetPlayerChargeDataList(params *PlayerChargeDataQueryParam) ([]*PlayerChargeInfoRecord, int64) {
	data := make([]*PlayerChargeInfoRecord, 0)
	var count int64
	sortOrder := "total_money desc"
	if params.Node == "" {
		DbCharge.Model(&PlayerChargeInfoRecord{}).Where(&PlayerChargeInfoRecord{
			PlatformId: params.PlatformId,
		}).Where("charge_count > 0 ").Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	} else {
		DbCharge.Model(&PlayerChargeInfoRecord{}).Where(&PlayerChargeInfoRecord{
			PlatformId: params.PlatformId,
		}).Where("charge_count > 0 ").Where("server_id in (?)", GetGameServerIdListByNode(params.Node)).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	}

	for _, e := range data {
		gameDb, err := GetGameDbByPlatformIdAndSid(e.PlatformId, e.ServerId)
		utils.CheckError(err)
		if err != nil {
			continue
		}
		defer gameDb.Close()
		player, err := GetPlayerByDb(gameDb, e.PlayerId)
		utils.CheckError(err)
		if err != nil {
			continue
		}
		e.PlayerName = player.ServerId + "." + player.Nickname
		e.Account = player.AccId
		e.LastLoginTime = player.LastLoginTime
		e.RegisterTime = player.RegTime
	}
	return data, count
}
