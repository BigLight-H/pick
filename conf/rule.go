package conf

//图书信息
type Books struct {
	Name string
	Year string
	Star string
	Author string
	Tags string
	Status string
	Intro string
	Types string
}
//图书章节链接
type Chapter struct {
	Name string
	Title string
	Link string
}

//网站源规则
type MainRule struct {
	Table  string //列表规则
	Body   string //页面body
	Name   string //图书名规则
	Year   string //年份规则
	Star   string //评分规则
	Author string //作者
	Tags   string //标签
	Status string //图书状态
	Intro  string //简介
	Types  string //类型
	Title  string //章节标题
	Link   string //章节链接
	Image  string //图片链接
	Detail string //图片详情页规则
	ImgSrc string //图片链接
	LTime  string //最后更新时间
	CTime  string //章节更新时间
	NTime  string //最新跟新时间
	NCTime string //章节最新跟新时间
}

//源一章节XPATH规则
func Choose(id int) *MainRule {
	switch id {
		case 1:
			return &MainRule{
				Table:   "//body/div[1]/div[1]/div[2]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[4]/div[1]/ul[1]/li",
				Body :   "//body",
				Name :   "//div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/ol[1]/li[3]",
				Year :   "//div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[3]/div[2]/div[1]/div[2]/div[1]/div[2]/a[1]",
				Star :   "//div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[3]/div[2]/div[1]/div[1]/div[2]/div[1]/span[1]",
				Author : "//div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[3]/div[2]/div[1]/div[1]/div[6]/div[2]/div[1]/a[1]",
				Tags :   "//div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[3]/div[2]/div[1]/div[1]/div[8]/div[2]/div[1]",
				Status : "//div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[3]/div[2]/div[1]/div[2]/div[2]/div[2]",
				Intro :  "//div[1]/div[1]/div[2]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[2]/div[1]/p[2]",
				Types :  "//div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[2]/h1[1]/span[1]",
				Image :  "//div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[3]/div[1]/a[1]/img[1]/@data-src",
				Title :  "//a[1]/text()",
				Link : "//a[1]/@href",
				Detail : "//body/div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[2]/div[1]/div[1]/div[2]/div[position()<last()]",
				ImgSrc : "//img/@data-src",
				LTime : "//div[1]/div[1]/div[2]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[4]/div[1]/ul[1]/li[1]/span",
				CTime : "//span",
				NTime : "//div[1]/div[1]/div[2]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[4]/div[1]/ul[1]/li[1]/span[1]/span[1]/a[1]/img[1]/@alt",
				NCTime : "//span[1]/span[1]/a[1]/img[1]/@alt",
			}
		case 2:
			return &MainRule{
				Table:   "//body/div[6]/div[1]/div[1]/div[3]/div[1]/p",
				Body :   "//body",
				Name :   "./div[6]/div[1]/ol[1]/li[3]/a[1]/span",
				Year :   "./div[2]/ul[1]/i",
				Star :   "//div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[3]/div[2]/div[1]/div[1]/div[2]/div[1]/span[1]",
				Author : "./div[2]/ul[1]/li[2]/small[1]",
				Tags :   "./div[2]/ul[1]/li[3]/small[1]",
				Status : "./div[2]/ul[1]/li[4]/a[1]",
				Intro :  " ./div[6]/div[1]/div[1]/div[2]/div[1]/div[2]/text()",
				Types :  "//div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[2]/h1[1]/span[1]",
				Image :  "./div[6]/div[1]/ol[1]/li[3]/a[1]/img/@src",
				Title :  "./span[1]/a/@href",
				Link  :  "./span[1]/a/@title",
				Detail : "//article[@id='content']/img",
				ImgSrc : "./@data-original",
				LTime : "./div[6]/div[1]/div[1]/div[3]/div[1]/p[1]/span[2]/i[1]/time[1]",
				CTime : "./span[2]/i[1]/time",
				NTime : "./div[6]/div[1]/div[1]/div[3]/div[1]/p[2]/span[2]/i[1]/time[1]",
				NCTime : "//span[1]/span[1]/a[1]/img[1]/@alt",
			}
		case 3:
			break
	}
	return &MainRule{}
}


