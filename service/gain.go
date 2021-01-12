package service

import (
	"fmt"
	"github.com/beego/beego/client/orm"
	"github.com/beego/beego/core/config"
	"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"os"
	"pick/conf"
	"pick/models"
	"pick/util"
	"strconv"
	"strings"
	"time"
)

/**
 * 获取源二总页数
 */
func BookTwoLists(domain string, rid int) {
	//图片信息
	c := colly.NewCollector()

	// Find and visit all links
	c.OnXML("//body/div[6]/div[1]/div[2]/div[1]/div[6]/ul[1]/li[7]", func(e *colly.XMLElement) {
		lastLink := e.ChildText("//a")
		allNum, _ := strconv.Atoi(lastLink)
		for i := 1; i <= allNum; i++ {
			//获取分页数据并存入数据库
			go GetLinks(domain + "page=" + strconv.Itoa(i) + "", rid)
		}
	})
	c.Visit(domain)
}

/**
 * 获取源一总页数
 */
func BookLists(domain string, rid int) {
	//图片信息
	c := colly.NewCollector()

	// Find and visit all links
	c.OnXML("//body/div[1]/div[1]/div[2]/div[2]/div[1]/div[1]/div[1]/div[1]/div[2]/div[1]/div[2]/div[2]/div[1]/div[2]/div[1]/a[@class='last']", func(e *colly.XMLElement) {
		lastLink := e.ChildText("//@href")
		allPage := util.ChapterListOrder(lastLink, "/", 2)
		allNum, _ := strconv.Atoi(allPage)
		for i := 1; i <= allNum; i++ {
			//获取分页数据并存入数据库
			go GetLinks(domain + "page/" + strconv.Itoa(i) + "/", rid)
		}
	})
	c.Visit(domain)
}

/**
 * 获取源一链接
 */
func GetLinks(pageDomain string, rid int) {
	d := colly.NewCollector()
	//var wg sync.WaitGroup
	switch rid {
		case 1:
			onePickLinks(d)
			break
		case 2:
			twoPickLinks(d)
			break
	}

	d.Visit(pageDomain)
	//wg.Wait()
}

func onePickLinks(d *colly.Collector) {
	d.OnXML("//body/div[1]/div[1]/div[2]/div[2]/div[1]/div[1]/div[1]/div[1]/div[2]/div[1]/div[2]/div[2]/div[1]/div[1]/div", func(f *colly.XMLElement) {
		for a := 1; a <= 2; a++ {
			//链接redis
			redisPool := models.ConnectRedisPool()
			defer redisPool.Close()
			link := f.ChildText("//div[1]/div[" + strconv.Itoa(a) + "]/div[1]/div[1]/a/@href")
			title := f.ChildText("//div[1]/div[" + strconv.Itoa(a) + "]/div[1]/div[1]/a/@title")
			lastCharpter := f.ChildText("//div[1]/div[" + strconv.Itoa(a) + "]/div[1]/div[2]/div[3]/div[1]/span[1]")
			t1 := f.ChildText("//div[1]/div[" + strconv.Itoa(a) + "]/div[1]/div[1]/a[1]/span[1]")
			t2 := f.ChildText("//div[1]/div[" + strconv.Itoa(a) + "]/div[1]/div[1]/a[1]/span[2]")
			isExist, _ := redisPool.Do("HEXISTS", "book_all_lists", link)
			spew.Dump(link)
			//创建协程
			if isExist == int64(0) {
				if link != "" {
					o := orm.NewOrm()
					lists := models.Links{}
					lists.BookLink = link
					lists.BookName = title
					lists.LastChapter = lastCharpter
					lists.Status = 0
					lists.Type = t1 + "," + t2
					lid, err := o.Insert(&lists)
					if err == nil {
						_, err2 := redisPool.Do("HSET", "book_all_lists", link, lid)
						if err2 != nil {
							spew.Dump("漫画链接存入错误")
						}
						go ComicsCopy(link, 1)
					}
				}
			} else {
				lid, err := redisPool.Do("HGET", "book_all_lists", link)
				lId, _ := redis.Int(lid, err)
				//查看更新
				o := orm.NewOrm()
				lists := models.Links{Id: lId}
				lists.LastChapter = lastCharpter
				lists.Type = t1 + "," + t2
				if num1, err1 := o.Update(&lists, "LastChapter", "Type"); err1 == nil {
					//有更新改变字段值
					if num1 == int64(1) {
						lists.Status = 1
						if num2, err2 := o.Update(&lists, "Status"); err2 != nil {
							spew.Dump(num2)
						}
					}
					go ComicsCopy(link, 1)
				}
			}

		}
		//os.Exit(2)
	})
}

func twoPickLinks(d *colly.Collector) {
	d.OnXML("//body/div[6]/div[1]/div[2]/div[1]/div[5]/div", func(f *colly.XMLElement) {
		link := f.ChildText("//div[1]/a[1]/@href")
		title := f.ChildText("//div[1]/div[1]/h3[1]/a/text()")
		lastCharpter := f.ChildText("//div[1]/div[1]/a")
		t1 := f.ChildText("/div[1]/div[1]/small[1]/a[1]/text()")
		t2 := ""
		if link != "" {
			o := orm.NewOrm()
			lists := models.Links{}
			lists.BookLink = link
			lists.BookName = title
			lists.LastChapter = lastCharpter
			lists.Status = 0
			lists.Type = t1 + "," + t2
			lid, err := o.Insert(&lists)
			if err == nil {
				spew.Dump(lid)
				os.Exit(1)
				//_, err2 := redisPool.Do("HSET", "book_all_lists", link, lid)
				//if err2 != nil {
				//	spew.Dump("漫画链接存入错误")
				//}
				//go ComicsCopy(link, 1)
			}
		}

	})
}

//操作图书
func ComicsCopy(domin string, rootId int) {
	//获取指定规则
	role := conf.Choose(rootId)
	//爬取图书
	bookInfo, chapterInfo := BookInfo(role, domin, false)
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
			book.BookProfile = util.SubString(v["intro"], 0, 300)
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
					DownloadJpg(v["image"], bookid+"/"+bookid+"_thumb.jpg", false)
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
			chapter.CreateTime = time.Now().Unix()
			chapter.UpdateTime = time.Now().Unix()
			//cId, _ := c.Insert(&chapter)
			cId := 0

			// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
			if created, cid, err := c.ReadOrCreate(&chapter, "Link"); err == nil {
				cId = int(cid)
				if created {
					if cId > 0 {
						//图片路径
						imgRole := bookid +"/"+ epid +"/"+ bookid +"0"+ epid +"_thumb.jpg"
						upimg := models.BookEpisode{Id: cId}
						upimg.EpisodeThumbnail = bookid +"/"+ epid +"/"+ bookid +"0"+ epid +"_thumb.jpg"
						num, _ := c.Update(&upimg, "EpisodeThumbnail")
						if num > int64(0) {
							//创建章节目录
							util.MKdirs(bookid + "/" + epid)
							//下载章节首张图片
							DownloadJpg(fImg, imgRole, false)
							//下载图片
							if s["imgs"] != "" {
								go DoWork(bookid+"/"+epid, s["imgs"], bookid, epid, s["link"], false)
							}
						}
					}
				} else {
					if cId > 0 {
						//修改章节最后更新时间
						oldChapter := models.BookEpisode{Id: cId}
						oldChapter.LastTime = s["ctime"]
						num, _ := c.Update(&oldChapter, "LastTime")
						if num > int64(0) {
							//下载图片
							if s["imgs"] != "" {
								go DoWork(bookid+"/"+epid, s["imgs"], bookid, epid, s["link"], false)
							}
						}
					}
				}
			}
		}
	}
}

func DownloadJpg(url string, file_name string, caches bool)  {
	comicHub, _ := config.String("comic_hub")
	bs := util.Exists(comicHub + file_name)
	if bs && !caches {
		return
	}
	transport := &http.Transport{IdleConnTimeout: 0}
	client := &http.Client{Transport: transport}
	if url != "" {
		req,err := http.NewRequest("GET",url,nil)
		if err != nil{
			fmt.Println(err)
			return
		}

		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
		req.Header.Add("Referer","https://www.webtoon.xyz/")

		req.Close = true

		resp, er := client.Do(req)

		if er != nil {
			fmt.Println(er)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Println(resp.StatusCode)
			return
		}
		byteCotent, e := ioutil.ReadAll(resp.Body)
		util.HandError(e)

		_ = ioutil.WriteFile(comicHub+file_name, byteCotent, 0777)
		spew.Dump("已下载图片链接"+url)
	}
}

func DoWork(dir string, imgs string, bid string, eid string, link string, caches bool) {
	imgArr := strings.Split(imgs, ",")
	//删除第最后一个元素
	if len(imgArr) > 0 {
		imgArr = imgArr[:len(imgArr)-1]
		for _, value := range imgArr{
			imgAr := strings.Split(value, "/")
			name := imgAr[len(imgAr)-1]
			go DownloadJpg(value, dir+"/"+bid+"0"+eid+name, caches)
		}
		//存储章节爬取记录
		redisPool := models.ConnectRedisPool()
		defer redisPool.Close()
		_, err := redisPool.Do("HSET", "chapter_links", link, 1)
		if err != nil {
			spew.Dump("存入章节链接到redis错误")
		}
	}
}