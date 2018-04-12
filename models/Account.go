package models

import (
	"github.com/chnzrb/myadmin/utils"
	"github.com/jinzhu/gorm"
	"fmt"
)

type Account struct {
	AccId             string `json:"accId"`
	ServerId          string `json:"serverId"`
	IsCreateRole      int    `json:"isCreateRole"`
	IsEnterGame       int    `json:"isEnterGame"`
	IsFinishFirstTask int    `json:"isFinishFirstTask"`
}


//获取总注册人数
func GetTotalRegister(db *gorm.DB) int {
	var count int
	db.Model(&Account{}).Count(&count)
	return count
}

//获取今日创角人数
func GetTodayRegister(db *gorm.DB) int {
	var data struct {
		Count int
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM account WHERE time between ? and ?`)
	err := db.Raw(sql, todayZeroTimestamp, todayZeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	//logs.Info("ppp:%v,%v", gameServer.Node, data.Count)
	return data.Count
}

//获取某日创角人数
func GetThatDayRegister(db *gorm.DB, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
	//todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM account WHERE time between ? and ?`)
	err := db.Raw(sql, zeroTimestamp, zeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	//logs.Info("ppp:%v,%v", gameServer.Node, data.Count)
	return data.Count
}
