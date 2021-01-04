package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"os"
	"pick/conf"
	"pick/models"
	"pick/service"
	"pick/util"
	"strconv"
	"strings"
	"sync"
	"time"
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
	go util.BookLists(domain)
	p.MsgBack("采集全部图书链接完成", 1)
}

//采集图书
func (p * PickController) Collection() {
	//源ID
	rootId, _ := strconv.Atoi(p.GetString("id"))
	//是否用缓存
	cacheId, _ := strconv.Atoi(p.GetString("cache"))
	if rootId < 1 {
		spew.Dump("规则ID不能为空")
		os.Exit(1)
	}

	caches := false
	if cacheId > 0 {
		caches = true
	}
	//图书列表
	rootList := p.GetString("list")
	if rootList == "" {
		p.Ctx.WriteString("rootList is empty")
		return
	}
	list := strings.Split(rootList, "\n")


	all := int(len(list))
	if all>1 {
		var wg sync.WaitGroup
		for _, val := range list {
			//创建协程
			wg.Add(1)
			go Comics(val, rootId, caches)
			wg.Done()
		}
		wg.Wait()
	}else {
		go Comics(rootList, rootId, caches)
	}
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
func Comics(domin string, rootId int, caches bool) {
	//获取指定规则
	role := conf.Choose(rootId)
	//爬取图书
	bookInfo, chapterInfo := service.BookInfo(role, domin, caches)
	//图书ID
	bId := 0
	o := orm.NewOrm()
	for _, v := range bookInfo {
		//链接redis
		redisPool := models.ConnectRedisPool()
		defer redisPool.Close()
		isExist, _ := redisPool.Do("HEXISTS", "comic_links", domin)
		if isExist == int64(0) {
			//存入数据库
			book := models.BookList{}
			book.DomainName = domin
			book.BookTitle = v["name"]
			book.BookTags = v["tags"]
			book.BookProfile = v["intro"]
			book.BookStat = 0
			if v["status"] == "Completed" {
				book.BookStat = 1
			}
			book.BookAuthor = v["author"]
			//book.CoverUrl = v["image"]
			book.NowStatus = v["types"]
			book.Star = v["star"]
			book.Year = v["year"]
			book.LastTime = v["ltime"]
			book.TimesCollect = util.RandomNum(4)
			book.TimesBuy = util.RandomNum(4)
			book.TimesRead = util.RandomNum(4)
			book.TimesSubscribed = util.RandomNum(4)
			book.UserBuy = util.RandomNum(4)
			book.UserRead = util.RandomNum(4)
			book.BookThumbnail = ""
			book.IsAgeLimit = 1
			book.Status = 0
			book.CreateTime = time.Now().Unix()
			book.UpdateTime = time.Now().Unix()
			id, _ := o.Insert(&book)
			bId = int(id)
			if id > int64(0) && bId > 0 {
				up := models.BookList{BookId: bId}
				bookid := strconv.Itoa(bId)
				up.BookThumbnail = bookid+"/"+bookid+"_thumb.jpg"
				num, _ := o.Update(&up, "BookThumbnail")
				if num > int64(0) {
					//新建文件目录
					util.MKdirs(bookid)
					//下载封面图片到目录
					util.DownloadJpg(v["image"], bookid+"/"+bookid+"_thumb.jpg", false)
				}
				_, err := redisPool.Do("HSET", "comic_links", domin, bId)
				if err != nil {
					spew.Dump("漫画ID存入错误")
				}
			}
		} else {
			comicId, err := redisPool.Do("HGET", "comic_links", domin)
			bookId, _ := redis.Int(comicId, err)
			bId = bookId
			//修改书本最新更新时间
			if bId > 0 {
				oldBook := models.BookList{BookId: bId}
				oldBook.LastTime = v["ltime"]
				if num, err := o.Update(&oldBook, "LastTime"); err != nil {
					fmt.Println(num)
				}
			}
		}
	}
	bookid := strconv.Itoa(bId)
	if bId > 0 {
		for _, s := range chapterInfo {
			//切割链接
			_, son := util.GetSubdirectory(s["link"])
			epid := util.ChapterOrder(son, "-", 1)
			//存入章节表
			c := orm.NewOrm()
			chapter := models.BookEpisode{}
			chapter.BookId = bId
			chapter.EpisodeTitle = s["title"]
			chapter.EpisodeId = epid
			imgLen, fImg := util.SumImgs(s["imgs"])
			chapter.EpisodeImgtotal = imgLen
			chapter.EpisodeThumbnail = ""
			chapter.LastTime = s["ctime"]
			chapter.Link = s["link"]
			//cId, _ := c.Insert(&chapter)
			cId := 0
			//图片路径
			imgRole := bookid +"/"+ epid +"/"+ bookid +"0"+ epid +"_thumb.jpg"
			// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
			if created, cid, err := c.ReadOrCreate(&chapter, "Link"); err == nil {
				cId = int(cid)
				if created {
					upimg := models.BookEpisode{Id: cId}
					upimg.EpisodeThumbnail = bookid +"/"+ epid +"/"+ bookid +"0"+ epid +"_thumb.jpg"
					num, _ := c.Update(&upimg, "EpisodeThumbnail")
					if num > int64(0) {
						//创建章节目录
						util.MKdirs(bookid + "/" + epid)
						//下载章节首张图片
						util.DownloadJpg(fImg, imgRole, false)
						if s["imgs"] != "" {
							go util.DoWork(bookid+"/"+epid, s["imgs"], bookid, epid, s["link"], caches)
						}
					}
				} else {
					//修改章节最后更新时间
					oldChapter := models.BookEpisode{Id: cId}
					oldChapter.LastTime = s["ctime"]
					num, _ := c.Update(&oldChapter, "LastTime")
					if num > int64(0) || caches {
						//创建章节目录
						util.MKdirs(bookid + "/" + epid)
						//下载章节首张图片
						util.DownloadJpg(fImg, imgRole, false)
						if s["imgs"] != "" {
							go util.DoWork(bookid+"/"+epid, s["imgs"], bookid, epid, s["link"],caches)
						}
					}
				}
			}
		}
	}
}

//初始化写入redis值
func (p *PickController) SaveRedis() {
	o := orm.NewOrm()
	//链接redis
	redisPool := models.ConnectRedisPool()
	defer redisPool.Close()
	//获取全部链接
	var links []*models.Links
	_, _ = o.QueryTable(new(models.Links).TableName()).All(&links)
	for _, value := range links {
		spew.Dump(value.Id)
		_, err := redisPool.Do("HSET", "book_all_lists", value.BookLink, value.Id)
		if err != nil {
			spew.Dump("redis初始化链接存入错误")
		}
	}
	//获取全部章节图书链接
	var lists []*models.BookList
	_, _ = o.QueryTable(new(models.BookList).TableName()).All(&lists)
	for _, value1 := range lists {
		_, err1 := redisPool.Do("HSET", "comic_links", value1.DomainName, value1.BookId)
		if err1 != nil {
			spew.Dump("初始化漫画ID存入错误")
		}
	}
	//获取全部章节链接
	var epLists []*models.BookEpisode
	_, _ = o.QueryTable(new(models.BookEpisode).TableName()).All(&epLists)
	for _, value2 := range epLists {
		_, err2 := redisPool.Do("HSET", "chapter_links", value2.Link, 1)
		if err2 != nil {
			spew.Dump("初始化存入章节链接到redis错误")
		}
	}

	p.MsgBack("初始化链接进入redis完成", 1)
}