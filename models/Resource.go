package models

import (
	"fmt"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

func (a *Resource) TableName() string {
	return ResourceTBName()
}

type ResourceQueryParam struct {
	BaseQueryParam
}

//Resource 权限控制资源表
type Resource struct {
	Id    int    `json:"id"`
	Name  string `orm:"size(64)" json:"name"`
	Parent          *Resource          `orm:"null;rel(fk) " json:"-"` // RelForeignKey relation
	ParentId        int                `orm:"-" json:"parentId"`             // RelForeignKey relation
	Children        []*Resource        `orm:"reverse(many)" json:"children"` // fk 的反向关系
	UrlFor          string             `orm:"size(256)" json:"urlFor"`
	Url          string             	`orm:"-" json:"url"`
	RoleResourceRel []*RoleResourceRel `orm:"reverse(many)" json:"-"` // 设置一对多的反向关系
}

func ResourceOne(id int) (*Resource, error) {
	o := orm.NewOrm()
	m := Resource{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	if m.Parent != nil {
		m.ParentId = m.Parent.Id
	}

	return &m, nil
}

//获取分页数据
func ResourceList() []*Resource {
	query := orm.NewOrm().QueryTable(ResourceTBName())
	data := make([]*Resource, 0)
	query.All(&data)
	for _,e:= range data {
		if e.Parent != nil {
			e.ParentId = e.Parent.Id
		}

	}
	logs.Debug("ResourceList:%+v", data)
	return data
}

//ResourceTreeGrid4Parent 获取可以成为某个节点父节点的列表
func ResourceTreeGrid4Parent(id int) []*Resource {
	list := ResourceList()
	tmpList := make([] *Resource, 0)
	if id > 0 {
		for _, e := range list {
			if CanParent(id, e.Id) {
				tmpList = append(tmpList, e)
			}
		}
	} else {
		tmpList = list
	}
	return tmpList
}

func CanParent(resourceId int, parentResourceId int) bool {
	parentResource, _ := ResourceOne(parentResourceId)
	if (parentResource.Parent != nil && parentResource.Parent.Id == resourceId) || resourceId == parentResourceId {
		//logs.Debug("CanParent:%+v %+v %+v", resourceId, thisResourceId, thisResource.Parent.Id)
		return false
	}
	if parentResource.Parent == nil || parentResource.Parent.Id == 0 {
		return true
	}
	return CanParent(resourceId, parentResource.Parent.Id)
}

//根据用户获取有权管理的资源列表
func GetResourceListByUserId(userId int) []*Resource {
	var list []*Resource
	o := orm.NewOrm()
	user, err := UserOne(userId)
	logs.Info("user:%+v", user)
	utils.CheckError(err)
	if err != nil || user == nil {
		return list
	}

	var sql string
	if user.IsSuper == 1 {
		//如果是管理员，则查出所有的
		sql = fmt.Sprintf(`SELECT * FROM %s  Order By seq asc,Id asc`, ResourceTBName())
		o.Raw(sql).QueryRows(&list)
	} else {
		//	//联查多张表，找出某用户有权管理的
		sql = fmt.Sprintf(`SELECT DISTINCT T2.*
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
		INNER JOIN %s AS T2 ON T2.id = T0.resource_id
		WHERE T1.user_id = ? `, RoleResourceRelTBName(), RoleUserRelTBName(), ResourceTBName())
		o.Raw(sql, userId).QueryRows(&list)
	}
	result := list
	for _,e:= range list {
		if e.Parent != nil {
			e.ParentId = e.Parent.Id
		}
	}
	return result
}

func TranResourceList2ResourceTree(resourceList []*Resource) []*Resource {
	resourceTree := make([]*Resource, 0)
	for _, item := range resourceList {
		if item.Parent == nil || item.Parent.Id == 0 {
			item = TranResourceList2ResourceTree_(item, resourceList)
			resourceTree = append(resourceTree, item)
		}
	}
	logs.Debug("TranResourceList2ResourceTree:%+v", resourceTree)
	return resourceTree
}
func TranResourceList2ResourceTree_(cur *Resource, list []*Resource) *Resource {
	for _, item := range list {
		if item.Parent != nil && item.Parent.Id == cur.Id {
			item = TranResourceList2ResourceTree_(item, list)
			cur.Children = append(cur.Children, item)
		}
	}
	return cur
}
