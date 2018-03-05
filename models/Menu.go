package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

func (a *Menu) TableName() string {
	return MenuTBName()
}

type MenuQueryParam struct {
	BaseQueryParam
}

//Resource 权限控制资源表
type Menu struct {
	Id       int     `json:"id"`
	Title    string  `orm:"size(64)" json:"title"` //标题
	Name     string  `orm:"size(64)" json:"name"`
	Parent   *Menu   `orm:"null;rel(fk) " json:"parent"`
	Seq      int     `json:"seq"`
	Children []*Menu `orm:"reverse(many)" json:"children"` // fk 的反向关系
	Icon     string  `orm:"size(32)" json:"icon"`
}

func MenuOne(id int) (*Menu, error) {
	o := orm.NewOrm()
	m := Menu{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
func MenuList() []*Menu {
	o := orm.NewOrm()
	query := o.QueryTable(MenuTBName())
	list := make([]*Menu, 0)

	query.All(&list)
	return list
}

//
////Resource 获取treegrid顺序的列表
//func GetMenuTree() []*Resource {
//}
//
////将资源列表转成treegrid格式
//func resourceList2TreeGrid(list []*Resource) []*Resource {
//	result := make([]*Resource, 0)
//	for _, item := range list {
//		if item.Parent == nil || item.Parent.Id == 0 {
//			item.Level = 0
//			result = append(result, item)
//			result = resourceAddSons(item, list, result)
//		}
//	}
//	return result
//}
//
////resourceAddSons 添加子菜单
//func resourceAddSons(cur *Resource, list, result []*Resource) []*Resource {
//	for _, item := range list {
//		if item.Parent != nil && item.Parent.Id == cur.Id {
//			cur.SonNum++
//			item.Level = cur.Level + 1
//			result = append(result, item)
//			result = resourceAddSons(item, list, result)
//		}
//	}
//	return result
//}

//将资源列表转成treegrid格式
func GetMenuTree() []*Menu {
	o := orm.NewOrm()
	query := o.QueryTable(MenuTBName())
	list := make([]*Menu, 0)

	query.All(&list)
	result := make([]*Menu, 0)
	for _, item := range list {
		if item.Parent == nil || item.Parent.Id == 0 {
			item = GetMenuTree_2(item, list)
			result = append(result, item)

		}
	}

	logs.Debug("A:%+v", result)
	return result
}

//resourceAddSons 添加子菜单
func GetMenuTree_2(cur *Menu, list []*Menu) *Menu {
	for _, item := range list {
		if item.Parent != nil && item.Parent.Id == cur.Id {
			item = GetMenuTree_2(item, list)
			cur.Children = append(cur.Children, item)
		}
	}
	return cur
}
