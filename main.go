package main

import (
	"github.com/astaxie/beego"
	_ "github.com/chnzrb/myadmin/routers"
	_ "github.com/chnzrb/myadmin/sysinit"

	//"github.com/astaxie/beego/plugins/cors"
	//"github.com/astaxie/beego/plugins/cors"
)

func main() {
	beego.Run()

}
