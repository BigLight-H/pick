package routers

import (
	"github.com/beego/beego/server/web"
	"pick/controllers"
)

func init() {
    //beego.Router("/", &controllers.MainController{},"get:Index")
    //beego.Router("/schedule", &controllers.MainController{},"post:Schedule")
    web.Router("/collection", &controllers.PickController{},"post:Collection")
	web.Router("/lists", &controllers.PickController{},"post:Lists")
	web.Router("/add/redis", &controllers.PickController{},"*:SaveRedis")
}
