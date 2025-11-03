package contracts

type CreateUserDto struct {
	Name string `json:"name"`
}

type GetUserResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CreateSongDto struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Lyrics string `json:"Lyrics"`
}

type GetSongResponse struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type CreateSongsReponse struct {
	Id        string `json:"id"`
	LineCount int    `json:"line_count"`
}

type CreateLessonDto struct {
	UserId string `json:"userId"`
}

type LessonItem struct {
	Type         string   `json:"type"` // "fillblanks" | "arrange"
	LineIndex    int      `json:"lineIndex"`
	RenderedLine string   `json:"renderedLine"` // only for fillblanks
	Words        []string `json:"words"`        // 4 options for fillblanks; CORRECT ORDER for arrange
}

type CreateLessonResponse struct {
	LessonId string       `json:"lessonId"`
	Items    []LessonItem `json:"items"`
}
