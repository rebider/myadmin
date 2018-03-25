package models

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/utils"
	"github.com/jinzhu/gorm"
	"github.com/astaxie/beego/logs"
)

type Result struct {
	Code enums.JsonResultCode `json:"code"`
	Msg  string               `json:"msg"`
	Data interface{}          `json:"data"`
}

type BaseQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

func GetDbByPlatformIdAndSid(platformId int, Sid string) (db *gorm.DB, err error) {
	gameServer, err := GetGameServer(platformId, Sid)
	utils.CheckError(err)
	serverNode, err := GetServerNode(gameServer.Node)
	utils.CheckError(err)
	logs.Debug(serverNode.Ip)
	dbArgs := "root:game1234@tcp(" + serverNode.Ip + ":3306)/game?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", dbArgs)
	return db, err
}

//func GetCenterDb() (db *gorm.DB, err error) {
//	dbArgs := "root:game1234@tcp(192.168.31.100:3306)/center?charset=utf8&parseTime=True&loc=Local"
//	db, err = gorm.Open("mysql", dbArgs)
//	return db, err
//}
