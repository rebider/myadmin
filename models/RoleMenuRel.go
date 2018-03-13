package models

import "time"

//角色与资源关系表
type RoleMenuRel struct {
	Id       int
	Role     *Role     `orm:"rel(fk)"`  //外键
	Menu *Menu `orm:"rel(fk)" ` // 外键
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
}

func (a *RoleMenuRel) TableName() string {
	return RoleMenuRelTBName()
}
