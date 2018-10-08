package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	"fmt"
)

type DailyStatistics struct {
	Node       string `json:"node" gorm:"primary_key"`
	PlatformId string `json:"platformId" gorm:"primary_key"`
	ServerId   string `json:"serverId" gorm:"primary_key"`
	Channel    string `json:"channel" gorm:"primary_key"`
	Time       int    `json:"time" gorm:"primary_key"`

	ChargeMoney            int     `json:"chargeMoney"`
	NewChargeMoney         int     `json:"newChargeMoney"`
	TotalChargeMoney       int     `json:"totalChargeMoney"`
	ChargePlayerCount      int     `json:"chargePlayerCount"`
	TotalChargePlayerCount int     `json:"totalChargePlayerCount"`
	ARPU                   float32 `json:"arpu" gorm:"-"`
	ActiveARPU             float32 `json:"active_arpu" gorm:"-"`
	NewChargePlayerCount   int     `json:"newChargePlayerCount"`
	//ActivePlayerCount    int     `json:"activePlayerCount" gorm:"-"`
	ActiveChargeRate float32 `json:"activeChargeRate" gorm:"-"`

	LoginTimes           int `json:"loginTimes"`
	LoginPlayerCount     int `json:"loginPlayerCount"`
	ActivePlayerCount    int `json:"activePlayerCount"`
	CreateRoleCount      int `json:"createRoleCount"`
	ShareCreateRoleCount int `json:"shareCreateRoleCount"`
	TotalCreateRoleCount int `json:"totalCreateRoleCount"`

	MaxOnlineCount int `json:"maxOnline"`
	MinOnlineCount int `json:"minOnline"`
	AvgOnlineCount int `json:"avgOnline"`
	AvgOnlineTime  int `json:"avgOnlineTime"`

	RegisterCount      int `json:"registerCount"`
	TotalRegisterCount int `json:"totalRegisterCount"`
	//CreateRoleCount int `json:"createRoleCount"`
	ValidRoleCount int `json:"validRoleCount"`
}

type DailyStatisticsQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string    `json:"serverId"`
	ChannelList [] string `json:"channelList"`
	StartTime   int
	EndTime     int
}

func GetDailyStatisticsList(params *DailyStatisticsQueryParam) []*DailyStatistics {
	if params.EndTime < params.StartTime {
		logs.Error("开始结束时间错误")
		return nil
	}
	data := make([]*DailyStatistics, 0, (params.EndTime-params.StartTime)/86400)
	channelLen := len(params.ChannelList)
	for i := params.StartTime; i <= params.EndTime; i = i + 86400 {
		tmpData := make([]*DailyStatistics, 0, channelLen)
		err := Db.Model(&DailyStatistics{}).Where(&DailyStatistics{PlatformId: params.PlatformId, ServerId: params.ServerId, Time: i}).Where("channel in (?)", params.ChannelList).Find(&tmpData).Error
		utils.CheckError(err)
		if len(tmpData) > 0 {
			tmpE := &DailyStatistics{
				PlatformId: params.PlatformId,
				ServerId:   params.ServerId,
				Time:       i,
			}
			for _, e := range tmpData {
				tmpE.ChargeMoney += e.ChargeMoney
				tmpE.TotalChargeMoney += e.TotalChargeMoney
				tmpE.NewChargeMoney += e.NewChargeMoney
				tmpE.ChargePlayerCount += e.ChargePlayerCount
				tmpE.TotalChargePlayerCount += e.TotalChargePlayerCount
				tmpE.NewChargePlayerCount += e.NewChargePlayerCount

				tmpE.LoginTimes += e.LoginTimes
				tmpE.LoginPlayerCount += e.LoginPlayerCount
				tmpE.ActivePlayerCount += e.ActivePlayerCount

				tmpE.MaxOnlineCount += e.MaxOnlineCount
				tmpE.MinOnlineCount += e.MinOnlineCount
				tmpE.AvgOnlineCount += e.AvgOnlineCount
				tmpE.AvgOnlineTime += e.AvgOnlineTime

				tmpE.RegisterCount += e.RegisterCount
				tmpE.CreateRoleCount += e.CreateRoleCount
				tmpE.ValidRoleCount += e.ValidRoleCount
				tmpE.TotalCreateRoleCount += e.TotalCreateRoleCount
				tmpE.ShareCreateRoleCount += e.ShareCreateRoleCount
				tmpE.TotalRegisterCount += e.TotalRegisterCount
			}
			data = append(data, tmpE)
		}
	}
	for _, e := range data {
		e.ARPU = CaclRate(e.ChargeMoney, e.ChargePlayerCount)
		e.ActiveARPU = CaclRate(e.ChargeMoney, e.LoginPlayerCount)
		e.ActiveChargeRate = CaclRate(e.ChargePlayerCount, e.ActivePlayerCount)
	}
	return data
}
func GetDailyStatisticsOne(platformId string, serverId string, channel string, timestamp int) (*DailyStatistics, error) {
	data := &DailyStatistics{}
	err := Db.Model(&DailyStatistics{}).Where(&DailyStatistics{
		PlatformId: platformId,
		ServerId:   serverId,
		Channel:    channel,
		Time:       timestamp,
	}).First(&data).Error
	return data, err
}
func UpdateDailyStatistics(platformId string, serverId string, channel string, timestamp int) error {
	logs.Info("UpdateDailyStatistics:%v, %v, %v, %v", platformId, serverId, channel, timestamp)
	serverNode, err := GetGameServerOne(platformId, serverId)
	if err != nil {
		return err
	}
	node := serverNode.Node
	gameDb, err := GetGameDbByNode2(serverNode.Node, platformId, serverId)
	if err != nil {
		return err
	}
	defer gameDb.Close()

	createRoleCount := GetCreateRoleCount(gameDb, serverId, channel, timestamp, timestamp+86400)
	registerRoleCount := GetRegisterRoleCount(gameDb, serverId, channel, timestamp, timestamp+86400)
	totalChargeMoney := GetTotalChargeMoney(platformId, serverId, channel, 0, timestamp+86400)
	chargeMoney := GetTotalChargeMoney(platformId, serverId, channel, timestamp, timestamp+86400)
	//logs.Info("%d: %d | %d", timestamp, chargeMoney, totalChargeMoney)
	m := &DailyStatistics{
		Node:                   serverNode.Node,
		PlatformId:             platformId,
		ServerId:               serverId,
		Channel:                channel,
		Time:                   timestamp,
		TotalChargeMoney:       totalChargeMoney,
		ChargeMoney:            chargeMoney,
		NewChargeMoney:         GetThatDayNewChargeMoney(platformId, serverId, channel, timestamp),
		ChargePlayerCount:      GetThatDayServerChargePlayerCount(platformId, serverId, channel, timestamp),
		TotalChargePlayerCount: GetThatDayChargePlayerCount(platformId, serverId, channel, timestamp+86400),
		NewChargePlayerCount:   GetThadDayServerFirstChargePlayerCount(platformId, serverId, channel, timestamp),

		LoginTimes:        GetThatDayLoginTimes(gameDb, serverId, channel, timestamp),
		LoginPlayerCount:  GetThatDayLoginPlayerCount(gameDb, serverId, channel, timestamp),
		ActivePlayerCount: GetThatDayActivePlayerCount(gameDb, serverId, channel, timestamp),

		TotalCreateRoleCount: GetHistoryCreateRoleCount(platformId, serverId, channel, timestamp) + createRoleCount,
		TotalRegisterCount:   GetHistoryRegisterRoleCount(platformId, serverId, channel, timestamp) + registerRoleCount,
		//TotalCreateRoleCount: GetCreateRoleCount(gameDb, serverId, channel, 0, timestamp+86400),
		//TotalRegisterCount: GetRegisterRoleCount(gameDb, serverId, channel, 0, timestamp+86400),

		MaxOnlineCount: GetThatDayMaxOnlineCount(platformId, serverId, [] string{channel}, timestamp, timestamp+86400),
		MinOnlineCount: GetThatDayMinOnlineCount(platformId, serverId, [] string{channel}, timestamp, timestamp+86400),
		AvgOnlineCount: GetThatDayAvgOnlineCount(platformId, serverId, [] string{channel}, timestamp, timestamp+86400),
		AvgOnlineTime:  GetOnlineTime(node, serverId, channel, timestamp),


		RegisterCount:        registerRoleCount,
		CreateRoleCount:      createRoleCount,
		ShareCreateRoleCount: GetThatDayShareCreateRoleCountByChannel(gameDb, serverId, channel, timestamp),
		ValidRoleCount:       GetThatDayValidCreateRoleCountByChannel(gameDb, serverId, channel, timestamp),
	}
	err = Db.Save(&m).Error
	return err
}

func GetHistoryCreateRoleCount(platformId string, serverId string, channel string, time int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT sum(create_role_count) as count FROM daily_statistics WHERE platform_id = '%s' and server_id = '%s' and channel = '%s' and time < %d `, platformId, serverId, channel, time)
	err := Db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}
func GetHistoryRegisterRoleCount(platformId string, serverId string, channel string, time int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT sum(register_count) as count FROM daily_statistics WHERE platform_id = '%s' and server_id = '%s' and channel = '%s' and time < %d `, platformId, serverId, channel, time)
	err := Db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//func UpdateDailyStatistics(node string, timestamp int) error {
//	logs.Info("更新DailyStatistics:%v, %v", node, timestamp)
//	m := &DailyStatistics{
//		Node:                 node,
//		Time:                 timestamp,
//		ChargeMoney:          GetThatDayServerTotalChargeMoney(node, timestamp),
//		ChargePlayerCount:    GetThatDayServerChargePlayerCount(node, timestamp),
//		NewChargePlayerCount: GetThadDayServerFirstChargePlayerCount(node, timestamp),
//	}
//	err := Db.Save(&m).Error
//	return err
//}
//func RepireCharge() {
//	logs.Info("修复每日充值 platform_id, server_id")
//	data := make([]*DailyStatistics, 0)
//	err := Db.Model(&DailyStatistics{}).Find(&data).Error
//	utils.CheckError(err)
//	for _, e := range data {
//
//		nodeName := strings.Split(e.Node, "@")[0]
//
//		platform := "qq"
//		serverId := ""
//		if strings.Contains(nodeName, "_") {
//			platform = strings.Split(nodeName, "_")[0]
//			serverId = strings.Split(nodeName, "_")[1]
//		} else {
//			serverId = nodeName
//		}
//		//e.PlatformId = platform
//		//e.ServerId = serverId
//		//err = Db.Save(&e).Error
//		sql := fmt.Sprintf("update daily_charge_statistics set platform_id = '%s', server_id = '%s', channel = '%s' where node = '%s';", platform, serverId, platform, e.Node)
//		logs.Info("sql:%s\n", sql)
//		err = Db.Debug().Exec(sql).Error
//		//err =Db.Model(&e).Updates(map[string]interface{}{"platform_id": platform, "server_id": serverId, "channel": platform}).Error
//		utils.CheckError(err)
//	}
//	logs.Info("修复每日充值完毕")
//}

//func RepireRemainTotal() {
//	logs.Info("修复 remain_total")
//	data := make([]*RemainTotal, 0)
//	err := Db.Model(&RemainTotal{}).Find(&data).Error
//	utils.CheckError(err)
//	for _, e := range data {
//
//		nodeName := strings.Split(e.Node, "@")[0]
//
//		platform := "qq"
//		serverId := ""
//		if strings.Contains(nodeName, "_") {
//			platform = strings.Split(nodeName, "_")[0]
//			serverId = strings.Split(nodeName, "_")[1]
//		} else {
//			serverId = nodeName
//		}
//		//e.PlatformId = platform
//		//e.ServerId = serverId
//		//err = Db.Save(&e).Error
//		sql := fmt.Sprintf("update remain_total set platform_id = '%s', server_id = '%s', channel = '%s' where node = '%s';", platform, serverId, platform, e.Node)
//		logs.Info("sql:%s\n", sql)
//		err = Db.Debug().Exec(sql).Error
//		//err =Db.Model(&e).Updates(map[string]interface{}{"platform_id": platform, "server_id": serverId, "channel": platform}).Error
//		utils.CheckError(err)
//	}
//	logs.Info("修复remain_total完毕")
//}

//func RepireTenMinuteStatistics() {
//	logs.Info("修复 RepireTenMinuteStatistics  2222222")
//	data := make([]*TenMinuteStatistics, 0)
//	err := Db.Model(&TenMinuteStatistics{}).Find(&data).Error
//	utils.CheckError(err)
//	for _, e := range data {
//
//		nodeName := strings.Split(e.Node, "@")[0]
//
//		platform := "qq"
//		serverId := ""
//		if strings.Contains(nodeName, "_") {
//			platform = strings.Split(nodeName, "_")[0]
//			serverId = strings.Split(nodeName, "_")[1]
//		} else {
//			serverId = nodeName
//		}
//		//e.PlatformId = platform
//		//e.ServerId = serverId
//		//err = Db.Save(&e).Error
//		if e.PlatformId == "" {
//			//go func() {
//				sql := fmt.Sprintf("update ten_minute_statistics set platform_id = '%s', server_id = '%s' where node = '%s';", platform, serverId, e.Node)
//				//logs.Info("sql:%s\n", sql)
//				err = Db.Exec(sql).Error
//				//err =Db.Model(&e).Updates(map[string]interface{}{"platform_id": platform, "server_id": serverId, "channel": platform}).Error
//				utils.CheckError(err)
//			//}()
//		}
//
//	}
//	logs.Info("修复RepireTenMinuteStatistics完毕")
//}

//func RepireRemainActive() {
//
//	logs.Info("修复 remain_active")
//	data := make([]*RemainActive, 0)
//	err := Db.Model(&RemainActive{}).Find(&data).Error
//	utils.CheckError(err)
//	for _, e := range data {
//
//		nodeName := strings.Split(e.Node, "@")[0]
//
//		platform := "qq"
//		serverId := ""
//		if strings.Contains(nodeName, "_") {
//			platform = strings.Split(nodeName, "_")[0]
//			serverId = strings.Split(nodeName, "_")[1]
//		} else {
//			serverId = nodeName
//		}
//		//e.PlatformId = platform
//		//e.ServerId = serverId
//		//err = Db.Save(&e).Error
//		sql := fmt.Sprintf("update remain_active set platform_id = '%s', server_id = '%s', channel = '%s' where node = '%s';", platform, serverId, platform, e.Node)
//		logs.Info("sql:%s\n", sql)
//		err = Db.Debug().Exec(sql).Error
//		//err =Db.Model(&e).Updates(map[string]interface{}{"platform_id": platform, "server_id": serverId, "channel": platform}).Error
//		utils.CheckError(err)
//	}
//	logs.Info("修复remain_active完毕")
//}
