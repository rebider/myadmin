package models

import (
	"github.com/chnzrb/myadmin/utils"
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

// 获取单个用户
func GetInventoryServerOne(id int) (*InventoryServer, error) {
	r := &InventoryServer{
		Id: id,
	}
	err := Db.First(&r).Error
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
	return tx.Commit().Error
}
