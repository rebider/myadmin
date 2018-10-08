package crons

import (
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
)

//更新所有游戏节点  DailyStatistics
//func TmpUpdateAllGameNodeDailyStatistics(Time int) {
//	logs.Info("TmpUpdateAllGameNodeDailyStatistics:%v", Time)
//	gameServerNodeList := models.GetAllGameServerNode()
//	for _, serverNode := range gameServerNodeList {
//		//err := models.UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
//		//utils.CheckError(err)
//		err := models.UpdateDailyChargeStatistics(serverNode.Node, Time)
//		utils.CheckError(err)
//	}
//}

//更新所有游戏节点  DailyStatistics
func UpdateAllGameNodeDailyStatistics() {
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	DoUpdateAllGameNodeDailyStatistics(todayZeroTimestamp - 86400)
}

//更新所有游戏节点  DailyStatistics
func UpdateAllGameNodeDailyLTV() {
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	DoUpdateAllGameNodeDailyLTV(todayZeroTimestamp - 86400)
}

//更新所有游戏节点  ReaminCharge
func UpdateAllGameNodeChargeRemain() {
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	DoUpdateAllGameNodeReaminCharge(todayZeroTimestamp - 86400)
}

//func Repire() {
//	//1532016000
//	//1536163200 退1536422400
//
//	for i := 1532016000; i <= 1538236800; i += 86400 {
//		DoUpdateAllGameNodeDailyLTV(i)
//	}
//}


func Repire() {
	//1532016000
	//1536163200 退1536422400

	for i := 1537200000; i <= 1538236800; i += 86400 {
		DoUpdateAllGameNodeReaminCharge(i)
	}
}



func DoUpdateAllGameNodeTenMinuteStatistics(timestamp int) {
	logs.Info("更新每10分钟统计:%v", timestamp)
	gameServerList, _ := models.GetAllGameServerDirty()

	for _, gameServer := range gameServerList {
		platformId := gameServer.PlatformId
		serverId := gameServer.Sid
		channelList := models.GetChannelListByPlatformId(platformId)
		if len(channelList) == 0 {
			logs.Error("渠道未配置:%v %+v", platformId, channelList)
		}

		if timestamp < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
			continue
		}
		for _, channel := range channelList {
			err := models.UpdateTenMinuteStatistics(platformId, serverId, channel.Channel, timestamp)
			utils.CheckError(err)
		}
	}
	logs.Info("更新每10分钟完毕.")
}


func DoUpdateAllGameNodeDailyStatistics(timestamp int) {
	//todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	logs.Info("更新每日统计:%v", timestamp)
	gameServerList, _ := models.GetAllGameServerDirty()

	for _, gameServer := range gameServerList {
		//err := models.UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
		//utils.CheckError(err)
		platformId := gameServer.PlatformId
		serverId := gameServer.Sid
		channelList := models.GetChannelListByPlatformId(platformId)
		if len(channelList) == 0 {
			logs.Error("渠道未配置:%v %+v", platformId, channelList)
		}

		if timestamp < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
			continue
		}
		for _, channel := range channelList {
			err := models.UpdateDailyStatistics(platformId, serverId, channel.Channel, timestamp)
			utils.CheckError(err)
		}
	}
	logs.Info("更新每日统计完毕.")
}


func DoUpdateAllGameNodeDailyLTV(timestamp int) {
	//todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	logs.Info("更新每日LTV:%v", timestamp)
	gameServerList, _ := models.GetAllGameServerDirty()
	for _, gameServer := range gameServerList {
		//err := models.UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
		//utils.CheckError(err)
		platformId := gameServer.PlatformId
		serverId := gameServer.Sid
		channelList := models.GetChannelListByPlatformId(platformId)
		if len(channelList) == 0 {
			logs.Error("渠道未配置:%v %+v", platformId, channelList)
		}

		if timestamp < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
			continue
		}
		for _, channel := range channelList {
			err := models.UpdateDailyLTV(platformId, serverId, channel.Channel, timestamp)
			utils.CheckError(err)
		}
	}
	logs.Info("更新每日LTV完毕.")
}


func DoUpdateAllGameNodeReaminCharge(timestamp int) {
	//todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	logs.Info("更新付费留存:%v", timestamp)
	gameServerList, _ := models.GetAllGameServerDirty()
	for _, gameServer := range gameServerList {
		//err := models.UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
		//utils.CheckError(err)
		platformId := gameServer.PlatformId
		//if platformId == "af" || platformId == "djs" {
			serverId := gameServer.Sid
			channelList := models.GetChannelListByPlatformId(platformId)
			if len(channelList) == 0 {
				logs.Error("渠道未配置:%v %+v", platformId, channelList)
			}

			if timestamp < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
				continue
			}
			for _, channel := range channelList {
				err := models.UpdateRemainCharge(platformId, serverId, channel.Channel, timestamp)
				utils.CheckError(err)
			}
		//}
	}
	logs.Info("更新付费留存完毕.")
}

////更新所有游戏节点  DailyStatistics
//func UpdateAllGameNodeDailyStatistics() {
//	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
//	logs.Info("UpdateAllGameNodeDailyStatistics:%v", todayZeroTimestamp)
//	gameServerNodeList := models.GetAllGameServerNode()
//	for _, serverNode := range gameServerNodeList {
//		//err := models.UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
//		//utils.CheckError(err)
//		err := models.UpdateDailyChargeStatistics(serverNode.Node, todayZeroTimestamp-86400)
//		utils.CheckError(err)
//		err = models.UpdateDailyOnlineStatistics(serverNode.Node, todayZeroTimestamp-86400)
//		utils.CheckError(err)
//		err = models.UpdateDailyRegisterStatistics(serverNode.Node, todayZeroTimestamp-86400)
//		utils.CheckError(err)
//		err = models.UpdateDailyActiveStatistics(serverNode.Node, todayZeroTimestamp-86400)
//		utils.CheckError(err)
//	}
//}

//更新所有游戏节点  RemainTotal
func UpdateAllGameNodeRemainTotal() {
	now := utils.GetTimestamp()
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	logs.Info("更新所有总体留存:%v", todayZeroTimestamp)
	gameServerList, _ := models.GetAllGameServerDirty()

	for _, gameServer := range gameServerList {
		serverNode, err := models.GetServerNode(gameServer.Node)
		utils.CheckError(err)
		platformId := gameServer.PlatformId
		serverId := gameServer.Sid
		channelList := models.GetChannelListByPlatformId(platformId)
		if len(channelList) == 0 {
			logs.Error("渠道未配置:%v %+v", platformId, channelList)
		}
		if now >= serverNode.OpenTime {
			for _, channel := range channelList {
				err := models.UpdateRemainTotal(platformId, serverId, channel.Channel, todayZeroTimestamp-86400)
				utils.CheckError(err)
			}
		}
	}
}


//func UpdateAllGameNodeLTV() {
//	now := utils.GetTimestamp()
//	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
//	logs.Info("更新所有LTV:%v", todayZeroTimestamp)
//	gameServerList, _ := models.GetAllGameServer()
//
//	for _, gameServer := range gameServerList {
//		serverNode, err := models.GetServerNode(gameServer.Node)
//		utils.CheckError(err)
//		platformId := gameServer.PlatformId
//		serverId := gameServer.Sid
//		channelList := models.GetChannelListByPlatformId(platformId)
//		if len(channelList) == 0 {
//			logs.Error("渠道未配置:%v %+v", platformId, channelList)
//		}
//		if now >= serverNode.OpenTime {
//			for _, channel := range channelList {
//				err := models.UpdateRemainTotal(platformId, serverId, channel.Channel, todayZeroTimestamp-86400)
//				utils.CheckError(err)
//			}
//		}
//	}
//}


//func UpdateAllGameNodeRemainTotal() {
//	now := utils.GetTimestamp()
//	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
//	logs.Info("更新所有总体留存:%v", todayZeroTimestamp)
//	gameServerNodeList := models.GetAllGameServerNode()
//	for _, serverNode := range gameServerNodeList {
//		if now >= serverNode.OpenTime {
//			err := models.UpdateRemainTotal(serverNode.Node, todayZeroTimestamp-86400)
//			utils.CheckError(err)
//		}
//
//	}
//}

//更新所有游戏节点  RemainActive
func UpdateAllGameNodeRemainActive() {
	now := utils.GetTimestamp()
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	logs.Info("更新所有活跃留存:%v", todayZeroTimestamp)
	gameServerList, _ := models.GetAllGameServerDirty()

	for _, gameServer := range gameServerList {
		serverNode, err := models.GetServerNode(gameServer.Node)
		utils.CheckError(err)
		platformId := gameServer.PlatformId
		serverId := gameServer.Sid
		channelList := models.GetChannelListByPlatformId(platformId)
		if len(channelList) == 0 {
			logs.Error("渠道未配置:%v %+v", platformId, channelList)
		}
		if now >= serverNode.OpenTime {
			for _, channel := range channelList {
				err := models.UpdateRemainActive(platformId, serverId, channel.Channel, todayZeroTimestamp-86400)
				utils.CheckError(err)
			}
		}
	}
}

//func TmpUpdateAllGameNodeRemainTotal(time int) {
//	now := utils.GetTimestamp()
//	logs.Info("更新所有总体留存:%v", time)
//	gameServerNodeList := models.GetAllGameServerNode()
//	for _, serverNode := range gameServerNodeList {
//		if now >= serverNode.OpenTime {
//			err := models.UpdateRemainTotal(serverNode.Node, time)
//			utils.CheckError(err)
//		}
//
//	}
//}
//
//func TmpUpdateAllGameNodeRemainActive(time int) {
//	now := utils.GetTimestamp()
//	logs.Info("更新所有活跃留存:%v", time)
//	gameServerNodeList := models.GetAllGameServerNode()
//	for _, serverNode := range gameServerNodeList {
//		if now >= serverNode.OpenTime {
//			err := models.UpdateRemainActive(serverNode.Node, time-86400)
//			utils.CheckError(err)
//		}
//
//	}
//}


