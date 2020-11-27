package models

type Photo struct {
	Id        int
	ChapterId int
	PicOrder  string
	ImgUrl  string
	BookId  int
}

func (m *Photo) TableName() string {
	return TableName("photo")
}