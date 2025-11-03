package models

type Song struct {
	Id     string     `bson:"_id,omitempty" json:"id"`
	Title  string     `bson:"title" json:"title"`
	Artist string     `bson:"artist" json:"artist"`
	Lyrics [][]string `bson:"lyrics" json:"lyrics"`
}
