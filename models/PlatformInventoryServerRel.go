package models

import "time"

//平台服务器关系表
type PlatformInventorySeverRel struct {
	Id                int
	PlatformId      string
	InventoryServerId int
	Created           time.Time
}

func (a *PlatformInventorySeverRel) TableName() string {
	return PlatformInventorySeverRelTBName()
}

func PlatformInventorySeverRelTBName() string {
	return "platform_inventory_server_rel"
}

// 删除平台服务器关系
func DeletePlatformInventorySeverRelByPlatformIdList(PlatformIdList [] string) (int, error) {
	var count int
	err := Db.Where("platform_id in (?)", PlatformIdList).Delete(&PlatformInventorySeverRel{}).Count(&count).Error
	return count, err
}

// 删除平台服务器关系
func DeletePlatformInventorySeverRelByInventoryServerIdList(InventoryServerIdList [] int) (int, error) {
	var count int
	err := Db.Where("inventory_server_id in (?)", InventoryServerIdList).Delete(&PlatformInventorySeverRel{}).Count(&count).Error
	return count, err
}
