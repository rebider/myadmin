package models

import (
	"github.com/chnzrb/myadmin/utils"
	"sort"
)

func (a *User) TableName() string {
	return UserTBName()
}

type UserQueryParam struct {
	BaseQueryParam
	Account string
}
type User struct {
	Id                 int            `json:"id"`
	Name               string         `orm:"size(32)" json:"name"`
	Account            string         `orm:"size(24);" json:"account"`
	Password           string         `json:"-"`
	IsSuper            int            `json:"isSuper"`
	ModifyPassword     string         `json:"Password" orm:"-" gorm:"-"`
	Status             int            `json:"status"`
	LoginTimes         int            `json:"loginTimes"`
	LastLoginTime      int            `json:"lastLoginTime"`
	LastLoginIp        string         `json:"lastLoginIp"`
	Mobile             string         `orm:"size(16)" json:"mobile"`
	RoleIds            []int          `orm:"-" json:"roleIds" gorm:"-"`
	RoleUserRel        []*RoleUserRel `json:"-" orm:"reverse(many)"` // 设置一对多的反向关系
	ResourceUrlForList []string       `orm:"-" gorm:"-"`
}

//获取分页数据
func GetUserList(params *UserQueryParam) ([]*User, int64) {
	data := make([]*User, 0)
	sortOrder := "Id"
	switch params.Sort {
	case "Id":
		sortOrder = "Id"
	}
	if params.Order == "descending" {
		sortOrder = "-" + sortOrder
	}
	//query = query.Filter("account__istartswith", params.Account)

	var count int64
	err := Db.Model(&User{}).Count(&count).Where(&User{
		Account:params.Account,
	}).Offset(params.Offset).Limit(params.Limit).Find(&data).Error
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

////获取分页数据
//func GetUserList(params *UserQueryParam) ([]*User, int64) {
//	query := orm.NewOrm().QueryTable(UserTBName())
//	data := make([]*User, 0)
//	//默认排序
//	sortOrder := "Id"
//	switch params.Sort {
//	case "Id":
//		sortOrder = "Id"
//	}
//	if params.Order == "descending" {
//		sortOrder = "-" + sortOrder
//	}
//	query = query.Filter("account__istartswith", params.Account)
//
//	query.RelatedSel().OrderBy(sortOrder).Limit(params.Limit, params.Offset).All(&data)
//	for _, v := range data {
//		_, error := orm.NewOrm().LoadRelated(v, "RoleUserRel")
//		utils.CheckError(error)
//		roleIds := make([] int, 0)
//		for _, e := range v.RoleUserRel {
//			roleIds = append(roleIds, e.Role.Id)
//		}
//		sort.Ints(roleIds)
//		v.RoleIds = roleIds
//	}
//	total, _ := query.Count()
//	return data, total
//}

// 根据id获取单条
func UserOne(id int) (user *User, err error) {
	user = &User{
		Id: id,
	}
	err = Db.First(&user).Error
	utils.CheckError(err)
	return user, err
}

// 根据用户名密码获取单条
func UserOneByAccount(account, password string) (*User, error) {
	user := &User{}
	err := Db.Where(&User{Account: account, Password: password}).First(&user).Error
	return user, err
}

// 删除用户列表
func DeleteUsers(userIdList [] int) (int, error) {
	var count int
	err := Db.Where(userIdList).Delete(&User{}).Count(&count).Error
	if err == nil {
		// 删除关联的 角色列表
		_, err = DeleteRoleUserRelByUserIdList(userIdList)
	}
	return  count, err
}
