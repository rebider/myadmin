package models

import (
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/astaxie/beego/logs"
	"strconv"
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
	ForbidType      int    `json:"forbidType"`
	ForbidTime      int    `json:"forbidTime"`
	RegTime         int    `json:"regTime"`
	LastLoginTime   int    `json:"lastLoginTime"`
	LastOfflineTime int    `json:"lastOfflineTime"`
	LastLoginIp     string `json:"lastLoginIp"`
	IsOnline        int    `json:"isOnline"`
	//DisableChatTime int    `json:"disableChatTime"`
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

func GetPlayerByPlatformIdAndSidAndNickname(platformId int, serverId string, nickname string) (*Player, error) {
	if nickname == "" {
		return nil, errors.New("nickname must not empty!")
	}
	logs.Debug("nickname:%v", nickname)
	db, err := GetDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	player := &Player{}
	err = db.Where(&Player{ServerId: serverId, Nickname: nickname}).First(&player).Error
	if err != nil {
		return nil, err
	}
	return player, err
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

	//var count int
	//db.Model(&Player{}).Count(&count)
	serverGeneralize.TotalRegister = GetTotalRegister(db)
	//db.Model(&Player{}).Where(&Player{IsOnline: 1}).Count(&count)
	serverGeneralize.OnlineCount = GetNowOnlineCount(db)
	return serverGeneralize, err
}

//获取总注册人数
func GetTotalRegister(db *gorm.DB) int {
	var count int
	db.Model(&Player{}).Count(&count)
	return count
}

//获取当前在线人数
func GetNowOnlineCount(db *gorm.DB) int {
	var count int
	db.Model(&Player{}).Where(&Player{IsOnline: 1}).Count(&count)
	return count
}

//获取当前在线人数
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

type CServerTraceLog struct {
	Node      string
	Time      int
	OnlineNum int
}

type ServerOnlineStatistics struct {
	PlatformId                int    `json:"platformId"`
	ServerId                  string `json:"serverId"`
	TotalRegister             int    `json:"totalRegister"`
	TodayRegister             int    `json:"todayRegister"`
	OnlineCount               int    `json:"onlineCount"`
	MaxOnlineCount            int    `json:"maxOnlineCount"`
	AverageOnlineCount        int    `json:"averageOnlineCount"`
	TodayOnlineList           [] string `json:"todayOnlineList"`
	YesterdayOnlineList       [] string `json:"yesterdayOnlineList"`
	BeforeYesterdayOnlineList [] string `json:"beforeYesterdayOnlineList"`
}

func GetServerOnlineStatistics(platformId int, serverId string) (*ServerOnlineStatistics, error) {
	db, err := GetDbByPlatformIdAndSid(platformId, serverId)
	defer db.Close()
	utils.CheckError(err)
	gameServer, err := GetGameServerOne(platformId, serverId)
	utils.CheckError(err)
	//todayOnlineList := make([]string, 0)
	//yesterdayTodayOnlineList := make([]int, 0)
	//beforeYesterdayTodayOnlineList := make([]int, 0)
	todayZeroTimestamp := int(utils.GetTodayZeroTimestamp())
	yesterdayZeroTimestamp := todayZeroTimestamp - 86400
	beforeYesterdayZeroTimestamp := yesterdayZeroTimestamp - 86400
	serverOnlineStatistics := &ServerOnlineStatistics{
		PlatformId:                platformId,
		ServerId:                  serverId,
		TotalRegister:             GetTotalRegister(db),
		OnlineCount:               GetNowOnlineCount(db),
		MaxOnlineCount:            GetMaxOnlineCount(platformId, serverId),
		TodayOnlineList:           get24hoursOnlineCount(gameServer.Node, todayZeroTimestamp),
		YesterdayOnlineList:       get24hoursOnlineCount(gameServer.Node, yesterdayZeroTimestamp),
		BeforeYesterdayOnlineList: get24hoursOnlineCount(gameServer.Node, beforeYesterdayZeroTimestamp),
	}
	return serverOnlineStatistics, nil
}

func get24hoursOnlineCount(node string, zeroTimestamp int) [] string {
	onlineCountList := make([] string, 0)
	for i := zeroTimestamp; i < zeroTimestamp + 86400; i = i + 10 * 60 {
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
	logs.Info("%+v", len(onlineCountList))
	return onlineCountList
}
