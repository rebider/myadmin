package models

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/chnzrb/myadmin/utils"
	"strconv"
	//"github.com/astaxie/beego/logs"
	//"github.com/zaaksam/dproxy/go/db"
	//"github.com/astaxie/beego/logs"
)

//获取总注册人数
func GetTotalCreateRole(db *gorm.DB) int {
	var count int
	db.Model(&Player{}).Count(&count)
	return count
}

//获取今日创角人数
func GetTodayCreateRole(db *gorm.DB) int {
	var data struct {
		Count int
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM player WHERE reg_time between ? and ?`)
	err := db.Raw(sql, todayZeroTimestamp, todayZeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	//logs.Info("ppp:%v,%v", gameServer.Node, data.Count)
	return data.Count
}

//获取某日创角人数
func GetThatDayCreateRole(db *gorm.DB, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM player WHERE reg_time between ? and ?`)
	err := db.Raw(sql, zeroTimestamp, zeroTimestamp+86400).Scan(&data).Error
	utils.CheckError(err)
	//logs.Info("ppp:%v,%v", gameServer.Node, data.Count)
	return data.Count
}

//获取某日创角的玩家Id列表
func GetThatDayCreateRolePlayerIdList(db *gorm.DB, zeroTimestamp int) [] int {
	var data [] struct {
		Id int
	}

	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
	sql := fmt.Sprintf(
		`SELECT id FROM player WHERE reg_time between ? and ?`)
	err := db.Raw(sql, zeroTimestamp, zeroTimestamp+86400).Find(&data).Error
	utils.CheckError(err)
	idList := make([] int, 0)
	for _, e := range data {
		idList = append(idList, e.Id)
	}
	//logs.Info("ppp:%v,%v", gameServer.Node, data.Count)
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
func GetMaxOnlineCount(platformId int, serverId string) int {
	gameServer, err := GetGameServerOne(platformId, serverId)
	utils.CheckError(err)
	var data struct {
		Count int
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
	sql := fmt.Sprintf(
		`SELECT max(online_num) as count FROM c_server_trace_log WHERE node = ? `)
	err = DbCenter.Raw(sql, gameServer.Node).Scan(&data).Error
	utils.CheckError(err)
	//logs.Info("ppp:%v,%v", gameServer.Node, data.Count)
	return data.Count
}

//获取当前在线人数
func GetMaxPlayerLevel(db *gorm.DB) int {
	var data struct {
		MaxLevel int
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
	sql := fmt.Sprintf(
		`SELECT max(level) as max_level FROM player_data `)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	//logs.Info("ppp:%v,%v", data.MaxLevel)
	return data.MaxLevel
}


// 获取当天 平均在线人数
func GetThatDayAverageOnlineCount(node string, zeroTimestamp int) float32 {
	var data struct {
		Count float32
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
	sql := fmt.Sprintf(
		`SELECT avg(online_num)  as count FROM c_server_trace_log where node = ? and time between ? and ? `)
	err := DbCenter.Raw(sql, node, zeroTimestamp, zeroTimestamp + 86400).Scan(&data).Error
	utils.CheckError(err)
	//logs.Info("ppp:%v,%v", data.MaxLevel)
	return data.Count
	//averageOnlineCount := 0
	//num := 0
	//for i := zeroTimestamp; i < zeroTimestamp+86400; i = i + 60*60 {
	//	cServerTraceLog := &CServerTraceLog{}
	//	err := DbCenter.Where(&CServerTraceLog{
	//		Node: node,
	//		Time: i,
	//	}).First(&cServerTraceLog).Error
	//	if err == nil {
	//		averageOnlineCount += cServerTraceLog.OnlineNum
	//		num += 1
	//	}
	//}
	////logs.Info("%v %v %v", num, averageOnlineCount, averageOnlineCount/num)
	//if num == 0 {
	//	return 0
	//}
	//return averageOnlineCount/num
}

// 获取当天 最高在线人数
func GetThatDayMaxOnlineCount(node string, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
	sql := fmt.Sprintf(
		`SELECT max(online_num)  as count FROM c_server_trace_log where node = ? and time between ? and ? `)
	err := DbCenter.Raw(sql, node, zeroTimestamp, zeroTimestamp + 86400).Scan(&data).Error
	utils.CheckError(err)
	//logs.Info("ppp:%v,%v", data.MaxLevel)
	return data.Count
}

// 获取当天 最低在线人数
func GetThatDayMinOnlineCount(node string, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
	sql := fmt.Sprintf(
		`SELECT min(online_num)  as count FROM c_server_trace_log where node = ? and time between ? and ? `)
	err := DbCenter.Raw(sql, node, zeroTimestamp, zeroTimestamp + 86400).Scan(&data).Error
	utils.CheckError(err)
	//logs.Info("ppp:%v,%v", data.MaxLevel)
	return data.Count
}

type RemainTask struct {
	TaskId int `json:"taskId"`
	Count int `json:"count"`
}

// 获取任务分布
func GetRemainTask(platformId int, serverId string) [] *RemainTask {
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
	gameDb, err:= GetGameDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	defer gameDb.Close()
	data := make([] *RemainTask, 0)
	sql := fmt.Sprintf(
		`SELECT task_id, count(*) as count FROM player_task group by task_id `)
	err = gameDb.Raw(sql).Find(&data).Error
	utils.CheckError(err)
	//logs.Info("ppp:%v,%v", data.MaxLevel)
	return data
}

func get24hoursOnlineCount(node string, zeroTimestamp int) [] string {
	onlineCountList := make([] string, 0)
	for i := zeroTimestamp; i < zeroTimestamp+86400; i = i + 10*60 {
		cServerTraceLog := &CServerTraceLog{}
		err := DbCenter.Where(&CServerTraceLog{
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

//获取当前在线人数
func GetPlayerName(db *gorm.DB, playerId int) string {
	var data struct {
		Name string
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)

	sql := fmt.Sprintf(
		`SELECT nickname as name FROM player where id = ? `)
	//logs.Info("GetPlayerName:%v", playerId)
	err := db.Raw(sql, playerId).Scan(&data).Error
	utils.CheckError(err)
	//logs.Info("ppp:%v,%v", data.MaxLevel)
	return data.Name
}
