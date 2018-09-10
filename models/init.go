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
	Db            *gorm.DB
	DbCenter      *gorm.DB
	DbCharge      *gorm.DB
	DbLoginServer *gorm.DB
)

// 现在是否在开服
//var IsNowOpenServer = false
//var IsNowOpenServerMap map[string]bool
var IsNowOpenServerMap = make( map[string]bool, 0)
//初始化
func init() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}
	initDb()
	initCenter()
	initCharge()
	initLoginServer()
}

func initDb() {
	//数据库名称
	dbName := beego.AppConfig.String("mysql" + "::db_name")
	//数据库用户名
	dbUser := beego.AppConfig.String("mysql" + "::db_user")
	//数据库密码
	dbPwd := beego.AppConfig.String("mysql" + "::db_password")
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
	Db.DB().SetMaxIdleConns(50)
}

func initCenter() {
	//数据库名称
	dbName := beego.AppConfig.String("center_db" + "::db_name")
	//数据库用户名
	dbUser := beego.AppConfig.String("center_db" + "::db_user")
	//数据库密码
	dbPwd := beego.AppConfig.String("center_db" + "::db_password")
	//数据库IP
	dbHost := beego.AppConfig.String("center_db" + "::db_host")
	//数据库端口
	dbPort := beego.AppConfig.String("center_db" + "::db_port")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd, dbHost, dbPort, dbName)
	var err error
	DbCenter, err = gorm.Open("mysql", dsn)
	DbCenter.SingularTable(true)
	utils.CheckError(err, "连接中心服失败")
}

func PingDb(db *gorm.DB) {
	sql := `show databases`
	err := db.Raw(sql).Error
	if err != nil {
		logs.Error("ping 数据库失败:%v", err)
	}
}

func initCharge() {
	//数据库名称
	dbName := beego.AppConfig.String("charge_db" + "::db_name")
	//数据库用户名
	dbUser := beego.AppConfig.String("charge_db" + "::db_user")
	//数据库密码
	dbPwd := beego.AppConfig.String("charge_db" + "::db_password")
	//数据库IP
	dbHost := beego.AppConfig.String("charge_db" + "::db_host")
	//数据库端口
	dbPort := beego.AppConfig.String("charge_db" + "::db_port")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd, dbHost, dbPort, dbName)
	var err error

	logs.Info("dbPwd:%v", dsn)
	DbCharge, err = gorm.Open("mysql", dsn)
	utils.CheckError(err, "连接充值数据库失败")
	//Db.LogMode(true)
	DbCharge.SetLogger(log.New(os.Stdout, "\r\n", 0))
	DbCharge.SingularTable(true)
}

func initLoginServer() {
	//数据库名称
	dbName := beego.AppConfig.String("login_server" + "::db_name")
	//数据库用户名
	dbUser := beego.AppConfig.String("login_server" + "::db_user")
	//数据库密码
	dbPwd := beego.AppConfig.String("login_server" + "::db_password")
	//数据库IP
	dbHost := beego.AppConfig.String("login_server" + "::db_host")
	//数据库端口
	dbPort := beego.AppConfig.String("login_server" + "::db_port")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd, dbHost, dbPort, dbName)
	var err error

	logs.Info("dbPwd:%v", dsn)
	DbLoginServer, err = gorm.Open("mysql", dsn)
	utils.CheckError(err, "连接登录服数据库失败")
	//Db.LogMode(true)
	DbLoginServer.SetLogger(log.New(os.Stdout, "\r\n", 0))
	DbLoginServer.SingularTable(true)
}

func TableName(name string) string {
	prefix := beego.AppConfig.String("db_dt_prefix")
	return prefix + name
}
