package models

import (
	//"encoding/json"
	//"io/ioutil"
	//"fmt"
	//"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
)

type Platform struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Time int `json:"time"`
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

//获取用户列表
func GetPlatformList() []*Platform {
	data := make([]*Platform, 0)
	err := Db.Model(&Platform{}).Find(&data).Error
	utils.CheckError(err)
	return data
}

// 获取单个用户
//func GetPlatformOne(id int) (*Platform, error) {
//	r := &Platform{
//		Id: id,
//	}
//	err := Db.First(&r).Error
//	return r, err
//}

// 删除用户列表
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
