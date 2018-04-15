package models

import (
	"github.com/chnzrb/myadmin/utils"
)

type ServerGeneralize struct {
	PlatformId              int     `json:"platformId"`
	ServerId                string  `json:"serverId"`
	OpenTime                int     `json:"openTime"`
	Version                 string  `json:"version"`
	TotalRegister           int     `json:"totalRegister"`
	TotalCreateRole         int     `json:"totalCreateRole"`
	TodayCreateRole         int     `json:"todayCreateRole"`
	TodayRegister           int     `json:"todayRegister"`
	TotalChargeIngot        int     `json:"totalChargeIngot"`
	TotalChargeMoney        int     `json:"totalChargeMoney"`
	TotalChargePlayerCount  int     `json:"totalChargePlayerCount"`
	SecondChargePlayerCount int     `json:"secondChargePlayerCount"`
	OnlineCount             int     `json:"onlineCount"`
	Status                  int     `json:"status"`
	ARPU                    float32 `json:"arpu"`
	ChargeRate              float32 `json:"chargeRate"`
	SecondChargeRate        float32 `json:"secondChargeRate"`
	MaxLevel                int     `json:"maxLevel"`
	MaxOnlineCount          int     `json:"maxOnlineCount"`
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
	serverGeneralize.TotalCreateRole = GetTotalCreateRoleCount(gameDb)
	//db.Model(&Player{}).Where(&Player{IsOnline: 1}).Count(&count)
	serverGeneralize.OnlineCount = GetNowOnlineCount(gameDb)
	serverGeneralize.MaxLevel = GetMaxPlayerLevel(gameDb)
	serverGeneralize.TodayCreateRole = GetTodayCreateRoleCount(gameDb)
	serverGeneralize.TodayRegister = GetTodayRegister(gameDb)
	serverGeneralize.MaxOnlineCount = GetMaxOnlineCount(platformId, serverId)

	serverGeneralize.TotalChargeIngot = GetServerTotalChargeIngot(platformId, serverId)
	serverGeneralize.TotalChargeMoney = GetServerTotalChargeMoney(platformId, serverId)
	serverGeneralize.TotalChargePlayerCount = GetServerChargePlayerCount(platformId, serverId)
	serverGeneralize.SecondChargePlayerCount = GetServerSecondChargePlayerCount(platformId, serverId)

	serverGeneralize.ARPU = CaclARPU(serverGeneralize.TotalChargeMoney, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeRate = CaclChargeRate(serverGeneralize.TotalChargePlayerCount, serverGeneralize.TotalCreateRole)
	serverGeneralize.SecondChargeRate = CaclSceondChargeRate(serverGeneralize.SecondChargePlayerCount, serverGeneralize.TotalChargePlayerCount)

	return serverGeneralize, err
}
