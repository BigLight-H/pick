package models

type BookEpisode struct {
	Id                  int
	BookId              int
	EpisodeId   	    string
	EpisodeTitle        string
	EpisodeThumbnail    string
	EpisodeImgtotal     int
	LastTime   			string
}

func (m *BookEpisode) TableName() string {
	return TableName("book_episode")
}