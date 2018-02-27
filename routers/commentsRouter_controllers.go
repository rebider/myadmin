package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/chnzrb/myadmin/controllers:GameServerController"] = append(beego.GlobalControllerRouter["github.com/chnzrb/myadmin/controllers:GameServerController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/staticblock/:key`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
