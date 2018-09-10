package models

import ()

type LoginNotice struct {
	PlatformId   string    `gorm:"primary_key" json:"platformId"`
	Notice      string `json:"notice"`
	Time         int64  `json:"time"`
	UserId       int    `json:"userId"`
	UserName     string `json:"userName" gorm:"-"`
}

type LoginNoticeQueryParam struct {
	BaseQueryParam
	PlatformId string
}

func GetAllLoginNotice() []*LoginNotice {
	data := make([]*LoginNotice, 0)
	Db.Model(&LoginNotice{}).Find(&data)
	for _, e := range data {
		u, err := GetUserOne(e.UserId)
		if err == nil {
			e.UserName = u.Name
		}
	}
	return data
}

func GetLoginNoticeListByPlatformIdList( platformIdList [] string) []*LoginNotice {
	data := make([]*LoginNotice, 0)
	if len(platformIdList) == 0 {
		return data
	}
	Db.Model(&LoginNotice{}).Where("platform_id in (?)", platformIdList).Find(&data)
	for _, e := range data {
		u, err := GetUserOne(e.UserId)
		if err == nil {
			e.UserName = u.Name
		}
	}
	return data
}


// 删除登录公告日志
func DeleteLoginNotice(ids [] int) error {
	err := Db.Where(ids).Delete(&LoginNotice{}).Error
	return err
}
