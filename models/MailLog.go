package models

import (
	"github.com/jinzhu/gorm"
)


type MailLog struct {
	Id             int    `json:"id"`
	PlatformId     int    `json:"platformId"`
	ServerIdList   string `json:"serverIdList"`
	PlayerNameList string `json:"playerNameList"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Time           int64  `json:"time"`
	UserId         int    `json:"userId"`
	ItemList       string `json:"itemList"`
	Status         int    `json:"status"`
	UserName       string `json:"userName" gorm:"-"`
}

type MailLogQueryParam struct {
	BaseQueryParam
	PlatformId int
	ServerId   string
	StartTime  int
	EndTime    int
	UserId     int
}

func GetMailLogList(params *MailLogQueryParam) ([]*MailLog, int64) {
	data := make([]*MailLog, 0)
	var count int64
	sortOrder := "id"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}
	f(Db.Model(&MailLog{}).Where(&MailLog{
		PlatformId: params.PlatformId,
		UserId:     params.UserId,
	}).Where("server_id_list LIKE ?", "%"+params.ServerId+"%")).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	for _, e := range data {
		u, err := GetUserOne(e.UserId)
		if err == nil {
			e.UserName = u.Name
		}
	}
	return data, count
}

// 删除邮件日志
func DeleteMailLog(ids [] int) error {
	err := Db.Where(ids).Delete(&MailLog{}).Error
	return err
}
