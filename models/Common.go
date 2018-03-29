package models

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/utils"
	"github.com/jinzhu/gorm"
)

type Result struct {
	Code enums.ResultCode `json:"code"`
	Msg  string           `json:"msg"`
	Data interface{}      `json:"data"`
}

type BaseQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

// 通过平台id和区服id获取 gorm.DB 实例
func GetDbByPlatformIdAndSid(platformId int, Sid string) (db *gorm.DB, err error) {
	gameServer, err := GetGameServerOne(platformId, Sid)
	utils.CheckError(err)
	serverNode, err := GetServerNode(gameServer.Node)
	utils.CheckError(err)
	//logs.Debug(serverNode.Ip)
	dbArgs := "root:game1234@tcp(" + serverNode.Ip + ":3306)/game?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", dbArgs)
	db.SingularTable(true)
	return db, err
}

//func GetCenterDb() (db *gorm.DB, err error) {
//	dbArgs := "root:game1234@tcp(192.168.31.100:3306)/center?charset=utf8&parseTime=True&loc=Local"
//	db, err = gorm.Open("mysql", dbArgs)
//	return db, err
//}
