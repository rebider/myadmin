package models

import (
	"github.com/chnzrb/myadmin/utils"
)

type GlobalAccount struct {
	PlatformId string `gorm:"primary_key" json:"platformId"`
	Account    string `gorm:"primary_key" json:"account"`
	Type       int `json:"type"`
	ForbidType int    `json:"forbidType"`
	ForbidTime int    `json:"forbidTime"`
}

func (a *GlobalAccount) TableName() string {
	return "global_account"
}

func GetGlobalAccount(platformId string, accId string) (*GlobalAccount, error) {
	globalAccount := &GlobalAccount{}
	err := DbCenter.Where(&GlobalAccount{PlatformId: platformId,
		Account: accId,
	}).First(&globalAccount).Error
	utils.CheckError(err)
	return globalAccount, err
}
