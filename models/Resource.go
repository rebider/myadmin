package models

import (
	"fmt"

	"github.com/chnzrb/myadmin/utils"

	"github.com/astaxie/beego/orm"
	//"github.com/astaxie/beego/logs"
	//"sort"
)

func (a *Resource) TableName() string {
	return ResourceTBName()
}

//Resource 权限控制资源表
type Resource struct {
	Id              int
	Title 			string		 `orm:"size(64)" json:"title"`//标题
	Name            string    `orm:"size(64)" json:"name"`
	Component            string    `orm:"size(64)" json:"componentName"`
	Parent          *Resource `orm:"null;rel(fk)" json:"-"` // RelForeignKey relation
	Rtype           int
	Seq             int
	Children            []*Resource        `orm:"reverse(many)" json:"children"` // fk 的反向关系
	SonNum          int                `orm:"-"`
	Icon            string             `orm:"size(32)" json:"icon"`
	LinkUrl         string             `orm:"-"`
	UrlFor          string             `orm:"size(256)" json:"path"`
	HtmlDisabled    int                `orm:"-"`             //在html里应用时是否可用
	Level           int                `orm:"-"`             //第几级，从0开始
	RoleResourceRel []*RoleResourceRel `orm:"reverse(many)"` // 设置一对多的反向关系
}

func ResourceOne(id int) (*Resource, error) {
	o := orm.NewOrm()
	m := Resource{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

//Resource 获取treegrid顺序的列表
func ResourceTreeGrid() []*Resource {
	o := orm.NewOrm()
	query := o.QueryTable(ResourceTBName())
	list := make([]*Resource, 0)

	query.All(&list)
	sonList := make([]*Resource, 0)
	realList := make([]*Resource, 0)

	for _, v := range list {
		if v.Parent ==  nil {
			//realList[v.Id] = v
			realList = append(realList, v)
		} else {
			sonList = append(sonList, v)
		}
	}

	for _, v := range sonList{
		for _, p := range realList{
			if p.Id == v.Parent.Id {
				p.Children = append(p.Children, v)
			}
		}
	}
	//logs.Info("%v", list)
	//return resourceList2TreeGrid(list)
	return realList
}

//ResourceTreeGrid4Parent 获取可以成为某个节点父节点的列表
func ResourceTreeGrid4Parent(id int) []*Resource {
	tree := ResourceTreeGrid()
	if id == 0 {
		return tree
	}
	var index = -1
	//找出当前节点所在索引
	for i, _ := range tree {
		if tree[i].Id == id {
			index = i
			break
		}
	}
	if index == -1 {
		return tree
	} else {
		tree[index].HtmlDisabled = 1
		for _, item := range tree[index+1:] {
			if item.Level > tree[index].Level {
				item.HtmlDisabled = 1
			} else {
				break
			}
		}
	}
	return tree
}

//根据用户获取有权管理的资源列表，并整理成teegrid格式
func ResourceTreeGridByUserId(backuserid, maxrtype int) []*Resource {
	cachekey := fmt.Sprintf("rms_ResourceTreeGridByUserId_%v_%v", backuserid, maxrtype)
	var list []*Resource
	if err := utils.GetCache(cachekey, &list); err == nil {
		return list
	}
	//o := orm.NewOrm()
	user, err := UserOne(backuserid)
	if err != nil || user == nil {
		return list
	}
	//var sql string
	//if user.IsSuper == true {
	//	//如果是管理员，则查出所有的
	//	sql = fmt.Sprintf(`SELECT id,name,parent_id,rtype,icon,seq,url_for FROM %s Where rtype <= ? Order By seq asc,Id asc`, ResourceTBName())
	//	o.Raw(sql, maxrtype).QueryRows(&list)
	//} else {
	//	//联查多张表，找出某用户有权管理的
	//	sql = fmt.Sprintf(`SELECT DISTINCT T0.resource_id,T2.id,T2.name,T2.parent_id,T2.rtype,T2.icon,T2.seq,T2.url_for
	//	FROM %s AS T0
	//	INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
	//	INNER JOIN %s AS T2 ON T2.id = T0.resource_id
	//	WHERE T1.backend_user_id = ? and T2.rtype <= ?  Order By T2.seq asc,T2.id asc`, RoleResourceRelTBName(), RoleUserRelTBName(), ResourceTBName())
	//	o.Raw(sql, backuserid, maxrtype).QueryRows(&list)
	//}
	result := resourceList2TreeGrid(list)
	utils.SetCache(cachekey, result, 30)
	return result
}

//将资源列表转成treegrid格式
func resourceList2TreeGrid(list []*Resource) []*Resource {
	result := make([]*Resource, 0)
	for _, item := range list {
		if item.Parent == nil || item.Parent.Id == 0 {
			item.Level = 0
			result = append(result, item)
			result = resourceAddSons(item, list, result)
		}
	}
	return result
}

//resourceAddSons 添加子菜单
func resourceAddSons(cur *Resource, list, result []*Resource) []*Resource {
	for _, item := range list {
		if item.Parent != nil && item.Parent.Id == cur.Id {
			cur.SonNum++
			item.Level = cur.Level + 1
			result = append(result, item)
			result = resourceAddSons(item, list, result)
		}
	}
	return result
}
