package models

import (
	"github.com/chnzrb/myadmin/utils"
	"strings"
)

type GameServerQueryParam struct {
	BaseQueryParam
	PlatformId string `json:"platformId"`
	ServerId   string `json:"serverId"`
	Node       string `json:"node"`
}

type GameServer struct {
	PlatformId string `gorm:"primary_key" json:"platformId"`
	Sid        string `gorm:"primary_key" json:"serverId"`
	Desc       string `json:"desc"`
	Node       string `json:"node"`
	State      int    `gorm:"-" json:"state"`
	IsShow     int    `json:"isShow"`
	OpenTime   int    `gorm:"-" json:"openTime"`
	ZoneNode   string `gorm:"-" json:"zoneNode"`
	IsAdd      int    `gorm:"-" json:"isAdd"`
	DbVersion     int    `json:"dbVersion" gorm:"-"`
	RunState   int    `json:"runState" gorm:"-"`
	StartTime  int    `json:"startTime" gorm:"-"`
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
	err := DbCenter.Model(&GameServer{}).Where(&GameServer{
		PlatformId: params.PlatformId,
		Sid:        params.ServerId,
		Node:       params.Node,
	}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		serverNode, err := GetServerNode(e.Node)
		e.DbVersion = GetDbVersion(e.Node)
		utils.CheckError(err)
		if err == nil {
			e.State = serverNode.State
			e.OpenTime = serverNode.OpenTime
			e.ZoneNode = serverNode.ZoneNode
			e.RunState = serverNode.RunState
			e.StartTime = GetNodeStartTime(e.Node)
		}

	}
	return data, count
}

// 获取单个游戏服
func GetGameServerOne(platformId string, id string) (*GameServer, error) {
	gameServer := &GameServer{
		Sid:        id,
		PlatformId: platformId,
	}
	err := DbCenter.First(&gameServer).Error
	return gameServer, err
}

func IsGameServerExists(platformId string, id string) bool {
	gameServer := &GameServer{
		Sid:        id,
		PlatformId: platformId,
	}
	return ! DbCenter.First(&gameServer).RecordNotFound()
}

// 获取该节点关联的所有游戏服
func GetGameServerByNode(node string) [] *GameServer {
	data := make([]*GameServer, 0)
	err := DbCenter.Model(&GameServer{}).Where(&GameServer{
		Node: node,
	}).Find(&data).Error
	utils.CheckError(err)
	return data
}
func GetGameServerIdListStringByNode(node string) string {
	serverIdList := GetGameServerIdListByNode(node)
	return "'" + strings.Join(serverIdList, "','") + "'"
}
func GetGameServerIdListByNode(node string) [] string {
	data := make([]*GameServer, 0)
	serverIdList := make([]string, 0)
	err := DbCenter.Model(&GameServer{}).Where(&GameServer{
		Node: node,
	}).Find(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		serverIdList = append(serverIdList, e.Sid)
	}
	return serverIdList
}
