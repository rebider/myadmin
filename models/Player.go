package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/chnzrb/myadmin/utils"
	//"github.com/astaxie/beego"
	//"fmt"
	//"database/sql"
	"github.com/jinzhu/gorm"
)

func (a *Player) TableName() string {
	return PlayerTBName()
}

type PlayerQueryParam struct {
	BaseQueryParam
	Account    string
	Ip         string
	PlayerId   string
	Nickname   string
	IsOnline   string
	PlatformId int
	ServerId   string
}
type Player struct {
	Id              int    `json:"id"`
	AccId           string `orm:"size(32)" json:"accId"`
	Nickname        string `orm:"size(24)" json:"nickname"`
	Sex             int    `json:"sex"`
	ServerId        string `json:"serverId"`
	DisableLogin    int    `json:"disableLogin"`
	RegTime         int    `json:"regTime"`
	LastLoginTime   int    `json:"lastLoginTime"`
	LastOfflineTime int    `json:"lastOfflineTime"`
	LastLoginIp     string `json:"lastLoginIp"`
	IsOnline        int    `json:"isOnline"`
	DisableChatTime int    `json:"disableChatTime"`
}

//获取分页数据
//func PlayerPageList(params *PlayerQueryParam) ([]*Player, int64) {
//	//initGame()
//	orm := orm.NewOrm()
//	err := orm.Using("game")
//	utils.CheckError(err)
//
//	query := orm.QueryTable(PlayerTBName())
//
//	data := make([]*Player, 0)
//	//默认排序
//	sortOrder := "Id"
//	switch params.Sort {
//	case "Id":
//		sortOrder = "Id"
//	case "LastLoginTime":
//		sortOrder = "last_login_time"
//	}
//	if params.Order == "descending" {
//		sortOrder = "-" + sortOrder
//	}
//	if params.Account != "" {
//		query = query.Filter("acc_id__contains", params.Account)
//	}
//	if params.Ip != "" {
//		query = query.Filter("last_login_ip", params.Ip)
//	}
//	if params.Nickname != "" {
//		query = query.Filter("nickname__contains", params.Nickname)
//	}
//	if params.IsOnline != "" {
//		query = query.Filter("is_online", params.IsOnline)
//	}
//	if params.PlayerId != "" {
//		query = query.Filter("id", params.PlayerId)
//	}
//	_, err = query.OrderBy(sortOrder).Limit(params.Limit, params.Offset).All(&data)
//	utils.CheckError(err)
//	total, _ := query.Count()
//	return data, total
//}

func GetDbByPlatformIdAndSid(platformId int, Sid string) (db *gorm.DB, err error){
	orm := orm.NewOrm()
	err1 := orm.Using("center")
	utils.CheckError(err1)
	gameServer, err := GetGameServer(platformId, Sid)
	utils.CheckError(err)
	serverNode, err := GetServerNode(gameServer.Node)
	utils.CheckError(err)
	dbArgs := "root:game1234@tcp(" +serverNode.Ip+ ":3306)/game?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", dbArgs)
	return db, err
}
//获取分页数据
func PlayerPageList(params *PlayerQueryParam) ([]*Player, int64) {
	//initGame()
	db, err := GetDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	defer db.Close()
	data := make([]*Player, 0)
	var count int64
	//默认排序
	sortOrder := "Id"
	switch params.Sort {
	case "Id":
		sortOrder = "Id"
	case "LastLoginTime":
		sortOrder = "last_login_time"
	}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}

	if params.Account != "" {
		db = db.Where("acc_id = ?", params.Account)
	}
	if params.Ip != "" {
		db = db.Where("last_login_ip = ?", params.Ip)
	}
	if params.Nickname != "" {
		db = db.Where("nickname LIKE ?", "%" + params.Nickname + "%")
	}
	if params.IsOnline != "" {
		db = db.Where("is_online = ?", params.IsOnline)
	}
	if params.PlayerId != "" {
		db = db.Where("id = ?", params.PlayerId)
	}

	db.Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Count(&count)
	return data, count
}

// 根据id获取单条
func PlayerOne(id int) (*Player, error) {
	o := orm.NewOrm()
	err := o.Using("center")
	utils.CheckError(err)
	m := Player{Id: id}
	err = o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
