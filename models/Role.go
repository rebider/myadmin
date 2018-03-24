package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/chnzrb/myadmin/utils"
	"sort"
)

func (a *Role) TableName() string {
	return RoleTBName()
}

type RoleQueryParam struct {
	BaseQueryParam
	NameLike string
}

//用户角色
type Role struct {
	Id              int                `form:"id" json:"id"`
	Name            string             `form:"name" json:"name"`
	ResourceIds     [] int             `orm:"-" json:"resourceIds"`
	MenuIds         [] int             `orm:"-" json:"menuIds"`
	RoleResourceRel []*RoleResourceRel `orm:"reverse(many)" json:"-"` // 设置一对多的反向关系
	RoleMenuRel     []*RoleMenuRel     `orm:"reverse(many)" json:"-"` // 设置一对多的反向关系
	//RoleUserRel     []*RoleUserRel     `orm:"reverse(many)" json:"-"` // 设置一对多的反向关系
}

//获取分页数据
func RolePageList(params *RoleQueryParam) ([]*Role, int64) {
	//query := orm.NewOrm().QueryTable(RoleTBName())
	data := make([]*Role, 0)
	//默认排序
	//sortorder := "Id"
	//switch params.Sort {
	//case "Id":
	//	sortorder = "Id"
	//case "Seq":
	//	sortorder = "Seq"
	//}
	//if params.Order == "desc" {
	//	sortorder = "-" + sortorder
	//}
	//query = query.Filter("name__istartswith", params.NameLike)
	//total, _ := query.Count()
	//query.RelatedSel().OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	var count int64
	err := Db.Model(&Role{}).Count(&count).Where(&Role{Name:params.NameLike}).Offset(params.Offset).Limit(params.Limit).Find(&data).Error
	utils.CheckError(err)
	for _, v := range data {
		err = Db.Model(&v).Related(&v.RoleResourceRel).Error
		//_, error := orm.NewOrm().LoadRelated(v, "RoleResourceRel")
		utils.CheckError(err)
		err = Db.Model(&v).Related(&v.RoleMenuRel).Error
		//_, error = orm.NewOrm().LoadRelated(v, "RoleMenuRel")
		utils.CheckError(err)

		resourceIds := make([] int, 0)
		for _, e := range v.RoleResourceRel {
			resourceIds = append(resourceIds, e.ResourceId)
		}
		sort.Ints(resourceIds)
		v.ResourceIds = resourceIds

		menuIds := make([] int, 0)
		for _, e := range v.RoleMenuRel {
			menuIds = append(menuIds, e.MenuId)
		}
		sort.Ints(menuIds)
		v.MenuIds = menuIds
	}
	return data, count
}

//获取角色列表
func RoleDataList(params *RoleQueryParam) []*Role {
	params.Limit = -1
	params.Sort = "Seq"
	params.Order = "asc"
	data, _ := RolePageList(params)
	return data
}

//批量删除
func RoleBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(RoleTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}
func RoleOne(id int) (*Role, error) {
	o := orm.NewOrm()
	m := Role{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
