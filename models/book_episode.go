package models

type BookEpisode struct {
	Id                  int
	BookId              int
	EpisodeId   	    string
	EpisodeTitle        string
	EpisodeThumbnail    string
	EpisodeImgtotal     int
	LastTime   			string
	Link   				string
	CreateTime   		int64
	UpdateTime   		int64
}

func (m *BookEpisode) TableName() string {
	return TableName("episode")
}