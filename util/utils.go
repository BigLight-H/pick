package util

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"os"
	"pick/models"
	"strings"
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
	orderId := chapterArr[len(chapterArr)-k]
	return orderId
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

func DoWork(dir string, imgs string, bid int, cid int) {
	imgArr := strings.Split(imgs, ",")
	//删除第最后一个元素
	if len(imgArr) > 0 {
		imgArr = imgArr[:len(imgArr)-1]
		for _, value := range imgArr{
			imgAr := strings.Split(value, "/")
			name := imgAr[len(imgAr)-1]
			//存入图片表
			c := orm.NewOrm()
			photo := models.Photo{}
			photo.ChapterId = cid
			photo.BookId = bid
			photo.PicOrder = ChapterOrder(name,".",2)
			photo.ImgUrl = value
			_, err := c.Insert(&photo)
			if err != nil {
				os.Exit(3)
			}
			//创建协程处理->获取图片并存储
			DownloadJpg(value, dir+"\\"+name)
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




