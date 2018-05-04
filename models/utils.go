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

// 获取该天登录次数
func GetThatDayLoginTimes(db *gorm.DB, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM player_login_log WHERE timestamp between ? and ?`)
	err := db.Raw(sql,zeroTimestamp, zeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

// 获取该天登录的玩家数量
func GetThatDayLoginPlayerCount(db *gorm.DB, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM player WHERE last_login_time between ? and ?`)
	err := db.Raw(sql,zeroTimestamp, zeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

// 获取该天活跃玩家数量
func GetThatDayActivePlayerCount(db *gorm.DB, zeroTimestamp int) int {
	count := 0
	data := make( [] *Player, 0)
	sql := fmt.Sprintf(
		`SELECT * FROM player WHERE last_login_time between ? and ?`)
	err := db.Raw(sql,zeroTimestamp, zeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		if e.LoginTimes >= (zeroTimestamp+86400 - e.RegTime) / 86400 {
			count ++
		}
	}
	return count
}

//获取当前在线人数
func GetNowOnlineCount(db *gorm.DB) int {
	var count int
	db.Model(&Player{}).Where(&Player{IsOnline: 1}).Count(&count)
	return count
}

//获取当前在线ip数
func GetNowOnlineIpCount(db *gorm.DB) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT COUNT( DISTINCT last_login_ip ) as count FROM player where is_online = 1;`)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
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
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	defer gameDb.Close()
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
		`SELECT player.id, player.is_online, player_task.task_id FROM player left join player_task on player.id = player_task.player_id`)
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
	totalCreateRole := GetTotalCreateRoleCount(gameDb)
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
func GetRemainLevel(node string) [] *RemainLevel {
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
		`SELECT player.id, player.is_online, player_data.level FROM player left join player_data on player.id = player_data.player_id`)
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
	totalCreateRole := GetTotalCreateRoleCount(gameDb)
	for _, e := range data {
		elementList := make([] *Element, 0)
		sql := fmt.Sprintf(
			`SELECT is_online, total_online_time FROM player where total_online_time >= ? and total_online_time < ? `)
		err = gameDb.Raw(sql, e.StartTime, e.EndTime).Find(&elementList).Error
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
	if err != nil {
		return  0
	}
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
	if totalChargePlayerCount > 0 {
		return float32(totalChargeValueNum) / float32(totalChargePlayerCount)
	}
	return float32(0)
}

//付费率
func CaclChargeRate(totalChargePlayerCount int, totalRoleCount int) float32 {
	if totalRoleCount == 0 {
		return 0
	}
	if totalRoleCount > 0 {
		return float32(totalChargePlayerCount) / float32(totalRoleCount)
	}
	return float32(0)
}

//二次付费率
func CaclSceondChargeRate(secondChargePlayerCount int, totalChargePlayerCount int) float32 {
	if totalChargePlayerCount == 0 {
		return 0
	}
	if totalChargePlayerCount > 0 {
		return float32(secondChargePlayerCount) / float32(totalChargePlayerCount)
	}
	return float32(0)
}

//获取区服二次付费人数
func GetServerSecondChargePlayerCount(node string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where server_id in (%s) and is_first = 0 and charge_type = 99;`, GetGameServerIdListStringByNode(node))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服首次付费人数
func GetThadDayServerFirstChargePlayerCount(node string, time int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where node = ? and is_first = 1 and charge_type = 99 and (record_time between ? and ?);`)
	err := DbCharge.Raw(sql, node, time, time+86400).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服总充值元宝
func GetServerTotalChargeIngot(node string) int {
	var data struct {
		Count int
	}
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

type ChargeTaskDistribution struct {
	TaskId int     `json:"taskId"`
	Count  int     `json:"count"`
	Rate   float32 `json:"rate"`
}
type ChargeTaskDistributionQueryParam struct {
	BaseQueryParam
	PlatformId int
	Node       string `json:"serverId"`
	StartTime  int
	EndTime    int
	IsFirst    int
}

// 获取充值任务分布
func GetChargeTaskDistribution(params ChargeTaskDistributionQueryParam) [] *ChargeTaskDistribution {
	data := make([] *ChargeTaskDistribution, 0)
	whereArray := make([] string, 0)
	whereArray = append(whereArray, fmt.Sprintf("charge_type = 99"))
	whereArray = append(whereArray, fmt.Sprintf(" part_id = %d", params.PlatformId))
	if params.IsFirst == 1 {
		whereArray = append(whereArray, fmt.Sprintf("is_first = 1"))
	}
	if params.Node != "" {
		whereArray = append(whereArray, fmt.Sprintf("server_id in (%s)", GetGameServerIdListStringByNode(params.Node)))
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

type ChargeMoneyDistribution struct {
	ValueString string  `json:"valueString"`
	Count       int     `json:"count"`
	Rate        float32 `json:"rate"`
	Min         int     `json:"-"`
	Max         int     `json:"-"`
}
type ChargeMoneyDistributionQueryParam struct {
	BaseQueryParam
	PlatformId int
	Node       string `json:"serverId"`
	StartTime  int
	EndTime    int
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
		whereArray = append(whereArray, fmt.Sprintf(" part_id = %d", params.PlatformId))
		whereArray = append(whereArray, fmt.Sprintf(" total_money > %d", e.Min*100))
		whereArray = append(whereArray, fmt.Sprintf(" total_money <= %d", e.Max*100))
		if params.Node != "" {
			whereArray = append(whereArray, fmt.Sprintf("server_id in (%s)", GetGameServerIdListStringByNode(params.Node)))
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
	PlatformId int
	Node       string `json:"serverId"`
	StartTime  int
	EndTime    int
	IsFirst    int
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
		whereArray = append(whereArray, fmt.Sprintf(" part_id = %d", params.PlatformId))
		if params.IsFirst == 1 {
			whereArray = append(whereArray, fmt.Sprintf("is_first = 1"))
		}
		if params.Node != "" {
			whereArray = append(whereArray, fmt.Sprintf("server_id in (%s)", GetGameServerIdListStringByNode(params.Node)))
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
