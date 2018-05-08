package models

import (
	"github.com/chnzrb/myadmin/utils"
	"sort"
	"errors"
)

func (a *User) TableName() string {
	return UserTBName()
}

func UserTBName() string {
	return TableName("user")
}

type UserQueryParam struct {
	BaseQueryParam
	Account string
}
type User struct {
	Id                      int            `json:"id"`
	Name                    string         `json:"name"`
	Account                 string         `json:"account"`
	Password                string         `json:"-"`
	IsSuper                 int            `json:"isSuper"`
	ModifyPassword          string         `json:"Password" gorm:"-"`
	Status                  int            `json:"status"`
	LoginTimes              int            `json:"loginTimes"`
	LastLoginTime           int            `json:"lastLoginTime"`
	LastLoginIp             string         `json:"lastLoginIp"`
	CanLoginTime            int            `json:"canLoginTime"`
	ContinueLoginErrorTimes int            `json:"-"`
	Mobile                  string         `json:"mobile"`
	RoleIds                 []int          `json:"roleIds" gorm:"-"`
	RoleUserRel             []*RoleUserRel `json:"-"`
	ResourceUrlForList      []string       `gorm:"-"`
}

//获取用户列表
func GetUserList(params *UserQueryParam) ([]*User, int64) {
	data := make([]*User, 0)
	sortOrder := "id"
	switch params.Sort {
	case "id":
		sortOrder = "id"
	}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := Db.Model(&User{}).Where(&User{
		Account: params.Account,
	}).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Error
	utils.CheckError(err)
	for _, v := range data {
		err = Db.Model(&v).Related(&v.RoleUserRel).Error
		utils.CheckError(err)
		roleIds := make([] int, 0)
		for _, e := range v.RoleUserRel {
			roleIds = append(roleIds, e.RoleId)
		}
		sort.Ints(roleIds)
		v.RoleIds = roleIds
	}
	return data, count
}

// 获取单个用户
func GetUserOne(id int) (*User, error) {
	user := &User{
		Id: id,
	}
	err := Db.First(&user).Error
	return user, err
}

// 根据用户名单条
func GetUserOneByAccount(account string) (*User, error) {
	user := &User{}
	isNotFound := Db.Where(&User{Account: account}).First(&user).RecordNotFound()
	if isNotFound {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// 根据用户名密码获取单条
func GetUserOneByAccountAndPassword(account, password string) (*User, error) {
	user := &User{}
	isNotFound := Db.Where(&User{Account: account, Password: password}).First(&user).RecordNotFound()
	if isNotFound {
		return nil, errors.New("用户名或者密码错误")
	}
	return user, nil
}

// 删除用户列表
func DeleteUsers(ids [] int) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&User{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if _, err := DeleteRoleUserRelByUserIdList(ids); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
