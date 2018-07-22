package models

import (
	"github.com/chnzrb/myadmin/utils"
	"github.com/jinzhu/gorm"
	"fmt"
	//"os/user"
	"github.com/astaxie/beego/logs"
	"strings"
)

type ChargeInfoRecord struct {
	OrderId       string `json:"orderId" gorm:"primary_key"`
	ChargeType    int    `json:"chargeType"`
	Ip            string `json:"ip"`
	PartId        string `json:"platformId"`
	ServerId      string `json:"serverId"`
	AccId         string `json:"accId"`
	IsFirst       int    `json:"isFirst"`
	CurrLevel     int    `json:"currLevel"`
	CurrTaskId    int    `json:"currTaskId"`
	RegTime       int    `json:"regTime"`
	FirstTime     int    `json:"firstTime"`
	CurrPower     int    `json:"currPower"`
	PlayerId      int    `json:"playerId"`
	PlayerName    string `json:"playerName" gorm:"-"`
	LastLoginTime int    `json:"lastLoginTime" gorm:"-"`
	Money         int    `json:"money"`
	Ingot         int    `json:"ingot"`
	RecordTime    int    `json:"recordTime"`
	ChargeItemId  int    `json:"chargeItemId"`
}

type ChargeInfoRecordQueryParam struct {
	BaseQueryParam
	PlatformId string
	Node       string `json:"serverId"`
	PlayerId   int
	PlayerName string
	OrderId    string
	AccId      string
	StartTime  int
	EndTime    int
}

func GetChargeInfoRecordList(params *ChargeInfoRecordQueryParam) ([]*ChargeInfoRecord, int64, int64, int64) {
	data := make([]*ChargeInfoRecord, 0)
	var count int64
	sortOrder := "record_time"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	} else if params.Order == "ascending" {
		sortOrder = sortOrder + " asc"
	} else {
		sortOrder = sortOrder + " desc"
	}
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("record_time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}

	var sumData struct {
		MoneyCount int64
		PlayerCount int64
	}
	whereArray := make([] string, 0)
	whereArray = append(whereArray, " charge_type = 99 ")
	if params.Node != "" {
		whereArray = append(whereArray, fmt.Sprintf(" server_id in (%s) ", GetGameServerIdListStringByNode(params.Node)))
	}
	if params.StartTime > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" record_time between %d and %d ", params.StartTime, params.EndTime))
	}
	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}
	sql := fmt.Sprintf(
		`select sum(money) as money_count, count(DISTINCT player_id) as player_count  from charge_info_record  %s;`, whereParam)
	err := DbCharge.Debug().Raw(sql).Scan(&sumData).Error
	utils.CheckError(err)

	if params.Node == "" {
		err := f(DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
			PartId:     params.PlatformId,
			AccId:      params.AccId,
			OrderId:    params.OrderId,
			PlayerId:   params.PlayerId,
			ChargeType: 99,
		})).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Error
		utils.CheckError(err)
	} else {
		err := f(DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
			PartId:     params.PlatformId,
			AccId:      params.AccId,
			OrderId:    params.OrderId,
			PlayerId:   params.PlayerId,
			ChargeType: 99,
		})).Where("server_id in(?)", GetGameServerIdListByNode(params.Node)).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Error
		utils.CheckError(err)
	}


	for _, e := range data {
		e.PlayerName = GetPlayerName_2(e.PartId, e.ServerId, e.PlayerId)
		e.LastLoginTime = GetPlayerLastLoginTime(e.PartId, e.ServerId, e.PlayerId)
		//e.ChargeItemId = GetChargeItemId(e.OrderId, e.PartId, e.ServerId)
	}
	return data, count, sumData.PlayerCount, sumData.MoneyCount
}

func Repair() {
	logs.Info("开始修复充值数据")
	data := make([]*ChargeInfoRecord, 0)
	//var count int64
	err := DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
		ChargeType: 99,
	}).Find(&data).Error
	utils.CheckError(err)
	if err != nil {
		return
	}
	for _, e := range data {
		value := GetChargeItemId(e.OrderId, e.PartId, e.ServerId)
		logs.Debug("value:%v, %v", e.OrderId, value)
		if value > 0 {
			err = DbCharge.Model(&e).Update("charge_item_id", value).Error
			utils.CheckError(err)
			if err != nil {
				return
			}
		}
	}
	logs.Info("修复充值数据成功")
}

//获取玩家最近登录时间
func GetChargeItemId(orderId string, platformId string, serverId string) int {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	defer gameDb.Close()
	var data struct {
		ChargeItemId int
	}
	sql := fmt.Sprintf(
		`SELECT charge_item_id FROM player_charge_record where order_id = ? `)
	err = gameDb.Raw(sql, orderId).Scan(&data).Error
	utils.CheckError(err)
	return data.ChargeItemId
}
