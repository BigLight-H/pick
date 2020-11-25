package routers

import (
	"github.com/astaxie/beego/orm"
	"pick/controllers"
	"github.com/astaxie/beego"
)

func init()  {
	orm.Debug = true
	//models.Init()
}

func init() {
    beego.Router("/", &controllers.PickController{})
}
