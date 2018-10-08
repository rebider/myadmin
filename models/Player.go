package models

import (
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"errors"
	"strings"
	"github.com/jinzhu/gorm"
)

type Player struct {
	Id               int    `json:"id"`
	AccId            string `json:"accId"`
	Nickname         string `json:"nickname"`
	Sex              int    `json:"sex"`
	ServerId         string `json:"serverId"`
	ForbidType       int    `json:"forbidType"`
	ForbidTime       int    `json:"forbidTime"`
	RegTime          int    `json:"regTime"`
	LoginTimes       int    `json:"loginTimes"`
	LastLoginTime    int    `json:"lastLoginTime"`
	LastOfflineTime  int    `json:"lastOfflineTime"`
	TotalOnlineTime  int    `json:"totalOnlineTime"`
	LastLoginIp      string `json:"lastLoginIp"`
	From             string `json:"from"`
	Channel          string `json:"channel"`
	IsOnline         int    `json:"isOnline"`
	Type             int    `json:"type" gorm:"-"`
	Level            int    `json:"level" gorm:"-" `
	Ingot            int    `json:"ingot" gorm:"-"`
	TotalChargeMoney int    `json:"totalChargeMoney" gorm:"-"`
	VipLevel         int    `json:"vipLevel" gorm:"-"`
	Power            int    `json:"power" gorm:"-"`
	FactionName      string `json:"factionName" gorm:"-"`
	FriendCode      string `json:"friendCode" gorm:"-"`
}

type PlayerQueryParam struct {
	BaseQueryParam
	Account     string
	Ip          string
	PlayerId    string
	Nickname    string
	IsOnline    string
	PlatformId  string
	Type        string
	ServerId    string    `json:"serverId"`
	ChannelList [] string `json:"channelList"`
}

func (a *Player) TableName() string {
	return "player"
}

//获取玩家列表
func GetPlayerList(params *PlayerQueryParam) ([]*Player, int64) {
	gameDb, err := GetGameDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	//gameDb, err := GetGameDbByNode(params.Node)
	utils.CheckError(err)
	if err != nil {
		return nil, 0
	}
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
	case "totalOnlineTime":
		sortOrder = "total_online_time"
	}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}

	whereArray := make([] string, 0)
	if params.Account != "" {
		whereArray = append(whereArray, fmt.Sprintf(" acc_id = '%s'", params.Account))
	}
	if params.Ip != "" {
		whereArray = append(whereArray, fmt.Sprintf(" last_login_ip = %s", params.Ip))
	}
	if params.Nickname != "" {
		//serverId, playerName, err := SplitPlayerName(params.Nickname)
		//utils.CheckError(err)
		whereArray = append(whereArray, fmt.Sprintf("nickname LIKE '%%%s%%' ", params.Nickname))
	}
	if params.IsOnline != "" {
		whereArray = append(whereArray, fmt.Sprintf(" is_online = %s", params.IsOnline))
	}
	if params.PlayerId != "" {
		whereArray = append(whereArray, fmt.Sprintf(" id = %s", params.PlayerId))
	}
	if len(params.ChannelList) > 0 {
		//whereArray = append(whereArray, fmt.Sprintf(" channel in  (%s) ", "'"+strings.Join(params.ChannelList, "','")+"'"))
		whereArray = append(whereArray, fmt.Sprintf(" channel in  (%s) ", GetSQLWhereParam(params.ChannelList)))
	}
	//if params.Type != "" {
	//	whereArray = append(whereArray, fmt.Sprintf(" type = %s", params.Type))
	//}

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
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	err = gameDb.Model(&Player{}).Raw("select count(1) from player " + whereParam).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		e.FactionName = GetPlayerFactionName(gameDb, e.Id)
		e.Nickname = e.ServerId + "." + e.Nickname
		e.Ingot = GetPlayerIngot(gameDb, e.Id)
		playerChargeData, err := GetPlayerChargeDataOne(e.Id)
		utils.CheckError(err)
		e.TotalChargeMoney = int(playerChargeData.TotalMoney)
		e.Type = GetAccountType(params.PlatformId, e.AccId)
		//e.LastLoginIp = e.LastLoginIp + "(" + utils.GetIpLocation(e.LastLoginIp) + ")"
	}
	return data, count
}

func GetAccountType(platfromId string, accId string) int {
	var data struct {
		Type int
	}
	sql := fmt.Sprintf(
		`SELECT type FROM global_account WHERE platform_id = '%s' and account = '%s'`, platfromId, accId)
	err := DbCenter.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Type
}

func GetPlayerFactionName(gameDb *gorm.DB, playerId int) string {
	var factionMember struct {
		FactionId int
	}
	var faction struct {
		Name string
	}
	sql := fmt.Sprintf(
		`SELECT faction_id  FROM faction_member WHERE player_id = %d`, playerId)
	isNotFound := gameDb.Raw(sql).Scan(&factionMember).RecordNotFound()
	//utils.CheckError(err)
	if isNotFound {
		return ""
	}

	//logs.Info("faction_id:%d", factionMember.FactionId)
	sql = fmt.Sprintf(
		`SELECT name  FROM faction WHERE id = %d`, factionMember.FactionId)
	err := gameDb.Raw(sql).Scan(&faction).Error
	utils.CheckError(err)
	if err != nil {
		return ""
	}
	//logs.Info("faction_name:%s", faction.Name)
	//logs.Info("ppp:%v,%v", gameServer.Node, data.Count)
	return faction.Name
}

// 获取单个玩家
func GetPlayerOneByNode(node string, id int) (*Player, error) {
	gameDb, err := GetGameDbByNode(node)
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

// 获取单个玩家
func GetPlayerOne(platformId string, serverId string, id int) (*Player, error) {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	if err != nil {
		return nil, err
	}
	defer gameDb.Close()
	player := &Player{
		Id: id,
	}
	err = gameDb.First(&player).Error
	if err == nil {
		player.Type = GetAccountType(platformId, player.AccId)
	}
	return player, err
}

func GetPlayerByDb(gameDb *gorm.DB, playerId int) (*Player, error) {
	player := &Player{
		Id: playerId,
	}
	err := gameDb.First(&player).Error
	return player, err
}

func GetPlayerDataByDb(gameDb *gorm.DB, playerId int) (*PlayerData, error) {
	playerData := &PlayerData{
		PlayerId: playerId,
	}
	err := gameDb.First(&playerData).Error
	return playerData, err
}

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
	FactionName string `json:"factionName"`
	TitleId     int
	//TotalChargeMoney int            `json:"totalChargeMoney"`
	Attack           int            `json:"attack"`
	MaxHp            int            `json:"maxHp"`
	Defense          int            `json:"defense"`
	Hit              int            `json:"hit"`
	Dodge            int            `json:"dodge"`
	Critical         int            `json:"critical"`
	Power            int            `json:"power"`
	LastWorldSceneId int            `json:"lastWorldSceneId"`
	PlayerPropList   [] *PlayerProp `json:"playerPropList"`
	EquipList        [] *PlayerProp `json:"equipList"`
}

type PlayerProp struct {
	PlayerId int `json:"playerId" gorm:"primary_key"`
	PropType int `json:"propType" gorm:"primary_key"`
	PropId   int `json:"propId" gorm:"primary_key"`
	Num      int `json:"num"`
}

func GetPlayerPropList(gameDb *gorm.DB, playerId int) ([]*PlayerProp, error) {
	playerPropList := make([]*PlayerProp, 0)
	sql := fmt.Sprintf(
		`SELECT * FROM player_prop WHERE player_id = %d `, playerId)
	err := gameDb.Raw(sql).Scan(&playerPropList).Error

	return playerPropList, err
}

func GetPlayerEquipList(gameDb *gorm.DB, playerId int) ([]*PlayerProp, error) {
	playerPropList := make([]*PlayerProp, 0)
	sql := fmt.Sprintf(
		`SELECT * FROM player_equip_pos WHERE player_id = %d `, playerId)
	err := gameDb.Raw(sql).Scan(&playerPropList).Error

	return playerPropList, err
}

func GetPlayerDetail(platformId string, serverId string, playerId int) (*PlayerDetail, error) {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("连接数据库失败:%v", serverId))
	}
	defer gameDb.Close()
	playerDetail := &PlayerDetail{}

	sql := fmt.Sprintf(
		`SELECT player.*, player_data.*, player_task.task_id, player_vip.level as vip_level FROM ((player LEFT JOIN player_data on player.id = player_data.player_id) LEFT JOIN player_task on player.id = player_task.player_id) LEFT JOIN player_vip on player_vip.player_id = player.id WHERE player.id = %d `, playerId)
	err = gameDb.Raw(sql).Scan(&playerDetail).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("查询玩家失败:%v, %v", serverId, playerId))
	}
	playerDetail.Player.Nickname = playerDetail.Player.ServerId + "." + playerDetail.Player.Nickname
	playerDetail.PlayerPropList, err = GetPlayerPropList(gameDb, playerId)
	utils.CheckError(err)
	playerDetail.FactionName = GetPlayerFactionName(gameDb, playerId)
	utils.CheckError(err)
	playerChargeData, err := GetPlayerChargeDataOne(playerId)
	utils.CheckError(err)
	playerDetail.TotalChargeMoney = int(playerChargeData.TotalMoney)
	playerDetail.Player.Type = GetAccountType(platformId, playerDetail.Player.AccId)
	//playerDetail.LastLoginIp = playerDetail.LastLoginIp + "(" + utils.GetIpLocation(playerDetail.LastLoginIp) + ")"
	return playerDetail, err
}

func GetPlayerByPlatformIdAndNickname(platformId string, nickname string) (*Player, error) {
	if nickname == "" {
		return nil, errors.New("角色名字不能为空!")
	}
	serverId, playerName, err := SplitPlayerName(nickname)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("非法角色名:%s", nickname))
	}
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	if err != nil {
		return nil, err
	}
	defer gameDb.Close()
	player := &Player{}
	isNotFound := gameDb.Where(&Player{ServerId: serverId, Nickname: playerName}).First(&player).RecordNotFound()
	if isNotFound {
		return nil, errors.New(fmt.Sprintf("角色不存在:%s", nickname))
	}
	player.Nickname = player.ServerId + "." + player.Nickname
	player.Type = GetAccountType(platformId, player.AccId)
	player.Ingot = GetPlayerIngot(gameDb, player.Id)
	return player, err
}

func GetPlayerIngot(gameDb *gorm.DB, playerId int) int {
	playerProp := &PlayerProp{
		PlayerId: playerId,
		PropType: 1,
		PropId:   2,
	}
	err := gameDb.FirstOrInit(&playerProp).Error
	utils.CheckError(err)
	return playerProp.Num
}

//func GetPlayerByNodeAndNickname(node string, serverId string, nickname string) (*Player, error) {
//	if nickname == "" {
//		return nil, errors.New("角色名字不能为空!")
//	}
//	logs.Debug("nickname:%v", nickname)
//	gameDb, err := GetGameDbByNode(node)
//	utils.CheckError(err)
//	if err != nil {
//		return nil, err
//	}
//	defer gameDb.Close()
//	player := &Player{}
//	err = gameDb.Where(&Player{ServerId: serverId, Nickname: nickname}).First(&player).Error
//	if err != nil {
//		return nil, err
//	}
//	return player, err
//}

type ServerOnlineStatistics struct {
	PlatformId string `json:"platformId"`0
	//ServerId                    string    `json:"serverId"`
	TodayCreateRole             int       `json:"todayCreateRole"`
	TodayRegister               int       `json:"todayRegister"`
	OnlineCount                 int       `json:"onlineCount"`
	OnlineIpCount               int       `json:"onlineIpCount"`
	MaxOnlineCount              int       `json:"maxOnlineCount"`
	AverageOnlineCount          int       `json:"averageOnlineCount"`
	TodayOnlineList             [] string `json:"todayOnlineList"`
	YesterdayOnlineList         [] string `json:"yesterdayOnlineList"`
	BeforeYesterdayOnlineList   [] string `json:"beforeYesterdayOnlineList"`
	TodayRegisterList           [] string `json:"todayRegisterList"`
	YesterdayRegisterList       [] string `json:"yesterdayRegisterList"`
	BeforeYesterdayRegisterList [] string `json:"beforeYesterdayRegisterList"`
}

func GetServerOnlineStatistics(platformId string, serverId string, channelList [] string) (*ServerOnlineStatistics, error) {

	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	yesterdayZeroTimestamp := todayZeroTimestamp - 86400
	beforeYesterdayZeroTimestamp := yesterdayZeroTimestamp - 86400

	todayOnlineList, nowOnline := get24hoursOnlineCount(platformId, serverId, channelList, todayZeroTimestamp)
	yesterdayOnlineList, _ := get24hoursOnlineCount(platformId, serverId, channelList, yesterdayZeroTimestamp)
	beforeYesterdayOnlineList, _ := get24hoursOnlineCount(platformId, serverId, channelList, beforeYesterdayZeroTimestamp)
	todayRegisterList, todayRegister := get24hoursRegisterCount(platformId, serverId, channelList, todayZeroTimestamp)
	yesterdayRegisterList, _ := get24hoursRegisterCount(platformId, serverId, channelList, yesterdayZeroTimestamp)
	beforeYesterdayRegisterList, _ := get24hoursRegisterCount(platformId, serverId, channelList, beforeYesterdayZeroTimestamp)
	serverOnlineStatistics := &ServerOnlineStatistics{
		PlatformId: platformId,
		//ServerId:                    serverId,
		OnlineCount: nowOnline,
		//TodayCreateRole: GetCreateRoleCountByChannelList(gameDb, serverId, channelList, todayZeroTimestamp, todayZeroTimestamp+86400),
		TodayRegister:   todayRegister,
		//MaxOnlineCount:              GetMaxOnlineCount(node),
		TodayOnlineList:             todayOnlineList,
		YesterdayOnlineList:         yesterdayOnlineList,
		BeforeYesterdayOnlineList:   beforeYesterdayOnlineList,
		TodayRegisterList:           todayRegisterList,
		YesterdayRegisterList:      yesterdayRegisterList,
		BeforeYesterdayRegisterList: beforeYesterdayRegisterList,
	}
	return serverOnlineStatistics, nil
}
