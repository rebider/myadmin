package models
//
//import (
//	"github.com/chnzrb/myadmin/utils"
//	"github.com/jinzhu/gorm"
//	"github.com/astaxie/beego/logs"
//	"fmt"
//	"regexp"
//	"github.com/astaxie/beego"
//	"strconv"
//	"time"
//	"strings"
//)
//
//type PlayerMissionLog struct {
//	Id          int    `json:"id"`
//	PlayerId    int    `json:"playerId"`
//	PlayerName  string `json:"playerName" gorm:"-"`
//	MissionType int    `json:"missionType"`
//	MissionId   int    `json:"missionId"`
//	Result      int    `json:"result"`
//}
//
//type PlayerMissionLogQueryParam struct {
//	BaseQueryParam
//	PlatformId  string
//	ServerId    string `json:"serverId"`
//	Ip          string
//	PlayerId    int
//	PlayerName  string
//	Datetime    int    `json:"datetime"`
//	StartTime   int
//	EndTime     int
//	MissionType int
//}
////
////func GetPlayerMissionLogList(params *PlayerMissionLogQueryParam) ([]*PlayerMissionLog, int) {
////	gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
////	utils.CheckError(err)
////	if err != nil {
////		return nil, 0
////	}
////	node := gameServer.Node
////	serverNode, err := GetServerNode(node)
////	utils.CheckError(err)
////	if err != nil {
////		return nil, 0
////	}
////	t := time.Unix(int64(params.Datetime), 0)
////	logDir := fmt.Sprintf("%d_%d_%d", t.Year(), t.Month(), t.Day())
////	logs.Debug("logDir:%s", logDir)
////
////	grepParam := ""
////	grepParam += fmt.Sprintf(" | /usr/bin/grep \\{p,%d\\} ", params.PlayerId)
////	if params.MissionType > 0 {
////		grepParam += fmt.Sprintf(" | /usr/bin/grep \\{mT,%d\\} ", params.MissionType)
////	}
////	sshKey := beego.AppConfig.String("ssh_key")
////	sshPort := beego.AppConfig.String("ssh_port")
////	nodeName := strings.Split(serverNode.Node, "@")[0]
////	nodeIp := strings.Split(serverNode.Node, "@")[1]
////	cmd := fmt.Sprintf("ssh -i %s -p%s %s ' /usr/bin/cat /data/log/game/%s/%s/player_prop_log.log %s'", sshKey, sshPort, nodeIp, nodeName, logDir, grepParam)
////	out, err := utils.ExecShell(cmd)
////	utils.CheckError(err)
////	//logs.Debug(out)
////
////	reg := regexp.MustCompile(`(\d+):(\d+):(\d+)\s+\[{pid,(\d+)},{mT,(\d+)},{mI,(\d+)},{r,(\d+)},{c,([-\d]+)},{n,(\d+)}\]`)
////	matchArray := reg.FindAllStringSubmatch(out, -1)
////	//logs.Debug("%+v", matchArray)
////	data := make([]*PlayerPropLog, 0)
////	for _, e := range matchArray {
////		//logs.Debug("%+v", e)
////		h, err := strconv.Atoi(e[1])
////		utils.CheckError(err)
////		m, err := strconv.Atoi(e[2])
////		utils.CheckError(err)
////		s, err := strconv.Atoi(e[3])
////		utils.CheckError(err)
////
////		time := h*60*60 + m*60 + s
////		if params.StartTime > 0 && params.EndTime > 0 {
////			if time < params.StartTime || time > params.EndTime {
////				continue
////			}
////		}
////		t := params.Datetime + time
////		playerId, err := strconv.Atoi(e[4])
////		utils.CheckError(err)
////		propType, err := strconv.Atoi(e[5])
////		utils.CheckError(err)
////		propId, err := strconv.Atoi(e[6])
////		utils.CheckError(err)
////		logType, err := strconv.Atoi(e[7])
////		utils.CheckError(err)
////		change, err := strconv.Atoi(e[8])
////		utils.CheckError(err)
////		new, err := strconv.Atoi(e[9])
////		utils.CheckError(err)
////		data = append(data, &PlayerPropLog{
////			PlayerId:    playerId,
////			PropType:    propType,
////			PropId:      propId,
////			OpType:      logType,
////			ChangeValue: change,
////			NewValue:    new,
////			OpTime:      t,
////		})
////	}
////	len := len(data)
////	limit := params.BaseQueryParam.Limit
////	start := params.BaseQueryParam.Offset
////	if start >= len {
////		return nil, len
////	}
////	if start+limit > len {
////		limit = len - start
////	}
////	logs.Debug(len, start, limit)
////	return data[start:start+limit], len
////	//exec_shell("ssh -i /root/.ssh/thyz_87 -p22 39.108.98.87 \"cat /data/log/game/qq_s0/2018_9_30/player_prop_log.log | grep 10130\"")
////	//cmdStr :=" /bin/bash -c 'ssh -i /root/.ssh/thyz_87 -p22 39.108.98.87 'cat /data/log/game/qq_s0/2018_9_30/player_prop_log.log | grep 10130''"
////	//cmdStr := " ssh -i /root/.ssh/thyz_qq -p6888 10.105.54.242 \"cat /data/log/game/s1/2018_10_2/player_prop_log.log |grep 15486\""
////	//commandArgs := strings.Split(cmdStr, " ")
////	//out, err := utils.Cmd(commandArgs[0], commandArgs[1:])
////	//utils.CheckError(err)
////	//logs.Debug(out)
////}
