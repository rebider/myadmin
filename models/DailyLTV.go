package models

import (
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
)

type DailyLTV struct {
	//Node         string  `json:"node" gorm:"primary_key"`
	PlatformId   string  `json:"platformId" gorm:"primary_key"`
	ServerId     string  `json:"serverId" gorm:"primary_key"`
	Channel      string  `json:"channel" gorm:"primary_key"`
	Time         int     `json:"time" gorm:"primary_key"`
	RegisterRole int     `json:"registerRole" gorm:"-"`
	CreateRole   int     `json:"createRole" gorm:"-"`
	C1           int     `json:"c1"  gorm:"column:c1"`
	C2           int     `json:"c2" gorm:"column:c2"`
	C3           int     `json:"c3" gorm:"column:c3"`
	C7           int     `json:"c7" gorm:"column:c7"`
	C14          int     `json:"c14" gorm:"column:c14"`
	C30          int     `json:"c30" gorm:"column:c30"`
	C60          int     `json:"c60" gorm:"column:c60"`
	C90          int     `json:"c90" gorm:"column:c90"`
	C120         int     `json:"c120" gorm:"column:c120"`
	LTV1         float32 `json:"ltv1"  gorm:"-"`
	LTV2         float32 `json:"ltv2" gorm:"-"`
	LTV3         float32 `json:"ltv3" gorm:"-"`
	LTV7         float32 `json:"ltv7" gorm:"-"`
	LTV14        float32 `json:"ltv14" gorm:"-"`
	LTV30        float32 `json:"ltv30" gorm:"-"`
	LTV60        float32 `json:"ltv60" gorm:"-"`
	LTV90        float32 `json:"ltv90" gorm:"-"`
	LTV120       float32 `json:"ltv120" gorm:"-"`
}

type DailyLTVQueryParam struct {
	BaseQueryParam
	PlatformId string
	//ServerId   string
	ServerId    string    `json:"serverId"`
	ChannelList [] string `json:"channelList"`
	Channel     string
	StartTime   int
	EndTime     int
}

// 获取总体留存
func GetDailyLTVList(params DailyLTVQueryParam) ([]*DailyLTV) {
	if params.EndTime < params.StartTime {
		logs.Error("开始结束时间错误")
		return nil
	}
	data := make([]*DailyLTV, 0, (params.EndTime-params.StartTime)/86400)
	channelLen := len(params.ChannelList)

	for i := params.StartTime; i <= params.EndTime; i = i + 86400 {
		tmpData := make([]*DailyLTV, 0, channelLen)
		err := Db.Model(&DailyLTV{}).Where(&DailyLTV{PlatformId: params.PlatformId, ServerId: params.ServerId, Time: i}).Where("channel in (?)", params.ChannelList).Find(&tmpData).Error
		//logs.Debug("ok")
		utils.CheckError(err)
		if len(tmpData) > 0 {
			tmpE := &DailyLTV{
				PlatformId: params.PlatformId,
				ServerId:   params.ServerId,
				Time:       i,
			}
			//logs.Debug("1")
			for _, e := range tmpData {
				dailyStatistics, err := GetDailyStatisticsOne(e.PlatformId, e.ServerId, e.Channel, e.Time)
				//logs.Debug("11")
				utils.CheckError(err)
				tmpE.RegisterRole += dailyStatistics.RegisterCount
				tmpE.CreateRole += dailyStatistics.CreateRoleCount
				if e.C1 > 0 {
					tmpE.C1 += e.C1
				}
				if e.C2 > 0 {
					tmpE.C2 += e.C2
				}
				if e.C3 > 0 {
					tmpE.C3 += e.C3
				}
				if e.C7 > 0 {
					tmpE.C7 += e.C7
				}
				if e.C14 > 0 {
					tmpE.C14 += e.C14
				}
				if e.C30 > 0 {
					tmpE.C30 += e.C30
				}
				if e.C60 > 0 {
					tmpE.C60 += e.C60
				}
				if e.C90 > 0 {
					tmpE.C90 += e.C90
				}
				if e.C120 > 0 {
					tmpE.C120 += e.C120
				}
			}
			//logs.Debug("2")
			data = append(data, tmpE)
		}
	}
	for _, e := range data {
		if e.CreateRole > 0 {
			e.LTV1 = float32(e.C1) / float32(e.CreateRole)
			e.LTV2 = float32(e.C2) / float32(e.CreateRole)
			e.LTV3 = float32(e.C3) / float32(e.CreateRole)
			e.LTV7 = float32(e.C7) / float32(e.CreateRole)
			e.LTV14 = float32(e.C14) / float32(e.CreateRole)
			e.LTV30 = float32(e.C30) / float32(e.CreateRole)
			e.LTV60 = float32(e.C60) / float32(e.CreateRole)
			e.LTV90 = float32(e.C90) / float32(e.CreateRole)
			e.LTV120 = float32(e.C120) / float32(e.CreateRole)

		}
	}
	return data
	//data := make([]*DailyLTV, 0)
	//var count int64
	//f := func(db *gorm.DB) *gorm.DB {
	//	if params.StartTime > 0 {
	//		return db.Where("time between ? and ?", params.StartTime, params.EndTime)
	//	}
	//	return db
	//}
	//err := f(Db.Model(&DailyLTV{}).Where(&DailyLTV{PlatformId: params.PlatformId, ServerId: params.ServerId, Channel: params.Channel})).Where(" channel in (?)", params.ChannelList).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data).Error
	//utils.CheckError(err)
	//for _, e := range data {
	//	dailyStatistics, err := GetDailyStatisticsOne(e.PlatformId, e.ServerId, e.Channel, e.Time)
	//	//logs.Info("dailyStatistics:%+v", dailyStatistics)
	//	utils.CheckError(err)
	//	e.CreateRole = dailyStatistics.CreateRoleCount
	//	e.RegisterRole = dailyStatistics.RegisterCount
	//}
	//return data, count
}

//更新 每日ltv
func UpdateDailyLTV(platformId string, serverId string, channelList [] * Channel, timestamp int) error {
	logs.Info("更新每日ltv:%v, %v, %v, %v", platformId, serverId, len(channelList), timestamp)
	gameServer, err := GetGameServerOne(platformId, serverId)
	if err != nil {
		return err
	}
	node := gameServer.Node
	serverNode, err := GetServerNode(node)
	if err != nil {
		return err
	}
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	defer gameDb.Close()
	openDayZeroTimestamp := utils.GetThatZeroTimestamp(int64(serverNode.OpenTime))

	for _, e := range channelList {
		channel := e.Channel
		for i := 1; i <= 120; i++ {
			if i == 1 || i == 2 || i == 3 || i == 7 || i == 14 || i == 30 || i == 60 || i == 90 || i == 120 {

			} else {
				continue
			}
			thatDayZeroTimestamp := timestamp - i*86400
			if openDayZeroTimestamp > thatDayZeroTimestamp {
				continue
			}

			//dailyStatistics, err := GetDailyStatisticsOne(platformId, serverId, channel, thatDayZeroTimestamp)

			//registerNum := dailyStatistics.RegisterCount
			//createNum := dailyStatistics.CreateRoleCount
			totalCharge := GetTotalChargeMoneyByRegisterTime(platformId, serverId, channel, 0, timestamp, thatDayZeroTimestamp)
			if totalCharge > 0 {
				rate := totalCharge
				//if createNum > 0 {
				//	rate = int(float32(totalCharge) / float32(createNum))
				//}

				m := &DailyLTV{
					//Node:       node,
					PlatformId: platformId,
					ServerId:   serverId,
					Channel:    channel,
					Time:       thatDayZeroTimestamp,
					C1:       0,
					C2:       0,
					C3:       0,
					C7:       0,
					C14:      0,
					C30:      0,
					C60:      0,
					C90:      0,
					C120:     0,
				}
				err = Db.FirstOrCreate(&m).Error
				if err != nil {
					return err
				}
				switch i {
				case 1:
					err = Db.Model(&m).Update("c1", rate).Error
				case 2:
					err = Db.Model(&m).Update("c2", rate).Error
				case 3:
					err = Db.Model(&m).Update("c3", rate).Error
				case 7:
					err = Db.Model(&m).Update("c7", rate).Error
				case 14:
					err = Db.Model(&m).Update("c14", rate).Error
				case 30:
					err = Db.Model(&m).Update("c30", rate).Error
				case 60:
					err = Db.Model(&m).Update("c60", rate).Error
				case 90:
					err = Db.Model(&m).Update("c90", rate).Error
				case 120:
					err = Db.Model(&m).Update("c120", rate).Error
				}
			}
			if err != nil {
				return err
			}
		}
	}
		return nil
}
