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
	Account      string //模糊查询
	//NameLike     string //模糊查询
	//Mobile       string //精确查询
	//SearchStatus string //为空不查询，有值精确查询
}
type User struct {
	Id             int
	Name           string `orm:"size(32)"`
	Account        string `orm:"size(24);unique"`
	Password       string `json:"-"`
	ModifyPassword string `json:"Password" orm:"-"`
	Status             int
	LoginTimes         int
	LastLoginTime      int
	LastLoginIp        string
	Mobile             string         `orm:"size(16)"`
	RoleIds            []int          `orm:"-"`
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
	//query = query.Filter("name__istartswith", params.NameLike)
	//if len(params.Mobile) > 0 {
	//	query = query.Filter("mobile", params.Mobile)
	//}
	//if len(params.SearchStatus) > 0 {
	//	query = query.Filter("status", params.SearchStatus)
	//}

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
func UserOne(id int) (*User, error) {
	o := orm.NewOrm()
	m := User{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// 根据用户名密码获取单条
func UserOneByAccount(account, password string) (*User, error) {
	m := User{}
	err := orm.NewOrm().QueryTable(UserTBName()).Filter("account", account).Filter("password", password).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
