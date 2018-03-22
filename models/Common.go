package models

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/utils"
	"github.com/jinzhu/gorm"
	"github.com/astaxie/beego/logs"
)

type JsonResult struct {
	Code enums.JsonResultCode `json:"code"`
	Msg  string               `json:"msg"`
	Obj  interface{}          `json:"obj"`
}

type Result struct {
	Code enums.JsonResultCode `json:"code"`
	Msg  string               `json:"msg"`
	Data  interface{}          `json:"data"`
}

type BaseQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int64  `json:"offset"`
	Limit  int    `json:"limit"`
}
//var (
//	DbCenter *gorm.DB
//)

//func init() {
//	dbArgs := "root:game1234@tcp(192.168.31.100:3306)/center?charset=utf8&parseTime=True&loc=Local"
//	var err error
//	DbCenter, err = gorm.Open("mysql", dbArgs)
//	utils.CheckError(err)
//}
//
//func Wheres(db *gorm.DB) *gorm.DB{
//
//}
func GetDbByPlatformIdAndSid(platformId int, Sid string) (db *gorm.DB, err error){
	gameServer, err := GetGameServer(platformId, Sid)
	logs.Debug("%v%v", platformId, Sid)
	logs.Debug("%v", gameServer)
	utils.CheckError(err)
	serverNode, err := GetServerNode(gameServer.Node)
	logs.Debug("%v", serverNode)
	utils.CheckError(err)
	logs.Debug(serverNode.Ip)
	dbArgs := "root:game1234@tcp(" +serverNode.Ip+ ":3306)/game?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", dbArgs)
	return db, err
}

func GetCenterDb() (db *gorm.DB, err error){
	dbArgs := "root:game1234@tcp(192.168.31.100:3306)/center?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", dbArgs)
	return db, err
}
