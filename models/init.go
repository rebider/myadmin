package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/chnzrb/myadmin/utils"
)
var (
	//Db *gorm.DB
	DbCenter *gorm.DB
)

//初始化
func init() {
	initCenter()
	orm.RegisterModel(
		new(User),
		new(Resource),
		new(Role),
		new(RoleResourceRel),
		new(RoleUserRel),
		new(GameServer),
		new(ServerNode),
		new(Menu),
		new(RoleMenuRel),
		new(Player),
		)
}

func initCenter() {
	dsn := "root:game1234@tcp(192.168.31.100:3306)/center?charset=utf8&parseTime=True&loc=Local"
	var err error
	DbCenter, err = gorm.Open("mysql", dsn)
	utils.CheckError(err, "连接中心服失败")
}

//下面是统一的表名管理
func TableName(name string) string {
	prefix := beego.AppConfig.String("db_dt_prefix")
	return prefix + name
}

//获取 User 对应的表名称
func UserTBName() string {
	return TableName("user")
}

//获取 Resource 对应的表名称
func ResourceTBName() string {
	return TableName("resource")
}

//获取 Menu 对应的表名称
func MenuTBName() string {
	return TableName("menu")
}

//获取 Role 对应的表名称
func RoleTBName() string {
	return TableName("role")
}

//角色与资源多对多关系表
func RoleResourceRelTBName() string {
	return TableName("role_resource_rel")
}

//角色与菜单多对多关系表
func RoleMenuRelTBName() string {
	return TableName("role_menu_rel")
}

//角色与用户多对多关系表
func RoleUserRelTBName() string {
	return TableName("role_user_rel")
}
