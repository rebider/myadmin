package models

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/utils"
	"time"
)

type LoginNotice struct {
	PlatformId string `gorm:"primary_key" json:"platformId"`
	Notice     string `json:"notice"`
	Time       int64  `json:"time"`
	UserId     int    `json:"userId"`
	UserName   string `json:"userName" gorm:"-"`
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

func GetLoginNoticeListByPlatformIdList(platformIdList [] string) []*LoginNotice {
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
func DeleteLoginNotice(ids [] string) error {
	err := Db.Where("platform_id in (?)", ids).Delete(&LoginNotice{}).Error
	return err
}

func UpdateAndPushLoginNotice(userId int, platformId string, notice string) error {
	var request struct {
		PlatformId string `json:"platformId"`
		Notice     string `json:"notice"`
	}
	request.PlatformId = platformId
	request.Notice = notice
	data, err := json.Marshal(request)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	url := utils.GetCenterURL() + "/set_login_notice"
	err = utils.HttpRequest(url, string(data))
	if err != nil {
		return err
	}
	noticeLog := &LoginNotice{
		PlatformId: platformId,
		Notice:     notice,
		Time:       time.Now().Unix(),
		UserId:     userId,
	}
	err = Db.Save(&noticeLog).Error
	return err
}
