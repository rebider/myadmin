package models

import (
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"github.com/astaxie/beego"
	//"github.com/zaaksam/dproxy/go/db"
	"github.com/astaxie/beego/logs"
)

type ServerNodeQueryParam struct {
	BaseQueryParam
	Type       int
	Node       string
	PlatformId string `json:"platformId"`
}

type ServerNode struct {
	Node          string `gorm:"primary_key" json:"node"`
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	WebPort       int    `json:"webPort"`
	DbHost        string `json:"dbHost"`
	DbPort        int    `json:"dbPort"`
	DbName        string `json:"dbName"`
	Type          int    `json:"type"`
	ZoneNode      string `json:"zoneNode"`
	ServerVersion int    `json:"serverVersion" gorm:"-"`
	DbVersion     int    `json:"dbVersion" gorm:"-"`
	IsAdd         int    `json:"isAdd" gorm:"-"`
	//ClientVersion string `json:"clientVersion"`
	OpenTime  int `json:"openTime"`
	StartTime int `json:"startTime" gorm:"-"`
	//IsTest        int    `json:"isTest"`
	PlatformId string `json:"platformId"`
	State      int    `json:"state"`
	RunState   int    `json:"runState"`
}

func (t *ServerNode) TableName() string {
	return "c_server_node"
}

//获取分页数据
func ServerNodePageList(params *ServerNodeQueryParam) ([]*ServerNode, int64) {
	data := make([]*ServerNode, 0)
	//默认排序
	sortOrder := "node"
	switch params.Sort {
	case "node":
		sortOrder = "node"
	}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := DbCenter.Model(&ServerNode{}).Where(&ServerNode{
		Type:       params.Type,
		Node:       params.Node,
		PlatformId: params.PlatformId,
	}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Error
	utils.CheckError(err)
	if err == nil {
		for _, e := range data {
			e.ServerVersion = GetNodeVersion(e.Node)
			e.DbVersion = GetDbVersion(e.Node)
			e.StartTime = GetNodeStartTime(e.Node)
		}
	}
	return data, count
}

func GetAllServerNodeList() ([]*ServerNode) {
	data := make([]*ServerNode, 0)
	err := DbCenter.Model(&ServerNode{}).Find(&data).Error
	utils.CheckError(err)
	return data
}

func GetServerNode(node string) (*ServerNode, error) {
	serverNode := &ServerNode{
		Node: node,
	}
	err := DbCenter.First(&serverNode).Error
	return serverNode, err
}

func IsServerNodeExists(node string) bool {
	serverNode := &ServerNode{
		Node: node,
	}
	return ! DbCenter.First(&serverNode).RecordNotFound()
}

// 获取节点版本
func GetNodeVersion(node string) int {
	//return "nullddd"
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	defer gameDb.Close()
	var data struct {
		Version int
	}

	sql := fmt.Sprintf(
		`SELECT int_data as version FROM server_data where id = 0 `)
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	return data.Version
}

//获取节点合服时间
func GetMergeTime(node string) int {
	//return "nullddd"
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	defer gameDb.Close()
	var data struct {
		Time int
	}

	sql := fmt.Sprintf(
		`SELECT int_data as time FROM server_data where id = 6 `)
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	return data.Time
}

func GetDbVersion(node string) int {
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return -1
	}
	defer gameDb.Close()
	var data struct {
		Version int
	}

	sql := fmt.Sprintf(
		`SELECT version FROM db_version`)
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	if err != nil {
		return -1
	}
	return data.Version
}

func BackDatabases() {
	list := GetAllServerNodeList()
	gameDbPwd := beego.AppConfig.String( "game_db_password")
	dbBackupDir := beego.AppConfig.String( "db_backup_dir")
	if dbBackupDir == "" {
		logs.Error("数据库备份路径未配置!!!!!!!!!!!!!!!!!!!!!!!")
		return
	}
	for _, e := range list {
		user := "root"
		dbName := e.DbName
		host := e.DbHost
		commandArgs := []string{
			" -u ",
			user,
			" -h ",
			host,
			//" -p ",
			//gameDbPwd,
			fmt.Sprintf(" -p%s ", gameDbPwd),
			dbName,
			" > ",
			fmt.Sprintf(" %s%s_%d ", dbBackupDir, dbName, utils.GetTimestamp()),
		}
		logs.Info("开始备份数据库:%s......", dbName)
		out, err := utils.Cmd("mysqldump", commandArgs)
		if err != nil {
			logs.Error("备份数据库失败:%s!!!!!, %s", dbName, out)
		} else {
			logs.Info("备份数据库成功:%s.", dbName)
		}
	}
}
func GetNodeStartTime(node string) int {
	//return "nullddd"
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	defer gameDb.Close()
	var data struct {
		Time int
	}

	sql := fmt.Sprintf(
		`SELECT int_data as time FROM server_data where id = 3 `)
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err, "获取节点启动时间")
	if err != nil {
		return 0
	}
	return data.Time
}

// 获取所有游戏节点
func GetAllGameServerNode() []*ServerNode {
	//data := make([]*ServerNode, 0)
	//err := DbCenter.Model(&ServerNode{}).Where(&ServerNode{
	//	Type: 1,
	//}).Find(&data).Error
	//utils.CheckError(err, "查询所有游戏节点失败")
	return GetAllServerNodeByType(1)
}

// 获取节点列表
func GetAllServerNodeByType(nodeType int) []*ServerNode {
	data := make([]*ServerNode, 0)
	err := DbCenter.Model(&ServerNode{}).Where(&ServerNode{
		Type: nodeType,
	}).Find(&data).Error
	utils.CheckError(err, "获取节点列表失败")
	return data
}


// 获取所有游戏节点
func GetAllGameServerNodeByPlatformId(platformId string) []*ServerNode {
	data := make([]*ServerNode, 0)
	err := DbCenter.Model(&ServerNode{}).Where(&ServerNode{
		Type:       1,
		PlatformId: platformId,
	}).Find(&data).Error
	utils.CheckError(err)
	return data
}

// 获取所有游戏节点
func GetAllGameNodeByPlatformId(platformId string) []string {
	data := make([]string, 0)
	serverNodeList := GetAllGameServerNodeByPlatformId(platformId)
	for _, e := range serverNodeList {
		data = append(data, e.Node)
	}
	return data
}

// 获取登录节点 14101 11101
//func GetLoginServerNode() (*ServerNode, error) {
//	serverNode := &ServerNode{}
//	err := DbCenter.Where(&ServerNode{
//		Type: 4,
//	}).First(&serverNode).Error
//	return serverNode, err
//}

func CheckAllGameNode() {
	isCheck := beego.AppConfig.DefaultBool("is_check_node", false)
	if isCheck {
		now := utils.GetTimestamp()
		serverNodeList := GetAllServerNodeList()
		num := 0
		for _, e := range serverNodeList {
			if e.Type == 1 && e.OpenTime + 60 > now {
				logs.Debug("忽略:%s", e.Node)
				continue
			}
			if e.RunState == 0 {
				logs.Info("节点未开启:~p", e.Node)
				num ++
			}
		}
		if num > 0 {
			logs.Error("%d个节点未开启!!!!!!!!!!!!", num)
			utils.ReportMsg("105138", "13616005067")
			utils.ReportMsg("105138", "19905929917")
		}
	}
}
