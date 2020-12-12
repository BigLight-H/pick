package util

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"github.com/gocolly/colly"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"pick/conf"
	"pick/models"
	"pick/service"
	"strconv"
	"strings"
	"time"
)

//建立图书目录
func ThisMkdir(domin string) string {
	pathArr := strings.SplitAfter(domin, "/")
	ps := pathArr[len(pathArr)-1]
	path := ""
	if ps == "" {
		ps = pathArr[len(pathArr)-2]
		path = strings.TrimRight(ps, "/")
	} else {
		path = ps
	}
	MKdirs(path)
	return path
}

//获取父子级目录名
func GetSubdirectory(domin string) (string, string)  {
	imgArr := strings.Split(domin, "/")
	fater := imgArr[len(imgArr)-3]
	son := imgArr[len(imgArr)-2]
	return fater, son
}


func ChapterListOrder(chapterName string, str string, k int) string  {
	chapterArr := strings.Split(chapterName, str)
	var orderId string
	orderId = chapterArr[len(chapterArr)-k]
	return orderId
}

//获取章节排序
func ChapterOrder(chapterName string, str string, k int) string  {
	chapterArr := strings.Split(chapterName, str)
	var orderId string
	if len(chapterArr) > 2 {
		orderId = chapterArr[len(chapterArr)-(k+1)]+"-"+chapterArr[len(chapterArr)-k]
	} else {
		orderId = chapterArr[len(chapterArr)-k]
	}
	return orderId
}

//计算图片张数
func SumImgs(imgs string) (int, string)  {
	chapterArr := strings.Split(imgs, ",")
	return len(chapterArr), chapterArr[0]
}

//建立指定目录
func MKdirs(path string) {
	ph := beego.AppConfig.String("comic_hub") + path
	if _, err := os.Stat(ph); os.IsNotExist(err) {
		os.Mkdir(ph, os.ModePerm)
	}
}

func DownloadJpg(url string, file_name string)  {
	transport := &http.Transport{IdleConnTimeout: 0}
	client := &http.Client{Transport: transport}
	if url != "" {
		req,err := http.NewRequest("GET",url,nil)
		if err != nil{
			fmt.Println(err)
			return
		}

		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.2222")
		req.Header.Add("Referer","https://www.webtoon.xyz/")

		req.Close = true

		resp, er := client.Do(req)

		if er != nil {
			fmt.Println(er)
			defer resp.Body.Close()
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Println(resp.StatusCode)
			return
		}
		byteCotent, e := ioutil.ReadAll(resp.Body)
		HandError(e)

		ioutil.WriteFile(beego.AppConfig.String("comic_hub") + file_name, byteCotent, 0777)
		spew.Dump("已下载图片链接"+url)
	}
}

func HandError(err error)  {
	if err != nil{
		fmt.Println("error",err)
	}
}

func DoWork(dir string, imgs string, bid string, eid string, link string) {
	imgArr := strings.Split(imgs, ",")
	//删除第最后一个元素
	if len(imgArr) > 0 {
		imgArr = imgArr[:len(imgArr)-1]
		for _, value := range imgArr{
			imgAr := strings.Split(value, "/")
			name := imgAr[len(imgAr)-1]
			go DownloadJpg(value, dir+"/"+bid+"0"+eid+name)
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

//数组平分
func SplitArray(arr []string, num int) ([][]string) {
	max := int(len(arr))
	if max < num {
		return nil
	}
	var segmens =make([][]string,0)
	quantity:=max/num
	end:=int(0)
	for i := int(1); i <= num; i++ {
		qu:=i*quantity
		if i != num {
			segmens= append(segmens,arr[i-1+end:qu])
		}else {
			segmens= append(segmens,arr[i-1+end:])
		}
		end=qu-i
	}
	return segmens
}

func BookLists(domain string) {
	//图片信息
	c := colly.NewCollector()

	// Find and visit all links
	c.OnXML("//body/div[1]/div[1]/div[2]/div[2]/div[1]/div[1]/div[1]/div[1]/div[2]/div[1]/div[2]/div[2]/div[1]/div[2]/div[1]/a[@class='last']", func(e *colly.XMLElement) {
		lastLink := e.ChildText("//@href")
		allPage := ChapterListOrder(lastLink, "/", 2)
		allNum, _ := strconv.Atoi(allPage)
		for i := 1; i <= allNum; i++ {
			//获取分页数据并存入数据库
			go GetLinks(domain + "page/" + strconv.Itoa(i) + "/")
		}
	})
	c.Visit(domain)
}


func GetLinks(pageDomain string) {
	d := colly.NewCollector()
	//var wg sync.WaitGroup
	d.OnXML("//body/div[1]/div[1]/div[2]/div[2]/div[1]/div[1]/div[1]/div[1]/div[2]/div[1]/div[2]/div[2]/div[1]/div[1]/div", func(f *colly.XMLElement) {
		for a := 1; a <= 2; a++ {
			//链接redis
			redisPool := models.ConnectRedisPool()
			defer redisPool.Close()
			link := f.ChildText("//div[1]/div["+strconv.Itoa(a)+"]/div[1]/div[1]/a/@href")
			title := f.ChildText("//div[1]/div["+strconv.Itoa(a)+"]/div[1]/div[1]/a/@title")
			lastCharpter := f.ChildText("//div[1]/div["+strconv.Itoa(a)+"]/div[1]/div[2]/div[3]/div[1]/span[1]")
			t1 := f.ChildText("//div[1]/div["+strconv.Itoa(a)+"]/div[1]/div[1]/a[1]/span[1]")
			t2 := f.ChildText("//div[1]/div["+strconv.Itoa(a)+"]/div[1]/div[1]/a[1]/span[2]")
			isExist, _ := redisPool.Do("HEXISTS", "book_all_lists", link)
			spew.Dump(link)
			//创建协程
			//wg.Add(1)
			if isExist != int64(1) {
				o := orm.NewOrm()
				lists := models.Links{}
				lists.BookLink = link
				lists.BookName = title
				lists.LastChapter = lastCharpter
				lists.Status = 0
				lists.Type = t1+","+t2
				lid, err := o.Insert(&lists)
				if err == nil {
					_, err := redisPool.Do("HSET", "book_all_lists", link, lid)
					if err != nil {
						spew.Dump("漫画链接存入错误")
					}
					go ComicsCopy(link, 1)
					//wg.Done()
				}
			} else {
				lid, err := redisPool.Do("HGET", "book_all_lists", link)
				lId, _ := redis.Int(lid, err)
				//查看更新
				o := orm.NewOrm()
				lists := models.Links{Id:lId}
				lists.LastChapter = lastCharpter
				lists.Type = t1+","+t2
				if num, err := o.Update(&lists, "LastChapter", "Type"); err == nil {
					//有更新改变字段值
					if num > 0 {
						u := orm.NewOrm()
						up := models.Links{Id:lId}
						up.Status = 1
						if num, err := u.Update(&lists, "Status"); err == nil {
							spew.Dump(num)
						}
						go ComicsCopy(link, 1)
						//wg.Done()
					}
				}
			}

		}
		//os.Exit(2)
	})
	d.Visit(pageDomain)
	//wg.Wait()
}

// 随机数字串
func RandomNum(length int) string {
	result := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		num := rand.Intn(10)
		result = result + strconv.Itoa(num)
	}
	return result
}

//操作图书
func ComicsCopy(domin string, rootId int) {
	//获取指定规则
	role := conf.Choose(rootId)
	//爬取图书
	bookInfo, chapterInfo := service.BookInfo(role, domin)
	//图书ID
	bId := 0
	o := orm.NewOrm()
	for _, v := range bookInfo {
		//链接redis
		redisPool := models.ConnectRedisPool()
		defer redisPool.Close()
		isExist, _ := redisPool.Do("HEXISTS", "comic_links", domin)
		if isExist != int64(1) {
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
			book.TimesCollect = RandomNum(4)
			book.TimesBuy = RandomNum(4)
			book.TimesRead = RandomNum(4)
			book.TimesSubscribed = RandomNum(4)
			book.UserBuy = RandomNum(4)
			book.BookThumbnail = ""
			book.IsAgeLimit = 1
			id, _ := o.Insert(&book)
			bId = int(id)
			if id > int64(0) {
				up := models.BookList{BookId: bId}
				bookid := strconv.Itoa(bId)
				up.BookThumbnail = bookid+"/"+bookid+"_thumb.jpg"
				num, _ := o.Update(&up, "BookThumbnail")
				if num > int64(0) {
					//新建文件目录
					MKdirs(bookid)
					//下载封面图片到目录
					DownloadJpg(v["image"], bookid+"/"+bookid+"_thumb.jpg")
				}
			}
			_, err := redisPool.Do("HSET", "comic_links", domin, bId)
			if err != nil {
				spew.Dump("漫画ID存入错误")
			}
		} else {
			comicId, err := redisPool.Do("HGET", "comic_links", domin)
			bookId, _ := redis.Int(comicId, err)
			bId = bookId
			//修改书本最新更新时间

			oldBook := models.BookList{BookId: bId}
			oldBook.LastTime = v["ltime"]
			if num, err := o.Update(&oldBook, "LastTime"); err != nil {
				fmt.Println(num)
			}
		}
	}
	bookid := strconv.Itoa(bId)
	for _, s := range chapterInfo {
		//切割链接
		_, son := GetSubdirectory(s["link"])
		epid := ChapterOrder(son, "-", 1)
		//存入章节表
		c := orm.NewOrm()
		chapter := models.BookEpisode{}
		chapter.BookId = bId
		chapter.EpisodeTitle = s["title"]
		chapter.EpisodeId = epid
		imgLen, fImg := SumImgs(s["imgs"])
		chapter.EpisodeImgtotal = imgLen
		chapter.EpisodeThumbnail = ""
		chapter.LastTime = s["ctime"]
		chapter.Link = s["link"]
		//cId, _ := c.Insert(&chapter)
		cId := 0

		// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
		if created, cid, err := c.ReadOrCreate(&chapter, "Link"); err == nil {
			if created {
				cId = int(cid)
				//图片路径
				imgRole := bookid +"/"+ epid +"/"+ bookid +"0"+ epid +"_thumb.jpg"
				upimg := models.BookEpisode{Id: cId}
				upimg.EpisodeThumbnail = bookid +"/"+ epid +"/"+ bookid +"0"+ epid +"_thumb.jpg"
				num, _ := c.Update(&upimg, "EpisodeThumbnail")
				if num > int64(0) {
					//创建章节目录
					MKdirs(bookid + "/" + epid)
					//下载章节首张图片
					DownloadJpg(fImg, imgRole)

				}
			} else {
				cId = int(cid)
				//修改章节最后更新时间
				oldChapter := models.BookEpisode{Id: cId}
				oldChapter.LastTime = s["ctime"]
				num, _ := c.Update(&oldChapter, "LastTime")
				if num > int64(0) {
					spew.Dump(num)
				}
			}

			//协程下载图片
			//if cId > 0 {
				if s["imgs"] != "" {
					go DoWork(bookid+"/"+epid, s["imgs"], bookid, epid, s["link"])
				}
			//}
		}


	}

}