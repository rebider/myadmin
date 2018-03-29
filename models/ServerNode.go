package models

import (
	"github.com/chnzrb/myadmin/utils"
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
	Type          int    `json:"type"`
	ZoneNode      string `json:"zoneNode"`
	ServerVersion string `json:"serverVersion"`
	ClientVersion string `json:"clientVersion"`
	OpenTime      int    `json:"openTime"`
	IsTest        int    `json:"isTest"`
	PlatformId    int    `json:"platformId"`
	State         int    `json:"state"`
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
	return data, count
}

func GetServerNode(node string) (*ServerNode, error) {
	serverNode := &ServerNode{
		Node: node,
	}
	err := DbCenter.First(&serverNode).Error
	return serverNode, err
}
