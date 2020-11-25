package controllers

import (
	"github.com/astaxie/beego"
	"github.com/davecgh/go-spew/spew"
	"pick/conf"
	"pick/service"
)

type PickController struct {
	beego.Controller
}

func (p * PickController) Get() {
	domin := "https://www.webtoon.xyz/read/ring-my-bell/"
	//获取指定规则
	role := conf.Choose(1)
	//爬取图书
	bookInfo, chapterInfo := service.BookInfo(role, domin)
	//spew.Dump(bookInfo, chapterInfo)
	for k, v := range bookInfo {
		spew.Dump(k, v)
	}
	p.Data["json"] = chapterInfo
	p.ServeJSON()
}