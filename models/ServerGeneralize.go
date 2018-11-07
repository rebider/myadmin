package models

import (
	"github.com/chnzrb/myadmin/utils"
	//"github.com/astaxie/beego/logs"
)

type ServerGeneralize struct {
	PlatformId              string  `json:"platformId"`
	ServerId                string  `json:"serverId"`
	OpenTime                int     `json:"openTime"`
	MergeTime               int     `json:"mergeTime"`
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
	OnlineCount      int     `json:"onlineCount"`
	Status           int     `json:"status"`
	ARPU             float32 `json:"arpu"`
	ChargeRate       float32 `json:"chargeRate"`
	SecondChargeRate float32 `json:"secondChargeRate"`
	MaxLevel         int     `json:"maxLevel"`
	MaxOnlineCount   int     `json:"maxOnlineCount"`
	TotalIngot       int     `json:"totalIngot"`
	TotalCoin        int     `json:"totalCoin"`
}

type ServerGeneralizeQueryParam struct {
	PlatformId  string
	ServerId    string    `json:"serverId"`
	ChannelList [] string `json:"channelList"`
}

func GetServerGeneralize(platformId string, serverId string, channelList [] string) (*ServerGeneralize, error) {
	gameServer, err := GetGameServerOne(platformId, serverId)
	if err != nil {
		return nil, err
	}
	node := gameServer.Node
	gameDb, err := GetGameDbByNode(node)
	if err != nil {
		return nil, err
	}
	defer gameDb.Close()
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	chargeTimesList := GetServerChargeCountList(platformId, serverId, channelList)
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
	todayZeroTime := utils.GetTodayZeroTimestamp()
	serverGeneralize := &ServerGeneralize{
		PlatformId:              platformId,
		ServerId:                GetGameServerIdListStringByNode(node),
		OpenTime:                serverNode.OpenTime,
		Version:                 GetNodeVersion(node),
		MergeTime:               GetMergeTime(node),
		Status:                  serverNode.State,
		TotalRegister:           GetTotalRegisterRoleCountByChannelList(gameDb, serverId, channelList),
		TotalCreateRole:         GetTotalCreateRoleCountByChannelList(gameDb, serverId, channelList),
		OnlineCount:             GetNowOnlineCount2(gameDb, serverId, channelList),
		MaxLevel:                GetMaxPlayerLevel(gameDb, serverId, channelList),
		TodayCreateRole:         GetCreateRoleCountByChannelList(gameDb, serverId, channelList, todayZeroTime, todayZeroTime+86400),
		TodayRegister:           GetRegisterRoleCountByChannelList(gameDb, serverId, channelList, todayZeroTime, todayZeroTime+86400),
		MaxOnlineCount:          GetThatDayMaxOnlineCount(platformId, serverId, channelList, todayZeroTime, todayZeroTime+86400),
		TotalChargeIngot:        GetServerTotalChargeIngot(platformId, serverId, channelList),
		TotalChargeMoney:        GetServerTotalChargeMoneyByChannelList(platformId, serverId, channelList),
		TotalChargePlayerCount:  GetServerChargePlayerCount(platformId, serverId, channelList),
		SecondChargePlayerCount: GetServerSecondChargePlayerCount(platformId, serverId, channelList),
		//ChargeCountMore:         chargeCountMore,
		TotalIngot: GetTotalProp(gameDb, 1, 2, channelList),
		TotalCoin:  GetTotalProp(gameDb, 1, 1, channelList),
	}

	serverGeneralize.ARPU = CaclRate(serverGeneralize.TotalChargeMoney, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeRate = CaclRate(serverGeneralize.TotalChargePlayerCount, serverGeneralize.TotalCreateRole)
	serverGeneralize.SecondChargeRate = CaclRate(serverGeneralize.SecondChargePlayerCount, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeCount2Rate = CaclRate(chargeCount2, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeCount3Rate = CaclRate(chargeCount3, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeCount5Rate = CaclRate(chargeCount5, serverGeneralize.TotalChargePlayerCount)

	return serverGeneralize, err
}
