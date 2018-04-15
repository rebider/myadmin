package models

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/utils"
	"github.com/jinzhu/gorm"
	"fmt"
	"golang.org/x/net/websocket"
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
func GetGameDbByPlatformIdAndSid(platformId int, Sid string) (gameDb *gorm.DB, err error) {
	gameServer, err := GetGameServerOne(platformId, Sid)
	utils.CheckError(err)
	serverNode, err := GetServerNode(gameServer.Node)
	utils.CheckError(err)
	//logs.Debug(serverNode.Ip)
	dbArgs := "root:game1234@tcp(" + serverNode.Ip + ":3306)/game?charset=utf8&parseTime=True&loc=Local"
	gameDb, err = gorm.Open("mysql", dbArgs)
	gameDb.SingularTable(true)
	return gameDb, err
}

// 通过ServerNode获取 gorm.DB 实例
func GetGameDbByServerNode(serverNode *ServerNode) (gameDb *gorm.DB, err error) {
	//logs.Debug(serverNode.Ip)
	dbArgs := "root:game1234@tcp(" + serverNode.Ip + ":3306)/game?charset=utf8&parseTime=True&loc=Local"
	gameDb, err = gorm.Open("mysql", dbArgs)
	gameDb.SingularTable(true)
	return gameDb, err
}

// 通过平台id和区服id 获取ip地址和端口
func GetIpAndPortByPlatformIdAndSid(platformId int, Sid string) (string, int, error) {
	gameServer, err := GetGameServerOne(platformId, Sid)
	utils.CheckError(err)
	serverNode, err := GetServerNode(gameServer.Node)
	utils.CheckError(err)
	return serverNode.Ip, serverNode.Port, err
}

func GetWsByPlatformIdAndSid(platformId int, Sid string) (*websocket.Conn, error){
	ip, port, err := GetIpAndPortByPlatformIdAndSid(platformId, Sid)
	if err != nil {
		return  nil, err
	}
	wsUrl := fmt.Sprintf("ws://%s:%d", ip, port)
	ws, err := websocket.Dial(wsUrl, "", wsUrl)
	return ws, err
}
