package models

type Photo struct {
	Id        int
	ChapterId int
	PicOrder  string
	ImgUrl  string
}

func (m *Photo) TableName() string {
	return TableName("photo")
}