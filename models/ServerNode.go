package models

import (
	"github.com/astaxie/beego/orm"
	//"fmt"
	"encoding/json"
	"fmt"
	//"github.com/chnzrb/admin/models"
	"github.com/astaxie/beego/logs"
)

type ServerNodeQueryParam struct {
	BaseQueryParam
	Type       int
	Node       string
	PlatformId int `json:"platform_id"`
}
type ServerNode struct {
	Node          string `orm:"pk" json:"node"`
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	Type          int    `json:"type"`
	TypeName      string `orm:"-"`
	ZoneNode      string `json:"zone_node"`
	ServerVersion string `json:"server_version"`
	ClientVersion string `json:"client_version"`
	OpenTime      int    `json:"open_time"`
	IsTest        int
	PlatformId    int    `json:"platform_id"`
	PlatformName  string `orm:"-"`
	State         int    `json:"state"`
}

func (t *ServerNode) TableName() string {
	return "c_server_node"
}
func GetServerNodeById(node string) (*ServerNode, error) {
	serverNode := &ServerNode{
		Node: node,
	}

	o := orm.NewOrm()
	o.Using("center")
	err := o.Read(serverNode)
	if err != nil {
		return nil, err
	}
	return serverNode.fill(), nil
}
func GetServerNodeList(filters ...interface{}) ([]*ServerNode, int64) {
	//offset := (page - 1) * pageSize

	serverNodeList := make([]*ServerNode, 0)
	o := orm.NewOrm()
	o.Using("center")
	//fmt.Println("111")
	query := o.QueryTable("c_server_node")
	//fmt.Println("111")
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("-Node").All(&serverNodeList)

	//fmt.Println(serverNodeList)
	return serverNodeList, total
}

//获取分页数据
func ServerNodePageList(params *ServerNodeQueryParam) ([]*ServerNode, int64) {
	logs.Debug("params:%+v", params)
	o := orm.NewOrm()
	err := o.Using("center")
	if err != nil {
		fmt.Println(err)
	}
	query := o.QueryTable("c_server_node")
	data := make([]*ServerNode, 0)
	//默认排序
	sortorder := "Node"
	switch params.Sort {
	case "Node":
		sortorder = "Node"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.Type > 0 {
		query = query.Filter("type", params.Type)
	}
	if params.Node != "" {
		query = query.Filter("node__contains", params.Node)
	}
	if params.PlatformId > 0 {
		query = query.Filter("platform_id", params.PlatformId)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return fillServerNodeList(data), total
}
func ShowGameNodeJson(Data map[interface{}]interface{}) map[interface{}]interface{} {
	gameServerNodeList, _ := GetServerNodeList("type", 1)
	out, err := json.Marshal(gameServerNodeList)
	fmt.Println("game_node_list error:", err, gameServerNodeList)
	Data["game_node_list"] = string(out)
	return Data
}
func ShowGameNodeList(Data map[interface{}]interface{}) map[interface{}]interface{} {
	gameServerNodeList, _ := GetServerNodeList("type", 1)
	Data["game_node_list"] = gameServerNodeList
	return Data
}
func GetTypeName(t int) string {
	platformMap := GetNodeTypeMap()
	platformName, ok := platformMap[t]
	if ok == true {
		return platformName
	}
	return "未定义"
}

func (t *ServerNode) fill() *ServerNode {
	t.TypeName = GetTypeName(t.Type)
	t.PlatformName = GetPlatformName(t.PlatformId)
	return t
}
func fillServerNodeList(ServerNodeList []*ServerNode) []*ServerNode {
	for _, g := range ServerNodeList {
		g.fill()
	}
	return ServerNodeList
}
func GetNodeTypeMap() map[int]string {
	nodeTypeMap := map[int]string{
		0: "中心节点",
		1: "游戏节点",
		2: "跨服节点",
		5: "唯一id节点",
	}
	return nodeTypeMap
}
func ShowNodeTypeList(Data map[interface{}]interface{}) map[interface{}]interface{} {
	nodeTypeMap := GetNodeTypeMap()
	nodeTypeList := make([]map[string]interface{}, 0, len(nodeTypeMap))
	for k, v := range nodeTypeMap {
		row := make(map[string]interface{})
		row["type_id"] = k
		row["type_name"] = v
		nodeTypeList = append(nodeTypeList, row)
	}

	//out,_ := json.Marshal(nodeTypeList)
	Data["node_type_list"] = nodeTypeList
	return Data
}

func ShowNodeTypeJsone(Data map[interface{}]interface{}) map[interface{}]interface{} {
	nodeTypeMap := GetNodeTypeMap()
	nodeTypeList := make([]map[string]interface{}, 0, len(nodeTypeMap))
	for k, v := range nodeTypeMap {
		row := make(map[string]interface{})
		row["type_id"] = k
		row["type_name"] = v
		nodeTypeList = append(nodeTypeList, row)
	}

	out, _ := json.Marshal(nodeTypeList)
	Data["node_type_list"] = string(out)
	return Data
}
