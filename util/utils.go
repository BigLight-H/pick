package util

import (
	"fmt"
	"github.com/beego/beego/core/config"
	"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"math/rand"
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
	ph, _ := config.String("comic_hub")
	path_ := ph + path
	if _, err := os.Stat(path_); os.IsNotExist(err) {
		_ = os.Mkdir(path_, os.ModePerm)
	}
}

//判断文件是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if os.IsNotExist(err) {
		return false
	}
	return true
}


func HandError(err error)  {
	if err != nil{
		fmt.Println("error",err)
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

//截取指定长度的字符串
func SubString(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start : end])
}

//增加进程数
func AddProcessNum() {
	redisPool := models.ConnectRedisPool()
	defer redisPool.Close()
	_, err2 := redisPool.Do("incr", "process-nums", 1)
	if err2 != nil {
		spew.Dump("增加进程数值错误")
	}
	val, err := redisPool.Do("get", "process-nums")
	num, _ := redis.Int(val, err)
	if num < 0 {
		//重置
		_, _ = redisPool.Do("set", "process-nums", 1)
	}
}

//减去进程数
func MinusProcessNum() {
	redisPool := models.ConnectRedisPool()
	defer redisPool.Close()
	_, err2 := redisPool.Do("incr", "process-nums", 1)
	if err2 != nil {
		spew.Dump("增加完成数值错误")
	}
}

func GetKeys(domain string, number int) int {
	RootLinks := make(map[int]string)
	for i := 1; i <= number; i++ {
		val, _ := config.String("source_"+strconv.Itoa(i))
		RootLinks[i] = val
	}
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率较高
	j := 1
	keys := 0
	for k := range RootLinks {
		if strings.Compare(domain, RootLinks[j]) == 0 {
			keys = k
			continue
		}
		j++
	}
	return keys
}

func StrToInt(str string) int {
	number, _ := strconv.Atoi(str)
	return number
}