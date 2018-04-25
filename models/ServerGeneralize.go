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
	Node   string `json:"serverId"`
}

func GetServerGeneralize(platformId int, node string) (*ServerGeneralize, error) {
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	defer gameDb.Close()
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	serverGeneralize := &ServerGeneralize{
		PlatformId:platformId,
		ServerId: GetGameServerIdListStringByNode(node),
		OpenTime: serverNode.OpenTime,
		Version: serverNode.ServerVersion,
		Status : serverNode.State,
		TotalRegister: GetTotalRegister(gameDb),
		TotalCreateRole: GetTotalCreateRoleCount(gameDb),
		OnlineCount:GetNowOnlineCount(gameDb),
		MaxLevel: GetMaxPlayerLevel(gameDb),
		TodayCreateRole:GetTodayCreateRoleCount(gameDb),
		TodayRegister:GetTodayRegister(gameDb),
		MaxOnlineCount:GetMaxOnlineCount(node),
		TotalChargeIngot:GetServerTotalChargeIngot(node),
		TotalChargeMoney:GetServerTotalChargeMoney(node),
		TotalChargePlayerCount:GetServerChargePlayerCount(node),
		SecondChargePlayerCount:GetServerSecondChargePlayerCount(node),
	}

	serverGeneralize.ARPU = CaclARPU(serverGeneralize.TotalChargeMoney, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeRate = CaclChargeRate(serverGeneralize.TotalChargePlayerCount, serverGeneralize.TotalCreateRole)
	serverGeneralize.SecondChargeRate = CaclSceondChargeRate(serverGeneralize.SecondChargePlayerCount, serverGeneralize.TotalChargePlayerCount)

	return serverGeneralize, err
}
