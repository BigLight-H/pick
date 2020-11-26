package controllers

import (
	"github.com/astaxie/beego"
	"pick/conf"
	"pick/service"
	"pick/util"
)

type PickController struct {
	beego.Controller
}

func (p * PickController) Get() {
	domin := "https://www.webtoon.xyz/read/ring-my-bell/"
	//新建目录,没有就新建->返回目录名
	dir := util.ThisMkdir(domin)
	//获取指定规则
	role := conf.Choose(1)
	//爬取图书
	bookInfo, chapterInfo := service.BookInfo(role, domin)
	//spew.Dump(bookInfo, chapterInfo)
	for _, v := range bookInfo {
		//下载封面图片
		util.DownloadJpg(v["image"], dir+"\\thumb.jpg")
	}
	for _, s := range chapterInfo {
		//切割链接
		_,son := util.GetSubdirectory(s["link"])
		//创建章节目录
		util.MKdirs(dir+"\\"+son)
		//创建协程
		//for i := 0; i <= 4; i++ {
			//创建协程处理
			util.DoWork(dir+"\\"+son, s["imgs"])
		//}












	}
	p.Data["json"] = chapterInfo
	p.ServeJSON()
}