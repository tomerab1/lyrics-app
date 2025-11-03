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
    Answers   []LessonAnswer `bson:"answers"       json:"-"`
	CreatedAt time.Time    `bson:"created_at"     json:"-"`
}

type LessonItem struct {
	Type         LessonType `bson:"type"           json:"type"`
	LineIndex    int        `bson:"line_index"     json:"lineIndex"`
	RenderedLine string     `bson:"rendered_line"  json:"renderedLine"`
	Words        []string   `bson:"words"          json:"words"`
	CorrectWord  string     `bson:"correct_word" json:"correct_word"`
}

type LessonAnswer struct {
    ItemIndex int    `bson:"item_index"`
    Type      string `bson:"type"`          // persist only if fillblanks
    UserInput string `bson:"user_input"`    // chosen word
    Correct   bool   `bson:"correct"`
}
