package models

type Book struct {
	Id        int
	UniqueId  string
	BookName  int
	Tags      string
	Summary   string
	Domain    string
	End       int
	AuthorName string
	CoverUrl  string
}

func (m *Book) TableName() string {
	return TableName("book")
}