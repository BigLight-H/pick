package main

import (
	"github.com/beego/beego/client/orm"
	"github.com/beego/beego/server/web"
	"github.com/beego/beego/server/web/filter/cors"
	_ "github.com/go-sql-driver/mysql"
	"pick/models"
	_ "pick/routers"
)

func init()  {
	orm.Debug = true
	models.Init()
}

func main() {
	//InsertFilter是提供一个过滤函数
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		//允许访问所有源
		AllowAllOrigins: true,
		//可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
		//其中Options跨域复杂请求预检
		AllowMethods:   []string{"*"},
		//指的是允许的Header的种类
		AllowHeaders: 	[]string{"*"},
		//公开的HTTP标头列表
		ExposeHeaders:	[]string{"Content-Length"},
		//如果设置，则允许共享身份验证凭据，例如cookie
		AllowCredentials: true,
	}))
	web.Run()
}

