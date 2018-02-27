package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/chnzrb/myadmin/utils"
)

type GameServerQueryParam struct {
	BaseQueryParam
	PlatformId string
	Node       string
}

type GameServer struct {
	PlatformId int
	Sid        string `orm:"pk"`
	Desc       string
	Node       string
	PlatformName string `orm:"-"`
}

func (t *GameServer) TableName() string {
	return "c_game_server"
}
func (t *GameServer) fill() *GameServer {
	t.PlatformName = GetPlatformName(t.PlatformId)
	return t
}
func fillGameServerList(gameServerList []*GameServer) []*GameServer{
	for _,g := range gameServerList{
		g.fill()
	}
	return gameServerList
}
//获取所有数据
func GetAllGameServer() ([]*GameServer, int64){
	var params GameServerQueryParam
	//获取数据列表和总数
	data, total := GetGameServerList(&params)
	return data, total
}
//获取数据列表
func GetGameServerList(params *GameServerQueryParam) ([]*GameServer, int64) {
	o := orm.NewOrm()
	err := o.Using("center")
	utils.CheckError(err)
	query := o.QueryTable("c_game_server")

	//默认排序
	sortorder := "Sid"
	switch params.Sort {
	case "Sid":
		sortorder = "Sid"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.PlatformId != "" && params.PlatformId != "-1"{
		query = query.Filter("platform_id", params.PlatformId)
	}
	if params.Node != ""{
		query = query.Filter("node__contains", params.Node)
	}
	total, _ := query.Count()
	data := make([]*GameServer, total)
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return fillGameServerList(data), total
}
func GetGameServer(platformId int, id string) (*GameServer, error) {
	gameServer := &GameServer{
		Sid:         id,
		PlatformId: platformId,
	}

	o := orm.NewOrm()
	o.Using("center")
	err := o.Read(gameServer)
	if err != nil {
		return nil, err
	}
	return gameServer.fill(), nil
}

