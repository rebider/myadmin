package models

import (
	//"encoding/json"
	//"io/ioutil"
	//"fmt"
	//"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"sort"
	//"github.com/astaxie/beego/logs"
)

type Platform struct {
	Id                        string                       `json:"id"`
	Name                      string                       `json:"name"`
	IsAutoOpenServer          int                          `json:"isAutoOpenServer"`
	InventoryDatabaseId       int                          `json:"inventoryDatabaseId"`
	ZoneInventoryServerId     int                          `json:"zoneInventoryServerId"`
	CreateRoleLimit           int                          `json:"createRoleLimit"`
	Version                   string                       `json:"version"`
	PlatformInventorySeverRel []*PlatformInventorySeverRel `json:"-"`
	ChannelList               []*Channel                   `json:"channelList"`
	InventorySeverIds         []int                        `json:"inventorySeverIds" gorm:"-"`
	Time                      int                          `json:"time"`
}

func (a *Platform) TableName() string {
	return PlatformDatabaseTBName()
}

func PlatformDatabaseTBName() string {
	return "platform"
}

type PlatformParam struct {
	BaseQueryParam
}

//type Platform struct {
//	Id   int    `json:"id"`
//	Name string `json:"name"`
//	//User string `json:"user"`
//	//Port int    `json:"port"`
//	//Host string `json:"host"`
//	//AddTime           int    `json:"addTime"`
//	Time        int    `json:"time"`
//}

//获取平台列表
func GetPlatformList() []*Platform {
	data := make([]*Platform, 0)
	err := Db.Model(&Platform{}).Find(&data).Error
	utils.CheckError(err)
	for _, v := range data {
		v.InventorySeverIds = make([] int, 0)

		err = Db.Model(&v).Related(&v.PlatformInventorySeverRel).Error
		utils.CheckError(err)

		for _, e := range v.PlatformInventorySeverRel {
			v.InventorySeverIds = append(v.InventorySeverIds, e.InventoryServerId)
		}
		sort.Ints(v.InventorySeverIds)
		v.ChannelList = GetChannelListByPlatformId(v.Id)
	}
	return data
}

//获取平台列表
func GetPlatformListByUserId(userId int) []*Platform {
	var list []*Platform
	var channelList []*Channel
	user, err := GetUserOne(userId)
	utils.CheckError(err)
	if user.IsSuper == 1 {
		list = GetPlatformList()
	} else {
		sql := fmt.Sprintf(`SELECT DISTINCT T2.*
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
		INNER JOIN %s AS T2 ON T2.id = T0.channel_id
		WHERE T1.user_id = ?`, RoleChannelRelTBName(), RoleUserRelTBName(), ChannelDatabaseTBName())
		rows, err := Db.Raw(sql, userId).Rows()
		defer rows.Close()
		utils.CheckError(err)
		for rows.Next() {
			var channel Channel
			Db.ScanRows(rows, &channel)
			//logs.Debug("channel:%+v", channel)
			channelList = append(channelList, &channel)
		}
		flag := make(map[string]bool)
		for _, v := range channelList {
			_, ok := flag[v.PlatformId]
			if ok {

			} else {
				flag[v.PlatformId] = true
				platform, err := GetPlatformOne(v.PlatformId)
				utils.CheckError(err)
				list = append(list, platform)
			}
		}
		for _, v := range list {
			for _, channel := range channelList {
				if channel.PlatformId == v.Id {
					v.ChannelList = append(v.ChannelList, channel)
				}
			}
		}
	}
	return list
}

////获取平台列表
//func GetPlatformListByUserId(userId int) []*Platform {
//	var list []*Platform
//	user, err := GetUserOne(userId)
//	utils.CheckError(err)
//	if user.IsSuper == 1 {
//		list = GetPlatformList()
//	} else {
//		sql := fmt.Sprintf(`SELECT DISTINCT T2.*
//		FROM %s AS T0
//		INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
//		INNER JOIN %s AS T2 ON T2.id = T0.platform_id
//		WHERE T1.user_id = ?`, RoleChannelRelTBName(), RoleUserRelTBName(), PlatformDatabaseTBName())
//		rows, err := Db.Raw(sql, userId).Rows()
//		defer rows.Close()
//		utils.CheckError(err)
//		for rows.Next() {
//			var platform Platform
//			Db.ScanRows(rows, &platform)
//			list = append(list, &platform)
//		}
//	}
//	return list
//}

//获取单个平台
func GetPlatformOne(id string) (*Platform, error) {
	r := &Platform{
		Id: id,
	}
	err := Db.First(&r).Error
	r.InventorySeverIds = make([] int, 0)

	err = Db.Model(&r).Related(&r.PlatformInventorySeverRel).Error
	utils.CheckError(err)

	for _, e := range r.PlatformInventorySeverRel {
		r.InventorySeverIds = append(r.InventorySeverIds, e.InventoryServerId)
	}
	sort.Ints(r.InventorySeverIds)
	return r, err
}

// 删除平台列表
func DeletePlatform(ids [] string) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&Platform{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//删除角色平台关系
	if _, err := DeleteRoleChannelRelByPlatformIdList(ids); err != nil {
		tx.Rollback()
		return err
	}
	//删除服务器平台关系
	if _, err := DeletePlatformInventorySeverRelByPlatformIdList(ids); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//
//func GetPlatformList()[] *Platform {
//	filename := "views/static/json/Platform.json"
//	bytes, err := ioutil.ReadFile(filename)
//	list := make([] *Platform, 0)
//	if err != nil {
//		fmt.Println("ReadFile: ", err.Error())
//		logs.Error("GetPlatformList:%v, %v", filename, err)
//		return nil, err
//	}
//
//	if err := json.Unmarshal(bytes, &list); err != nil {
//		logs.Error("Unmarshal json:%v, %v", filename, err)
//		return nil, err
//	}
//	return list, nil
//}
