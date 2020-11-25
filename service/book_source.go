package service

import (
	"fiction_web/models"
	"fiction_web/util"
	"github.com/astaxie/beego/orm"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

func GetBook(s string) []map[string]string {
	o := orm.NewOrm()
	c := colly.NewCollector()
	info := []map[string]string{}
	hub := []models.Hub{}
	o.QueryTable(new(models.Hub).TableName()).All(&hub)
	for _, v := range hub {
		c.OnXML(v.Root, func(e *colly.XMLElement) {
			link := ""
			new_list := "" //最新章节
			new_list_link := ""  //最新章节链接
			renew_time := ""  //最后更新时间
			status := ""  //书本状态:全本,连载
			img := ""  //书本封面图片
			if strings.Contains(e.ChildText(v.Link), "http") == false {
				link = v.BookHub + e.ChildText(v.Link)
			} else {
				link = e.ChildText(v.Link)
			}
			if v.NewList != ""{
				new_list = e.ChildText(v.NewList)
			}
			if v.NewListLink != ""{
				if strings.Contains(e.ChildText(v.NewListLink), "http") == false {
					new_list_link = v.BookHub + e.ChildText(v.NewListLink)
				} else {
					new_list_link = e.ChildText(v.NewListLink)
				}
			}
			if v.RenewTime != ""{
				renew_time = e.ChildText(v.RenewTime)
			}
			if v.RenewTime != ""{
				status = e.ChildText(v.Status)
			}
			if v.Image != "" {
				img = e.ChildText(v.Image)
			}
			name := e.ChildText(v.Name)
			uname := e.ChildText(v.Author)
			id := strconv.Itoa(v.Id)
			info = util.BackInfoMap(info, link, name, uname, id, new_list, new_list_link, renew_time, status, img)
		})
		c.Visit(v.Suffix + s)
	}
	return Deduplication(info)
}

//数据去重
func Deduplication(data []map[string]string) []map[string]string {
	src := make(map[string]interface{})
	info := []map[string]string{}
	for _, v := range data {
		if _,ok:=src[v["name"]];ok {
			continue
		} else {
			src[v["name"]] = v["name"]
			name := v["name"]
			link := v["link"]
			uname := v["uname"]
			new_list := v["new_list"]
			new_list_link := v["new_list_link"]
			renew_time := v["renew_time"]
			status := v["status"]
			id := v["id"]
			img := v["img"]
			info = util.BackInfoMap(info, link, name, uname, id, new_list, new_list_link, renew_time, status, img)
		}
	}
	return info
}
