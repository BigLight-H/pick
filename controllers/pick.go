package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"github.com/gocolly/colly"
	"pick/conf"
	"pick/models"
	"pick/service"
	"pick/util"
	"strconv"
	"strings"
	"sync"
)

type PickController struct {
	beego.Controller
}

//Json结构体
type Json struct {
	Msg string
	Status int
}

func (p *PickController) Lists() {
	domain := "https://www.webtoon.xyz/webtoons/"
	//图片信息
	c := colly.NewCollector()
	d := c.Clone()
	// Find and visit all links
	c.OnXML("//body/div[1]/div[1]/div[2]/div[2]/div[1]/div[1]/div[1]/div[1]/div[2]/div[1]/div[2]/div[2]/div[1]/div[2]/div[1]/a[@class='last']", func(e *colly.XMLElement) {
		lastLink := e.ChildText("//@href")
		allPage := util.ChapterOrder(lastLink, "/", 2)
		allNum, _ := strconv.Atoi(allPage)
		for i := 1; i <= allNum; i++ {
			//获取分页数据并存入数据库
			pageDomain := domain+"/page/"+strconv.Itoa(i)+"/"
			d.OnXML("//body/div[1]/div[1]/div[2]/div[2]/div[1]/div[1]/div[1]/div[1]/div[2]/div[1]/div[2]/div[2]/div[1]/div[1]/div", func(f *colly.XMLElement) {
				for a := 1; a <= 2; a++ {
					//链接redis
					redisPool := service.ConnectRedis()
					defer redisPool.Close()
					link := f.ChildText("//div[1]/div["+strconv.Itoa(a)+"]/div[1]/div[1]/a/@href")
					title := f.ChildText("//div[1]/div["+strconv.Itoa(a)+"]/div[1]/div[1]/a/@title")
					lastCharpter := f.ChildText("//div[1]/div["+strconv.Itoa(a)+"]/div[1]/div[2]/div[3]/div[1]/span[1]")
					//spew.Dump(lastCharpter)
					isExist, _ := redisPool.Do("HEXISTS", "book_all_lists", link)
					if isExist != int64(1) {
						o := orm.NewOrm()
						lists := models.Links{}
						lists.BookLink = link
						lists.BookName = title
						lists.LastChapter = lastCharpter
						lists.Status = 1
						lid, err := o.Insert(&lists)
						if err == nil {
							_, err := redisPool.Do("HSET", "book_all_lists", link, lid)
							if err != nil {
								spew.Dump("漫画链接存入错误")
							}
						}
					} else {
						lid, err := redisPool.Do("HGET", "book_all_lists", link)
						lId, _ := redis.Int(lid, err)
						//查看更新
						o := orm.NewOrm()
						lists := models.Links{Id:lId}
						lists.LastChapter = lastCharpter
						if num, err := o.Update(&lists, "LastChapter"); err == nil {
							//有更新改变字段值
							if num > 0 {
								u := orm.NewOrm()
								up := models.Links{Id:lId}
								up.Status = 1
								if num, err := u.Update(&lists, "Status"); err == nil {
									spew.Dump(num)
								}
							}
						}
					}
				}
			})
			d.Visit(pageDomain)
		}
	})
	c.Visit(domain)
	p.MsgBack("采集全部图书链接完成", 1)
}

func (p * PickController) Collection() {
	//源ID
	rootId, _ := strconv.Atoi(p.GetString("id"))
	//图书列表
	rootList := p.GetString("list")
	if rootList == "" {
		p.Ctx.WriteString("rootList is empty")
		return
	}
	list := strings.Split(rootList, " ")
	//all := int(len(list))
	//if all > 5 {
	//	a := float64(all / 10)
	//	num := int(math.Floor(a + 0/5))
	//	listArr := util.SplitArray(list, num)
	//	for _, val := range listArr {
	//		var wg sync.WaitGroup
	//		for _, value := range val {
	//			wg.Add(1)
	//			go Comics(value, rootId)
	//			wg.Done()
	//		}
	//		wg.Wait()
	//	}
	//} else {
		var wg sync.WaitGroup
		for _, val := range list {
			//创建协程
			wg.Add(1)
			go Comics(val, rootId)
			wg.Done()
		}
		wg.Wait()
	//}

	p.MsgBack("采集资源完成", 1)
}

//返回json信息
func (p *PickController) MsgBack(msg string, status int)  {
	data := &Json{
		msg,
		status,
	}
	p.Data["json"] = data
	p.ServeJSON()
}

//操作图书
func Comics(domin string, rootId int) {

	//新建目录,没有就新建->返回目录名
	dir := util.ThisMkdir(domin)
	//获取指定规则
	role := conf.Choose(rootId)
	//爬取图书
	bookInfo, chapterInfo := service.BookInfo(role, domin)
	//图书ID
	bId := 0
	for _, v := range bookInfo {
		//链接redis
		redisPool := service.ConnectRedis()
		defer redisPool.Close()
		isExist, _ := redisPool.Do("HEXISTS", "comic_link", domin)
		if isExist != int64(1) {
			//下载封面图片
			util.DownloadJpg(v["image"], dir+"\\thumb.jpg")
			//存入数据库
			o := orm.NewOrm()
			book := models.Book{}
			book.UniqueId = domin
			book.BookName = v["name"]
			book.Tags = v["tags"]
			book.Summary = v["intro"]
			book.End = 0
			if v["status"] == "Completed" {
				book.End = 1
			}
			book.AuthorName = v["author"]
			book.CoverUrl = v["image"]
			book.Type = v["types"]
			book.Star = v["star"]
			book.Year = v["year"]
			book.LastTime = v["ltime"]
			id, _ := o.Insert(&book)
			//if err != nil {
			//	os.Exit(1)
			//}
			bId = int(id)
			_, err := redisPool.Do("HSET", "comic_link", domin, bId)
			if err != nil {
				spew.Dump("漫画ID存入错误")
			}
		} else {
			comicId, err := redisPool.Do("HGET", "comic_link", domin)
			bookId, _ := redis.Int(comicId, err)
			bId = bookId
			//修改书本最新更新时间
			b := orm.NewOrm()
			oldBook := models.Book{Id: bId}
			oldBook.LastTime = v["ltime"]
			if num, err := b.Update(&oldBook, "LastTime"); err != nil {
				fmt.Println(num)
			}
		}
	}
	for _, s := range chapterInfo {
		//切割链接
		_, son := util.GetSubdirectory(s["link"])
		//存入章节表
		c := orm.NewOrm()
		chapter := models.Chapter{}
		chapter.ChapterName = s["title"]
		chapter.BookId = bId
		chapter.ChapterOrder = util.ChapterOrder(son, "-", 1)
		chapter.ChapterLink = s["link"]
		chapter.LastTime = s["ctime"]
		//cId, _ := c.Insert(&chapter)
		cId := 0

		// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
		if created, cid, err := c.ReadOrCreate(&chapter, "ChapterLink"); err == nil {
			if created {
				cId = int(cid)
			} else {
				cId = int(cid)
				//修改章节最后更新时间
				up := orm.NewOrm()
				oldChapter := models.Chapter{Id: cId}
				oldChapter.LastTime = s["ctime"]
				if num, err := up.Update(&oldChapter, "LastTime"); err != nil {
					fmt.Println(num)
				}
			}
		}

		//创建章节目录
		util.MKdirs(dir + "\\" + son)

		if cId > 0 {
			go util.DoWork(dir+"\\"+son, s["imgs"], bId, cId)
		}

	}

}