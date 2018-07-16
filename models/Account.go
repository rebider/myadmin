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
	PlayerId          int    `json:"playerId"`
	IsFinishFirstTask int    `json:"isFinishFirstTask"`
	Time              int    `json:"time"`
}

//获取总注册人数
func GetTotalRegister(db *gorm.DB) int {
	var count int
	db.Model(&Account{}).Count(&count)
	return count
}

//获取今日注册人数
func GetTodayRegister(db *gorm.DB) int {
	return GetThatDayRegister(db, utils.GetTodayZeroTimestamp())
}

//获取某日注册人数
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

//获取总创角色人数
func GetTotalCreateRoleCount(db *gorm.DB) int {
	var count int
	db.Model(&Account{}).Where(&Account{
		IsCreateRole: 1,
	}).Count(&count)
	return count
}

//获取今日创角人数
func GetTodayCreateRoleCount(db *gorm.DB) int {
	return GetThatDayCreateRoleCount(db, utils.GetTodayZeroTimestamp())
}

//获取某日创角人数
func GetThatDayCreateRoleCount(db *gorm.DB, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM account WHERE is_create_role = 1 and (time between ? and ?) `)
	err := db.Raw(sql, zeroTimestamp, zeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取当前全服资源总数
func GetTotalProp(db *gorm.DB, propType int, propId int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT sum(num) as count FROM player_prop WHERE prop_type = '%d' and prop_id = '%d'`, propType, propId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取总有效角色人数
func GetTotalValidCreateRoleCount(db *gorm.DB) int {
	var count int
	db.Model(&Account{}).Where(&Account{
		IsCreateRole: 1,
	}).Count(&count)
	return count
}

//获取今日有效角色人数
func GetTodayValidCreateRoleCount(db *gorm.DB) int {
	return GetThatDayValidCreateRoleCount(db, utils.GetTodayZeroTimestamp())
}

//获取某日有效角色人数
func GetThatDayValidCreateRoleCount(db *gorm.DB, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM account WHERE is_create_role = 1 and is_finish_first_task = 1 and (time between ? and ?) `)
	err := db.Raw(sql, zeroTimestamp, zeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}
