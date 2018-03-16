package models

import (
	"fmt"
	"github.com/chnzrb/myadmin/utils"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"sort"
)

func (a *Menu) TableName() string {
	return MenuTBName()
}

type MenuQueryParam struct {
	BaseQueryParam
}

//Menu 权限控制表
type Menu struct {
	Id       int     `json:"id,string"`
	Title    string  `orm:"size(64)" json:"title"` //标题
	Name     string  `orm:"size(64)" json:"name"`
	Parent   *Menu   `orm:"null;rel(fk) " json:"-"` // RelForeignKey relation
	ParentId int     `orm:"-" json:"parentId,string"`      // RelForeignKey relation
	Seq      int     `json:"seq"`
	Children []*Menu `orm:"reverse(many)" json:"children"` // fk 的反向关系
	Icon     string  `orm:"size(32)" json:"icon"`
	RoleMenuRel []*RoleMenuRel `orm:"reverse(many)" json:"-"` // 设置一对多的反向关系
}

type menuSlice [] *Menu

func (s menuSlice) Len() int { return len(s) }
func (s menuSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s menuSlice) Less(i, j int) bool { return s[i].Seq < s[j].Seq }


func sortMenuTree(list []*Menu) []*Menu {
	sort.Sort(menuSlice(list))
	for _, item := range list {
		if item.Children != nil  {
			sortMenuTree(item.Children)
		}
	}
	return list
}

func MenuOne(id int) (*Menu, error) {
	o := orm.NewOrm()
	m := Menu{Id: id}
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
func MenuList() []*Menu {
	query := orm.NewOrm().QueryTable(MenuTBName())
	data := make([]*Menu, 0)
	query.All(&data)
	for _, e := range data {
		if e.Parent != nil {
			e.ParentId = e.Parent.Id
		}
	}
	return data
}

func MenuTreeGrid4Parent(id int) []*Menu {
	list := MenuList()
	tmpList := make([] *Menu, 0)
	if id > 0 {
		for _, e := range list {
			if CanParentMenu(id, e.Id) {
				tmpList = append(tmpList, e)
			}
		}
	} else {
		tmpList = list
	}
	return tmpList
}

func CanParentMenu(menuId int, parentMenuId int) bool {
	parentMenu, _ := MenuOne(parentMenuId)
	if (parentMenu.Parent != nil && parentMenu.Parent.Id == menuId) || menuId == parentMenuId {
		return false
	}
	if parentMenu.Parent == nil || parentMenu.Parent.Id == 0 {
		return true
	}
	return CanParentMenu(menuId, parentMenu.Parent.Id)
}

//根据用户获取有权管理的菜单列表
func GetMenuListByUserId(userId int) []*Menu {
	var list []*Menu
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
		sql = fmt.Sprintf(`SELECT * FROM %s  Order By seq asc,Id asc`, MenuTBName())
		o.Raw(sql).QueryRows(&list)
	} else {
		//	//联查多张表，找出某用户有权管理的
		sql = fmt.Sprintf(`SELECT DISTINCT T2.*
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
		INNER JOIN %s AS T2 ON T2.id = T0.menu_id
		WHERE T1.user_id = ?  Order By T2.seq asc,T2.id asc`, RoleMenuRelTBName(), RoleUserRelTBName(), MenuTBName())
		o.Raw(sql, userId).QueryRows(&list)
	}
	result := list
	for _, e := range list {
		if e.Parent != nil {
			e.ParentId = e.Parent.Id
		}
	}
	return result
}

func TranMenuList2MenuTree(menuList []*Menu) []*Menu {
	menuTree := make([]*Menu, 0)
	for _, item := range menuList {
		if item.Parent == nil || item.Parent.Id == 0 {
			item = TranMenuList2MenuTree_(item, menuList)
			menuTree = append(menuTree, item)
		}
	}
	sortMenuTree(menuTree)
	return menuTree
}

func TranMenuList2MenuTree_(cur *Menu, list []*Menu) *Menu {
	for _, item := range list {
		if item.Parent != nil && item.Parent.Id == cur.Id {
			item = TranMenuList2MenuTree_(item, list)
			cur.Children = append(cur.Children, item)
		}
	}
	return cur
}
