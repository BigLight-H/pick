package service

import (
	"github.com/gocolly/colly"
	"pick/conf"
	"pick/models"
)

//书目录
func BookInfo(role *conf.MainRule, domin string, cache bool) ([]map[string]string, []map[string]string) {
	//图书信息
	var bookInfo []map[string]string
	//章节信息
	var chapterInfo []map[string]string
	//图片信息
	c := colly.NewCollector()
	// Find and visit all links
	//链接redis
	redisPool := models.ConnectRedisPool()
	defer redisPool.Close()
	c.OnXML(role.Table, func(e *colly.XMLElement) {
		//章节名
		title := e.ChildText(role.Link)
		//创建章节目录
		//util.MKdirs(dir+"\\"+title)
		//章节链接
		link := e.ChildText(role.Title)
		isExist, _ := redisPool.Do("HEXISTS", "chapter_links", link)
		//章节链接不存在redis里面采集
		if isExist != int64(1) {
			//章节更新时间
			ctime := e.ChildText(role.CTime)
			if ctime == "" {
				ctime = e.ChildText(role.NCTime)
			}
			img := GetDetail(role, link)
			chapterInfo = append(
				chapterInfo,
				map[string]string{"link": link, "title": title, "imgs": img, "ctime":ctime})
		}
		//不走缓存
		if cache {
			ctime := e.ChildText(role.CTime)
			if ctime == "" {
				ctime = e.ChildText(role.NCTime)
			}
			img := GetDetail(role, link)
			chapterInfo = append(
				chapterInfo,
				map[string]string{"link": link, "title": title, "imgs": img, "ctime":ctime})
		}
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
		//最后更新时间
		ltime := e.ChildText(role.LTime)
		if ltime == "" {
			ltime = e.ChildText(role.NTime)
		}
		bookInfo = append(
			bookInfo,
			map[string]string{"name": name, "year": year, "star": star, "author": author, "status": status, "intro": intro, "types": types, "tags": tags, "image": image, "ltime":ltime})
	})
	c.Visit(domin)
	return bookInfo, chapterInfo
}

func GetDetail(role *conf.MainRule, domin string) string {
	cs := colly.NewCollector()
	var img string
	cs.OnXML(role.Detail, func(e *colly.XMLElement) {
		//章节链接
		imgLink := e.ChildText(role.ImgSrc)
		//存入指定图片目录
		//imgArr := strings.Split(imgLink, "/")
		//name := imgArr[len(imgArr)-1]
		//util.DownloadJpg(imgLink, file_name+"\\"+name)
		img += imgLink+","
	})
	cs.Visit(domin)
	return  img
}

