package models

import (
	"github.com/chnzrb/myadmin/utils"
	"fmt"
)

type PlayerQueryParam struct {
	BaseQueryParam
	Account    string
	Ip         string
	PlayerId   string
	Nickname   string
	IsOnline   string
	PlatformId int
	ServerId   string
}

type Player struct {
	Id              int    `json:"id"`
	AccId           string `json:"accId"`
	Nickname        string `json:"nickname"`
	Sex             int    `json:"sex"`
	ServerId        string `json:"serverId"`
	DisableLogin    int    `json:"disableLogin"`
	RegTime         int    `json:"regTime"`
	LastLoginTime   int    `json:"lastLoginTime"`
	LastOfflineTime int    `json:"lastOfflineTime"`
	LastLoginIp     string `json:"lastLoginIp"`
	IsOnline        int    `json:"isOnline"`
	DisableChatTime int    `json:"disableChatTime"`
}

func (a *Player) TableName() string {
	return "player"
}

//获取玩家列表
func GetPlayerList(params *PlayerQueryParam) ([]*Player, int64) {
	db, err := GetDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	defer db.Close()
	data := make([]*Player, 0)
	var count int64
	sortOrder := "id"
	switch params.Sort {
	case "id":
		sortOrder = "id"
	case "lastLoginTime":
		sortOrder = "last_login_time"
	}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}

	if params.Account != "" {
		db = db.Where("acc_id = ?", params.Account)
	}
	if params.Ip != "" {
		db = db.Where("last_login_ip = ?", params.Ip)
	}
	if params.Nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+params.Nickname+"%")
	}
	if params.IsOnline != "" {
		db = db.Where("is_online = ?", params.IsOnline)
	}
	if params.PlayerId != "" {
		db = db.Where("id = ?", params.PlayerId)
	}
	db.Model(&Player{}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	return data, count
}

// 获取单个玩家
func GetPlayerOne(platformId int, serverId string, id int) (*Player, error) {
	db, err := GetDbByPlatformIdAndSid(platformId, serverId)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	player := &Player{
		Id: id,
	}
	err = db.First(&player).Error
	return player, err
}


type PlayerDetail struct {
	Player
	VipLevel         int `json:"vipLevel"`
	Exp              int `json:"exp"`
	Level            int `json:"level"`
	TitleId          int
	Attack           int `json:"attack"`
	MaxHp            int `json:"maxHp"`
	Defense          int `json:"defense"`
	Hit              int `json:"hit"`
	Dodge            int `json:"dodge"`
	Critical         int `json:"critical"`
	Power            int `json:"power"`
	LastWorldSceneId int `json:"lastWorldSceneId"`
	PlayerPropList   [] *PlayerProp
}

type PlayerProp struct {
	PropType int
	PropId   int
	Num      int
}

func GetPlayerPropList(platformId int, serverId string, playerId int) ([]*PlayerProp, error) {
	db, err := GetDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	defer db.Close()
	playerPropList := make([]*PlayerProp, 0)
	sql := fmt.Sprintf(
		`SELECT * FROM player_prop WHERE player_id = ? `)
	err = db.Raw(sql, playerId).Scan(&playerPropList).Error

	return playerPropList, err
}

type PlayerLoginLog struct {
	Id        int    `json:"id"`
	PlayerId  int    `json:"playerId"`
	Ip        string `json:"ip"`
	Timestamp int    `json:"time"`
}

type PlayerLoginLogQueryParam struct {
	BaseQueryParam
	PlatformId int
	ServerId   string
	Ip         string
	PlayerId   string
	StartTime  int
	EndTime    int
}

func GetPlayerLoginLogList(params *PlayerLoginLogQueryParam) ([]*PlayerLoginLog, int64) {
	db, err := GetDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	defer db.Close()
	data := make([]*PlayerLoginLog, 0)
	var count int64
	sortOrder := "id"
	//switch params.Sort {
	//case "id":
	//	sortOrder = "id"
	//case "lastLoginTime":
	//	sortOrder = "last_login_time"
	//}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	if params.Ip != "" {
		db = db.Where("ip = ?", params.Ip)
	}
	if params.PlayerId != "" {
		db = db.Where("player_id = ?", params.PlayerId)
	}
	if params.StartTime != 0 {
		db = db.Where("timestamp >= ?", params.StartTime)
	}
	if params.EndTime != 0 {
		db = db.Where("timestamp <= ?", params.EndTime)
	}
	db.Model(&PlayerLoginLog{}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	return data, count
}

func GetPlayerDetail(platformId int, serverId string, playerId int) (*PlayerDetail, error) {
	db, err := GetDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	defer db.Close()
	playerDetail := &PlayerDetail{}

	//m := make(map[interface{}]interface{}, 0)
	sql := fmt.Sprintf(
		`SELECT player.*, player_data.* FROM player LEFT JOIN player_data on player.id = player_data.player_id WHERE player.id = ? `)
	err = db.Raw(sql, playerId).Scan(&playerDetail).Error

	playerDetail.PlayerPropList, err = GetPlayerPropList(platformId, serverId, playerId)
	utils.CheckError(err)
	return playerDetail, err
}

type PlayerOnlineLog struct {
	Id          int `json:"id"`
	PlayerId    int `json:"playerId"`
	LoginTime   int `json:"loginTime"`
	OfflineTime int `json:"offlineTime"`
	OnlineTime  int `json:"onlineTime"`
}

type PlayerOnlineLogQueryParam struct {
	BaseQueryParam
	PlatformId int
	ServerId   string
	PlayerId   string
	StartTime  int
	EndTime    int
}

func GetPlayerOnlineLogList(params *PlayerOnlineLogQueryParam) ([]*PlayerOnlineLog, int64) {
	db, err := GetDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	defer db.Close()
	data := make([]*PlayerOnlineLog, 0)
	var count int64
	sortOrder := "id"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	if params.PlayerId != "" {
		db = db.Where("player_id = ?", params.PlayerId)
	}
	if params.StartTime != 0 {
		db = db.Where("offline_time >= ?", params.StartTime)
	}
	if params.EndTime != 0 {
		db = db.Where("offline_time <= ?", params.EndTime)
	}
	db.Model(&PlayerOnlineLog{}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	return data, count
}

type ServerGeneralize struct {
	PlatformId    int    `json:"platformId"`
	ServerId      string `json:"serverId"`
	OpenTime      int    `json:"openTime"`
	Version       string `json:"version"`
	TotalRegister int    `json:"totalRegister"`
	OnlineCount   int    `json:"onlineCount"`
	Status        int    `json:"status"`
}

type ServerGeneralizeQueryParam struct {
	PlatformId int
	ServerId   string
}

func GetServerGeneralize(platformId int, serverId string) (*ServerGeneralize, error) {
	db, err := GetDbByPlatformIdAndSid(platformId, serverId)
	defer db.Close()
	utils.CheckError(err)
	gameServer, err := GetGameServerOne(platformId, serverId)
	utils.CheckError(err)
	serverNode, err := GetServerNode(gameServer.Node)
	utils.CheckError(err)
	serverGeneralize := &ServerGeneralize{}
	serverGeneralize.PlatformId = platformId
	serverGeneralize.ServerId = serverId
	serverGeneralize.OpenTime = serverNode.OpenTime
	serverGeneralize.Version = serverNode.ServerVersion
	serverGeneralize.Status = serverNode.State

	var count int
	db.Model(&Player{}).Count(&count)
	serverGeneralize.TotalRegister = count
	db.Model(&Player{}).Where(&Player{IsOnline: 1}).Count(&count)
	serverGeneralize.OnlineCount = count
	return serverGeneralize, err
}

//func GetPlayerDetail(platformId int, serverId string, playerId int) (*map[string] string, error) {
//	db, err := GetDbByPlatformIdAndSid(platformId, serverId)
//	utils.CheckError(err)
//	defer db.Close()
//	//playerDetail := &PlayerDetail{}
//
//	//m := make(map[interface{}]interface{}, 0)
//	sql := fmt.Sprintf(
//		`SELECT player.*, player_data.* FROM player LEFT JOIN player_data on player.id = player_data.player_id WHERE player.id = ? `)
//	rows, err := db.Raw(sql, playerId).Rows()
//	utils.CheckError(err)
//	columns, _ := rows.Columns()
//	scanArgs := make([]interface{}, len(columns))
//	values := make([]interface{}, len(columns))
//	for i := range values {
//		scanArgs[i] = &values[i]
//	}
//	for rows.Next() {
//		err = rows.Scan(scanArgs...)
//		record := make(map[string] string)
//		for i, col := range values {
//			if col != nil {
//				logs.Debug("%+v", col)
//				logs.Debug("%+v", reflect.TypeOf(col))
//				switch reflect.TypeOf(col).String() {
//				case "int":
//					record[columns[i]] = strconv.Itoa(col.(int))
//				case "int32":
//					record[columns[i]] = strconv.Itoa(int(col.(int32)))
//				case "int64":
//					record[columns[i]] = strconv.Itoa(int(col.(int64)))
//				case "[]uint8":
//					record[columns[i]] = string(col.([]byte))
//				}
//
//			}
//		}
//		return &record, err
//	}
//	return nil, err
//}
