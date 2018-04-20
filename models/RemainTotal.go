package models

import (
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
)

type RemainTotal struct {
	Node         string `json:"node" gorm:"primary_key"`
	Time         int    `json:"time" gorm:"primary_key"`
	RegisterRole int    `json:"registerRole" gorm:"-"`
	CreateRole   int    `json:"createRole" gorm:"-"`
	Remain2      int    `json:"remain2"`
	Remain3      int    `json:"remain3"`
	Remain4      int    `json:"remain4"`
	Remain5      int    `json:"remain5"`
	Remain6      int    `json:"remain6"`
	Remain7      int    `json:"remain7"`
}

type TotalRemainQueryParam struct {
	BaseQueryParam
	PlatformId int
	//ServerId   string
	Node string `json:"serverId"`
	StartTime  int
	EndTime    int
}

func GetRemainTotalList(params *TotalRemainQueryParam) ([]*RemainTotal, int64) {
	data := make([]*RemainTotal, 0)
	var count int64
	//gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
	//utils.CheckError(err)
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}
	err := f(Db.Model(&RemainTotal{}).Where(&RemainTotal{Node: params.Node})).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		dailyRegisterStatistics, err := GetDailyRegisterStatisticsOne(e.Node, e.Time)
		utils.CheckError(err)
		e.CreateRole = dailyRegisterStatistics.CreateRoleCount
		e.RegisterRole = dailyRegisterStatistics.RegisterCount
	}
	return data, count
}

//更新 总体留存
func UpdateRemainTotal(node string, timestamp int) error {
	logs.Info("更新总体留存:%v, %v", node, timestamp)
	serverNode, err := GetServerNode(node)
	if err != nil {
		return err
	}
	gameDb, err := GetGameDbByServerNode(serverNode)
	openDayZeroTimestamp := utils.GetThatZeroTimestamp(int64(serverNode.OpenTime))
	defer gameDb.Close()
	if err != nil {
		return err
	}

	for i := 1; i < 7; i++ {
		thatDayZeroTimestamp := timestamp - i*86400
		if openDayZeroTimestamp > thatDayZeroTimestamp {
			continue
		}
		createRolePlayerIdList := GetThatDayCreateRolePlayerIdList(gameDb, thatDayZeroTimestamp)
		createRoleNum := len(createRolePlayerIdList)

		logs.Info("createRolePlayerIdList:%v", createRolePlayerIdList)
		logs.Info("createRoleNum:%v", createRoleNum)
		rate := 0
		if createRoleNum > 0 {
			loginNum := 0
			for _, playerId := range createRolePlayerIdList {
				if IsThatDayPlayerLogin(gameDb, timestamp, playerId) {
					loginNum += 1
				}
			}
			logs.Info("loginNum:%v", loginNum)
			logs.Info("createRoleNum:%v", createRoleNum)
			rate = int(float32(loginNum) / float32(createRoleNum) * 10000)
			logs.Info("rate:%v", rate)
			//logs.Info("rate:%v", float32(loginNum) / float32(createRoleNum) * 100)
			//logs.Info("qqqqq:%v", 40 / 50)
			//logs.Info("wwww:%v", 40.0 / 50.0)
		}
		//m := &RemainTotal{}
		m := &RemainTotal{
			Node:    node,
			Time:    thatDayZeroTimestamp,
			Remain2: -1,
			Remain3: -1,
			Remain4: -1,
			Remain5: -1,
			Remain6: -1,
			Remain7: -1,
		}
		err = Db.Debug().FirstOrCreate(&m).Error
		switch i {
		case 1:
			err = Db.Debug().Model(&m).Update("Remain2", rate).Error
		case 2:
			err = Db.Debug().Model(&m).Update("Remain3", rate).Error
		case 3:
			err = Db.Debug().Model(&m).Update("Remain4", rate).Error
		case 4:
			err = Db.Debug().Model(&m).Update("Remain5", rate).Error
		case 5:
			err = Db.Debug().Model(&m).Update("Remain6", rate).Error
		case 6:
			err = Db.Debug().Model(&m).Update("Remain7", rate).Error
		}
		//err = Db.Debug().Save(&m).Error
	}

	return err
}
