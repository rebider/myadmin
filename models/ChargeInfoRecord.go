package models

import (
	"github.com/chnzrb/myadmin/utils"
	//"github.com/jinzhu/gorm"
	"fmt"
	//"os/user"
	//"github.com/astaxie/beego/logs"
	"strings"
	//"github.com/zaaksam/dproxy/go/db"
	"github.com/astaxie/beego/logs"
)

type ChargeInfoRecord struct {
	OrderId       string  `json:"orderId" gorm:"primary_key"`
	ChargeType    int     `json:"chargeType"`
	Ip            string  `json:"ip"`
	PartId        string  `json:"platformId"`
	ServerId      string  `json:"serverId"`
	AccId         string  `json:"accId"`
	IsFirst       int     `json:"isFirst"`
	CurrLevel     int     `json:"currLevel"`
	CurrTaskId    int     `json:"currTaskId"`
	RegTime       int     `json:"regTime"`
	FirstTime     int     `json:"firstTime"`
	CurrPower     int     `json:"currPower"`
	PlayerId      int     `json:"playerId"`
	PlayerName    string  `json:"playerName" gorm:"-"`
	LastLoginTime int     `json:"lastLoginTime" gorm:"-"`
	Money         float32 `json:"money"`
	Ingot         int     `json:"ingot"`
	RecordTime    int     `json:"recordTime"`
	ChargeItemId  int     `json:"chargeItemId"`
	Channel       string  `json:"channel"`
}

type ChargeInfoRecordQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string    `json:"serverId"`
	ChannelList [] string `json:"channelList"`
	PlayerId    int
	PlayerName  string
	OrderId     string
	AccId       string
	StartTime   int
	EndTime     int
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
	//f := func(db *gorm.DB) *gorm.DB {
	//	if params.StartTime > 0 {
	//		return db.Where("record_time between ? and ?", params.StartTime, params.EndTime)
	//	}
	//	return db
	//}

	var sumData struct {
		MoneyCount  float32
		PlayerCount int64
	}
	whereArray := make([] string, 0)
	whereArray = append(whereArray, " charge_type = 99 ")
	whereArray = append(whereArray, " part_id =  '"+params.PlatformId+"' ")
	if params.ServerId != "" {
		whereArray = append(whereArray, " server_id =  '"+params.ServerId+"' ")
		//whereArray = append(whereArray, fmt.Sprintf(" server_id in (%s) ", GetGameServerIdListStringByNode(params.Node)))
	}
	if len(params.ChannelList) > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" channel in (%s) ", GetSQLWhereParam(params.ChannelList)))
	}
	//if params.StartTime > 0 {
	//	whereArray = append(whereArray, fmt.Sprintf(" record_time between %d and %d ", params.StartTime, params.EndTime))
	//}
	if params.StartTime > 0 {
		whereArray = append(whereArray, fmt.Sprintf("record_time between %d and %d", params.StartTime, params.EndTime))
	}

	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}
	sql := fmt.Sprintf(
		`select sum(money) as money_count, count(DISTINCT player_id) as player_count  from charge_info_record  %s;`, whereParam)
	logs.Debug("sql:", sql)
	err := DbCharge.Raw(sql).Scan(&sumData).Error
	utils.CheckError(err)

	//if params.Node == "" {
	if params.StartTime > 0 {
		err = DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
			PartId:     params.PlatformId,
			AccId:      params.AccId,
			OrderId:    params.OrderId,
			PlayerId:   params.PlayerId,
			ServerId:   params.ServerId,
			ChargeType: 99,
		}).Where("record_time between ? and ? ", params.StartTime, params.EndTime).Where("channel in(?)", params.ChannelList).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Error
		utils.CheckError(err)
	} else {
		err = DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
			PartId:     params.PlatformId,
			AccId:      params.AccId,
			OrderId:    params.OrderId,
			PlayerId:   params.PlayerId,
			ServerId:   params.ServerId,
			ChargeType: 99,
		}).Where("channel in(?)", params.ChannelList).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Error
		utils.CheckError(err)
	}

	//} else {
	//	err := f(DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
	//		PartId:     params.PlatformId,
	//		AccId:      params.AccId,
	//		OrderId:    params.OrderId,
	//		PlayerId:   params.PlayerId,
	//		ChargeType: 99,
	//	})).Where("server_id in(?)", GetGameServerIdListByNode(params.Node)).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Error
	//	utils.CheckError(err)
	//}

	for _, e := range data {
		e.PlayerName = GetPlayerName_2(e.PartId, e.ServerId, e.PlayerId)
		e.LastLoginTime = GetPlayerLastLoginTime(e.PartId, e.ServerId, e.PlayerId)
		//e.ChargeItemId = GetChargeItemId(e.OrderId, e.PartId, e.ServerId)
	}
	return data, count, sumData.PlayerCount, int64(sumData.MoneyCount)
}

//func Repair() {
//	logs.Info("开始修复充值数据")
//	data := make([]*ChargeInfoRecord, 0)
//	//var count int64
//	err := DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
//		ChargeType: 99,
//	}).Find(&data).Error
//	utils.CheckError(err)
//	if err != nil {
//		return
//	}
//	for _, e := range data {
//		value := GetChargeItemId(e.OrderId, e.PartId, e.ServerId)
//		logs.Debug("value:%v, %v", e.OrderId, value)
//		if value > 0 {
//			err = DbCharge.Model(&e).Update("charge_item_id", value).Error
//			utils.CheckError(err)
//			if err != nil {
//				return
//			}
//		}
//	}
//	logs.Info("修复充值数据成功")
//}
//
////获取玩家最近登录时间
//func GetChargeItemId(orderId string, platformId string, serverId string) int {
//	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
//	utils.CheckError(err)
//	if err != nil {
//		return 0
//	}
//	defer gameDb.Close()
//	var data struct {
//		ChargeItemId int
//	}
//	sql := fmt.Sprintf(
//		`SELECT charge_item_id FROM player_charge_record where order_id = %d `, orderId)
//	err = gameDb.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.ChargeItemId
//}

type ChargeStatistics struct {
	PlatformId string `json:"platformId"`
	//ServerId                    string    `json:"serverId"`
	TodayCharge                      int `json:"todayCharge"`
	TodayChargePlayerCount           int `json:"todayChargePlayerCount"`
	YesterdayCharge                  int `json:"yesterdayCharge"`
	YesterdayChargePlayerCount       int `json:"yesterdayChargePlayerCount"`
	BeforeYesterdayCharge            int `json:"beforeYesterdayCharge"`
	BeforeYesterdayChargePlayerCount int `json:"beforeYesterdayChargePlayerCount"`

	ChargeData            [] map[string]string `json:"chargeData"`
	ChargePlayerCountData [] map[string]string `json:"chargePlayerCountData"`

	TodayChargeList           [] string `json:"todayChargeList"`
	YesterdayChargeList       [] string `json:"yesterdayChargeList"`
	BeforeYesterdayChargeList [] string `json:"beforeYesterdayChargeList"`

	TodayChargePlayerCountList           [] string `json:"todayChargePlayerCountList"`
	YesterdayChargePlayerCountList       [] string `json:"yesterdayChargePlayerCountList"`
	BeforeYesterdayChargePlayerCountList [] string `json:"beforeYesterdayChargePlayerCountList"`
}

func GetChargeStatistics(platformId string, serverId string, channelList [] string) (*ChargeStatistics, error) {

	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	yesterdayZeroTimestamp := todayZeroTimestamp - 86400
	beforeYesterdayZeroTimestamp := yesterdayZeroTimestamp - 86400

	todayOnlineList, todayTotalCharge := get24hoursChargeCount(platformId, serverId, channelList, todayZeroTimestamp)
	yesterdayOnlineList, yesterdayTotalCharge := get24hoursChargeCount(platformId, serverId, channelList, yesterdayZeroTimestamp)
	beforeYesterdayOnlineList, beforeYesterdayTotalCharge := get24hoursChargeCount(platformId, serverId, channelList, beforeYesterdayZeroTimestamp)

	todayChargePlayerCountList, todayChargePlayerCount := get24hoursChargePlayerCount(platformId, serverId, channelList, todayZeroTimestamp)
	yesterdayChargePlayerCountList, yesterdayChargePlayerCount := get24hoursChargePlayerCount(platformId, serverId, channelList, yesterdayZeroTimestamp)
	beforeYesterdayChargePlayerCountList, beforeYesterdayChargePlayerCount := get24hoursChargePlayerCount(platformId, serverId, channelList, beforeYesterdayZeroTimestamp)

	chargeData := make([]map[string]string, 0, 144)
	//logs.Info("len:%d", len(todayOnlineList))
	for i := 0; i < 6*24; i = i + 1 {
		m := make(map[string]string, 4)
		m["时间"] = utils.FormatTime(i * 10 * 60)
		m["今日充值"] = todayOnlineList[i]
		m["昨日充值"] = yesterdayOnlineList[i]
		m["前日充值"] = beforeYesterdayOnlineList[i]
		//logs.Info(i)
		chargeData = append(chargeData, m)
	}

	chargePlayerCountData := make([]map[string]string, 0, 144)
	for i := 0; i < 6*24; i = i + 1 {
		m := make(map[string]string, 4)
		m["时间"] = utils.FormatTime(i * 10 * 60)
		m["今日充值人数"] = todayChargePlayerCountList[i]
		m["昨日充值人数"] = yesterdayChargePlayerCountList[i]
		m["前日充值人数"] = beforeYesterdayChargePlayerCountList[i]
		chargePlayerCountData = append(chargePlayerCountData, m)
	}

	chargeStatistics := &ChargeStatistics{
		PlatformId:             platformId,
		TodayCharge:            todayTotalCharge,
		TodayChargePlayerCount: todayChargePlayerCount,
		//TodayCreateRole: GetCreateRoleCountByChannelList(gameDb, serverId, channelList, todayZeroTimestamp, todayZeroTimestamp+86400),
		YesterdayCharge:            yesterdayTotalCharge,
		YesterdayChargePlayerCount: yesterdayChargePlayerCount,
		//MaxOnlineCount:              GetMaxOnlineCount(node),
		BeforeYesterdayCharge:            beforeYesterdayTotalCharge,
		BeforeYesterdayChargePlayerCount: beforeYesterdayChargePlayerCount,

		TodayChargeList:           todayOnlineList,
		YesterdayChargeList:       yesterdayOnlineList,
		BeforeYesterdayChargeList: beforeYesterdayOnlineList,

		TodayChargePlayerCountList:           todayChargePlayerCountList,
		YesterdayChargePlayerCountList:       yesterdayChargePlayerCountList,
		BeforeYesterdayChargePlayerCountList: beforeYesterdayChargePlayerCountList,
		ChargeData:                           chargeData,
		ChargePlayerCountData:                chargePlayerCountData,
	}
	return chargeStatistics, nil
}
