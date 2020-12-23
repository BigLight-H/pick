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
    beego.Router("/", &controllers.MainController{},"get:Index")
    beego.Router("/schedule", &controllers.MainController{},"post:Schedule")
    beego.Router("/collection", &controllers.PickController{},"post:Collection")
    beego.Router("/lists", &controllers.PickController{},"post:Lists")
    beego.Router("/save", &controllers.PickController{},"get:SaveRedis")
}
