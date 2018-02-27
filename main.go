package main

import (
	"github.com/astaxie/beego"
	_ "github.com/chnzrb/myadmin/routers"
	_ "github.com/chnzrb/myadmin/sysinit"

	//"github.com/astaxie/beego/plugins/cors"
)

func main() {
	//beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
	//	AllowAllOrigins:  "localhost:9527",
	//	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	//	AllowHeaders:     []string{"Origin", "Authorization", "Cookie", "Host", "Referer", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "x-token"},
	//	ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	//	AllowCredentials: true,
	//}))
	beego.Run()


}
