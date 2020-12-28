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
}

func (m *BookEpisode) TableName() string {
	return TableName("episode")
}