package models

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/chnzrb/myadmin/utils"
	"strconv"
	"sort"
	"strings"
	"errors"
	"github.com/astaxie/beego/logs"
	"os"
)

//获取某日创角的玩家Id列表
func GetThatDayCreateRolePlayerIdList(db *gorm.DB, serverId string, channel string, zeroTimestamp int) [] int {
	var data [] struct {
		Id int
	}

	sql := fmt.Sprintf(
		`SELECT id FROM player WHERE reg_time between %d and %d and channel = '%s' and server_id = '%s'`, zeroTimestamp, zeroTimestamp+86400, channel, serverId)
	err := db.Raw(sql).Find(&data).Error
	utils.CheckError(err)
	idList := make([] int, 0)
	for _, e := range data {
		idList = append(idList, e.Id)
	}
	return idList
}

// 是否该玩家某天登录过
func IsThatDayPlayerLogin(db *gorm.DB, zeroTimestamp int, playerId int) bool {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM player_login_log WHERE player_id = %d and timestamp between %d and %d`, playerId, zeroTimestamp, zeroTimestamp+86400)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	if data.Count == 0 {
		return false
	}
	return true
}

// 是否该玩家连续登录过
func IsPlayerContinueLogin(db *gorm.DB, zeroTimestamp int, continueDay int, playerId int) bool {
	for i := 0; i <= continueDay; i ++ {
		if IsThatDayPlayerLogin(db, zeroTimestamp+86400*i, playerId) == false {
			return false
		}
	}
	return true
}

func GetSQLWhereParam(args [] string) string {
	return "'" + strings.Join(args, "','") + "'"
}

// 获取该天登录次数
func GetThatDayLoginTimes(db *gorm.DB, serverId string, channel string, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM player_login_log WHERE timestamp between %d and %d and player_id in (select id from player where channel = '%s' and server_id = '%s')`, zeroTimestamp, zeroTimestamp+86400, channel, serverId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

// 获取该天登录的玩家数量
func GetThatDayLoginPlayerCount(db *gorm.DB, serverId string, channel string, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM player WHERE last_login_time between %d and %d and channel = '%s' and server_id = '%s'`, zeroTimestamp, zeroTimestamp+86400, channel, serverId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

// 获取该天活跃玩家数量
func GetThatDayActivePlayerCount(db *gorm.DB, serverId string, channel string, zeroTimestamp int) int {
	count := 0
	data := make([] *Player, 0)
	sql := fmt.Sprintf(
		`SELECT * FROM player WHERE last_login_time between %d and %d and channel = '%s' and server_id = '%s'`, zeroTimestamp, zeroTimestamp+86400, channel, serverId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		if e.LoginTimes >= (zeroTimestamp+86400-e.RegTime)/86400 {
			count ++
		}
	}
	return count
}

//获取当前在线人数
func GetNowOnlineCountByNode(node string) int {
	gameDb, err := GetGameDbByNode(node)

	utils.CheckError(err)
	if err != nil {
		return -1
	}
	defer gameDb.Close()
	return GetNowOnlineCount(gameDb)
}

//获取当前在线人数
func GetNowOnlineCount(db *gorm.DB) int {
	var count int
	db.Model(&Player{}).Where(&Player{IsOnline: 1}).Count(&count)
	return count
}

//获取当前在线人数
func GetNowOnlineCount2(db *gorm.DB, serverId string, channelList [] string) int {
	var count int
	db.Model(&Player{}).Where(&Player{ServerId: serverId, IsOnline: 1}).Where(" channel in(?)", channelList).Count(&count)
	return count
}

//获取当前在线ip数
func GetNowOnlineIpCount(db *gorm.DB, serverId string, channelList [] string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT COUNT( DISTINCT last_login_ip ) as count FROM player where is_online = 1 and server_id = '%s' and channel in(%s);`, serverId, GetSQLWhereParam(channelList))
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取最高在线人数
//func GetMaxOnlineCount(node string) int {
//	var data struct {
//		Count int
//	}
//	sql := fmt.Sprintf(
//		`SELECT max(online_num) as count FROM c_ten_minute_statics WHERE node = '%s' `, node)
//	err := DbCenter.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}

//获取时间段内最高在线人数
func GetThatDayMaxOnlineCount(platformId string, serverId string, channelList [] string, startTime int, endTime int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT max(online_count) as count FROM ten_minute_statistics WHERE platform_id = '%s' and server_id = '%s' and channel in(%s) and time between %d and %d`, platformId, serverId, GetSQLWhereParam(channelList), startTime, endTime)
	err := Db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取时间段内最低在线人数
func GetThatDayMinOnlineCount(platformId string, serverId string, channelList [] string, startTime int, endTime int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT min(online_count) as count FROM ten_minute_statistics WHERE platform_id = '%s' and server_id = '%s' and channel in(%s) and time between %d and %d`, platformId, serverId, GetSQLWhereParam(channelList), startTime, endTime)
	err := Db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取时间段内平均在线人数
func GetThatDayAvgOnlineCount(platformId string, serverId string, channelList [] string, startTime int, endTime int) int {
	var data struct {
		Count float32
	}
	sql := fmt.Sprintf(
		`SELECT avg(online_count) as count FROM ten_minute_statistics WHERE platform_id = '%s' and server_id = '%s' and channel in(%s) and time between %d and %d`, platformId, serverId, GetSQLWhereParam(channelList), startTime, endTime)
	err := Db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return int(data.Count)
}

//// 获取当天 平均在线人数
//func GetThatDayAverageOnlineCountByChannel(node string, channel string, zeroTimestamp int) float32 {
//	var data struct {
//		Count float32
//	}
//	sql := fmt.Sprintf(
//		`SELECT avg(online_num)  as count FROM c_ten_minute_statics where node = '%s' and time between %d and %d and channel = '%s'`, node, zeroTimestamp, zeroTimestamp+86400, channel)
//	err := DbCenter.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}

//// 获取当天 平均在线人数
//func GetThatDayAverageOnlineCount(node string, zeroTimestamp int) float32 {
//	var data struct {
//		Count float32
//	}
//	sql := fmt.Sprintf(
//		`SELECT avg(online_num)  as count FROM c_ten_minute_statics where node = '%s' and time between %d and %d `, node, zeroTimestamp, zeroTimestamp+86400)
//	err := DbCenter.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}

//// 获取当天 最高在线人数
//func GetThatDayMaxOnlineCount(node string, channel string, zeroTimestamp int) int {
//	var data struct {
//		Count int
//	}
//	sql := fmt.Sprintf(
//		`SELECT max(online_num)  as count FROM c_ten_minute_statics where node = '%s' and time between %d and %d and channel = '%s'`, node, zeroTimestamp, zeroTimestamp+86400, channel)
//	err := DbCenter.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}

//// 获取当天 最低在线人数
//func GetThatDayMinOnlineCount(node string, channel string, zeroTimestamp int) int {
//	var data struct {
//		Count int
//	}
//	sql := fmt.Sprintf(
//		`SELECT min(online_num)  as count FROM c_ten_minute_statics where node = '%s' and time between %d and %d and channel = '%s'`, node, zeroTimestamp, zeroTimestamp+86400, channel)
//	err := DbCenter.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}

//获取该服最高等级
func GetMaxPlayerLevel(db *gorm.DB, serverId string, channelList [] string) int {
	var data struct {
		MaxLevel int
	}
	sql := fmt.Sprintf(
		`SELECT max(level) as max_level FROM player_data where player_id in (select id from player where server_id = '%s' and channel in (%s))`, serverId, GetSQLWhereParam(channelList))
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.MaxLevel
}

//获取那天在线时长
func GetOnlineTime(node string, serverId string, channel string, zeroTimestamp int) int {
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	defer gameDb.Close()
	var data struct {
		Time int
	}
	sql := fmt.Sprintf(
		`SELECT sum(online_time) as time FROM player_online_log where login_time between %d and %d and player_id in (select id from player where channel = '%s' and server_id = '%s')`, zeroTimestamp, zeroTimestamp+86400, channel, serverId)
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return int(data.Time)
}

////获取那天平均在线时长
//func GetAvgOnlineTime(node string, channel string, zeroTimestamp int) int {
//	gameDb, err := GetGameDbByNode(node)
//	utils.CheckError(err)
//	if err != nil {
//		return 0
//	}
//	defer gameDb.Close()
//	var data struct {
//		Time float32
//	}
//	sql := fmt.Sprintf(
//		`SELECT avg(online_time) as time FROM player_online_log where login_time between %d and %d and player_id in (select id from player where channel = '%s')`, zeroTimestamp, zeroTimestamp+86400, channel)
//	err = gameDb.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return int(data.Time)
//}

type RemainTask struct {
	TaskId     int     `json:"taskId"`
	Count      int     `json:"count"`
	LeaveCount int     `json:"leaveCount"`
	Rate       float32 `json:"rate"`
}

// 获取任务分布
func GetRemainTask(platformId string, serverId string, channelList [] string) [] *RemainTask {
	gameServer, err := GetGameServerOne(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return nil
	}
	node := gameServer.Node
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil
	}
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
		`SELECT player.id, player.is_online, player_task.task_id FROM player left join player_task on player.id = player_task.player_id where player.channel in (%s)`, GetSQLWhereParam(channelList))
	err = gameDb.Raw(sql).Find(&elementList).Error
	utils.CheckError(err)
	if err != nil {
		return nil
	}
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
	totalCreateRole := GetTotalCreateRoleCountByChannelList(gameDb, serverId, channelList)
	for key, _ := range mapData {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, id := range keys {
		e := mapData[id]
		if totalCreateRole > 0 {
			e.Rate = float32(e.LeaveCount) / float32(totalCreateRole) * 100
		}
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
func GetRemainLevel(platformId string, serverId string, channelList [] string) [] *RemainLevel {
	gameServer, err := GetGameServerOne(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return nil
	}
	node := gameServer.Node
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil
	}
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
		`SELECT player.id, player.is_online, player_data.level FROM player left join player_data on player.id = player_data.player_id where player.channel in (%s)`, GetSQLWhereParam(channelList))
	err = gameDb.Raw(sql).Find(&elementList).Error
	utils.CheckError(err)
	if err != nil {
		return nil
	}
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
	totalCreateRole := GetTotalCreateRoleCountByChannelList(gameDb, serverId, channelList)
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
func GetRemainTime(platformId string, serverId string, channelList [] string) [] *RemainTime {
	gameServer, err := GetGameServerOne(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return nil
	}
	node := gameServer.Node
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil
	}
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
	totalCreateRole := GetTotalCreateRoleCountByChannelList(gameDb, serverId, channelList)
	for _, e := range data {
		elementList := make([] *Element, 0)
		sql := fmt.Sprintf(
			`SELECT is_online, total_online_time FROM player where total_online_time >= %d and total_online_time < %d and channel in (%s) `, e.StartTime, e.EndTime, GetSQLWhereParam(channelList))
		err = gameDb.Raw(sql).Find(&elementList).Error
		utils.CheckError(err)
		if err != nil {
			return nil
		}
		for _, ee := range elementList {
			e.Count ++
			if ee.IsOnline == 0 {
				e.LeaveCount ++
			}
		}
		if totalCreateRole > 0 {
			e.Rate = float32(e.LeaveCount) / float32(totalCreateRole) * 100
		}

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

func get24hoursOnlineCount(platformId string, serverId string, channelList [] string, zeroTimestamp int) ([] string, int) {
	onlineCountList := make([] string, 0, 144)
	now := utils.GetTimestamp()
	//gameServer, _ := GetGameServerOne(platformId, serverId)
	nowOnline := 0

	whereArray := make([] string, 0)
	//whereArray = append(whereArray, fmt.Sprintf("time = %d", i))
	whereArray = append(whereArray, fmt.Sprintf("platform_id = '%s'", platformId))
	whereArray = append(whereArray, fmt.Sprintf("online_count > 0"))
	if serverId != "" {
		whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", serverId))
	}
	whereArray = append(whereArray, fmt.Sprintf("channel in(%s)", GetSQLWhereParam(channelList)))
	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	for i := zeroTimestamp; i < zeroTimestamp+86400; i = i + 10*60 {
		if i < now {
			var data struct {
				Sum int
			}
			sql := fmt.Sprintf(
				`SELECT sum(online_count) as sum from ten_minute_statistics %s and time = %d`, whereParam, i)
			err := Db.Raw(sql).Scan(&data).Error
			if err == nil {
				nowOnline = data.Sum
				onlineCountList = append(onlineCountList, strconv.Itoa(data.Sum))
			} else {
				onlineCountList = append(onlineCountList, "null")
			}
		} else {
			onlineCountList = append(onlineCountList, "null")
		}
	}
	return onlineCountList, nowOnline
}

func get24hoursRegisterCount(platformId string, serverId string, channelList [] string, zeroTimestamp int) ([] string, int) {
	onlineCountList := make([] string, 0, 144)
	now := utils.GetTimestamp()
	totalCount := 0

	whereArray := make([] string, 0)
	//whereArray = append(whereArray, fmt.Sprintf("time = %d", i))
	whereArray = append(whereArray, fmt.Sprintf("platform_id = '%s'", platformId))
	whereArray = append(whereArray, fmt.Sprintf("register_count > 0"))
	if serverId != "" {
		whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", serverId))
	}
	whereArray = append(whereArray, fmt.Sprintf("channel in(%s)", GetSQLWhereParam(channelList)))
	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}
	for i := zeroTimestamp; i < zeroTimestamp+86400; i = i + 10*60 {
		if i < now {
			var data struct {
				Sum int
			}

			sql := fmt.Sprintf(
				`SELECT sum(register_count) as sum from ten_minute_statistics %s and time = %d`, whereParam, i)
			err := Db.Raw(sql).Scan(&data).Error
			utils.CheckError(err)
			if err == nil {
				totalCount += data.Sum
				onlineCountList = append(onlineCountList, strconv.Itoa(data.Sum))
			} else {
				onlineCountList = append(onlineCountList, "null")
			}
		} else {
			onlineCountList = append(onlineCountList, "null")
		}
	}
	return onlineCountList, totalCount
}

func get24hoursChargeCount(platformId string, serverId string, channelList [] string, zeroTimestamp int) ([] string, int) {
	chargeCountList := make([] string, 0, 144)
	now := utils.GetTimestamp()
	totalCount := 0
	//gameServer, _ := GetGameServerOne(platformId, serverId)
	for i := zeroTimestamp + 600; i <= zeroTimestamp+86400; i = i + 10*60 {
		if i < now {
			var data struct {
				Sum int
			}
			whereArray := make([] string, 0)
			whereArray = append(whereArray, fmt.Sprintf("time = %d", i))
			whereArray = append(whereArray, fmt.Sprintf("platform_id = '%s'", platformId))
			if serverId != "" {
				whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", serverId))
			}
			whereArray = append(whereArray, fmt.Sprintf("channel in(%s)", GetSQLWhereParam(channelList)))
			whereParam := strings.Join(whereArray, " and ")
			if whereParam != "" {
				whereParam = " where " + whereParam
			}
			sql := fmt.Sprintf(
				`SELECT sum(charge_count) as sum from ten_minute_statistics %s`, whereParam)
			err := Db.Raw(sql).Scan(&data).Error
			utils.CheckError(err)
			if err == nil {
				totalCount += data.Sum
				chargeCountList = append(chargeCountList, strconv.Itoa(totalCount))
			} else {
				chargeCountList = append(chargeCountList, "null")
			}
		} else {
			chargeCountList = append(chargeCountList, "null")
		}
	}
	//logs.Info("%+v", len(onlineCountList))
	return chargeCountList, totalCount
}

func get24hoursChargePlayerCount(platformId string, serverId string, channelList [] string, zeroTimestamp int) ([] string, int) {
	chargePlayerCountList := make([] string, 0, 144)
	now := utils.GetTimestamp()
	totalCount := 0
	//gameServer, _ := GetGameServerOne(platformId, serverId)
	for i := zeroTimestamp + 600; i <= zeroTimestamp+86400; i = i + 10*60 {
		if i < now {
			var data struct {
				Sum int
			}
			whereArray := make([] string, 0)
			whereArray = append(whereArray, fmt.Sprintf("time = %d", i))
			whereArray = append(whereArray, fmt.Sprintf("platform_id = '%s'", platformId))
			if serverId != "" {
				whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", serverId))
			}
			whereArray = append(whereArray, fmt.Sprintf("channel in(%s)", GetSQLWhereParam(channelList)))
			whereParam := strings.Join(whereArray, " and ")
			if whereParam != "" {
				whereParam = " where " + whereParam
			}
			sql := fmt.Sprintf(
				`SELECT sum(charge_player_count) as sum from ten_minute_statistics %s`, whereParam)
			err := Db.Raw(sql).Scan(&data).Error
			utils.CheckError(err)
			if err == nil {
				if data.Sum > totalCount {
					totalCount = data.Sum
				}
				chargePlayerCountList = append(chargePlayerCountList, strconv.Itoa(totalCount))
			} else {
				chargePlayerCountList = append(chargePlayerCountList, "null")
			}
		} else {
			chargePlayerCountList = append(chargePlayerCountList, "null")
		}
	}
	//logs.Info("%+v", len(onlineCountList))
	return chargePlayerCountList, totalCount
}

//获取玩家名字
func GetPlayerName(db *gorm.DB, playerId int) string {
	var data struct {
		ServerId string
		Name     string
	}

	sql := fmt.Sprintf(
		`SELECT server_id, nickname as name FROM player where id = %d `, playerId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.ServerId + "." + data.Name
}

//获取玩家最近登录时间
func GetPlayerLastLoginTime(platformId string, serverId string, playerId int) int {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	defer gameDb.Close()
	var data struct {
		Time int
	}
	sql := fmt.Sprintf(
		`SELECT last_login_time as time FROM player where id = %d `, playerId)
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Time
}

//获取玩家名字
func GetPlayerName_2(platformId string, serverId string, playerId int) string {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return ""
	}
	defer gameDb.Close()
	var data struct {
		ServerId string
		Name     string
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)

	sql := fmt.Sprintf(
		`SELECT server_id, nickname as name FROM player where id = %d `, playerId)
	//logs.Info("GetPlayerName:%v", playerId)
	err = gameDb.Raw(sql).Scan(&data).Error
	//utils.CheckError(err)
	if err != nil {
		return fmt.Sprintf("角色不存在:%s_%d", serverId, playerId)
	}
	//logs.Info("ppp:%v,%v", data.MaxLevel)
	return data.ServerId + "." + data.Name
}

//获取区服付费人数
func GetServerChargePlayerCount(platformId string, serverId string, channelList [] string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel in(%s) and charge_type = 99;`, platformId, serverId, GetSQLWhereParam(channelList))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服付费人数
func GetThatDayServerChargePlayerCount(platformId string, serverId string, channel string, time int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and ( record_time between %d and %d);`, platformId, serverId, channel, time, time+86400)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服付费人数
func GetThatDayChargePlayerCount(platformId string, serverId string, channel string, time int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and  record_time < %d ;`, platformId, serverId, channel, time)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//// arpu
//func CaclARPU(totalChargeValueNum int, totalChargePlayerCount int) float32 {
//	if totalChargePlayerCount == 0 {
//		return 0
//	}
//	if totalChargePlayerCount > 0 {
//		return float32(totalChargeValueNum) / float32(totalChargePlayerCount)
//	}
//	return float32(0)
//}

//率
func CaclRate(totalChargePlayerCount int, totalRoleCount int) float32 {
	if totalRoleCount == 0 {
		return 0
	}
	if totalRoleCount > 0 {
		return float32(totalChargePlayerCount) / float32(totalRoleCount)
	}
	return float32(0)
}

func GetNodeListByServerIdList(platformId string, serverIdList [] string) [] string {
	nodeList := make([] string, 0)
	for _, serverId := range serverIdList {
		gameServer, err := GetGameServerOne(platformId, serverId)
		if err != nil {
			logs.Error("获取节点列表失败!!!! %+v, %+v, %+v", platformId, serverId, err)
			return nodeList
		}
		isContain := false
		for _, node := range nodeList {
			if node == gameServer.Node {
				isContain = true
			}
		}
		if isContain == false {
			nodeList = append(nodeList, gameServer.Node)
		}
	}
	return nodeList
}

////二次付费率
//func CaclChargeRate(secondChargePlayerCount int, totalChargePlayerCount int) float32 {
//	if totalChargePlayerCount == 0 {
//		return 0
//	}
//	if totalChargePlayerCount > 0 {
//		return float32(secondChargePlayerCount) / float32(totalChargePlayerCount)
//	}
//	return float32(0)
//}

//获取区服二次付费人数
//func GetServerSecondChargePlayerCount(node string) int {
//	var data struct {
//		Count int
//	}
//	sql := fmt.Sprintf(
//		`select count(DISTINCT player_id) as count from charge_info_record where server_id in (%s) and is_first = 0 and charge_type = 99;`, GetGameServerIdListStringByNode(node))
//	err := DbCharge.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}
func GetServerSecondChargePlayerCount(platformId string, serverId string, channelList [] string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel in(%s) and is_first = 0 and charge_type = 99;`, platformId, serverId, GetSQLWhereParam(channelList))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服充值次数列表
func GetServerChargeCountList(platformId string, serverId string, channelList [] string) [] int {
	var data [] struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(player_id) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel in (%s)  and charge_type = 99 group by player_id;`, platformId, serverId, GetSQLWhereParam(channelList))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	r := make([] int, 0)
	//logs.Debug("data:%+v", data)
	for _, e := range data {
		r = append(r, e.Count)
	}
	//logs.Debug("r:%+v", r)
	return r
}

//获取区服首次付费人数
func GetThadDayServerFirstChargePlayerCount(platformId string, serverId string, channel string, time int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and is_first = 1 and charge_type = 99 and (record_time between %d and %d);`, platformId, serverId, channel, time, time+86400)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服新增付费人数
func GetThadDayNewChargePlayerCount(platformId string, serverId string, channel string, time int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and is_first = 1 and charge_type = 99 and (record_time between %d and %d) and (reg_time between %d and %d);`, platformId, serverId, channel, time, time+86400, time, time+86400)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取时间内充值的玩家id列表
func GetChargePlayerIdList(platformId string, serverId string, channel string, startTime int, endTime int) [] int {
	var data [] struct {
		Id int
	}
	sql := fmt.Sprintf(
		`select DISTINCT player_id as id from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s'  and charge_type = 99 and (record_time between %d and %d) ;`, platformId, serverId, channel, startTime, endTime)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	idList := make([] int, 0)
	for _, e := range data {
		idList = append(idList, e.Id)
	}
	return idList
}

//获取区服总充值元宝
func GetServerTotalChargeIngot(platformId string, serverId string, channelList [] string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select sum(ingot) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel in(%s) and charge_type = 99;`, platformId, serverId, GetSQLWhereParam(channelList))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取该时间内区服总充值人民币
func GetTotalChargeMoney(platformId string, serverId string, channel string, startTime int, endTime int) int {
	var data struct {
		Count float32
	}
	sql := fmt.Sprintf(
		`select sum(money) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and record_time between %d and %d;`, platformId, serverId, channel, startTime, endTime)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return int(data.Count)
}

//获取该时间内区服充值人数
func GetTotalChargePlayerCount(platformId string, serverId string, channel string, startTime int, endTime int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and record_time between %d and %d;`, platformId, serverId, channel, startTime, endTime)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取某天注册的玩家 该时间区间内总充值人民币
func GetTotalChargeMoneyByRegisterTime(platformId string, serverId string, channel string, startTime int, endTime int, registerTime int) int {
	var data struct {
		Count float32
	}
	sql := fmt.Sprintf(
		`select sum(money) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and record_time between %d and %d and reg_time between %d and %d;`, platformId, serverId, channel, startTime, endTime, registerTime, registerTime+86400)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return int(data.Count)
}

//获取区服总充值人民币
func GetServerTotalChargeMoneyByChannelList(platformId string, serverId string, channelList [] string) int {
	var data struct {
		Count float32
	}
	sql := fmt.Sprintf(
		`select sum(money) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel in(%s) and charge_type = 99;`, platformId, serverId, GetSQLWhereParam(channelList))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return int(data.Count)
}

//获取区服总充值人民币
//func GetTotalChargeMoney(platformId string, serverId string, channel string, time int) int {
//	var data struct {
//		Count float32
//	}
//	sql := fmt.Sprintf(
//		`select sum(money) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and record_time < %d;`, platformId, serverId, channel, time)
//	err := DbCharge.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return int(data.Count)
//}

func GetThatDayNewChargeMoney(platformId string, serverId string, channel string, time int) int {
	var data struct {
		Count float32
	}
	sql := fmt.Sprintf(
		`select sum(money) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and (record_time between %d and %d) and (reg_time between %d and %d);`, platformId, serverId, channel, time, time+86400, time, time+86400)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return int(data.Count)
}

func SplitPlayerName(playerName string) (string, string, error) {
	array := strings.Split(playerName, ".")
	if len(array) == 2 {
		return array[0], array[1], nil
	}
	return "", "", errors.New(fmt.Sprintf("解析玩家名字失败:%s", playerName))
}

type ChargeTaskDistribution struct {
	TaskId int     `json:"taskId"`
	Count  int     `json:"count"`
	Rate   float32 `json:"rate"`
}
type ChargeTaskDistributionQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string    `json:"serverId"`
	ChannelList [] string `json:"channelList"`
	StartTime   int
	EndTime     int
	IsFirst     int
}

// 获取充值任务分布
func GetChargeTaskDistribution(params ChargeTaskDistributionQueryParam) [] *ChargeTaskDistribution {
	data := make([] *ChargeTaskDistribution, 0)
	whereArray := make([] string, 0)
	whereArray = append(whereArray, fmt.Sprintf("charge_type = 99"))
	whereArray = append(whereArray, fmt.Sprintf(" part_id = '%s'", params.PlatformId))
	whereArray = append(whereArray, fmt.Sprintf(" channel in (%s) ", GetSQLWhereParam(params.ChannelList)))
	if params.IsFirst == 1 {
		whereArray = append(whereArray, fmt.Sprintf("is_first = 1"))
	}
	if params.ServerId != "" {
		whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", params.ServerId))
	}
	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	sql := fmt.Sprintf(
		`SELECT curr_task_id as task_id, count(*) as count FROM charge_info_record %s group by task_id `, whereParam)
	err := DbCharge.Raw(sql).Find(&data).Error
	utils.CheckError(err)

	//if params.Node == "" {
	//	sql := fmt.Sprintf(
	//		`SELECT curr_task_id as task_id, count(*) as count FROM charge_info_record where charge_type = 99 group by task_id `)
	//	err := DbCharge.Raw(sql).Find(&data).Error
	//	utils.CheckError(err)
	//} else {
	//	serverIdList := GetGameServerIdListStringByNode(params.Node)
	//	sql := fmt.Sprintf(
	//		`SELECT curr_task_id as task_id, count(*) as count FROM charge_info_record where server_id in (%s) and charge_type = 99 group by task_id `, serverIdList)
	//	err := DbCharge.Raw(sql).Find(&data).Error
	//	utils.CheckError(err)
	//}

	sum := 0
	for _, e := range data {
		sum += e.Count
	}
	if sum > 0 {
		for _, e := range data {
			e.Rate = float32(e.Count) / float32(sum) * 100
		}
	}

	return data
}

type ChargeActivityDistribution struct {
	ChargeItemId int     `json:"chargeItemId"`
	Count        int     `json:"count"`
	Rate         float32 `json:"rate"`
	Money        float32 `json:"money"`
	MoneyRate    float32 `json:"moneyRate"`
}
type ChargeActivityDistributionQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string    `json:"serverId"`
	ChannelList [] string `json:"channelList"`
	StartTime   int
	EndTime     int
	IsFirst     int
}

// 获取充值任务分布
func GetChargeActivityDistribution(params ChargeActivityDistributionQueryParam) [] *ChargeActivityDistribution {
	data := make([] *ChargeActivityDistribution, 0)
	whereArray := make([] string, 0)
	whereArray = append(whereArray, fmt.Sprintf("charge_type = 99"))
	whereArray = append(whereArray, fmt.Sprintf(" part_id = '%s'", params.PlatformId))
	whereArray = append(whereArray, fmt.Sprintf(" channel in (%s) ", GetSQLWhereParam(params.ChannelList)))
	if params.IsFirst == 1 {
		whereArray = append(whereArray, fmt.Sprintf("is_first = 1"))
	}
	if params.StartTime > 0 {
		whereArray = append(whereArray, fmt.Sprintf("record_time between %d and %d", params.StartTime, params.EndTime))
	}
	if params.ServerId != "" {
		whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", params.ServerId))
	}
	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	sql := fmt.Sprintf(
		`SELECT charge_item_id, count(*) as count, sum(money) as money FROM charge_info_record %s group by charge_item_id `, whereParam)
	err := DbCharge.Raw(sql).Find(&data).Error
	utils.CheckError(err)

	//if params.Node == "" {
	//	sql := fmt.Sprintf(
	//		`SELECT curr_task_id as task_id, count(*) as count FROM charge_info_record where charge_type = 99 group by task_id `)
	//	err := DbCharge.Raw(sql).Find(&data).Error
	//	utils.CheckError(err)
	//} else {
	//	serverIdList := GetGameServerIdListStringByNode(params.Node)
	//	sql := fmt.Sprintf(
	//		`SELECT curr_task_id as task_id, count(*) as count FROM charge_info_record where server_id in (%s) and charge_type = 99 group by task_id `, serverIdList)
	//	err := DbCharge.Raw(sql).Find(&data).Error
	//	utils.CheckError(err)
	//}

	sum := 0
	var moneySum float32
	for _, e := range data {
		sum += e.Count
		moneySum += e.Money
	}
	if sum > 0 {
		for _, e := range data {
			e.Rate = float32(e.Count) / float32(sum) * 100
		}
	}
	if moneySum > 0 {
		for _, e := range data {
			e.MoneyRate = float32(e.Money) / float32(moneySum) * 100
		}
	}
	return data
}

type ChargeMoneyDistribution struct {
	ValueString string  `json:"valueString"`
	Count       int     `json:"count"`
	Rate        float32 `json:"rate"`
	Min         int     `json:"-"`
	Max         int     `json:"-"`
}
type ChargeMoneyDistributionQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string    `json:"serverId"`
	ChannelList [] string `json:"channelList"`
	StartTime   int
	EndTime     int
}

// 获取充值金额分布
func GetChargeMoneyDistribution(params ChargeMoneyDistributionQueryParam) [] *ChargeMoneyDistribution {
	var data = [] *ChargeMoneyDistribution{
		&ChargeMoneyDistribution{
			Min:         0,
			Max:         1,
			ValueString: "1元",
		},
		&ChargeMoneyDistribution{
			Min:         1,
			Max:         5,
			ValueString: "1~5元",
		},
		&ChargeMoneyDistribution{
			Min:         5,
			Max:         10,
			ValueString: "5~10元",
		},
		&ChargeMoneyDistribution{
			Min:         10,
			Max:         20,
			ValueString: "10~20元",
		},
		&ChargeMoneyDistribution{
			Min:         20,
			Max:         50,
			ValueString: "20~50元",
		},
		&ChargeMoneyDistribution{
			Min:         50,
			Max:         100,
			ValueString: "50~100",
		},
		&ChargeMoneyDistribution{
			Min:         100,
			Max:         200,
			ValueString: "100~200元",
		},
		&ChargeMoneyDistribution{
			Min:         200,
			Max:         500,
			ValueString: "200~500元",
		},
		&ChargeMoneyDistribution{
			Min:         500,
			Max:         1000,
			ValueString: "500~1000元",
		},
		&ChargeMoneyDistribution{
			Min:         1000,
			Max:         2000,
			ValueString: "1000~2000元",
		},
		&ChargeMoneyDistribution{
			Min:         2000,
			Max:         5000,
			ValueString: "2000~5000元",
		},
		&ChargeMoneyDistribution{
			Min:         5000,
			Max:         10000,
			ValueString: "5000~10000元",
		},
		&ChargeMoneyDistribution{
			Min:         10000,
			Max:         20000,
			ValueString: "1万~2万元",
		},
		&ChargeMoneyDistribution{
			Min:         20000,
			Max:         50000,
			ValueString: "2万~5万元",
		},
		&ChargeMoneyDistribution{
			Min:         50001,
			Max:         100000,
			ValueString: "5万~10万元",
		},
		&ChargeMoneyDistribution{
			Min:         100001,
			Max:         100000000000,
			ValueString: "大于10万元",
		},
	}
	maxCount := 0
	for _, e := range data {
		whereArray := make([] string, 0)
		whereArray = append(whereArray, fmt.Sprintf(" part_id = '%s'", params.PlatformId))
		whereArray = append(whereArray, fmt.Sprintf(" total_money > %d", e.Min))
		whereArray = append(whereArray, fmt.Sprintf(" total_money <= %d", e.Max))
		whereArray = append(whereArray, fmt.Sprintf(" channel in (%s) ", GetSQLWhereParam(params.ChannelList)))
		if params.ServerId != "" {
			whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", params.ServerId))
		}
		whereParam := strings.Join(whereArray, " and ")
		if whereParam != "" {
			whereParam = " where " + whereParam
		}

		sql := fmt.Sprintf(`SELECT count(*) as count FROM player_charge_info_record %s`, whereParam)
		err := DbCharge.Raw(sql).Find(&e).Error
		utils.CheckError(err)
	}
	for _, e := range data {
		maxCount += e.Count
	}
	if maxCount > 0 {
		for _, e := range data {
			e.Rate = float32(e.Count) / float32(maxCount) * 100
		}
	}

	return data
}

type ChargeLevelDistribution struct {
	LevelString string  `json:"levelString"`
	Count       int     `json:"count"`
	Rate        float32 `json:"rate"`
	Min         int     `json:"-"`
	Max         int     `json:"-"`
}
type ChargeLevelDistributionQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string    `json:"serverId"`
	ChannelList [] string `json:"channelList"`
	StartTime   int
	EndTime     int
	IsFirst     int
}

// 获取充值等级分布
func GetChargeLevelDistribution(params ChargeLevelDistributionQueryParam) [] *ChargeLevelDistribution {
	var data = [] *ChargeLevelDistribution{
		&ChargeLevelDistribution{
			Min:         1,
			Max:         5,
			LevelString: "1~5级",
		},
		&ChargeLevelDistribution{
			Min:         6,
			Max:         10,
			LevelString: "6~10级",
		},
		&ChargeLevelDistribution{
			Min:         11,
			Max:         20,
			LevelString: "11~20级",
		},
		&ChargeLevelDistribution{
			Min:         21,
			Max:         30,
			LevelString: "21~30级",
		},
		&ChargeLevelDistribution{
			Min:         31,
			Max:         40,
			LevelString: "31~40级",
		},
		&ChargeLevelDistribution{
			Min:         41,
			Max:         50,
			LevelString: "41~50级",
		},
		&ChargeLevelDistribution{
			Min:         51,
			Max:         60,
			LevelString: "51~60级",
		},
		&ChargeLevelDistribution{
			Min:         61,
			Max:         70,
			LevelString: "61~70级",
		},
		&ChargeLevelDistribution{
			Min:         71,
			Max:         80,
			LevelString: "71~80级",
		},
		&ChargeLevelDistribution{
			Min:         81,
			Max:         90,
			LevelString: "81~90级",
		},
		&ChargeLevelDistribution{
			Min:         91,
			Max:         100,
			LevelString: "91~100级",
		},
		&ChargeLevelDistribution{
			Min:         101,
			Max:         110,
			LevelString: "101~110级",
		},
		&ChargeLevelDistribution{
			Min:         111,
			Max:         120,
			LevelString: "111~120级",
		}, &ChargeLevelDistribution{
			Min:         121,
			Max:         130,
			LevelString: "121~130级",
		},
		&ChargeLevelDistribution{
			Min:         131,
			Max:         140,
			LevelString: "131~140级",
		},
		&ChargeLevelDistribution{
			Min:         141,
			Max:         150,
			LevelString: "141~150级",
		},
		&ChargeLevelDistribution{
			Min:         151,
			Max:         160,
			LevelString: "151~160级",
		},
		&ChargeLevelDistribution{
			Min:         161,
			Max:         170,
			LevelString: "161~170级",
		},
		&ChargeLevelDistribution{
			Min:         171,
			Max:         180,
			LevelString: "171~180级",
		},
		&ChargeLevelDistribution{
			Min:         181,
			Max:         190,
			LevelString: "181~190级",
		},
		&ChargeLevelDistribution{
			Min:         191,
			Max:         200,
			LevelString: "191~200级",
		},
		&ChargeLevelDistribution{
			Min:         201,
			Max:         250,
			LevelString: "201~250级",
		},
		&ChargeLevelDistribution{
			Min:         251,
			Max:         300,
			LevelString: "251~300级",
		},
		&ChargeLevelDistribution{
			Min:         301,
			Max:         400,
			LevelString: "301~400级",
		},
		&ChargeLevelDistribution{
			Min:         401,
			Max:         500,
			LevelString: "401~500级",
		},
		&ChargeLevelDistribution{
			Min:         501,
			Max:         10000,
			LevelString: "大于501级",
		},
	}
	maxCount := 0
	for _, e := range data {
		whereArray := make([] string, 0)
		whereArray = append(whereArray, fmt.Sprintf(" curr_level >= %d", e.Min))
		whereArray = append(whereArray, fmt.Sprintf(" curr_level <= %d", e.Max))
		whereArray = append(whereArray, fmt.Sprintf("charge_type = 99"))
		whereArray = append(whereArray, fmt.Sprintf(" part_id = '%s'", params.PlatformId))
		whereArray = append(whereArray, fmt.Sprintf(" channel in (%s) ", GetSQLWhereParam(params.ChannelList)))
		if params.IsFirst == 1 {
			whereArray = append(whereArray, fmt.Sprintf("is_first = 1"))
		}
		if params.ServerId != "" {
			whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", params.ServerId))
		}
		whereParam := strings.Join(whereArray, " and ")
		if whereParam != "" {
			whereParam = " where " + whereParam
		}
		sql := fmt.Sprintf(`SELECT count(*) as count FROM charge_info_record %s `, whereParam)
		err := DbCharge.Raw(sql).Find(&e).Error
		utils.CheckError(err)
	}
	for _, e := range data {
		maxCount += e.Count
	}
	if maxCount > 0 {
		for _, e := range data {
			e.Rate = float32(e.Count) / float32(maxCount) * 100
		}
	}

	return data
}

//获取该ip节点数量
func GetIpNodeCount(ip string) int {
	l := GetAllServerNodeList()
	count := 0
	for _, e := range l {
		thisIp := strings.Split(e.Node, "@")[1]
		//logs.Debug("thisIp:%+v", thisIp)
		if thisIp == ip {
			count ++
		}
	}
	return count
}

//获取该ip在线人数
func GetIpOnlinePlayerCount(ip string) int {
	l := GetAllServerNodeList()
	count := 0
	for _, e := range l {
		thisIp := strings.Split(e.Node, "@")[1]
		//logs.Debug("thisIp:%+v", thisIp)
		if thisIp == ip {
			gameDb, err := GetGameDbByNode(e.Node)
			utils.CheckError(err)
			if err != nil {
				continue
			}
			defer gameDb.Close()
			count += GetNowOnlineCount(gameDb)
		}
	}
	return count
}

// 向中心服添加game_server
func AddGameServer(PlatformId string, sid string, desc string, node string, zoneNode string, state int, openTime int, isShow int) (string, error) {
	logs.Info("向中心服添加game_server:%v", PlatformId, sid, desc, node, zoneNode, state, openTime, isShow)
	out, err := utils.CenterNodeTool(
		"mod_server_mgr",
		"add_game_server",
		PlatformId,
		sid,
		desc,
		node,
		zoneNode,
		strconv.Itoa(state),
		strconv.Itoa(openTime),
		strconv.Itoa(isShow),
	)
	return out, err
}

func AddServerNode(node string, ip string, port int, webPort int, serverType int, platformId string, dbHost string, dbPort int, dbName string) (string, error) {
	logs.Info("向中心服添加server_node:%v", node, ip, port, webPort, serverType, platformId, dbHost, dbPort, dbName)
	out, err := utils.CenterNodeTool(
		"mod_server_mgr",
		"add_server_node",
		node,
		ip,
		strconv.Itoa(port),
		strconv.Itoa(webPort),
		strconv.Itoa(serverType),
		platformId,
		dbHost,
		strconv.Itoa(dbPort),
		dbName,
	)
	return out, err
}

func InstallNode(node string) error {
	logs.Info("开始部署节点:%s......", node)
	var commandArgs []string
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	app := ""
	switch serverNode.Type {
	case 0:
		app = "center"
	case 1:
		app = "game"
	case 2:
		app = "zone"
	case 4:
		app = "login_server"
	case 5:
		app = "unique_id"
	case 6:
		app = "charge"
	case 7:
		app = "war"
	case 8:
		app = "web"
	}

	version := ""
	if serverNode.PlatformId == "" {
		logs.Warning("部署节点(%s)没有对应平台， 取默认版本库server!!!!!!!!!!!!!!!!!", node)
		version = "server"
	} else {
		platform, err := GetPlatformOne(serverNode.PlatformId)
		utils.CheckError(err)
		version = platform.Version
	}
	logs.Info("版本库:%s", version)
	commandArgs = []string{"/data/tool/ansible/do-install.sh", serverNode.Node, app, serverNode.DbName, serverNode.DbHost, strconv.Itoa(serverNode.DbPort), "root", version}
	out, err := utils.Cmd("sh", commandArgs)
	utils.CheckError(err, fmt.Sprintf("部署节点失败:%v %v", node, out))
	if err != nil {
		return err
	}
	logs.Info("部署节点成功:%v!!!", node)
	return nil
}

func NodeAction(nodes [] string, action string) error {
	logs.Info("节点操作:nodes->%v, action->%v", nodes, action)
	curDir := utils.GetCurrentDirectory()
	defer os.Chdir(curDir)
	toolDir := utils.GetToolDir()
	err := os.Chdir(toolDir)
	utils.CheckError(err)
	if err != nil {
		return err
	}

	var commandArgs []string
	for _, node := range nodes {
		switch action {
		case "start":
			commandArgs = []string{"node_tool.sh", node, action,}
		case "stop":
			commandArgs = []string{"node_tool.sh", node, action,}
		case "pull":
			commandArgs = []string{"node_tool.sh", node, action,}
		case "hot_reload":
			commandArgs = []string{"node_hot_reload.sh", node, "server",}
		case "cold_reload":
			commandArgs = []string{"node_cold_reload.sh", node, "server",}
		}
		out, err := utils.Cmd("sh", commandArgs)
		utils.CheckError(err, fmt.Sprintf("操作节点失败:%v %v", action, out))
		if err != nil {
			return err
		}
	}
	logs.Info("节点操作成功:nodes->%v, action->%v!", nodes, action)
	return nil
}

func AfterAddGameServer() error {
	out, err := utils.CenterNodeTool(
		"mod_server_sync",
		"after_add_game_node",
	)
	utils.CheckError(err, out)
	return err
}

func RefreshGameServer() error {
	out, err := utils.CenterNodeTool(
		"mod_server_sync",
		"push_all_login_server_node",
	)
	utils.CheckError(err, out)
	return err
	//var result struct {
	//	ErrorCode int
	//}
	//logs.Info("刷新区服入口:...")
	//
	//data := fmt.Sprintf("time=%d", utils.GetTimestamp())
	//sign := utils.String2md5(data + enums.GmSalt)
	//base64Data := base64.URLEncoding.EncodeToString([]byte(data))
	//
	//baseUrl := beego.AppConfig.String("login_server" + "::url")
	//url := fmt.Sprintf("%s?data=%s&sign=%s", baseUrl, base64Data, sign)
	////url := "http://192.168.31.100:16667/refresh?" + "data=" + base64Data+ "&sign=" + sign
	//
	//logs.Info("url:%s", url)
	//resp, err := http.Get(url)
	//
	//utils.CheckError(err)
	//if err != nil {
	//	return err
	//}
	//
	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//
	//utils.CheckError(err)
	//if err != nil {
	//	return err
	//}
	//
	//logs.Info("刷新区服入口 result:%v", string(body))
	//
	//err = json.Unmarshal(body, &result)
	//utils.CheckError(err)
	//if err != nil {
	//	return err
	//}
	//
	//if result.ErrorCode != 0 {
	//	logs.Error("刷新区服入口失败!!!!")
	//	return errors.New("刷新区服入口失败!!!!")
	//}
	//logs.Info("刷新区服入口成功")
	//return nil
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

//func S() {
//	logs.Info("开始统计")
//	gameServerList, _ := GetAllGameServer()
//	context := ""
//	for _, e := range gameServerList {
//		var data [] struct {
//			PlayerId int
//		}
//		gameDb, err := GetGameDbByNode(e.Node)
//		utils.CheckError(err)
//		if err != nil {
//			return
//		}
//		defer gameDb.Close()
//		sql := fmt.Sprintf(
//			`select a.player_id as player_id,b.acc_id from player_platform_award AS a,player AS b where a.player_id = b.id and a.id = 1601`)
//		err = gameDb.Raw(sql).Find(&data).Error
//		utils.CheckError(err)
//		for _, i := range data {
//			player, err := GetPlayerOne(e.PlatformId, e.Sid, i.PlayerId)
//			utils.CheckError(err)
//			is2 := IsThatDayPlayerLogin(gameDb, utils.GetThatZeroTimestamp(int64(player.RegTime))+86400, player.Id)
//			is3 := IsThatDayPlayerLogin(gameDb, utils.GetThatZeroTimestamp(int64(player.RegTime))+86400*2, player.Id)
//			var moneyList [] struct {
//				Money        int
//				ChargeItemId int
//			}
//			sql := fmt.Sprintf(
//				`select money , charge_item_id from charge_info_record where player_id = %d;`, i.PlayerId)
//			err = DbCharge.Raw(sql).Find(&moneyList).Error
//			m := make([] string, 0)
//			//logs.Info("moneyList:%+v", moneyList)
//			for _, e := range moneyList {
//				m = append(m, strconv.Itoa(e.Money))
//			}
//
//			m1 := make([] string, 0)
//			//logs.Info("moneyList:%+v", moneyList)
//			for _, e := range moneyList {
//				m1 = append(m1, strconv.Itoa(e.ChargeItemId))
//			}
//			context += fmt.Sprintf("%s, %d, %s, %t, %t, [%s], [%s]\n", e.Sid, player.Id, player.AccId, is2, is3, strings.Join(m, " "), strings.Join(m1, " "))
//		}
//	}
//	utils.FilePutContext("data.txt", context)
//	logs.Info("统计完毕")
//}

// 获取服务器整形数据
func GetServerDataInt(db *gorm.DB, serverDataId int) int {
	var data struct {
		Data int
	}
	sql := fmt.Sprintf(
		`SELECT int_data as data FROM server_data WHERE id =  %d`, serverDataId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Data
}

// 获取服务器字符串数据
func GetServerDataStr(db *gorm.DB, serverDataId int) string {
	var data struct {
		Data string
	}
	sql := fmt.Sprintf(
		`SELECT str_data as data FROM server_data WHERE id =  %d`, serverDataId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Data
}

func RRRR() {
	logs.Info("go")
	accountList := make([] string, 0)
	var data [] struct {
		Id int
	}
	content := ""
	sql := fmt.Sprintf(
		`SELECT player_id as id FROM player_charge_info_record WHERE total_money  >=  200 and part_id = 'qq'`)
	err := DbCharge.Raw(sql).Find(&data).Error
	utils.CheckError(err)
	if err != nil {
		return
	}
	for _, e := range data {
		globalPlayer, err := GetGlobalPlayerOne(e.Id)
		utils.CheckError(err)
		if err != nil {
			return
		}
		if inArray(globalPlayer.Account, accountList) == false {
			accountList = append(accountList, globalPlayer.Account)
		}

	}
	for _, e := range accountList {

		content += fmt.Sprintf("%s\n", strings.Replace(strings.Replace(e, "ios_", "", -1), "android_", "", -1))
		//content += "nodes="
		//content += "\"[" + strings.Join(nodes, ", ") + "]\""
		//content += "\n\n"

	}
	err = utils.FilePutContext("account.txt", content)
	utils.CheckError(err)
	logs.Info("go finish")
}

func inArray(v string, array [] string) bool {
	for _, e := range array {
		if e == v {
			return true
		}
	}
	return false
}

func SSS() {
	logs.Info("go")
	var data [] struct {
		Id string
	}

	//content := ""
	sql := fmt.Sprintf(
		`select distinct account as id from  global_player`)
	err := DbCenter.Raw(sql).Find(&data).Error
	utils.CheckError(err)
	logs.Info("go1")
	if err != nil {
		return
	}
	args := make([] string, 0, len(data))
	for _, e := range data {
		args = append(args, e.Id)

	}
	logs.Info("go2")
	content := strings.Join(args, "\n")
	logs.Info("go3")
	err = utils.FilePutContext("all_account.txt", content)
	logs.Info("go4")
	utils.CheckError(err)
	logs.Info("go finish")
}

func QQQQ() {
	logs.Info("go")
	args := make([] string, 0, 100000)
	sql := fmt.Sprintf(
		`select acc_id from player , player_data where player.id = player_data.player_id and player_data.vip_level = 0 and player_data.level > 50 limit 1000;`)
	for i := 1; i <= 300; i++ {
		var data [] struct {
			AccId string
		}
		gameDb, err := GetGameDbByPlatformIdAndSid("wx", fmt.Sprintf("s%d", i))
		utils.CheckError(err)
		if err != nil {
			return
		}
		defer gameDb.Close()
		err = gameDb.Raw(sql).Find(&data).Error
		utils.CheckError(err)
		for _, e := range data {
			args = append(args, e.AccId)
			if len(args) == 100000 {
				break
			}
		}
		if len(args) == 100000 {
			break
		}

	}
	logs.Info("go2")
	content := strings.Join(args, "\n")
	logs.Info("go3")
	err := utils.FilePutContext("account_20181115.txt", content)
	logs.Info("go4")
	utils.CheckError(err)
	logs.Info("go finish")
}
