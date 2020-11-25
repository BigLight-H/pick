package controllers

import "C"
import (
	"github.com/astaxie/beego"
	"github.com/davecgh/go-spew/spew"
	"github.com/gocolly/colly"
)

type PickController struct {
	beego.Controller
}

func (p * PickController) Get() {
	domin := "https://www.webtoon.xyz/read/ring-my-bell/"
	spew.Dump(domin)
	c := colly.NewCollector()
	// Find and visit all links
	c.OnXML("//body/div[1]/div[1]/div[2]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[4]/div[1]/ul[1]/li", func(e *colly.XMLElement) {
		//e.Request.Visit(e.Attr("li"))
		link := e.ChildText("//a[1]/@href")
		title := e.ChildText("//a[1]/text()")
		//var ele  = NewXMLElement(e)
		//content := ele.ChildHtml(con.Content)
		spew.Dump(link,title)
	})
	c.OnXML("//body/div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[2]", func(e *colly.XMLElement) {
		//e.Request.Visit(e.Attr("li"))
		links := e.ChildText("//h1")
		//var ele  = NewXMLElement(e)
		//content := ele.ChildHtml(con.Content)
		spew.Dump(links)
	})
	c.Visit("https://www.webtoon.xyz/read/ring-my-bell/")
	p.ServeJSON()
}