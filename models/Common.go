package models

import (
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/utils"
	"github.com/jinzhu/gorm"
	"fmt"
	"golang.org/x/net/websocket"
	"strings"
	"github.com/astaxie/beego/logs"
	"errors"
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
	if err != nil {
		return nil, err
	}
	serverNode, err := GetServerNode(gameServer.Node)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	return GetGameDbByNode(serverNode.Node)
}

// 通过平台id和区服id获取 gorm.DB 实例
func GetGameURLByNode(node string) string {
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	if err != nil {
		return ""
	}
	url := fmt.Sprintf("http://%s:%d", serverNode.Ip, serverNode.WebPort)
	return url
}


// 通过平台id和区服id获取 gorm.DB 实例
func GetGameURLByPlatformIdAndSid(platformId int, Sid string) string {
	gameServer, err := GetGameServerOne(platformId, Sid)
	utils.CheckError(err)
	if err != nil {
		return ""
	}
	serverNode, err := GetServerNode(gameServer.Node)
	utils.CheckError(err)
	if err != nil {
		return ""
	}
	url := fmt.Sprintf("http://%s:%d", serverNode.Ip, serverNode.WebPort)
	return url
}


// 通过node获取 gorm.DB 实例
func GetGameDbByNode(node string) (gameDb *gorm.DB, err error) {
	if node == "" {
		return nil, errors.New("节点不能未空")
	}
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	array := strings.Split(serverNode.Node, "@")
	if len(array) != 2 {
		return nil, errors.New("解析节点名字失败:" + serverNode.Node)
	}
	//gameDbName := "game_" + array[0]
	gameDbName := "game_weixin"
	dbArgs := "root:game1234@tcp(" + serverNode.Ip + ":3306)/" + gameDbName +"?charset=utf8&parseTime=True&loc=Local"
	gameDb, err = gorm.Open("mysql", dbArgs)
	if err != nil {
		return nil, err
	}
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

func GetWsByPlatformIdAndSid(platformId int, Sid string) (*websocket.Conn, error) {
	ip, port, err := GetIpAndPortByPlatformIdAndSid(platformId, Sid)
	if err != nil {
		return nil, err
	}
	wsUrl := fmt.Sprintf("ws://%s:%d", ip, port)
	ws, err := websocket.Dial(wsUrl, "", wsUrl)
	return ws, err
}

func GetWsByNode(node string) (*websocket.Conn, error) {
	serverNode, err := GetServerNode(node)
	if err != nil {
		return nil, err
	}
	wsUrl := fmt.Sprintf("ws://%s:%d", serverNode.Ip, serverNode.Port)
	ws, err := websocket.Dial(wsUrl, "", wsUrl)
	return ws, err
}


type Server struct {
	PlatformId int    `json:"platformId"`
	Sid        string `json:"serverId"`
	Desc       string `json:"desc"`
	//Node       string `json:"node"`
}

func GetServerList() [] *Server {
	serverList := make([] *Server, 0)
	gameServerNodeList := GetAllGameServerNode()
	for _, serverNode := range gameServerNodeList {
		gameServerList := GetGameServerByNode(serverNode.Node)
		if len(gameServerList) > 0 {
			serverIdList := make([] string, 0)
			for _, gameServer := range gameServerList {
				serverIdList = append(serverIdList, gameServer.Sid)
			}
			var desc string
			if len(serverIdList) > 1 {
				// 合服
				desc = "合:[" + strings.Join(serverIdList, ", ") + "]"
			} else {
				desc = strings.Join(serverIdList, ", ")
			}
			serverList = append(serverList, &Server{
				PlatformId: gameServerList[0].PlatformId,
				Sid:        serverNode.Node,
				Desc:       desc,
			})
		} else {
			logs.Warning("节点没有对应的游戏服:%v", serverNode.Node)
		}
	}
	return serverList
}

func TranPlayerNameSting2PlayerIdList(platformId int, playerName string) ([] int, error) {
	playerIdList := make([] int, 0)
	nameList := strings.Split(playerName, ",")
	for _, name := range nameList {
		player, err :=  GetPlayerByPlatformIdAndNickname(platformId, name)
		if err != nil {
			return nil, err
		}
		playerIdList = append(playerIdList, player.Id)
	}
	return playerIdList, nil
}
