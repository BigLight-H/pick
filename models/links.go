package models

type Links struct {
	Id        int
	BookLink  string
	BookName  string
	LastChapter  string
	Status    int
	Type      string
}

func (m *Links) TableName() string {
	return TableName("links")
}