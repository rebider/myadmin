package models

import (
	"github.com/chnzrb/myadmin/utils"
)

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

func (a *Player) TableName() string {
	return "player"
}


//获取分页数据
func GetPlayerList(params *PlayerQueryParam) ([]*Player, int64) {
	//initGame()
	db, err := GetDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	defer db.Close()
	data := make([]*Player, 0)
	var count int64
	sortOrder := "id"
	switch params.Sort {
	case "id":
		sortOrder = "id"
	case "lastLoginTime":
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
	db.Scopes()
	//logs.Debug("Player:%+v", params)
	db.Model(&Player{}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	return data, count
}

//// 根据id获取单条
//func PlayerOne(id int) (*Player, error) {
//	o := orm.NewOrm()
//	err := o.Using("center")
//	utils.CheckError(err)
//	m := Player{Id: id}
//	err = o.Read(&m)
//	if err != nil {
//		return nil, err
//	}
//	return &m, nil
//}
