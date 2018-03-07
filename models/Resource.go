package models

import (
	"fmt"

	"github.com/chnzrb/myadmin/utils"

	"github.com/astaxie/beego/orm"
	//"github.com/astaxie/beego/logs"
	//"sort"
	//"sort"
	//"github.com/astaxie/beego/logs"
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
	Title string `orm:"size(64)" json:"title"` //标题
	Name  string `orm:"size(64)" json:"name"`
	//Component            string    `orm:"size(64)" json:"component"`
	Parent          *Resource          `orm:"null;rel(fk) " json:"-"` // RelForeignKey relation
	ParentId        int                `orm:"-" json:"parentId"`             // RelForeignKey relation
	Type            int                `json:"type"`
	Seq             int                `json:"seq"`
	Children        []*Resource        `orm:"reverse(many)" json:"children"` // fk 的反向关系
	SonNum          int                `orm:"-" json:"-"`
	Icon            string             `orm:"size(32)" json:"icon"`
	LinkUrl         string             `orm:"-" json:"-"`
	UrlFor          string             `orm:"size(256)" json:"path"`
	HtmlDisabled    int                `orm:"-" json:"-"`             //在html里应用时是否可用
	Level           int                `orm:"-" json:"-"`             //第几级，从0开始
	RoleResourceRel []*RoleResourceRel `orm:"reverse(many)" json:"-"` // 设置一对多的反向关系
}

func ResourceOne(id int) (*Resource, error) {
	o := orm.NewOrm()
	m := Resource{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	m.ParentId = m.Parent.Id
	return &m, nil
}

//获取分页数据
func ResourceList() []*Resource {
	query := orm.NewOrm().QueryTable(ResourceTBName())
	data := make([]*Resource, 0)
	query.All(&data)
	for _,e:= range data {
		e.ParentId = e.Parent.Id
	}
	//total, _ := query.Count()
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

	//return tmpList
	//logs.Debug("before:%+v ", tmpList)
	//logs.Debug("after:%+v ", TranResourceList2ResourceTree(tmpList))
	return tmpList
	//return TranResourceList2ResourceTree(tmpList)
}

func CanParent(resourceId int, thisResourceId int) bool {
	thisResource, _ := ResourceOne(thisResourceId)
	if thisResource.Parent.Id == resourceId || resourceId == thisResourceId {
		//logs.Debug("CanParent:%+v %+v %+v", resourceId, thisResourceId, thisResource.Parent.Id)
		return false
	}
	if thisResource.Parent.Id == 0 {
		return true
	}
	return CanParent(resourceId, thisResource.Parent.Id)
}

//根据用户获取有权管理的资源列表
func GetResourceListByUserId(userId, maxrtype int) []*Resource {
	var list []*Resource
	o := orm.NewOrm()
	user, err := UserOne(userId)
	logs.Info("user:%+v", user)
	utils.CheckError(err)
	if err != nil || user == nil {
		return list
	}

	var sql string
	//if user.IsSuper == true {
	//	//如果是管理员，则查出所有的
	//	sql = fmt.Sprintf(`SELECT id,name,parent_id,rtype,icon,seq,url_for FROM %s Where rtype <= ? Order By seq asc,Id asc`, ResourceTBName())
	//	o.Raw(sql, maxrtype).QueryRows(&list)
	//} else {
	//	//联查多张表，找出某用户有权管理的
	sql = fmt.Sprintf(`SELECT DISTINCT T2.*
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
		INNER JOIN %s AS T2 ON T2.id = T0.resource_id
		WHERE T1.user_id = ? and T2.type <= ?  Order By T2.seq asc,T2.id asc`, RoleResourceRelTBName(), RoleUserRelTBName(), ResourceTBName())
	o.Raw(sql, userId, maxrtype).QueryRows(&list)
	result := list
	for _,e:= range list {
		e.ParentId = e.Parent.Id
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
