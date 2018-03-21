package sysinit

import (
	_ "github.com/chnzrb/myadmin/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
)

//初始化数据连接
func InitDatabase() {
	//读取配置文件，设置数据库参数
	//数据库类别
	dbType := beego.AppConfig.String("db_type")
	//连接名称
	dbAlias := beego.AppConfig.String(dbType + "::db_alias")
	//数据库名称
	dbName := beego.AppConfig.String(dbType + "::db_name")
	//数据库连接用户名
	dbUser := beego.AppConfig.String(dbType + "::db_user")
	//数据库连接用户名
	dbPwd := beego.AppConfig.String(dbType + "::db_pwd")
	//数据库IP（域名）
	dbHost := beego.AppConfig.String(dbType + "::db_host")
	//数据库端口
	dbPort := beego.AppConfig.String(dbType + "::db_port")
	switch dbType {
	case "sqlite3":
		orm.RegisterDataBase(dbAlias, dbType, dbName)
	case "mysql":
		dbCharset := beego.AppConfig.String(dbType + "::db_charset")
		orm.RegisterDataBase(dbAlias, dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+
			dbPort+")/"+dbName+"?charset="+dbCharset, 30)
	}
	//如果是开发模式，则显示命令信息
	//isDev := (beego.AppConfig.String("runmode") == "dev")
	//自动建表
	//orm.RunSyncdb("default", false, isDev)
	initNode()
	//initGame()
	//if isDev {
	//	orm.Debug = isDev
	//}
}
func initNode(){
	dbhost := beego.AppConfig.String("center_db::db_host")
	dbport := beego.AppConfig.String("center_db::db_port")
	dbuser := beego.AppConfig.String("center_db::db_user")
	dbpassword := beego.AppConfig.String("center_db::db_password")
	dbname := beego.AppConfig.String("center_db::db_name")
	//dsn = "root:gamehome1234@tcp(192.168.31.100:3306)/h5_center?charset=utf8"
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	err := orm.RegisterDataBase("center", "mysql", dsn)
	if err != nil {
		fmt.Println(dsn, err)
	}
}

func initGame(){
	dbhost := beego.AppConfig.String("center_db::db_host")
	dbport := beego.AppConfig.String("center_db::db_port")
	dbuser := beego.AppConfig.String("center_db::db_user")
	dbpassword := beego.AppConfig.String("center_db::db_password")
	dbname := "game"
	//dsn = "root:gamehome1234@tcp(192.168.31.100:3306)/h5_center?charset=utf8"
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	err := orm.RegisterDataBase("game", "mysql", dsn)
	if err != nil {
		fmt.Println(dsn, err)
	}
}
