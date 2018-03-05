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
	Id              int				`json:"id"`
	Title 			string		 `orm:"size(64)" json:"title"`//标题
	Name            string    `orm:"size(64)" json:"name"`
	//Component            string    `orm:"size(64)" json:"component"`
	Parent          *Resource `orm:"null;rel(fk) " json:"parent"` // RelForeignKey relation
	Type           int			`json:"type"`
	Seq             int			`json:"seq"`
	Children            []*Resource        `orm:"reverse(many)" json:"children"` // fk 的反向关系
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
	return &m, nil
}

//Resource 获取treegrid顺序的列表
func MenuTree() []*Resource {
	o := orm.NewOrm()
	query := o.QueryTable(ResourceTBName()).Filter("type", 0)
	list := make([]*Resource, 0)

	query.All(&list)
	sonList := make([]*Resource, 0)
	realList := make([]*Resource, 0)

	for _, v := range list {
		//logs.Info("ResourceTreeGrid:%+v", v.Parent)
		if v.Parent ==  nil || v.Parent.Id == 0{
			realList = append(realList, v)
		} else {
			sonList = append(sonList, v)
		}
	}

	//logs.Info("realList:%+v", realList)
	//logs.Info("sonList:%+v", sonList)
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

//Resource 获取treegrid顺序的列表
func ResourceTreeGrid() []*Resource {
	o := orm.NewOrm()
	query := o.QueryTable(ResourceTBName())
	list := make([]*Resource, 0)

	query.All(&list)
	sonList := make([]*Resource, 0)
	realList := make([]*Resource, 0)

	for _, v := range list {
		//logs.Info("ResourceTreeGrid:%+v", v.Parent)
		if v.Parent ==  nil || v.Parent.Id == 0{
			realList = append(realList, v)
		} else {
			sonList = append(sonList, v)
		}
	}

	//logs.Info("realList:%+v", realList)
	//logs.Info("sonList:%+v", sonList)
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

//获取分页数据
func ResourcePageList(params *ResourceQueryParam) ([]*Resource, int64) {
	query := orm.NewOrm().QueryTable(ResourceTBName())
	data := make([]*Resource, 0)
	query.RelatedSel().Limit(params.Limit, params.Offset).All(&data)
	total, _ := query.Count()
	return data, total
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
func ResourceTreeGridByUserId(userId, maxrtype int) []*Resource {
	cachekey := fmt.Sprintf("rms_ResourceTreeGridByUserId_%v_%v", userId, maxrtype)
	var list []*Resource
	if err := utils.GetCache(cachekey, &list); err == nil {
		return list
	}
	o := orm.NewOrm()
	user, err := UserOne(userId)
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
		sql = fmt.Sprintf(`SELECT DISTINCT T0.resource_id,T2.id,T2.name,T2.parent_id,T2.type,T2.icon,T2.seq,T2.url_for
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
		INNER JOIN %s AS T2 ON T2.id = T0.resource_id
		WHERE T1.user_id = ?   Order By T2.seq asc,T2.id asc`, RoleResourceRelTBName(), RoleUserRelTBName(), ResourceTBName())
		o.Raw(sql, userId).QueryRows(&list)
	//}
	for _, e := range list{
		logs.Debug("ResourceTreeGridByUserId:%+v", e)
	}

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
