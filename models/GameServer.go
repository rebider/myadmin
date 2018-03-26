package models

import (
	"github.com/chnzrb/myadmin/utils"
)

type GameServerQueryParam struct {
	BaseQueryParam
	PlatformId int    `json:"platformId"`
	ServerId   string `json:"serverId"`
	Node       string `json:"node"`
}

type GameServer struct {
	PlatformId int    `gorm:"primary_key" json:"platformId"`
	Sid        string `gorm:"primary_key" json:"serverId"`
	Desc       string `json:"desc"`
	Node       string `json:"node"`
}

func (t *GameServer) TableName() string {
	return "c_game_server"
}

//获取所有数据
func GetAllGameServer() ([]*GameServer, int64) {
	var params GameServerQueryParam
	params.Limit = -1
	//获取数据列表和总数
	data, total := GetGameServerList(&params)
	return data, total
}

//获取游戏服列表
func GetGameServerList(params *GameServerQueryParam) ([]*GameServer, int64) {
	sortOrder := "Sid"
	switch params.Sort {
	case "Sid":
		sortOrder = "Sid"
	}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	data := make([]*GameServer, 0)
	var count int64
	err := DbCenter.Model(&GameServer{}).Count(&count).Where(&GameServer{
		PlatformId: params.PlatformId,
		Sid:        params.ServerId,
		Node:       params.Node,
	}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Error
	utils.CheckError(err)
	return data, count
}

// 获取单个游戏服
func GetGameServerOne(platformId int, id string) (*GameServer, error) {
	gameServer := &GameServer{
		Sid:        id,
		PlatformId: platformId,
	}
	err := DbCenter.First(&gameServer).Error
	return gameServer, err
}
