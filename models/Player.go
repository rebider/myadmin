package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/chnzrb/myadmin/utils"
)

func (a *Player) TableName() string {
	return PlayerTBName()
}

type PlayerQueryParam struct {
	BaseQueryParam
}
type Player struct {
	Id              int    `json:"id"`
	AccId           string `orm:"size(32)" json:"acc_id"`
	Nickname        string `orm:"size(24)" json:"nick_name"`
	Sex             int    `json:"sex"`
	ServerId        string `json:"server_id"`
	DisableLogin    int    `json:"disable_login"`
	RegTime         int    `json:"reg_time"`
	LastLoginTime   int    `json:"last_login_time"`
	LastOfflineTime int    `json:"last_offline_time"`
	LastLoginIp     string    `json:"last_login_ip"`
	IsOnline        int    `json:"is_online"`
	DisableChatTime int    `json:"disable_chat_time"`
}

//获取分页数据
func PlayerPageList(params *PlayerQueryParam) ([]*Player, int64) {
	orm := orm.NewOrm()
	err := orm.Using("game")
	utils.CheckError(err)

	query := orm.QueryTable(PlayerTBName())
	data := make([]*Player, 0)
	//默认排序
	sortOrder := "Id"
	switch params.Sort {
	case "Id":
		sortOrder = "Id"
	}
	if params.Order == "descending" {
		sortOrder = "-" + sortOrder
	}

	_, err = query.OrderBy(sortOrder).Limit(params.Limit, params.Offset).All(&data)
	utils.CheckError(err)
	total, _ := query.Count()
	return data, total
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
