package models

type Chapter struct {
	Id        int
	ChapterName  string
	BookId    int
	ChapterOrder  string
}

func (m *Chapter) TableName() string {
	return TableName("chapter")
}