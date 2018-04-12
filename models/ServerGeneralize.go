package models

import (
	"github.com/chnzrb/myadmin/utils"
)

type ServerGeneralize struct {
	PlatformId      int    `json:"platformId"`
	ServerId        string `json:"serverId"`
	OpenTime        int    `json:"openTime"`
	Version         string `json:"version"`
	TotalRegister   int    `json:"totalRegister"`
	TotalCreateRole int    `json:"totalCreateRole"`
	TodayCreateRole int    `json:"todayCreateRole"`
	TodayRegister   int    `json:"todayRegister"`
	OnlineCount     int    `json:"onlineCount"`
	Status          int    `json:"status"`
	MaxLevel        int    `json:"maxLevel"`
	MaxOnlineCount  int    `json:"maxOnlineCount"`
}

type ServerGeneralizeQueryParam struct {
	PlatformId int
	ServerId   string
}

func GetServerGeneralize(platformId int, serverId string) (*ServerGeneralize, error) {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	defer gameDb.Close()
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
	serverGeneralize.TotalRegister = GetTotalRegister(gameDb)
	serverGeneralize.TotalCreateRole = GetTotalCreateRole(gameDb)
	//db.Model(&Player{}).Where(&Player{IsOnline: 1}).Count(&count)
	serverGeneralize.OnlineCount = GetNowOnlineCount(gameDb)
	serverGeneralize.MaxLevel = GetMaxPlayerLevel(gameDb)
	serverGeneralize.TodayCreateRole = GetTodayCreateRole(gameDb)
	serverGeneralize.TodayRegister = GetTodayRegister(gameDb)
	serverGeneralize.MaxOnlineCount = GetMaxOnlineCount(platformId, serverId)
	return serverGeneralize, err
}
