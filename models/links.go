package models

type Links struct {
	Id        int
	BookLink  string
	BookName  string
	LastChapter  string
	Status    int
}

func (m *Links) TableName() string {
	return TableName("links")
}