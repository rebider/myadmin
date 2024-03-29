package models

import (
	"github.com/chnzrb/myadmin/utils"
	"strings"
	"strconv"
	"github.com/astaxie/beego/logs"
	"fmt"
	"github.com/astaxie/beego"
	"time"
	"github.com/linclin/gopub/src/github.com/pkg/errors"
	"sort"
)

type GameServerQueryParam struct {
	BaseQueryParam
	PlatformId string `json:"platformId"`
	ServerId   string `json:"serverId"`
	Node       string `json:"node"`
}

type GameServer struct {
	PlatformId      string `gorm:"primary_key" json:"platformId"`
	Sid             string `gorm:"primary_key" json:"serverId"`
	Desc            string `json:"desc"`
	Node            string `json:"node"`
	State           int    `gorm:"-" json:"state"`
	IsShow          int    `json:"isShow"`
	Ip            string `json:"ip" gorm:"-"`
	OpenTime        int    `gorm:"-" json:"openTime"`
	ZoneNode        string `gorm:"-" json:"zoneNode"`
	IsAdd           int    `gorm:"-" json:"isAdd"`
	DbVersion       int    `json:"dbVersion" gorm:"-"`
	RunState        int    `json:"runState" gorm:"-"`
	StartTime       int    `json:"startTime" gorm:"-"`
	OnlineCount     int    `gorm:"-" json:"onlineCount"`
	CreateRoleCount int    `gorm:"-" json:"createRoleCount"`
}

func (t *GameServer) TableName() string {
	return "c_game_server"
}

//获取所有数据
func GetAllGameServerDirty() ([]*GameServer, int64) {
	data := make([]*GameServer, 0)
	var count int64
	err := DbCenter.Model(&GameServer{}).Count(&count).Find(&data).Error
	utils.CheckError(err)
	return data, count
}

func GetAllGameServerDirtyByPlatfomrId(platformId string) ([]*GameServer, int64) {
	data := make([]*GameServer, 0)
	var count int64
	err := DbCenter.Model(&GameServer{}).Where(&GameServer{PlatformId:platformId}).Count(&count).Find(&data).Error
	utils.CheckError(err)
	return data, count
}

//获取所有数据
//func GetAllGameServer() ([]*GameServer, int64) {
//	var params GameServerQueryParam
//	params.Limit = -1
//	//获取数据列表和总数
//	data, total := GetGameServerList(&params)
//	return data, total
//}

type gameServerSlice [] *GameServer

func (s gameServerSlice) Len() int           { return len(s) }
func (s gameServerSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s gameServerSlice) Less(i, j int) bool {
	iId, err := strconv.Atoi(SubString(s[i].Sid, 1, len(s[i].Sid)-1))
	utils.CheckError(err)
	jId, err := strconv.Atoi(SubString(s[j].Sid, 1, len(s[j].Sid)-1))
	utils.CheckError(err)
	return iId > jId
	}

func sortGameServer(list []*GameServer) []*GameServer {
	sort.Sort(gameServerSlice(list))
	return list
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
	tmpData := make([]*GameServer, 0)
	var count int64
	err := DbCenter.Model(&GameServer{}).Where(&GameServer{
		PlatformId: params.PlatformId,
		Sid:        params.ServerId,
		Node:       params.Node,
	}).Count(&count).Offset(params.Offset).Find(&tmpData).Error
	utils.CheckError(err)
	sortGameServer(tmpData)
	data := make([]*GameServer, 0, params.Limit)
	i:= 0
	j := 0
	for _, e := range tmpData {
		if i >= params.Offset {
			if j < params.Limit {
				data = append(data, e)
				j++
			}
		}
		i++
	}
	for _, e := range data {
		serverNode, err := GetServerNode(e.Node)
		e.DbVersion = GetDbVersion(e.Node)
		utils.CheckError(err)

		if err == nil {
			e.State = serverNode.State
			e.OpenTime = serverNode.OpenTime
			e.Ip = serverNode.Ip
			e.ZoneNode = serverNode.ZoneNode
			e.RunState = serverNode.RunState
			e.StartTime = GetNodeStartTime(e.Node)
			e.OnlineCount = GetNowOnlineCountByNode(e.Node)
			e.CreateRoleCount = GetTotalCreateRoleCountByNode(e.Node)
		}

	}
	return data, count
}
//func GetGameServerList(params *GameServerQueryParam) ([]*GameServer, int64) {
//	sortOrder := "Sid"
//	switch params.Sort {
//	case "Sid":
//		sortOrder = "Sid"
//	}
//	if params.Order == "descending" {
//		sortOrder = sortOrder + " desc"
//	}
//	data := make([]*GameServer, 0)
//	var count int64
//	err := DbCenter.Model(&GameServer{}).Where(&GameServer{
//		PlatformId: params.PlatformId,
//		Sid:        params.ServerId,
//		Node:       params.Node,
//	}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Error
//	utils.CheckError(err)
//	for _, e := range data {
//		serverNode, err := GetServerNode(e.Node)
//		e.DbVersion = GetDbVersion(e.Node)
//		utils.CheckError(err)
//
//		if err == nil {
//			e.State = serverNode.State
//			e.OpenTime = serverNode.OpenTime
//			e.Ip = serverNode.Ip
//			e.ZoneNode = serverNode.ZoneNode
//			e.RunState = serverNode.RunState
//			e.StartTime = GetNodeStartTime(e.Node)
//			e.OnlineCount = GetNowOnlineCountByNode(e.Node)
//			e.CreateRoleCount = GetTotalCreateRoleCountByNode(e.Node)
//		}
//
//	}
//	return data, count
//}

//获取游戏服列表
func GetAllGameServerList() []*GameServer {
	data := make([]*GameServer, 0)
	err := DbCenter.Model(&GameServer{}).Find(&data).Error
	utils.CheckError(err)
	return data
}


func GetMaxGameServerByPlatformId(platformId string) (*GameServer, int, error) {
	l := GetAllGameServerList()
	maxId := -1
	var maxGameServer *GameServer
	for _, e := range l {
		//logs.Info("GetMaxGameServerByPlatformId:%+v, %+v", e, platformId)
		//sid := strings.Split(e.Sid, "@")[0]
		if e.PlatformId != platformId {
			continue
		}
		sid := e.Sid
		//logs.Debug("sid:%s", sid)

		//logs.Debug("2222:%s", SubString(sid, 1, len(sid)-1))
		id, err := strconv.Atoi(SubString(sid, 1, len(sid)-1))
		utils.CheckError(err)
		if err != nil {
			return nil, 0, err
		}
		if id > maxId {
			maxId = id
			maxGameServer = e
		}
	}
	return maxGameServer, maxId, nil
}

func GetThisIpMaxPort(ip string) int {
	l := GetAllServerNodeList()
	maxPort := 30000
	for _, e := range l {
		nodeIp := strings.Split(e.Node, "@")[1]
		if nodeIp == ip {
			if e.Port > maxPort {
				maxPort = e.Port
			}
			if e.WebPort > maxPort {
				maxPort = e.WebPort
			}
		}
	}
	return maxPort
}

//获取最新的跨服节点
func GetLatestZone() (*ServerNode, int, error) {
	l := GetAllServerNodeByType(2)
	maxId := -1
	var maxNode *ServerNode
	for _, e := range l {
		nodeName := strings.Split(e.Node, "@")[0]

		id, err := strconv.Atoi(SubString(nodeName, 1, len(nodeName)-1))
		utils.CheckError(err)
		if err != nil {
			return nil, 0, err
		}
		if id > maxId {
			maxId = id
			maxNode = e
		}
	}
	return maxNode, maxId, nil
}

//获取最新的跨服节点
func GetLatestZoneByPlatformId(platformId string) (*ServerNode, int, error) {
	l := GetAllServerNodeByType(2)
	maxId := -1
	var maxNode *ServerNode
	for _, e := range l {
		if e.PlatformId != platformId {
			continue
		}
		nodeName := strings.Split(e.Node, "@")[0]

		if strings.Contains(nodeName, "_") {
			nodeName2 :=strings.Split(nodeName, "_")[1]
			id, err := strconv.Atoi(SubString(nodeName2, 1, len(nodeName2)-1))
			utils.CheckError(err)
			if err != nil {
				return nil, 0, err
			}
			if id > maxId {
				maxId = id
				maxNode = e
			}
		} else {
			id, err := strconv.Atoi(SubString(nodeName, 1, len(nodeName)-1))
			utils.CheckError(err)
			if err != nil {
				return nil, 0, err
			}
			if id > maxId {
				maxId = id
				maxNode = e
			}
		}

	}
	return maxNode, maxId, nil
}

func GetFreeZoneByPlatformId(platformId  string) (string, error) {
	platform, err := GetPlatformOne(platformId)
	utils.CheckError(err, "获取平台失败:" + platformId)
	if err != nil {
		return "", err
	}

	inventoryDatabase, err:= GetInventoryDatabaseOne(platform.InventoryDatabaseId)
	utils.CheckError(err, "获取数据库配置失败")
	if err != nil {
		return "", err
	}

	inventoryServer, err:= GetInventoryServerOneDirty(platform.ZoneInventoryServerId)
	utils.CheckError(err, "获取跨服配置失败")
	if err != nil {
		return "", err
	}

	serverNode, intZid, err := GetLatestZoneByPlatformId(platformId)
	if serverNode == nil || intZid == -1{
		return "", errors.New("没有对应的跨服节点")
	}
	utils.CheckError(err)
	if err != nil {
		return "", err
	}

	connectCount := GetZoneConnectNodeCount(serverNode.Node)
	logs.Info("最新的跨服节点:%s, 连接的游戏节点个数:%d", serverNode.Node, connectCount)
	if connectCount <= 2 {
		return serverNode.Node, nil
	}
	//configDbHost := beego.AppConfig.String("config_db_host")
	newIntZid := intZid + 1
	newNode := fmt.Sprintf("%s_z%d@%s", platformId, newIntZid, serverNode.Ip)
	logs.Info("新跨服节点:%s", newNode)
	out, err := AddServerNode(newNode, inventoryServer.InnerIp, 0, 0, 2, platformId, inventoryDatabase.Host, 3306, fmt.Sprintf("db_%s_zone_z%d", platformId, newIntZid))
	utils.CheckError(err, "新增节点失败:"+out)
	if err != nil {
		return "", err
	}

	//time.Sleep(time.Duration(5) * time.Second)

	for i := 0; i < 30; i++ {
		logs.Info("等待跨服节点(%s)数据写入[%d]......", newNode, i)
		time.Sleep(time.Duration(1) * time.Second)
		isExists := IsServerNodeExists(newNode)
		if isExists == true {
			break
		}
	}
	logs.Info("跨服节点(%s)数据写入成功.", newNode)

	err = InstallNode(newNode)
	utils.CheckError(err, "部署节点失败")
	if err != nil {
		return "", err
	}

	err = NodeAction([] string{newNode}, "start")
	utils.CheckError(err, "启动节点失败")
	if err != nil {
		return "", err
	}

	return newNode, nil

}

func GetFreeZone() (string, error) {
	serverNode, intZid, err := GetLatestZone()
	utils.CheckError(err)
	if err != nil {
		return "", err
	}

	connectCount := GetZoneConnectNodeCount(serverNode.Node)
	logs.Info("最新的跨服节点:%s, 连接的游戏节点个数:%d", serverNode.Node, connectCount)
	if connectCount <= 2 {
		return serverNode.Node, nil
	}
	configDbHost := beego.AppConfig.String("config_db_host")
	newIntZid := intZid + 1
	newNode := fmt.Sprintf("z%d@%s", newIntZid, serverNode.Ip)
	logs.Info("新跨服节点:%s", newNode)
	out, err := AddServerNode(newNode, serverNode.Ip, 0, 0, 2, "", configDbHost, 3306, fmt.Sprintf("db_zone_z%d", newIntZid))
	utils.CheckError(err, "新增节点失败:"+out)
	if err != nil {
		return "", err
	}

	//time.Sleep(time.Duration(5) * time.Second)

	for i := 0; i < 30; i++ {
		logs.Info("等待跨服节点(%s)数据写入[%d]......", newNode, i)
		time.Sleep(time.Duration(1) * time.Second)
		isExists := IsServerNodeExists(newNode)
		if isExists == true {
			break
		}
	}
	logs.Info("跨服节点(%s)数据写入成功.", newNode)

	err = InstallNode(newNode)
	utils.CheckError(err, "部署节点失败")
	if err != nil {
		return "", err
	}

	err = NodeAction([] string{newNode}, "start")
	utils.CheckError(err, "启动节点失败")
	if err != nil {
		return "", err
	}

	return newNode, nil

}

// 获取跨服节点连接的数量
func GetZoneConnectNodeCount(node string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM c_server_node WHERE type = 1 and zone_node = '%s'`, node)
	err := DbCenter.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}
func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}
	// 返回子串
	return string(rs[begin:end])
}

func AutoCreateAndOpenServer(platformId string,isCheck bool, openServerTime int) error {
	//logs.Info("1")
	//logs.Info("%+v", IsNowOpenServerMap)
	isNowOpenServer, ok := IsNowOpenServerMap[platformId]
	if ok {
		if isNowOpenServer == true {
			logs.Warning("正在开服中:%s !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", platformId)
			return nil
		}
	}
	IsNowOpenServerMap[platformId] = true
	//logs.Info("2")
	//logs.Info("%+v", IsNowOpenServerMap)
	defer func() {
		IsNowOpenServerMap[platformId] = false
	}()
	t0 := time.Now()

	platform, err := GetPlatformOne(platformId)
	utils.CheckError(err, "获取平台失败:" + platformId)
	if err != nil {
		return err
	}


	//logs.Info("3")
	//openServerRoleCount, err := beego.AppConfig.Int("open_server_create_role_count")
	//utils.CheckError(err, "读取自动开服人数配置失败")
	//if err != nil {
	//	return err
	//}

	//configDbHost := beego.AppConfig.String("config_db_host")
	//if configDbHost == "" {
	//	logs.Error("读取配置游戏服连接的数据库配置失败")
	//	return err
	//}

	if isCheck {
		if platform.IsAutoOpenServer == 0 {
			//不自动开服
			return nil
		}
		//logs.Info("检测自动开服:%s......", platformId)
		// 检测是否满足开服条件
		//isAutoOpenServer, err := beego.AppConfig.Bool("is_auto_open_server")
		//utils.CheckError(err, "读取是否开启自动开服配置失败")
		//if err != nil {
		//	return err
		//}
		//if isAutoOpenServer == false {
		//	return err
		//}

		if time.Now().Hour() >= 22 {
			//logs.Info("晚上10点后不自动开服")
			return nil
		}
	} else {
		logs.Info("立即开服:%s......", platformId)
	}
	//logs.Debug("platform:%+v", platform)

	if platform.CreateRoleLimit == 0 {
		err = errors.New(fmt.Sprintf("自动开服人数配置错误:%s, %d", platform.Id, platform.CreateRoleLimit))
		utils.CheckError(err)
		return err
	}


	inventoryDatabase, err:= GetInventoryDatabaseOne(platform.InventoryDatabaseId)
	utils.CheckError(err, "获取数据库配置失败")
	if err != nil {
		return err
	}



	maxGameServer, intSid, err := GetMaxGameServerByPlatformId(platformId)
	utils.CheckError(err, "获取最大区服失败")
	if err != nil {
		return err
	}
	//logs.Info("最大区服:%+v", maxGameServer)
	//logs.Info("最大区服ID:%+v(%d)", maxGameServer.Sid, intSid)
	gameDb, err := GetGameDbByNode(maxGameServer.Node)
	utils.CheckError(err, "连接游戏服数据库失败")
	if err != nil {
		return err
	}
	defer gameDb.Close()
	createRoleCount := GetTotalCreateRoleCount(gameDb)
	//logs.Info("最新区服:%s, 当前创角:%d, 创角临界值:%d", maxGameServer.Sid, createRoleCount, platform.CreateRoleLimit)
	newIntSid := intSid + 1
	newSid := fmt.Sprintf("s%d", newIntSid)
	if isCheck == false || createRoleCount >= platform.CreateRoleLimit {
		logs.Info("*************************** 开服部署新服 %s - %s *****************************\n", platformId, newSid)

		//serverNode, err := GetServerNode(maxGameServer.Node)
		//utils.CheckError(err, "获取节点失败!!")
		//if err != nil {
		//	return err
		//}
		maxFreeServer, err := GetMaxFreeServerByPlatformId(platformId)
		utils.CheckError(err)
		if err != nil {
			return err
		}

		logs.Info("最空闲的服务器:%+v", maxFreeServer)
		logs.Info("新区服id:%s", newSid)
		newNode := fmt.Sprintf("%s_%s@%s", platformId, newSid, maxFreeServer.InnerIp)
		logs.Info("新区服节点:%s", newNode)

		//maxPort := Max(serverNode.Port, serverNode.WebPort)
		maxPort := GetThisIpMaxPort(maxFreeServer.InnerIp)
		//logs.Info("最大端口:%d", maxPort)
		out, err := AddServerNode(newNode, maxFreeServer.Host, maxPort+1, maxPort+2, 1, platformId, inventoryDatabase.Host, 3306, fmt.Sprintf("db_%s_game_%s", platformId, newSid))
		utils.CheckError(err, "新增节点失败:"+out)
		if err != nil {
			return err
		}
		logs.Info("添加节点成功:%s", newNode)
		zoneNode, err := GetFreeZoneByPlatformId(platformId)
		utils.CheckError(err, "获取空闲跨服节点失败:"+out)
		if err != nil {
			return err
		}

		out, err = AddGameServer(platformId, newSid, fmt.Sprintf("%d区", newIntSid), newNode, zoneNode, 3, openServerTime, 1)

		utils.CheckError(err, "新增游戏服失败:"+out)
		if err != nil {
			return err
		}
		logs.Info("添加game_server成功:%s", newSid)
		//time.Sleep(time.Duration(15) * time.Second)

		for i := 0; i < 30; i++ {
			logs.Info("等待节点(%s)数据写入[%d]......", newNode, i)
			time.Sleep(time.Duration(1) * time.Second)
			isExists := IsServerNodeExists(newNode)
			if isExists == true {
				break
			}
		}
		logs.Info("节点(%s)数据写入成功.", newNode)

		err = InstallNode(newNode)
		utils.CheckError(err, "部署节点失败")
		if err != nil {
			return err
		}

		err = NodeAction([] string{newNode}, "start")
		utils.CheckError(err, "启动节点失败")
		if err != nil {
			return err
		}


		err = AfterAddGameServer()
		utils.CheckError(err, "同步登录充值战区失败")
		if err != nil {
			return err
		}

		err = NodeAction([] string{zoneNode}, "pull")
		utils.CheckError(err, "更新跨服节点数据")
		if err != nil {
			return err
		}



		err = CreateAnsibleInventory()
		utils.CheckError(err, "生成ansible inventory失败")
		if err != nil {
			return err
		}
		usedTime := time.Since(t0)
		logs.Info("************************ 自动开服成功:%s - %s 耗时:%s **********************", platformId, newSid, usedTime.String())
	} else {
		//logs.Info("不满足开服条件.")
	}
	return nil
}

func OpenServerNow(platformId string) {
	logs.Info("现在开服:%s", platformId)
	err := AutoCreateAndOpenServer(platformId,false, utils.GetTimestamp())
	utils.CheckError(err, "开服失败!!!!!!" + platformId)
}
//func AutoCreateAndOpenServer(isCheck bool) error {
//	if IsNowOpenServer == true {
//		logs.Warning("正在开服中!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
//		return nil
//	}
//	t0 := time.Now()
//
//
//	openServerRoleCount, err := beego.AppConfig.Int("open_server_create_role_count")
//	utils.CheckError(err, "读取自动开服人数配置失败")
//	if err != nil {
//		return err
//	}
//
//	configDbHost := beego.AppConfig.String("config_db_host")
//	if configDbHost == "" {
//		logs.Error("读取配置游戏服连接的数据库配置失败")
//		return err
//	}
//
//	if isCheck {
//		logs.Info("检测自动开服......")
//		// 检测是否满足开服条件
//		isAutoOpenServer, err := beego.AppConfig.Bool("is_auto_open_server")
//		utils.CheckError(err, "读取是否开启自动开服配置失败")
//		if err != nil {
//			return err
//		}
//		if isAutoOpenServer == false {
//			return err
//		}
//
//		if time.Now().Hour() >= 22 {
//			logs.Info("晚上10点后不自动开服")
//			return err
//		}
//	} else {
//		logs.Info("立即开服......")
//	}
//
//	now := utils.GetTimestamp()
//
//	maxGameServer, intSid, err := GetMaxGameServer()
//	utils.CheckError(err, "获取最大区服失败")
//	if err != nil {
//		return err
//	}
//	//logs.Info("最大区服:%+v", maxGameServer)
//	//logs.Info("最大区服ID:%+v(%d)", maxGameServer.Sid, intSid)
//	gameDb, err := GetGameDbByNode(maxGameServer.Node)
//	utils.CheckError(err, "连接游戏服数据库失败")
//	if err != nil {
//		return err
//	}
//	defer gameDb.Close()
//	createRoleCount := GetTotalCreateRoleCount(gameDb)
//	logs.Info("最新区服:%s, 当前创角:%d, 创角临界值:%d", maxGameServer.Sid, createRoleCount, openServerRoleCount)
//
//	if isCheck == false || createRoleCount >= openServerRoleCount {
//		IsNowOpenServer = true
//		defer func() {
//			IsNowOpenServer = false
//		}()
//		logs.Info("*************************** 开服部署新服 *****************************\n")
//		newIntSid := intSid + 1
//		newSid := fmt.Sprintf("s%d", newIntSid)
//		serverNode, err := GetServerNode(maxGameServer.Node)
//		utils.CheckError(err, "获取节点失败!!")
//		if err != nil {
//			return err
//		}
//		maxFreeServer, err := GetMaxFreeServer()
//		utils.CheckError(err)
//		if err != nil {
//			return err
//		}
//
//		logs.Info("最空闲的服务器:%+v", maxFreeServer)
//		logs.Info("新服id:%s", newSid)
//		newNode := fmt.Sprintf("%s_%s@%s", serverNode.PlatformId, newSid, maxFreeServer.InnerIp)
//		logs.Info("新节点:%s", newNode)
//
//		//maxPort := Max(serverNode.Port, serverNode.WebPort)
//		maxPort := GetThisIpMaxPort(maxFreeServer.InnerIp)
//		logs.Info("最大端口:%d", maxPort)
//		out, err := AddServerNode(newNode, maxFreeServer.Host, maxPort+1, maxPort+2, 1, serverNode.PlatformId, configDbHost, 3306, fmt.Sprintf("db_%s_game_%s", serverNode.PlatformId, newSid))
//		utils.CheckError(err, "新增节点失败:"+out)
//		if err != nil {
//			return err
//		}
//
//		zoneNode, err := GetFreeZone()
//		utils.CheckError(err, "获取空闲跨服节点失败:"+out)
//		if err != nil {
//			return err
//	}
//
//		out, err = AddGameServer(maxGameServer.PlatformId, newSid, fmt.Sprintf("%d区", newIntSid), newNode, zoneNode, 3, now, 1)
//
//		utils.CheckError(err, "新增游戏服失败:"+out)
//		if err != nil {
//			return err
//		}
//
//		//time.Sleep(time.Duration(15) * time.Second)
//
//		for i := 0; i < 30; i++ {
//			logs.Info("等待节点(%s)数据写入[%d]......", newNode, i)
//			time.Sleep(time.Duration(1) * time.Second)
//			isExists := IsServerNodeExists(newNode)
//			if isExists == true {
//				break
//			}
//		}
//		logs.Info("节点(%s)数据写入成功.", newNode)
//
//		err = InstallNode(newNode)
//		utils.CheckError(err, "部署节点失败")
//		if err != nil {
//			return err
//		}
//
//		err = NodeAction([] string{newNode}, "start")
//		utils.CheckError(err, "启动节点失败")
//		if err != nil {
//			return err
//		}
//
//
//		err = RefreshGameServer()
//		utils.CheckError(err, "刷新区服入口失败")
//		if err != nil {
//			return err
//		}
//
//		err = NodeAction([] string{zoneNode}, "pull")
//		utils.CheckError(err, "更新跨服节点数据")
//		if err != nil {
//			return err
//		}
//
//
//
//		err = CreateAnsibleInventory()
//		utils.CheckError(err, "生成ansible inventory失败")
//		if err != nil {
//			return err
//		}
//		usedTime := time.Since(t0)
//		logs.Info("************************ 自动开服成功:%s 耗时:%s **********************", newSid, usedTime.String())
//	} else {
//		logs.Info("不满足开服条件.")
//	}
//	return nil
//}

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
// 获取该平台所有游戏服
func GetGameServerByPlatformId( platformId string) [] *GameServer {
	data := make([]*GameServer, 0)
	err := DbCenter.Model(&GameServer{}).Where(&GameServer{
		PlatformId: platformId,
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
