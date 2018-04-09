package models

import (
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/chnzrb/myadmin/utils"
	"fmt"
	"log"
	"os"
	"github.com/astaxie/beego/logs"
)
var (
	Db *gorm.DB
	DbCenter *gorm.DB
)

//初始化
func init() {
	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return defaultTableName
	}
	initDb()
	initCenter()
}

func initDb() {
	//数据库名称
	dbName := beego.AppConfig.String("mysql" + "::db_name")
	//数据库用户名
	dbUser := beego.AppConfig.String("mysql" + "::db_user")
	//数据库密码
	dbPwd := beego.AppConfig.String("mysql" + "::db_pwd")
	//数据库IP
	dbHost := beego.AppConfig.String("mysql" + "::db_host")
	//数据库端口
	dbPort := beego.AppConfig.String("mysql" + "::db_port")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd, dbHost, dbPort, dbName)
	//dsn := "root:game1234@tcp(192.168.31.100:3306)/center?charset=utf8&parseTime=True&loc=Local"
	var err error

	logs.Info("dbPwd:%v", dsn)
	Db, err = gorm.Open("mysql", dsn)
	utils.CheckError(err, "连接后台数据库失败")
	//Db.LogMode(true)
	Db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	Db.SingularTable(true)
}

func initCenter() {
	dsn := "root:game1234@tcp(192.168.31.100:3306)/center?charset=utf8&parseTime=True&loc=Local"
	var err error
	DbCenter, err = gorm.Open("mysql", dsn)
	DbCenter.SingularTable(true)
	utils.CheckError(err, "连接中心服失败")
}

func TableName(name string) string {
	prefix := beego.AppConfig.String("db_dt_prefix")
	return prefix + name
}
