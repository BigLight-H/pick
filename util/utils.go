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
	"pick/models"
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

func DownloadJpg(url string,file_name string)  {
	client := &http.Client{}

	req,err := http.NewRequest("GET",url,nil)
	if err != nil{
		fmt.Println(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.2222")
	req.Header.Add("Referer","https://www.webtoon.xyz/")
	resp,err := client.Do(req)
	// 先判断是否有错误
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	byteCotent, err := ioutil.ReadAll(resp.Body)
	HandError(err)

	ioutil.WriteFile(beego.AppConfig.String("comic_hub") + file_name, byteCotent, 0777)

}

func HandError(err error)  {
	if err != nil{
		fmt.Println("error",err)
	}
}

func DoWork(dir string, imgs string, bid string, eid string) {
	imgArr := strings.Split(imgs, ",")
	//删除第最后一个元素
	if len(imgArr) > 0 {
		imgArr = imgArr[:len(imgArr)-1]
		for _, value := range imgArr{
			imgAr := strings.Split(value, "/")
			name := imgAr[len(imgAr)-1]
			////存入图片表
			//c := orm.NewOrm()
			//photo := models.Photo{}
			//photo.ChapterId = cid
			//photo.BookId = bid
			//photo.PicOrder = ChapterOrder(name,".",2)
			//photo.ImgUrl = value
			//_, err := c.Insert(&photo)
			//if err != nil {
			//	os.Exit(3)
			//}
			//创建协程处理->获取图片并存储
			DownloadJpg(value, dir+"\\"+bid+"0"+eid+name)
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
		allPage := ChapterOrder(lastLink, "/", 2)
		allNum, _ := strconv.Atoi(allPage)
		for i := 1; i <= allNum; i++ {
			//获取分页数据并存入数据库
			pageDomain := domain + "page/" + strconv.Itoa(i) + "/"
			GetLinks(pageDomain)
		}
	})
	c.Visit(domain)
}


func GetLinks(pageDomain string) {
	d := colly.NewCollector()
	d.OnXML("//body/div[1]/div[1]/div[2]/div[2]/div[1]/div[1]/div[1]/div[1]/div[2]/div[1]/div[2]/div[2]/div[1]/div[1]/div", func(f *colly.XMLElement) {
		for a := 1; a <= 2; a++ {
			//链接redis
			redisPool := models.ConnectRedis()
			defer redisPool.Close()
			link := f.ChildText("//div[1]/div["+strconv.Itoa(a)+"]/div[1]/div[1]/a/@href")
			title := f.ChildText("//div[1]/div["+strconv.Itoa(a)+"]/div[1]/div[1]/a/@title")
			lastCharpter := f.ChildText("//div[1]/div["+strconv.Itoa(a)+"]/div[1]/div[2]/div[3]/div[1]/span[1]")

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
		//os.Exit(2)
	})
	d.Visit(pageDomain)
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
