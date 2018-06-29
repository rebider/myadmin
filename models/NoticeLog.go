package models

import ()

type NoticeLog struct {
	Id           int    `json:"id"`
	PlatformId   string    `json:"platformId"`
	NodeList     string `json:"serverIdList"`
	Content      string `json:"content"`
	NoticeType   int    `json:"noticeType"`
	NoticeTime   int    `json:"noticeTime"`
	Status       int    `json:"status"`
	IsAllServer  int    `json:"isAllServer"`
	Time         int64  `json:"time"`
	UserId       int    `json:"userId"`
	UserName     string `json:"userName" gorm:"-"`
	LastSendTime int    `json:"-"`
}

type NoticeLogQueryParam struct {
	BaseQueryParam
	PlatformId string
	StartTime  int
	EndTime    int
	UserId     int
	NoticeType int
}

func GetNoticeLogOne(id int) (*NoticeLog, error) {
	noticeLog := &NoticeLog{
		Id: id,
	}
	err := Db.First(&noticeLog).Error
	return noticeLog, err
}

func GetAllNoticeLog() []*NoticeLog {
	data := make([]*NoticeLog, 0)
	Db.Model(&NoticeLog{}).Find(&data)
	return data
}

func GetNoticeLogList(params *NoticeLogQueryParam) ([]*NoticeLog, int64) {
	data := make([]*NoticeLog, 0)
	var count int64
	sortOrder := "time"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	Db.Model(&NoticeLog{}).Where(&NoticeLog{
		PlatformId: params.PlatformId,
		NoticeType:params.NoticeType,
	}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	for _, e := range data {
		u, err := GetUserOne(e.UserId)
		if err == nil {
			e.UserName = u.Name
		}
	}
	return data, count
}

// 删除公告日志
func DeleteNoticeLog(ids [] int) error {
	err := Db.Where(ids).Delete(&NoticeLog{}).Error
	return err
}
