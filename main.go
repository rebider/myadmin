package main

import (
	"github.com/astaxie/beego"
	_ "github.com/chnzrb/myadmin/routers"
	_ "github.com/chnzrb/myadmin/sysinit"
)

func main() {
	beego.Run()
}
