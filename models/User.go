package models

import (
	"github.com/astaxie/beego/orm"
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
	ModifyPassword     string         `json:"Password" orm:"-"`
	Status             int            `json:"status"`
	LoginTimes         int            `json:"loginTimes"`
	LastLoginTime      int            `json:"lastLoginTime"`
	LastLoginIp        string         `json:"lastLoginIp"`
	Mobile             string         `orm:"size(16)" json:"mobile"`
	RoleIds            []int          `orm:"-" json:"roleIds"`
	RoleUserRel        []*RoleUserRel `json:"-" orm:"reverse(many)"` // 设置一对多的反向关系
	ResourceUrlForList []string       `orm:"-"`
}

//获取分页数据
func UserPageList(params *UserQueryParam) ([]*User, int64) {
	query := orm.NewOrm().QueryTable(UserTBName())
	data := make([]*User, 0)
	//默认排序
	sortOrder := "Id"
	switch params.Sort {
	case "Id":
		sortOrder = "Id"
	}
	if params.Order == "descending" {
		sortOrder = "-" + sortOrder
	}
	query = query.Filter("account__istartswith", params.Account)

	query.RelatedSel().OrderBy(sortOrder).Limit(params.Limit, params.Offset).All(&data)
	for _, v := range data {
		_, error := orm.NewOrm().LoadRelated(v, "RoleUserRel")
		utils.CheckError(error)
		roleIds := make([] int, 0)
		for _, e := range v.RoleUserRel {
			roleIds = append(roleIds, e.Role.Id)
		}
		sort.Ints(roleIds)
		v.RoleIds = roleIds
	}
	total, _ := query.Count()
	return data, total
}

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
func UserOneByAccount(account, password string) (user *User, err error) {
	//m := User{}
	//err := orm.NewOrm().QueryTable(UserTBName()).Filter("account", account).Filter("password", password).One(&m)
	//if err != nil {
	//	return nil, err
	//}
	//return &m, nil
	err = Db.Where(User{Account:account, Password:password}).First(&user).Error
	utils.CheckError(err)
	return user, err
}
