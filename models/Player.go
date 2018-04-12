package models

import (
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"errors"
	"github.com/astaxie/beego/logs"
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
}

func (a *Player) TableName() string {
	return "player"
}

//获取玩家列表
func GetPlayerList(params *PlayerQueryParam) ([]*Player, int64) {
	gameDb, err := GetGameDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	defer gameDb.Close()
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
		gameDb = gameDb.Where("acc_id = ?", params.Account)
	}
	if params.Ip != "" {
		gameDb = gameDb.Where("last_login_ip = ?", params.Ip)
	}
	if params.Nickname != "" {
		gameDb = gameDb.Where("nickname LIKE ?", "%"+params.Nickname+"%")
	}
	if params.IsOnline != "" {
		gameDb = gameDb.Where("is_online = ?", params.IsOnline)
	}
	if params.PlayerId != "" {
		gameDb = gameDb.Where("id = ?", params.PlayerId)
	}
	gameDb.Model(&Player{}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	return data, count
}

// 获取单个玩家
func GetPlayerOne(platformId int, serverId string, id int) (*Player, error) {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	if err != nil {
		return nil, err
	}
	defer gameDb.Close()
	player := &Player{
		Id: id,
	}
	err = gameDb.First(&player).Error
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
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	defer gameDb.Close()
	playerPropList := make([]*PlayerProp, 0)
	sql := fmt.Sprintf(
		`SELECT * FROM player_prop WHERE player_id = ? `)
	err = gameDb.Raw(sql, playerId).Scan(&playerPropList).Error

	return playerPropList, err
}

func GetPlayerDetail(platformId int, serverId string, playerId int) (*PlayerDetail, error) {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	defer gameDb.Close()
	playerDetail := &PlayerDetail{}

	//m := make(map[interface{}]interface{}, 0)
	sql := fmt.Sprintf(
		`SELECT player.*, player_data.* FROM player LEFT JOIN player_data on player.id = player_data.player_id WHERE player.id = ? `)
	err = gameDb.Raw(sql, playerId).Scan(&playerDetail).Error

	playerDetail.PlayerPropList, err = GetPlayerPropList(platformId, serverId, playerId)
	utils.CheckError(err)
	return playerDetail, err
}

func GetPlayerByPlatformIdAndSidAndNickname(platformId int, serverId string, nickname string) (*Player, error) {
	if nickname == "" {
		return nil, errors.New("nickname must not empty!")
	}
	logs.Debug("nickname:%v", nickname)
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	defer gameDb.Close()
	player := &Player{}
	err = gameDb.Where(&Player{ServerId: serverId, Nickname: nickname}).First(&player).Error
	if err != nil {
		return nil, err
	}
	return player, err
}

type CServerTraceLog struct {
	Node      string
	Time      int
	OnlineNum int
}

type ServerOnlineStatistics struct {
	PlatformId                int       `json:"platformId"`
	ServerId                  string    `json:"serverId"`
	TodayCreateRole           int       `json:"todayCreateRole"`
	TodayRegister             int       `json:"todayRegister"`
	OnlineCount               int       `json:"onlineCount"`
	MaxOnlineCount            int       `json:"maxOnlineCount"`
	AverageOnlineCount        float32       `json:"averageOnlineCount"`
	TodayOnlineList           [] string `json:"todayOnlineList"`
	YesterdayOnlineList       [] string `json:"yesterdayOnlineList"`
	BeforeYesterdayOnlineList [] string `json:"beforeYesterdayOnlineList"`
}

func GetServerOnlineStatistics(platformId int, serverId string) (*ServerOnlineStatistics, error) {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	defer gameDb.Close()
	utils.CheckError(err)
	gameServer, err := GetGameServerOne(platformId, serverId)
	utils.CheckError(err)
	//todayOnlineList := make([]string, 0)
	//yesterdayTodayOnlineList := make([]int, 0)
	//beforeYesterdayTodayOnlineList := make([]int, 0)
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	yesterdayZeroTimestamp := todayZeroTimestamp - 86400
	beforeYesterdayZeroTimestamp := yesterdayZeroTimestamp - 86400
	serverOnlineStatistics := &ServerOnlineStatistics{
		PlatformId:                platformId,
		ServerId:                  serverId,
		TodayCreateRole:           GetTodayCreateRole(gameDb),
		TodayRegister:             GetTodayRegister(gameDb),
		OnlineCount:               GetNowOnlineCount(gameDb),
		MaxOnlineCount:            GetMaxOnlineCount(platformId, serverId),
		TodayOnlineList:           get24hoursOnlineCount(gameServer.Node, todayZeroTimestamp),
		YesterdayOnlineList:       get24hoursOnlineCount(gameServer.Node, yesterdayZeroTimestamp),
		BeforeYesterdayOnlineList: get24hoursOnlineCount(gameServer.Node, beforeYesterdayZeroTimestamp),
		AverageOnlineCount:        GetThatDayAverageOnlineCount(gameServer.Node, todayZeroTimestamp),
	}
	return serverOnlineStatistics, nil
}
