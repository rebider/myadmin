package models

import (
	//"github.com/astaxie/beego/orm"
	//"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
)

type ServerNodeQueryParam struct {
	BaseQueryParam
	Type       int
	Node       string
	PlatformId int `json:"platformId"`
}
type ServerNode struct {
	Node          string `orm:"pk" json:"node"`
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	Type          int    `json:"type"`
	ZoneNode      string `json:"zoneNode"`
	ServerVersion string `json:"serverVersion"`
	ClientVersion string `json:"clientVersion"`
	OpenTime      int    `json:"openTime"`
	IsTest        int    `json:"isTest"`
	PlatformId    int    `json:"platformId"`
	State int `json:"state"`
}

func (t *ServerNode) TableName() string {
	return "c_server_node"
}

//获取分页数据
func ServerNodePageList(params *ServerNodeQueryParam) ([]*ServerNode, int64) {
	logs.Debug("params:%+v", params)
	db, err := GetCenterDb()
	utils.CheckError(err)
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
	if params.Type > 0 {
		db = db.Where("type = ?", params.Type)
	}
	if params.Node != "" {
		db = db.Where("node = ?", params.Node)
	}
	if params.PlatformId > 0 {
		db = db.Where("platform_id = ?", params.PlatformId)
	}
	var count int64
	err = db.Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Count(&count).Error
	utils.CheckError(err)
	return data, count
}

func GetServerNode(node string) (*ServerNode, error) {
	serverNode := &ServerNode{
		Node: node,
	}
	DbCenter.First(&serverNode)
	return serverNode, nil
}
