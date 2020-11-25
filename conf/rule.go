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
}

func Choose(id int) *MainRule {
	if id == 1 {
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
			Title :  "//a[1]/@href",
			Link : "//a[1]/text()",
			Detail : "//body/div[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[2]/div[1]/div[1]/div[2]/div[position()<last()]",
			ImgSrc : "//img/@data-src",
		}
	}
	return &MainRule{}
}

