package service

import (
	"fiction_web/models"
	"github.com/gocolly/colly"
)

func HubSource(source models.Source) map[string]interface{} {
	c := colly.NewCollector()
	info := []map[string]string{}
	list := []map[string]string{}
	list2 := []map[string]string{}
	list3 := []map[string]string{}
	list4 := []map[string]string{}
	list5 := []map[string]string{}
	lists := map[string]interface{}{}
	c.OnXML(source.Root, func(e *colly.XMLElement) {
		title := e.ChildText(source.TypeTitle)
		more := e.ChildText(source.MoreList)
		img := e.ChildText(source.Img)
		book_name := e.ChildText(source.BookName)
		book_name_link := e.ChildText(source.BookNameLink)
		book_author := e.ChildText(source.BookAuthor)
		book_author_link := e.ChildText(source.BookAuthorLink)
		book_mark := e.ChildText(source.BookMark)
		if title != "" {
			info = append(
				info,
				map[string]string{
					"title":            title,
					"more":             more,
					"book_name":        book_name,
					"img":              img,
					"book_name_link":   source.PcDomain+book_name_link,
					"book_author":      book_author,
					"book_author_link": source.PcDomain+book_author_link,
					"book_mark":        book_mark,
				})
		}
	})

	c.OnXML(source.ListTypeRoot, func(e *colly.XMLElement) {
		list_type_name := e.ChildText(source.ListTypeName)
		list_book_name := e.ChildText(source.ListBookName)
		list_book_author := e.ChildText(source.ListBookAuthor)
		list_type_link := e.ChildText(source.ListTypeLink)
		list_book_name_link := e.ChildText(source.ListBookNameLink)
		list_book_author_link := e.ChildText(source.ListBookAuthorLink)
		if list_book_name != "" {
			list = append(
				list,
				map[string]string{
					"list_type_name":        list_type_name,
					"list_book_name":        list_book_name,
					"list_book_author":      list_book_author,
					"list_type_link":        source.PcDomain+list_type_link,
					"list_book_name_link":   source.PcDomain+list_book_name_link,
					"list_book_author_link": source.PcDomain+list_book_author_link,
				})
		}
	})

	c.OnXML("//div[6]//ul[1]/li", func(e *colly.XMLElement) {
		list_type_name := e.ChildText(source.ListTypeName)
		list_book_name := e.ChildText(source.ListBookName)
		list_book_author := e.ChildText(source.ListBookAuthor)
		list_type_link := e.ChildText(source.ListTypeLink)
		list_book_name_link := e.ChildText(source.ListBookNameLink)
		list_book_author_link := e.ChildText(source.ListBookAuthorLink)
		if list_book_name != "" {
			list2 = append(
				list2,
				map[string]string{
					"list_type_name":        list_type_name,
					"list_book_name":        list_book_name,
					"list_book_author":      list_book_author,
					"list_type_link":        source.PcDomain+list_type_link,
					"list_book_name_link":   source.PcDomain+list_book_name_link,
					"list_book_author_link": source.PcDomain+list_book_author_link,
				})
		}
	})
	c.OnXML("//div[7]//ul[1]/li", func(e *colly.XMLElement) {
		list_type_name := e.ChildText(source.ListTypeName)
		list_book_name := e.ChildText(source.ListBookName)
		list_book_author := e.ChildText(source.ListBookAuthor)
		list_type_link := e.ChildText(source.ListTypeLink)
		list_book_name_link := e.ChildText(source.ListBookNameLink)
		list_book_author_link := e.ChildText(source.ListBookAuthorLink)
		if list_book_name != "" {
			list3 = append(
				list3,
				map[string]string{
					"list_type_name":        list_type_name,
					"list_book_name":        list_book_name,
					"list_book_author":      list_book_author,
					"list_type_link":        source.PcDomain+list_type_link,
					"list_book_name_link":   source.PcDomain+list_book_name_link,
					"list_book_author_link": source.PcDomain+list_book_author_link,
				})
		}
	})
	c.OnXML("//div[8]//ul[1]/li", func(e *colly.XMLElement) {
		list_type_name := e.ChildText(source.ListTypeName)
		list_book_name := e.ChildText(source.ListBookName)
		list_book_author := e.ChildText(source.ListBookAuthor)
		list_type_link := e.ChildText(source.ListTypeLink)
		list_book_name_link := e.ChildText(source.ListBookNameLink)
		list_book_author_link := e.ChildText(source.ListBookAuthorLink)
		if list_book_name != "" {
			list4 = append(
				list4,
				map[string]string{
					"list_type_name":        list_type_name,
					"list_book_name":        list_book_name,
					"list_book_author":      list_book_author,
					"list_type_link":        source.PcDomain+list_type_link,
					"list_book_name_link":   source.PcDomain+list_book_name_link,
					"list_book_author_link": source.PcDomain+list_book_author_link,
				})
		}
	})
	c.OnXML("//div[9]//ul[1]/li", func(e *colly.XMLElement) {
		list_type_name := e.ChildText(source.ListTypeName)
		list_book_name := e.ChildText(source.ListBookName)
		list_book_author := e.ChildText(source.ListBookAuthor)
		list_type_link := e.ChildText(source.ListTypeLink)
		list_book_name_link := e.ChildText(source.ListBookNameLink)
		list_book_author_link := e.ChildText(source.ListBookAuthorLink)
		if list_book_name != "" {
			list5 = append(
				list5,
				map[string]string{
					"list_type_name":        list_type_name,
					"list_book_name":        list_book_name,
					"list_book_author":      list_book_author,
					"list_type_link":        source.PcDomain+list_type_link,
					"list_book_name_link":   source.PcDomain+list_book_name_link,
					"list_book_author_link": source.PcDomain+list_book_author_link,
				})
		}
	})

	c.Visit(source.Domain)

	lists["info"] = info
	lists["list"] = list
	lists["list2"] = list2
	lists["list3"] = list3
	lists["list4"] = list4
	lists["list5"] = list5
	return lists
}


func BookType(t models.Type, link string) map[string]interface{} {
	c := colly.NewCollector()
	info := []map[string]string{}
	page := map[string]string{}
	data := map[string]interface{}{}
	c.OnXML(t.Root, func(e *colly.XMLElement) {
		book_name := e.ChildText(t.BookName)
		book_name_link := t.DomainPc + e.ChildText(t.BookNameLink)
		book_author := e.ChildText(t.BookAuthor)
		info = append(
			info,
			map[string]string{
				"book_name":      book_name,
				"book_name_link": book_name_link,
				"book_author":    book_author,
			})
	})
	c.OnXML(t.BookPageRoot, func(e *colly.XMLElement) {
		book_first_page := ""
		book_previous := ""
		book_next_page := ""
		book_last_page := ""
		if e.ChildText(t.BookNextPage) != "" {
			book_first_page = t.Domain + e.ChildText(t.BookFirstPage)
			book_previous = t.Domain + e.ChildText(t.BookPrevious)
			book_next_page = t.Domain + e.ChildText(t.BookNextPage)
			book_last_page = t.Domain + e.ChildText(t.BookLastPage)
		} else {
			book_next_page = t.Domain + e.ChildText(t.BookFirstPage)
			book_last_page = t.Domain + e.ChildText(t.BookPrevious)
		}
		page = map[string]string{
			"book_first_page": book_first_page,
			"book_previous":   book_previous,
			"book_next_page":  book_next_page,
			"book_last_page":  book_last_page,
		}
	})

	c.Visit(link)
	data["list"] = info
	data["page"] = page
	return data
}

func GetBoard(t models.Leaderboard, link string) map[string]interface{} {
	c := colly.NewCollector()
	info := []map[string]string{}
	page := map[string]string{}
	data := map[string]interface{}{}
	c.OnXML(t.Root, func(e *colly.XMLElement) {
		book_name := e.ChildText(t.BookName)
		book_name_link := t.DomainPc + e.ChildText(t.BookNameLink)
		book_author := e.ChildText(t.BookAuthor)
		book_type := e.ChildText(t.BookType)
		info = append(
			info,
			map[string]string{
				"book_name":      book_name,
				"book_name_link": book_name_link,
				"book_author":    book_author,
				"book_type":      book_type,
			})
	})
	c.OnXML(t.BookPageRoot, func(e *colly.XMLElement) {
		book_first_page := ""
		book_previous := ""
		book_next_page := ""
		book_last_page := ""
		if e.ChildText(t.BookNextPage) != "" {
			book_first_page = t.Domain + e.ChildText(t.BookFirstPage)
			book_previous = t.Domain + e.ChildText(t.BookPrevious)
			book_next_page = t.Domain + e.ChildText(t.BookNextPage)
			book_last_page = t.Domain + e.ChildText(t.BookLastPage)
		} else {
			book_next_page = t.Domain + e.ChildText(t.BookFirstPage)
			book_last_page = t.Domain + e.ChildText(t.BookPrevious)
		}
		page = map[string]string{
			"book_first_page": book_first_page,
			"book_previous":   book_previous,
			"book_next_page":  book_next_page,
			"book_last_page":  book_last_page,
		}
	})

	c.Visit(link)
	data["list"] = info
	data["page"] = page
	return data
}

func GetBookEnd(t models.Completed, link string) map[string]interface{} {
	c := colly.NewCollector()
	info := []map[string]string{}
	page := map[string]string{}
	data := map[string]interface{}{}
	c.OnXML(t.Root, func(e *colly.XMLElement) {
		book_name := e.ChildText(t.BookName)
		book_name_link := t.DomainPc + e.ChildText(t.BookNameLink)
		book_author := e.ChildText(t.BookAuthor)
		book_type := e.ChildText(t.BookType)
		info = append(
			info,
			map[string]string{
				"book_name":      book_name,
				"book_name_link": book_name_link,
				"book_author":    book_author,
				"book_type":      book_type,
			})
	})
	c.OnXML(t.BookPageRoot, func(e *colly.XMLElement) {
		book_first_page := ""
		book_previous := ""
		book_next_page := ""
		book_last_page := ""
		if e.ChildText(t.BookNextPage) != "" {
			book_first_page = t.Domain + e.ChildText(t.BookFirstPage)
			book_previous = t.Domain + e.ChildText(t.BookPrevious)
			book_next_page = t.Domain + e.ChildText(t.BookNextPage)
			book_last_page = t.Domain + e.ChildText(t.BookLastPage)
		} else {
			book_next_page = t.Domain + e.ChildText(t.BookFirstPage)
			book_last_page = t.Domain + e.ChildText(t.BookPrevious)
		}
		page = map[string]string{
			"book_first_page": book_first_page,
			"book_previous":   book_previous,
			"book_next_page":  book_next_page,
			"book_last_page":  book_last_page,
		}
	})

	c.Visit(link)
	data["list"] = info
	data["page"] = page
	return data
}
