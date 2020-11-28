package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
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

//获取全部书籍链接
func (p *PickController) Lists() {
	domain := beego.AppConfig.String("source_1")
	util.BookLists(domain)
	p.MsgBack("采集全部图书链接完成", 1)
}

//采集图书
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