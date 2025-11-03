// models/lesson.go
package models

import "time"

type LessonType = string

const (
	LessonTypeFillBlanks LessonType = "fillblanks"
	LessonTypeArrange    LessonType = "arrange"
)

type Lesson struct {
	Id        string       `bson:"_id,omitempty"  json:"lessonId"`
	UserId    string       `bson:"user_id"        json:"-"`
	SongId    string       `bson:"song_id"        json:"-"`
	Items     []LessonItem `bson:"items"          json:"items"`
	CreatedAt time.Time    `bson:"created_at"     json:"-"`
}

type LessonItem struct {
	Type         LessonType `bson:"type"           json:"type"`
	LineIndex    int        `bson:"line_index"     json:"lineIndex"`
	RenderedLine string     `bson:"rendered_line"  json:"renderedLine"`
	Words        []string   `bson:"words"          json:"words"`
	CorrectWord  string     `bson:"correct_word" json:"correct_word"`
}
