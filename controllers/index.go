package controllers

import (
	"github.com/beego/beego/server/web"
)

type MainController struct {
	web.Controller
}

func (c *MainController) Index() {
	//c.Data["Website"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.html"
}

func (c *MainController) Schedule() {
	c.TplName = "index.tpl"
}
