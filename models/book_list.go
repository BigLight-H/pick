package models

type BookList struct {
	BookId              int    `orm:"pk;auto"`  //设置主键自增长 字段名为 id
	BookAuthor   	    string
	BookThumbnail       string
	BookTitle       	string
	BookProfile    		string
	BookTags        	string
	BookStat 			int
	IsAgeLimit   		int
	TimesCollect   		string
	TimesBuy   			string
	TimesRead   		string
	TimesSubscribed   	string
	UserBuy   			string
	UserRead   			string
	Year       			string
	Star       			string
	NowStatus       	string
	LastTime   			string
	DomainName   		string
	Status   		    int
	CreateTime   		string
	UpdateTime   		string
}

func (m *BookList) TableName() string {
	return TableName("list")
}