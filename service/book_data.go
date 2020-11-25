package service

import (
	"github.com/gocolly/colly"
	"pick/conf"
)

//书目录
func BookInfo(role *conf.MainRule, domin string) ([]map[string]string, []map[string]string) {
	//图书信息
	bookInfo := []map[string]string{}
	//章节信息
	chapterInfo := []map[string]string{}
	//图片信息
	c := colly.NewCollector()
	// Find and visit all links
	c.OnXML(role.Table, func(e *colly.XMLElement) {
		//章节链接
		link := e.ChildText(role.Title)
		img := GetDetail(role,link)
		//章节名
		title := e.ChildText(role.Link)
		chapterInfo = append(
			chapterInfo,
			map[string]string{"link": link, "title": title, "imgs": img})

	})
	c.OnXML(role.Body, func(e *colly.XMLElement) {
		//图书名
		name := e.ChildText(role.Name)
		//年代
		year := e.ChildText(role.Year)
		//分数
		star := e.ChildText(role.Star)
		//作者
		author := e.ChildText(role.Author)
		//标签
		tags := e.ChildText(role.Tags)
		//状态
		status := e.ChildText(role.Status)
		//简介
		intro := e.ChildText(role.Intro)
		//图书类型
		types := e.ChildText(role.Types)
		//图书封面
		image := e.ChildText(role.Image)

		bookInfo = append(
			bookInfo,
			map[string]string{"name": name, "year": year, "star": star, "author": author, "status": status, "intro": intro, "types": types, "tags": tags, "image":image})
	})
	c.Visit(domin)
	return bookInfo, chapterInfo
}

func GetDetail(role *conf.MainRule, domin string) string {
	cs := colly.NewCollector()
	var img string
	cs.OnXML(role.Detail, func(e *colly.XMLElement) {
		//章节链接
		img += e.ChildText(role.ImgSrc)+","
	})
	cs.Visit(domin)
	return  img
}

//书正文
//func BookContent(con models.Content, links string) []map[string]string {
//	c := colly.NewCollector()
//	info := []map[string]string{}
//	c.OnXML(con.Root, func(e *colly.XMLElement) {
//		name := e.ChildText(con.Name)
//		var ele  = NewXMLElement(e)
//		content := ele.ChildHtml(con.Content)
//		s_page := ""
//		x_page := ""
//		list := ""
//		switch con.Id {
//			case 3:
//				s_page = con.Domain + e.ChildText(con.SPage)
//				x_page = con.Domain + e.ChildText(con.XPage)
//				list = con.Domain + e.ChildText(con.List)
//				break
//			default:
//				s_page = e.ChildText(con.SPage)
//				x_page = e.ChildText(con.XPage)
//				list = e.ChildText(con.List)
//				break
//		}
//		info = append(
//			info,
//			map[string]string{"name": name, "content": content, "s_page": s_page, "x_page": x_page, "list": list})
//	})
//	c.Visit(links)
//	return info
//}
////书图片,简介
//func BookSynosis(con models.Synopsis, links string) []map[string]string  {
//	c := colly.NewCollector()
//	info := []map[string]string{}
//	c.OnXML(con.Root, func(e *colly.XMLElement) {
//		name := e.ChildText(con.Name)
//		writer := e.ChildText(con.Writer)
//		img := e.ChildText(con.Img)
//		synopsis := e.ChildText(con.Synopsis)
//		renew_time := e.ChildText(con.RenewTime)
//		info = append(
//			info,
//			map[string]string{"name": name, "writer": writer, "img": img, "synopsis": synopsis, "renew_time": renew_time})
//	})
//	c.Visit(links)
//	return info
//}
//
////检测书本是否更新
//func BookSynosisCheck(con models.Synopsis, links string) string  {
//	c := colly.NewCollector()
//	info := ""
//	c.OnXML(con.Root, func(e *colly.XMLElement) {
//		renew_time := e.ChildText(con.RenewTime)
//		info = renew_time
//	})
//	c.Visit(links)
//	return info
//}