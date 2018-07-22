package models

import (
	"github.com/chnzrb/myadmin/utils"
	//"github.com/astaxie/beego/logs"
)

type ServerGeneralize struct {
	PlatformId              string  `json:"platformId"`
	ServerId                string  `json:"serverId"`
	OpenTime                int     `json:"openTime"`
	Version                 int     `json:"version"`
	TotalRegister           int     `json:"totalRegister"`
	TotalCreateRole         int     `json:"totalCreateRole"`
	TodayCreateRole         int     `json:"todayCreateRole"`
	TodayRegister           int     `json:"todayRegister"`
	TotalChargeIngot        int     `json:"totalChargeIngot"`
	TotalChargeMoney        int     `json:"totalChargeMoney"`
	TotalChargePlayerCount  int     `json:"totalChargePlayerCount"`
	SecondChargePlayerCount int     `json:"secondChargePlayerCount"`
	ChargeCount2Rate        float32 `json:"chargeCount2Rate"`
	ChargeCount3Rate        float32 `json:"chargeCount3Rate"`
	ChargeCount5Rate        float32 `json:"chargeCount5Rate"`
	//ChargeCountMore         int     `json:"chargeCountMore"`
	OnlineCount             int     `json:"onlineCount"`
	Status                  int     `json:"status"`
	ARPU                    float32 `json:"arpu"`
	ChargeRate              float32 `json:"chargeRate"`
	SecondChargeRate        float32 `json:"secondChargeRate"`
	MaxLevel                int     `json:"maxLevel"`
	MaxOnlineCount          int     `json:"maxOnlineCount"`
	TotalIngot              int     `json:"totalIngot"`
	TotalCoin               int     `json:"totalCoin"`
}

type ServerGeneralizeQueryParam struct {
	PlatformId string
	Node       string `json:"serverId"`
}

func GetServerGeneralize(platformId string, node string) (*ServerGeneralize, error) {
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	defer gameDb.Close()
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	chargeTimesList := GetServerChargeCountList(node)
	//logs.Debug("chargeTimesList:%+v", chargeTimesList)
	chargeCount2 := 0
	chargeCount3 := 0
	chargeCount5 := 0
	//chargeCountMore := 0
	for _, e := range chargeTimesList {
		if e >= 2 {
			chargeCount2++
		}
		if e >= 3 {
			chargeCount3++
		}
		if e >= 5 {
			chargeCount5++
		}
		//if e >= 6 {
		//	chargeCountMore++
		//}
	}
	serverGeneralize := &ServerGeneralize{
		PlatformId:              platformId,
		ServerId:                GetGameServerIdListStringByNode(node),
		OpenTime:                serverNode.OpenTime,
		Version:                 GetNodeVersion(node),
		Status:                  serverNode.State,
		TotalRegister:           GetTotalRegister(gameDb),
		TotalCreateRole:         GetTotalCreateRoleCount(gameDb),
		OnlineCount:             GetNowOnlineCount(gameDb),
		MaxLevel:                GetMaxPlayerLevel(gameDb),
		TodayCreateRole:         GetTodayCreateRoleCount(gameDb),
		TodayRegister:           GetTodayRegister(gameDb),
		MaxOnlineCount:          GetMaxOnlineCount(node),
		TotalChargeIngot:        GetServerTotalChargeIngot(node),
		TotalChargeMoney:        GetServerTotalChargeMoney(node),
		TotalChargePlayerCount:  GetServerChargePlayerCount(node),
		SecondChargePlayerCount: GetServerSecondChargePlayerCount(node),
		//ChargeCountMore:         chargeCountMore,
		TotalIngot:              GetTotalProp(gameDb, 1, 2),
		TotalCoin:               GetTotalProp(gameDb, 1, 1),
	}

	serverGeneralize.ARPU = CaclARPU(serverGeneralize.TotalChargeMoney, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeRate = CaclChargeRate(serverGeneralize.TotalChargePlayerCount, serverGeneralize.TotalCreateRole)
	serverGeneralize.SecondChargeRate = CaclChargeRate(serverGeneralize.SecondChargePlayerCount, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeCount2Rate = CaclChargeRate(chargeCount2, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeCount3Rate = CaclChargeRate(chargeCount3, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeCount5Rate = CaclChargeRate(chargeCount5, serverGeneralize.TotalChargePlayerCount)

	return serverGeneralize, err
}
