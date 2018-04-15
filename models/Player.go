package models

import (
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"errors"
	"github.com/astaxie/beego/logs"
	//"github.com/jinzhu/gorm"
	"strings"
	"github.com/jinzhu/gorm"
	//"github.com/zaaksam/dproxy/go/db"
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
	TotalOnlineTime int    `json:"totalOnlineTime"`
	LastLoginIp     string `json:"lastLoginIp"`
	IsOnline        int    `json:"isOnline"`
	Level           int    `json:"level" gorm:"-"`
	VipLevel        int    `json:"vipLevel" gorm:"-"`
	Power           int    `json:"power" gorm:"-"`
	FactionName     string `json:"factionName" gorm:"-"`
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
	case "level":
		sortOrder = "level"
	case "vipLevel":
		sortOrder = "vip_level"
	case "power":
		sortOrder = "power"
	}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}

	err = gameDb.Model(&Player{}).Count(&count).Error
	utils.CheckError(err)
	whereArray := make([] string, 0)
	if params.Account != "" {
		whereArray = append(whereArray, fmt.Sprintf(" acc_id = %s", params.Account))
	}
	if params.Ip != "" {
		whereArray = append(whereArray, fmt.Sprintf(" last_login_ip = %s", params.Ip))
	}
	if params.Nickname != "" {
		whereArray = append(whereArray, fmt.Sprintf(" nickname LIKE '%%%s%%' ", params.Nickname))
	}
	if params.IsOnline != "" {
		whereArray = append(whereArray, fmt.Sprintf(" is_online = %s", params.IsOnline))
	}
	if params.PlayerId != "" {
		whereArray = append(whereArray, fmt.Sprintf(" id = %s", params.PlayerId))
	}

	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}
	sql := fmt.Sprintf(
		" select player.*, player_data.level, player_vip.level as vip_level, player_data.power from ( player left join player_data on player.id = player_data.player_id) left join player_vip on player.id = player_vip.player_id  %s order by %s limit %d,%d; ",
		whereParam,
		sortOrder,
		params.Offset,
		params.Limit,
	)
	err = gameDb.Debug().Raw(sql).Scan(&data).Error
	utils.CheckError(err)

	for _, e := range data {
		e.FactionName = GetPlayerFactionName(gameDb, e.Id)
	}
	return data, count
}

func GetPlayerFactionName(gameDb *gorm.DB, playerId int) string {
	var factionMember struct {
		FactionId int
	}
	var faction struct {
		Name string
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
	//todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	sql := fmt.Sprintf(
		`SELECT faction_id  FROM faction_member WHERE player_id = ?`)
	err := gameDb.Raw(sql, playerId).Scan(&factionMember).Error
	utils.CheckError(err)
	if err != nil {
		return ""
	}

	//logs.Info("faction_id:%d", factionMember.FactionId)
	sql = fmt.Sprintf(
		`SELECT name  FROM faction WHERE id = ?`)
	err = gameDb.Raw(sql, factionMember.FactionId).Scan(&faction).Error
	utils.CheckError(err)
	if err != nil {
		return ""
	}
	//logs.Info("faction_name:%s", faction.Name)
	//logs.Info("ppp:%v,%v", gameServer.Node, data.Count)
	return faction.Name
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

//func GetPlayerData(gameDb *gorm.DB, playerId int) ( *PlayerData, error) {
//	playerData := &PlayerData{
//		PlayerId: playerId,
//	}
//	err := gameDb.First(&playerData).Error
//	return playerData, err
//}

type PlayerData struct {
	PlayerId int `json:"playerId" gorm:"primary_key"`
	VipLevel int `json:"vipLevel"`
	Level    int `json:"level"`
	Power    int `json:"power"`
}

type PlayerDetail struct {
	Player
	VipLevel int `json:"vipLevel"`
	Exp      int `json:"exp"`
	Level    int `json:"level"`
	TaskId   int `json:"taskId"`
	//FactionId        int `json:"-"`
	FactionName      string `json:"factionName"`
	TitleId          int
	Attack           int `json:"attack"`
	MaxHp            int `json:"maxHp"`
	Defense          int `json:"defense"`
	Hit              int `json:"hit"`
	Dodge            int `json:"dodge"`
	Critical         int `json:"critical"`
	Power            int `json:"power"`
	LastWorldSceneId int `json:"lastWorldSceneId"`
	PlayerPropList   [] *PlayerProp `json:"playerPropList"`
}

type PlayerProp struct {
	PropType int `json:"propType"`
	PropId   int `json:"propId"`
	Num      int `json:"num"`
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
		`SELECT player.*, player_data.*, player_task.task_id, player_vip.level as vip_level FROM ((player LEFT JOIN player_data on player.id = player_data.player_id) LEFT JOIN player_task on player.id = player_task.player_id) LEFT JOIN player_vip on player_vip.player_id = player.id WHERE player.id = ? `)
	err = gameDb.Raw(sql, playerId).Scan(&playerDetail).Error

	playerDetail.PlayerPropList, err = GetPlayerPropList(platformId, serverId, playerId)
	playerDetail.FactionName = GetPlayerFactionName(gameDb, playerId)
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

//`c_ten_minute_statics`
type CTenMinuteStatics struct {
	Node          string
	Time          int
	OnlineNum     int
	RegisterCount int
}

type ServerOnlineStatistics struct {
	PlatformId                int       `json:"platformId"`
	ServerId                  string    `json:"serverId"`
	TodayCreateRole           int       `json:"todayCreateRole"`
	TodayRegister             int       `json:"todayRegister"`
	OnlineCount               int       `json:"onlineCount"`
	MaxOnlineCount            int       `json:"maxOnlineCount"`
	AverageOnlineCount        float32   `json:"averageOnlineCount"`
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
		TodayCreateRole:           GetTodayCreateRoleCount(gameDb),
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

type PropConsumeStatistics struct {
	OpType int     `json:"opType"`
	Count  int     `json:"count"`
	Rate   float32 `json:"rate"`
}
type PropConsumeStatisticsQueryParam struct {
	PlayerName string
	PlayerId   int
	PlatformId int
	ServerId   string
	StartTime  int
	EndTime    int
	PropType   int
	PropId     int
	Type       int
}

func GetPropConsumeStatistics(param *PropConsumeStatisticsQueryParam) ([]*PropConsumeStatistics, error) {
	if param.PropType == 0 || param.PropId == 0 {
		return nil, errors.New("请选择道具")
	}
	gameDb, err := GetGameDbByPlatformIdAndSid(param.PlatformId, param.ServerId)
	defer gameDb.Close()
	utils.CheckError(err)
	list := make([]*PropConsumeStatistics, 0)
	//if param.EndTime == 0 {
	//	param.EndTime = 9999999999999
	//}
	var changeValue string
	if param.Type == 0 {
		changeValue = "change_value < 0"
	} else {
		changeValue = "change_value > 0"
	}
	var timeRange string
	if param.StartTime > 0 {
		timeRange = fmt.Sprintf("and op_time between %d and %d", param.StartTime, param.EndTime)
	}

	var selectPlayer string
	if param.PlayerId > 0 {
		timeRange = fmt.Sprintf("and player_id = %d", param.PlayerId)
	}

	sql := fmt.Sprintf(
		` select op_type, sum(change_value) as count from player_prop_log where %s and prop_type = ? and prop_id = ? %s %s group by op_type; `, changeValue, timeRange, selectPlayer)
	err = gameDb.Debug().Raw(sql, param.PropType, param.PropId).Scan(&list).Error
	utils.CheckError(err)
	var sum = 0
	for _, e := range list {
		sum += e.Count
	}
	for _, e := range list {
		e.Rate = float32(e.Count) / float32(sum) * 100
	}
	return list, nil
}
