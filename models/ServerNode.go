package models

import (
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	//"github.com/zaaksam/dproxy/go/db"
)

type ServerNodeQueryParam struct {
	BaseQueryParam
	Type       int
	Node       string
	PlatformId int `json:"platformId"`
}

type ServerNode struct {
	Node          string `gorm:"primary_key" json:"node"`
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	WebPort       int    `json:"webPort"`
	Type          int    `json:"type"`
	ZoneNode      string `json:"zoneNode"`
	ServerVersion string `json:"serverVersion" gorm:"-"`
	IsAdd int `json:"isAdd" gorm:"-"`
	//ClientVersion string `json:"clientVersion"`
	OpenTime int `json:"openTime"`
	//IsTest        int    `json:"isTest"`
	PlatformId int `json:"platformId"`
	State      int `json:"state"`
	RunState   int `json:"runState"`
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
		}
	}
	return data, count
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

func GetNodeVersion(node string) string {
	//return "nullddd"
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return "null"
	}
	defer gameDb.Close()
	var data struct {
		Version string
	}

	sql := fmt.Sprintf(
		`SELECT str_data as version FROM server_data where id = 0 `)
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	if err != nil {
		return "null"
	}
	return data.Version
}

// 获取所有游戏节点
func GetAllGameServerNode() []*ServerNode {
	data := make([]*ServerNode, 0)
	err := DbCenter.Model(&ServerNode{}).Where(&ServerNode{
		Type: 1,
	}).Find(&data).Error
	utils.CheckError(err)
	return data
}

// 获取所有游戏节点
func GetAllGameServerNodeByPlatformId(platformId int) []*ServerNode {
	data := make([]*ServerNode, 0)
	err := DbCenter.Model(&ServerNode{}).Where(&ServerNode{
		Type:       1,
		PlatformId: platformId,
	}).Find(&data).Error
	utils.CheckError(err)
	return data
}

// 获取所有游戏节点
func GetAllGameNodeByPlatformId(platformId int) []string {
	data := make([]string, 0)
	serverNodeList := GetAllGameServerNodeByPlatformId(platformId)
	for _, e := range serverNodeList {
		data = append(data, e.Node)
	}
	return data
}
