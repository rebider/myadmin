package crons

import (
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/logs"
)

//更新所有游戏节点  DailyStatistics
func TmpUpdateAllGameNodeDailyStatistics(Time int) {
	logs.Info("TmpUpdateAllGameNodeDailyStatistics:%v", Time)
	gameServerNodeList := models.GetAllGameServerNode()
	for _, serverNode := range gameServerNodeList {
		//err := models.UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
		//utils.CheckError(err)
		err := models.UpdateDailyChargeStatistics(serverNode.Node, Time)
		utils.CheckError(err)
	}
}

//更新所有游戏节点  DailyStatistics
func UpdateAllGameNodeDailyStatistics() {
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	logs.Info("UpdateAllGameNodeDailyStatistics:%v", todayZeroTimestamp)
	gameServerNodeList := models.GetAllGameServerNode()
	for _, serverNode := range gameServerNodeList {
		//err := models.UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
		//utils.CheckError(err)
		err := models.UpdateDailyChargeStatistics(serverNode.Node, todayZeroTimestamp - 86400)
		utils.CheckError(err)
		err = models.UpdateDailyOnlineStatistics(serverNode.Node, todayZeroTimestamp - 86400)
		utils.CheckError(err)
		err = models.UpdateDailyRegisterStatistics(serverNode.Node, todayZeroTimestamp - 86400)
		utils.CheckError(err)
		err = models.UpdateDailyActiveStatistics(serverNode.Node, todayZeroTimestamp - 86400)
		utils.CheckError(err)
	}
}

//更新所有游戏节点  RemainTotal
func UpdateAllGameNodeRemainTotal() {
	now := utils.GetTimestamp()
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	logs.Info("更新所有总体留存:%v", todayZeroTimestamp)
	gameServerNodeList := models.GetAllGameServerNode()
	for _, serverNode := range gameServerNodeList {
		if now >= serverNode.OpenTime {
			err := models.UpdateRemainTotal(serverNode.Node, todayZeroTimestamp - 86400)
			utils.CheckError(err)
		}

	}
}

//更新所有游戏节点  RemainActive
func UpdateAllGameNodeRemainActive() {
	now := utils.GetTimestamp()
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	logs.Info("更新所有活跃留存:%v", todayZeroTimestamp)
	gameServerNodeList := models.GetAllGameServerNode()
	for _, serverNode := range gameServerNodeList {
		if now >= serverNode.OpenTime {
			err := models.UpdateRemainActive(serverNode.Node, todayZeroTimestamp - 86400)
			utils.CheckError(err)
		}

	}
}

func TmpUpdateAllGameNodeRemainTotal(time int) {
	now := utils.GetTimestamp()
	logs.Info("更新所有总体留存:%v", time)
	gameServerNodeList := models.GetAllGameServerNode()
	for _, serverNode := range gameServerNodeList {
		if now >= serverNode.OpenTime {
			err := models.UpdateRemainActive(serverNode.Node, time)
			utils.CheckError(err)
		}

	}
}
