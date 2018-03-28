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
	VipLevel int `json:"vipLevel"`
	Exp      int `json:"exp"`
	Level    int `json:"level"`
	TitleId  int
	Attack   int `json:"attack"`
	MaxHp    int `json:"maxHp"`
	Defense  int `json:"defense"`
	Hit      int `json:"hit"`
	Dodge    int `json:"dodge"`
	Critical int `json:"critical"`
	Power    int `json:"power"`
}

func GetPlayerDetail(platformId int, serverId string, playerId int) (*PlayerDetail, error) {
	db, err := GetDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)

	playerDetail := &PlayerDetail{}
	defer db.Close()
	sql := fmt.Sprintf(
		`SELECT player.*, player_data.* FROM player LEFT JOIN player_data on player.id = player_data.player_id WHERE player.id = ? `)
	err = db.Raw(sql, playerId).Scan(&playerDetail).Error
	return playerDetail, err
}
