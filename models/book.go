package models

type Book struct {
	Id        int
	UniqueId  string
	BookName  string
	Tags      string
	Summary   string
	End       int
	AuthorName string
	CoverUrl  string
	Year  string
	Star  string
	Type  string
	LastTime  string
}

func (m *Book) TableName() string {
	return TableName("book")
}