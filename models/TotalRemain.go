package models

import (
	"github.com/chnzrb/myadmin/utils"
)

type TotalRemain struct {
	Node         string `json:"node"`
	Date         int    `json:"date"`
	RegisterRole int    `json:"registerRole"`
	CreateRole   int    `json:"createRole"`
	Remain2      int    `json:"remain2"`
	Remain3      int    `json:"remain3"`
	Remain4      int    `json:"remain4"`
	Remain5      int    `json:"remain5"`
	Remain6      int    `json:"remain6"`
	Remain7      int    `json:"remain7"`
}

type TotalRemainQueryParam struct {
	BaseQueryParam
	PlatformId int
	ServerId   string
	StartTime  int
	EndTime    int
}

func GetTotalRemainList(params *TotalRemainQueryParam) ([]*TotalRemain, int64) {
	data := make([]*TotalRemain, 0)
	var count int64
	gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	err = Db.Model(&TotalRemain{}).Where(&TotalRemain{Node: gameServer.Node}).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data).Error
	utils.CheckError(err)
	return data, count
}
