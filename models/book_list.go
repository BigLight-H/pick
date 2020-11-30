package models

type BookList struct {
	Id                  int
	BookId              int
	BookAuthor   	    string
	BookThumbnail       string
	BookProfile    		string
	BookTags        	int
	BookStat 			string
	IsAgeLimit   		string
	TimesCollect   		int
	TimesBuy   			int
	TimesRead   		int
	TimesSubscribed   	int
	UserBuy   			int
	UpdateChapter   	string
	Year       			string
	Star       			string
	NowStatus       	string
	LastTime   			string
}

func (m *BookList) TableName() string {
	return TableName("book_list")
}