package models

import "time"

//角色与用户关系
type RoleUserRel struct {
	Id          int
	//Role        *Role        `orm:"rel(fk)"`  //外键
	//User *User `orm:"rel(fk)" ` // 外键
	UserId  int `orm:"-" ` // 外键
	RoleId  int `orm:"-" ` // 外键
	Created     time.Time    `orm:"auto_now_add;type(datetime)"`
}

func (a *RoleUserRel) TableName() string {
	return RoleUserRelTBName()
}

// 删除用户列表
func DeleteRoleUserRelByUserIdList(userIdList [] int) (int, error) {
	var count int
	err := Db.Where("user_id in (?)", userIdList).Delete(&RoleUserRel{}).Count(&count).Error
	return  count, err
}
