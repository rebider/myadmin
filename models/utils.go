package models

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/chnzrb/myadmin/utils"
	"strconv"
	//"time"
	"sort"
	"strings"
	"errors"
)

//获取某日创角的玩家Id列表
func GetThatDayCreateRolePlayerIdList(db *gorm.DB, zeroTimestamp int) [] int {
	var data [] struct {
		Id int
	}

	sql := fmt.Sprintf(
		`SELECT id FROM player WHERE reg_time between ? and ?`)
	err := db.Raw(sql, zeroTimestamp, zeroTimestamp+86400).Find(&data).Error
	utils.CheckError(err)
	idList := make([] int, 0)
	for _, e := range data {
		idList = append(idList, e.Id)
	}
	return idList
}

// 该玩家某天是否登录过
func IsThatDayPlayerLogin(db *gorm.DB, zeroTimestamp int, playerId int) bool {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM player_login_log WHERE player_id = ? and timestamp between ? and ?`)
	err := db.Raw(sql, playerId, zeroTimestamp, zeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	//logs.Info("IsThatDayPlayerLogin:%v", data.Count)
	if data.Count == 0 {
		return false
	}
	return true
}

//获取当前在线人数
func GetNowOnlineCount(db *gorm.DB) int {
	var count int
	db.Model(&Player{}).Where(&Player{IsOnline: 1}).Count(&count)
	return count
}

//获取最高在线人数
func GetMaxOnlineCount(node string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT max(online_num) as count FROM c_ten_minute_statics WHERE node = ? `)
	err := DbCenter.Raw(sql, node).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取该服最高等级
func GetMaxPlayerLevel(db *gorm.DB) int {
	var data struct {
		MaxLevel int
	}
	sql := fmt.Sprintf(
		`SELECT max(level) as max_level FROM player_data `)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.MaxLevel
}

//获取那天平均在线时长
func GetAvgOnlineTime(node string, zeroTimestamp int) int {
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	gameDb, err := GetGameDbByServerNode(serverNode)
	utils.CheckError(err)
	var data struct {
		Time float32
	}
	sql := fmt.Sprintf(
		`SELECT avg(online_time) as time FROM player_online_log where login_time between ? and ? `)
	err = gameDb.Raw(sql, zeroTimestamp, zeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	return int(data.Time)
}

// 获取当天 平均在线人数
func GetThatDayAverageOnlineCount(node string, zeroTimestamp int) float32 {
	var data struct {
		Count float32
	}
	sql := fmt.Sprintf(
		`SELECT avg(online_num)  as count FROM c_ten_minute_statics where node = ? and time between ? and ? `)
	err := DbCenter.Raw(sql, node, zeroTimestamp, zeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

// 获取当天 最高在线人数
func GetThatDayMaxOnlineCount(node string, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT max(online_num)  as count FROM c_ten_minute_statics where node = ? and time between ? and ? `)
	err := DbCenter.Raw(sql, node, zeroTimestamp, zeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

// 获取当天 最低在线人数
func GetThatDayMinOnlineCount(node string, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT min(online_num)  as count FROM c_ten_minute_statics where node = ? and time between ? and ? `)
	err := DbCenter.Raw(sql, node, zeroTimestamp, zeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

type RemainTask struct {
	TaskId     int     `json:"taskId"`
	Count      int     `json:"count"`
	LeaveCount int     `json:"leaveCount"`
	Rate       float32 `json:"rate"`
}

// 获取任务分布
func GetRemainTask(node string) [] *RemainTask {
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	defer gameDb.Close()

	type Element struct {
		TaskId   int
		PlayerId int
		IsOnline int
	}
	mapData := make(map[int]*RemainTask, 0)
	data := make([] *RemainTask, 0)
	elementList := make([] *Element, 0)
	sql := fmt.Sprintf(
		`SELECT player.id, player.is_online, player_task.task_id FROM player left join player_task on player.id = player_task.player_id`)
	err = gameDb.Raw(sql).Find(&elementList).Error
	utils.CheckError(err)

	for _, e := range elementList {
		if r, ok := mapData[e.TaskId]; ok == true {
			if e.IsOnline == 1 {
				r.Count ++
			} else {
				r.Count ++
				r.LeaveCount ++
			}
		} else {
			if e.IsOnline == 1 {
				mapData[e.TaskId] = &RemainTask{
					TaskId:     e.TaskId,
					Count:      1,
					LeaveCount: 0,
				}
			} else {
				mapData[e.TaskId] = &RemainTask{
					TaskId:     e.TaskId,
					Count:      1,
					LeaveCount: 1,
				}
			}
		}

	}

	var keys [] int
	totalCreateRole := GetTotalCreateRoleCount(gameDb)
	for key, _ := range mapData {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, id := range keys {
		e := mapData[id]
		e.Rate = float32(e.LeaveCount) / float32(totalCreateRole) * 100
		data = append(data, e)
	}
	return data
}

//// 获取任务分布
//func GetRemainTask(platformId int, serverId string) [] *RemainTask{
//	gameDb, err:= GetGameDbByPlatformIdAndSid(platformId, serverId)
//	utils.CheckError(err)
//	defer gameDb.Close()
//	data := make([] *RemainTask, 0)
//	sql := fmt.Sprintf(
//		`SELECT task_id, count(*) as count FROM player_task group by task_id `)
//	err = gameDb.Raw(sql).Find(&data).Error
//
//	totalCreateRole := GetTotalCreateRoleCount(gameDb)
//	for _,e:= range data {
//		e.Rate = float32(e.Count) / float32(totalCreateRole) * 100
//	}
//	return data
//}

type RemainLevel struct {
	Level      int     `json:"level"`
	Count      int     `json:"count"`
	LeaveCount int     `json:"leaveCount"`
	Rate       float32 `json:"rate"`
}

// 获取等级分布
func GetRemainLevel(node string) [] *RemainLevel {
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	defer gameDb.Close()

	type Element struct {
		Level    int
		PlayerId int
		IsOnline int
	}
	mapData := make(map[int]*RemainLevel, 0)
	data := make([] *RemainLevel, 0)
	elementList := make([] *Element, 0)
	sql := fmt.Sprintf(
		`SELECT player.id, player.is_online, player_data.level FROM player left join player_data on player.id = player_data.player_id`)
	err = gameDb.Raw(sql).Find(&elementList).Error
	utils.CheckError(err)

	for _, e := range elementList {
		if r, ok := mapData[e.Level]; ok == true {
			if e.IsOnline == 1 {
				r.Count ++
			} else {
				r.Count ++
				r.LeaveCount ++
			}
		} else {
			if e.IsOnline == 1 {
				mapData[e.Level] = &RemainLevel{
					Level:      e.Level,
					Count:      1,
					LeaveCount: 0,
				}
			} else {
				mapData[e.Level] = &RemainLevel{
					Level:      e.Level,
					Count:      1,
					LeaveCount: 1,
				}
			}
		}

	}

	var keys [] int
	totalCreateRole := GetTotalCreateRoleCount(gameDb)
	for key, _ := range mapData {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, id := range keys {
		e := mapData[id]
		e.Rate = float32(e.LeaveCount) / float32(totalCreateRole) * 100
		data = append(data, e)
	}
	return data
	//gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	//utils.CheckError(err)
	//defer gameDb.Close()
	//data := make([] *RemainLevel, 0)
	//sql := fmt.Sprintf(
	//	`SELECT level, count(*) as count FROM player_data group by level `)
	//err = gameDb.Raw(sql).Find(&data).Error
	//utils.CheckError(err)
	//
	//totalCreateRole := GetTotalCreateRoleCount(gameDb)
	//for _, e := range data {
	//	e.Rate = float32(e.Count) / float32(totalCreateRole) * 100
	//}
	//return data
}

type RemainTime struct {
	StartTime  int     `json:"-"`
	EndTime    int     `json:"-"`
	TimeString string  `json:"timeString"`
	Count      int     `json:"count"`
	LeaveCount int     `json:"leaveCount"`
	Rate       float32 `json:"rate"`
}

// 获取时长分布
func GetRemainTime(node string) [] *RemainTime {
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	defer gameDb.Close()
	var data = [] *RemainTime{
		&RemainTime{
			StartTime:  0,
			EndTime:    60,
			TimeString: "小于1分钟",
		},
		&RemainTime{
			StartTime:  60,
			EndTime:    300,
			TimeString: "1~5分钟",
		},
		&RemainTime{
			StartTime:  300,
			EndTime:    600,
			TimeString: "5~10分钟",
		},
		&RemainTime{
			StartTime:  600,
			EndTime:    1800,
			TimeString: "10~30分钟",
		},
		&RemainTime{
			StartTime:  1800,
			EndTime:    3600,
			TimeString: "30~60分钟",
		},
		&RemainTime{
			StartTime:  3600,
			EndTime:    3600 * 2,
			TimeString: "1~2小时",
		},
		&RemainTime{
			StartTime:  3600 * 2,
			EndTime:    3600 * 3,
			TimeString: "2~3小时",
		},
		&RemainTime{
			StartTime:  3600 * 3,
			EndTime:    3600 * 4,
			TimeString: "3~4小时",
		},
		&RemainTime{
			StartTime:  3600 * 4,
			EndTime:    3600 * 5,
			TimeString: "4~5小时",
		},
		&RemainTime{
			StartTime:  3600 * 5,
			EndTime:    3600 * 6,
			TimeString: "5~6小时",
		},
		&RemainTime{
			StartTime:  3600 * 6,
			EndTime:    3600 * 9,
			TimeString: "6~9小时",
		},
		&RemainTime{
			StartTime:  3600 * 9,
			EndTime:    3600 * 12,
			TimeString: "9~12小时",
		},
		&RemainTime{
			StartTime:  3600 * 12,
			EndTime:    3600 * 24,
			TimeString: "12~24小时",
		},
		&RemainTime{
			StartTime:  3600 * 24,
			EndTime:    3600 * 48,
			TimeString: "1~2天",
		},
		&RemainTime{
			StartTime:  3600 * 48,
			EndTime:    3600 * 72,
			TimeString: "2~3天",
		},
		&RemainTime{
			StartTime:  3600 * 72,
			EndTime:    3600 * 999999,
			TimeString: ">3天",
		},
	}
	type Element struct {
		OnlineTime int
		//PlayerId int
		IsOnline int
	}
	//elementList := make([] *Element, 0)
	//sql := fmt.Sprintf(
	//	`SELECT is_online, total_online_time FROM player`)
	//err = gameDb.Raw(sql).Find(&elementList).Error
	totalCreateRole := GetTotalCreateRoleCount(gameDb)
	for _, e := range data {
		elementList := make([] *Element, 0)
		sql := fmt.Sprintf(
			`SELECT is_online, total_online_time FROM player where total_online_time >= ? and total_online_time < ? `)
		err = gameDb.Raw(sql, e.StartTime, e.EndTime).Find(&elementList).Error
		utils.CheckError(err)
		for _, ee := range elementList {
			e.Count ++
			if ee.IsOnline == 0 {
				e.LeaveCount ++
			}
		}
		e.Rate = float32(e.LeaveCount) / float32(totalCreateRole) * 100
	}
	return data
	//totalCreateRole := GetTotalCreateRoleCount(gameDb)
	//for _, e := range data {
	//	sql := fmt.Sprintf(
	//		`SELECT count(*) as count FROM player where total_online_time >= ? and total_online_time < ? `)
	//	err = gameDb.Raw(sql, e.StartTime, e.EndTime).Find(&e).Error
	//	utils.CheckError(err)
	//	e.Rate = float32(e.Count) / float32(totalCreateRole) * 100
	//}
	//return data
}

func get24hoursOnlineCount(node string, zeroTimestamp int) [] string {
	onlineCountList := make([] string, 0)
	for i := zeroTimestamp; i < zeroTimestamp+86400; i = i + 10*60 {
		cServerTraceLog := &CTenMinuteStatics{}
		err := DbCenter.Where(&CTenMinuteStatics{
			Node: node,
			Time: i,
		}).First(&cServerTraceLog).Error
		if err == nil {
			onlineCountList = append(onlineCountList, strconv.Itoa(cServerTraceLog.OnlineNum))
		} else {
			onlineCountList = append(onlineCountList, "null")
		}
	}
	//logs.Info("%+v", len(onlineCountList))
	return onlineCountList
}

func get24hoursRegisterCount(node string, zeroTimestamp int) [] string {
	onlineCountList := make([] string, 0)
	for i := zeroTimestamp; i < zeroTimestamp+86400; i = i + 10*60 {
		cServerTraceLog := &CTenMinuteStatics{}
		err := DbCenter.Where(&CTenMinuteStatics{
			Node: node,
			Time: i,
		}).First(&cServerTraceLog).Error
		if err == nil {
			onlineCountList = append(onlineCountList, strconv.Itoa(cServerTraceLog.RegisterCount))
		} else {
			onlineCountList = append(onlineCountList, "null")
		}
	}
	//logs.Info("%+v", len(onlineCountList))
	return onlineCountList
}

//获取玩家名字
func GetPlayerName(db *gorm.DB, playerId int) string {
	var data struct {
		ServerId string
		Name     string
	}

	sql := fmt.Sprintf(
		`SELECT server_id, nickname as name FROM player where id = ? `)
	err := db.Raw(sql, playerId).Scan(&data).Error
	utils.CheckError(err)
	return data.ServerId + "." + data.Name
}

//获取玩家最近登录时间
func GetPlayerLastLoginTime(platformId int, serverId string, playerId int) int {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	defer gameDb.Close()
	var data struct {
		Time int
	}
	sql := fmt.Sprintf(
		`SELECT last_login_time as time FROM player where id = ? `)
	err = gameDb.Raw(sql, playerId).Scan(&data).Error
	utils.CheckError(err)
	return data.Time
}

//获取玩家名字
func GetPlayerName_2(platformId int, serverId string, playerId int) string {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	defer gameDb.Close()
	var data struct {
		ServerId string
		Name string
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)

	sql := fmt.Sprintf(
		`SELECT server_id, nickname as name FROM player where id = ? `)
	//logs.Info("GetPlayerName:%v", playerId)
	err = gameDb.Raw(sql, playerId).Scan(&data).Error
	utils.CheckError(err)
	//logs.Info("ppp:%v,%v", data.MaxLevel)
	return data.ServerId + "." + data.Name
}

//获取区服付费人数
func GetServerChargePlayerCount(node string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where server_id in(%s) and charge_type = 99;`, GetGameServerIdListStringByNode(node))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服付费人数
func GetThatDayServerChargePlayerCount(node string, time int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where node = ? and charge_type = 99 and ( record_time between ? and ?);`)
	err := DbCharge.Raw(sql, node, time, time+86400).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

// arpu
func CaclARPU(totalChargeValueNum int, totalChargePlayerCount int) float32 {
	if totalChargePlayerCount == 0 {
		return 0
	}
	return float32(totalChargeValueNum) / float32(totalChargePlayerCount) / 100
}

//付费率
func CaclChargeRate(totalChargePlayerCount int, totalRoleCount int) float32 {
	if totalRoleCount == 0 {
		return 0
	}
	return float32(totalChargePlayerCount) / float32(totalRoleCount)
}

//二次付费率
func CaclSceondChargeRate(secondChargePlayerCount int, totalChargePlayerCount int) float32 {
	if totalChargePlayerCount == 0 {
		return 0
	}
	return float32(secondChargePlayerCount) / float32(totalChargePlayerCount)
}

//获取区服二次付费人数
func GetServerSecondChargePlayerCount(node string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where server_id in (%s) and is_first = 0 and charge_type = 99;`, GetGameServerIdListStringByNode(node))
	err := DbCharge.Debug().Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服首次付费人数
func GetThadDayServerSecondChargePlayerCount(node string, time int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where node = ? and is_first = 0 and charge_type = 99 and (record_time between ? and ?);`)
	err := DbCharge.Raw(sql, node, time, time+86400).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服总充值元宝
func GetServerTotalChargeIngot(node string) int {
	var data struct {
		Count int
	}
	//gameServer, err := GetGameServerOne(platformId, serverId)
	//utils.CheckError(err)
	//if err != nil {
	//	return 0
	//}
	sql := fmt.Sprintf(
		`select sum(ingot) as count from charge_info_record where server_id in (%s) and charge_type = 99;`, GetGameServerIdListStringByNode(node))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服总充值人民币
func GetThatDayServerTotalChargeMoney(node string, time int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select sum(money) as count from charge_info_record where node = ? and charge_type = 99 and (record_time between ? and ?);`)
	err := DbCharge.Raw(sql, node, time, time+86400).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服总充值人民币
func GetServerTotalChargeMoney(node string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select sum(money) as count from charge_info_record where server_id in (%s) and charge_type = 99;`, GetGameServerIdListStringByNode(node))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

func SplitPlayerName(playerName string) (string, string, error) {
	array := strings.Split(playerName, ".")
	if len(array) == 2 {
		return array[0], array[1], nil
	}
	return "", "", errors.New(fmt.Sprintf("解析玩家名字失败:%s", playerName))
}
