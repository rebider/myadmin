package models

import (
	"github.com/chnzrb/myadmin/utils"
	"github.com/linclin/gopub/src/github.com/pkg/errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strings"
	"fmt"
	"time"
)

func (a *InventoryServer) TableName() string {
	return InventoryServerTBName()
}

func InventoryServerTBName() string {
	return TableName("inventory_server")
}

type InventoryServerParam struct {
	BaseQueryParam
	Type int
}
type InventoryServer struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	OutIp             string `json:"outIp"`
	InnerIp           string `json:"innerIp"`
	Host              string `json:"host"`
	Type              int    `json:"type"`
	MaxNodeCount      int    `json:"maxNodeCount"`
	NodeCount         int    `json:"nodeCount" gorm:"-"`
	OnlinePlayerCount int    `json:"onlinePlayerCount" gorm:"-"`
	AddTime           int    `json:"addTime"`
	UpdateTime        int    `json:"updateTime"`
}

//获取用户列表
func GetInventoryServerList(params *InventoryServerParam) ([]*InventoryServer, int64) {
	data := make([]*InventoryServer, 0)
	sortOrder := "id"
	switch params.Sort {
	case "id":
		sortOrder = "id"
	}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := Db.Model(&InventoryServer{}).Where(&InventoryServer{
		Type: params.Type,
	}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		e.NodeCount = GetIpNodeCount(e.InnerIp)
		e.OnlinePlayerCount = GetIpOnlinePlayerCount(e.InnerIp)
	}
	return data, count
}

func GetAllServerListOfGame() []*InventoryServer {
	data := make([]*InventoryServer, 0)
	err := Db.Model(&InventoryServer{}).Where(&InventoryServer{
		Type: 4,
	}).Find(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		e.NodeCount = GetIpNodeCount(e.InnerIp)
		e.OnlinePlayerCount = GetIpOnlinePlayerCount(e.InnerIp)
	}
	return data
}

func GetAllServerList() []*InventoryServer {
	data := make([]*InventoryServer, 0)
	err := Db.Model(&InventoryServer{}).Find(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		e.NodeCount = GetIpNodeCount(e.InnerIp)
		e.OnlinePlayerCount = GetIpOnlinePlayerCount(e.InnerIp)
	}
	return data
}

func GetMaxFreeServerByPlatformId(platformId string) (*InventoryServer, error) {
	platform, err := GetPlatformOne(platformId)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}

	err = Db.Model(&platform).Related(&platform.PlatformInventorySeverRel).Error
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	l := make([] *InventoryServer, 0)
	for _, v := range platform.PlatformInventorySeverRel {
		thisInventoryServer, err := GetInventoryServerOne(v.InventoryServerId)
		utils.CheckError(err)
		if err != nil {
			return nil, err
		}
		l = append(l, thisInventoryServer)
	}
	if len(l) == 0 {
		return nil, errors.New("没有空闲服务器")
	}
	minCount := -1
	inventoryServer := &InventoryServer{}
	for _, e := range l {
		logs.Debug("server:%s, nodeCount:%d, onlineCount:%d", e.Name, e.NodeCount, e.OnlinePlayerCount, e.MaxNodeCount)
		if e.NodeCount >= e.MaxNodeCount {
			// 一个服务器最多安装50个节点
			continue
		}
		// 一个节点当作10个在线来计算
		thisValue := e.OnlinePlayerCount + e.NodeCount*10

		if minCount == -1 {
			minCount = thisValue
			inventoryServer = e
		} else {

			if thisValue < minCount {
				minCount = thisValue
				inventoryServer = e
			}
		}
	}
	if minCount == -1 {
		return nil, errors.New("没有空闲的服务器")
	}
	return inventoryServer, nil
}

func GetMaxFreeServer() (*InventoryServer, error) {
	l := GetAllServerListOfGame()
	if len(l) == 0 {
		return nil, errors.New("没有空闲服务器")
	}
	minCount := -1
	inventoryServer := &InventoryServer{}
	for _, e := range l {
		logs.Debug("server:%s, nodeCount:%d, onlineCount:%d", e.Name, e.NodeCount, e.OnlinePlayerCount)
		if e.NodeCount >= 33 {
			// 一个服务器最多安装33个节点
			continue
		}
		// 一个节点当作10个在线来计算
		thisValue := e.OnlinePlayerCount + e.NodeCount*10

		if minCount == -1 {
			minCount = thisValue
			inventoryServer = e
		} else {

			if thisValue < minCount {
				minCount = thisValue
				inventoryServer = e
			}
		}
	}
	if minCount == -1 {
		return nil, errors.New("没有空闲的服务器")
	}
	return inventoryServer, nil
}

// 获取单个服务器
func GetInventoryServerOneDirty(id int) (*InventoryServer, error) {
	r := &InventoryServer{
		Id: id,
	}
	err := Db.First(&r).Error
	return r, err
}

// 获取单个服务器
func GetInventoryServerOne(id int) (*InventoryServer, error) {
	r := &InventoryServer{
		Id: id,
	}
	err := Db.First(&r).Error
	r.NodeCount = GetIpNodeCount(r.InnerIp)
	r.OnlinePlayerCount = GetIpOnlinePlayerCount(r.InnerIp)
	return r, err
}

// 删除用户列表
func DeleteInventoryServers(ids [] int) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&InventoryServer{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if _, err := DeletePlatformInventorySeverRelByInventoryServerIdList(ids); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func CreateAnsibleInventory() error {
	logs.Info("开始创建ansible inventory 文件.....")
	ansibleInventoryFile := beego.AppConfig.String("ansible_inventory_file")
	if ansibleInventoryFile == "" {
		logs.Error("读取配置ansible_inventory_file失败")
		return errors.New("读取配置ansible_inventory_file失败")
	}
	mapList := make(map[string][]string, 0)
	serverNodeList := GetAllServerNodeList()
	for _, e := range serverNodeList {
		array := strings.Split(e.Node, "@")
		if len(array) != 2 {
			return errors.New("解析节点名字失败:" + e.Node)
		}
		nodeName := array[0]
		nodeIp := array[1]
		//logs.Info("nodeName:%+v", nodeName)
		if v, ok := mapList[nodeIp]; ok {
			v = append(v, "'"+nodeName+"'")
			mapList[nodeIp] = v
		} else {
			mapList[nodeIp] = append(make([] string, 0), "'"+nodeName+"'")
		}
	}
	serverOfGameList := GetAllServerList()
	for _, e := range serverOfGameList {
		if _, ok := mapList[e.InnerIp]; ok {
		} else {
			logs.Info("e:%+v", e)
			mapList[e.InnerIp] = make([] string, 0)
		}
	}
	logs.Info("serverOfGameList:%+v", mapList)
	logs.Info("mapList:%+v", mapList)
	content := "## Generated automatically, no need to modify.\n"
	content += fmt.Sprintf("## Auto Created :%s\n\n", time.Now().String())
	//for ip, _ := range mapList {
	//	content += fmt.Sprintf("%s\n", ip)
	//}
	//content += "\n\n\n"

	for ip, nodes := range mapList {
		content += fmt.Sprintf("%s    ", ip)
		content += "nodes="
		content += "\"[" + strings.Join(nodes, ", ") + "]\""
		content += "\n\n"

	}

	err := utils.FilePutContext(ansibleInventoryFile, content)

	if err != nil {
		return err
	}
	logs.Info("创建ansible inventory 文件(%s)成功.", ansibleInventoryFile)
	return nil
}
//func CreateAnsibleInventory2() error {
//	logs.Info("开始创建ansible inventory2 文件.....")
//	ansibleInventoryDir := beego.AppConfig.String("ansible_inventory_dir")
//	if ansibleInventoryDir == "" {
//		logs.Error("ansible_inventory_dir")
//		return errors.New("读取配置ansible_inventory_dir失败")
//	}
//	mapList := make(map[string][]string, 0)
//	serverNodeList := GetAllServerNodeList()
//	typeList := [] struct{
//		nodeType int
//		name string
//	}
//	for _, e := range typeList {
//
//	}
//	for _, e := range serverNodeList {
//		array := strings.Split(e.Node, "@")
//		if len(array) != 2 {
//			return errors.New("解析节点名字失败:" + e.Node)
//		}
//		nodeName := array[0]
//		nodeIp := array[1]
//		//logs.Info("nodeName:%+v", nodeName)
//		if v, ok := mapList[nodeIp]; ok {
//			v = append(v, "'"+nodeName+"'")
//			mapList[nodeIp] = v
//		} else {
//			mapList[nodeIp] = append(make([] string, 0), "'"+nodeName+"'")
//		}
//	}
//	serverOfGameList := GetAllServerList()
//	for _, e := range serverOfGameList {
//		if _, ok := mapList[e.InnerIp]; ok {
//		} else {
//			logs.Info("e:%+v", e)
//			mapList[e.InnerIp] = make([] string, 0)
//		}
//	}
//	logs.Info("serverOfGameList:%+v", mapList)
//	logs.Info("mapList:%+v", mapList)
//	content := "## Generated automatically, no need to modify.\n"
//	content += fmt.Sprintf("## Auto Created :%s\n\n", time.Now().String())
//	//for ip, _ := range mapList {
//	//	content += fmt.Sprintf("%s\n", ip)
//	//}
//	//content += "\n\n\n"
//
//	for ip, nodes := range mapList {
//		content += fmt.Sprintf("%s    ", ip)
//		content += "nodes="
//		content += "\"[" + strings.Join(nodes, ", ") + "]\""
//		content += "\n\n"
//
//	}
//
//	err := utils.FilePutContext(ansibleInventoryFile, content)
//
//	if err != nil {
//		return err
//	}
//	logs.Info("创建ansible inventory 文件(%s)成功.", ansibleInventoryFile)
//	return nil
//}


//func CreateAnsibleInventory() error{
//	logs.Info("开始创建ansible inventory 文件.....")
//	ansibleInventoryFile := beego.AppConfig.String("ansible_inventory_file")
//	if ansibleInventoryFile == "" {
//		logs.Error("读取配置ansible_inventory_file失败")
//		return errors.New("读取配置ansible_inventory_file失败")
//	}
//	mapList := make(map[string] []string , 0)
//	serverNodeList := GetAllServerNodeList()
//	for _, e := range serverNodeList {
//		array := strings.Split(e.Node, "@")
//		if len(array) != 2 {
//			return errors.New("解析节点名字失败:" + e.Node)
//		}
//		nodeName := array[0]
//		nodeIp := array[1]
//		//logs.Info("nodeName:%+v", nodeName)
//		if v, ok := mapList[nodeIp]; ok {
//			v = append(v, "'" + nodeName + "'")
//			mapList[nodeIp] = v
//		} else {
//			mapList[nodeIp] = append(make([] string, 0), "'" + nodeName + "'")
//		}
//	}
//	//logs.Info("mapList:%+v", mapList)
//	content := "## Generated automatically, no need to modify.\n"
//	content += fmt.Sprintf("## Auto Created :%s\n\n", time.Now().String())
//	for ip, _ := range mapList {
//		content += fmt.Sprintf("%s\n", ip)
//	}
//	content += "\n\n\n"
//
//	for ip, nodes := range mapList {
//		content += fmt.Sprintf("[%s:vars]\n", ip)
//		content += "nodes="
//		content += "[" + strings.Join(nodes, ", ") + "]"
//		content += "\n\n"
//
//	}
//
//	err := utils.FilePutContext(ansibleInventoryFile, content)
//
//	if err != nil {
//		return err
//	}
//	logs.Info("创建ansible inventory 文件(%s)成功.", ansibleInventoryFile)
//	return nil
//}
