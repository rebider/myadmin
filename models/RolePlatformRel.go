package models

import "time"

//角色与平台关系表
type RolePlatformRel struct {
	Id         int
	RoleId     int
	PlatformId string
	Created    time.Time
}

func (a *RolePlatformRel) TableName() string {
	return RolePlatformRelTBName()
}

func RolePlatformRelTBName() string {
	return TableName("role_platform_rel")
}

// 删除角色平台关系
func DeleteRolePlatformRelByRoleIdList(roleIdList [] int) (int, error) {
	var count int
	err := Db.Where("role_id in (?)", roleIdList).Delete(&RolePlatformRel{}).Count(&count).Error
	return count, err
}

// 删除角色平台关系
func DeleteRolePlatformRelByPlatformIdList(platformIdList [] string) (int, error) {
	var count int
	err := Db.Where("platform_id in (?)", platformIdList).Delete(&RolePlatformRel{}).Count(&count).Error
	return count, err
}
