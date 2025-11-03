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
	Lyrics string `json:"lyrics"`
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
	CorrectWord  string   `json:"correct_word"`
}

type CreateLessonResponse struct {
	LessonId string       `json:"lessonId"`
	Items    []LessonItem `json:"items"`
}

// --- Answers & Summary ---

type SubmitAnswerDto struct {
	LessonId  string `json:"lessonId"`
	ItemIndex int    `json:"itemIndex"`
	Type      string `json:"type"`      // "fillblanks" | "arrange"
	Correct   *bool  `json:"correct"`   // optional, ignored for persistence rules
	UserInput string `json:"userInput"` // the chosen word for fillblanks
}

type SubmitAnswerResponse struct {
	Ok      bool `json:"ok"`
	Correct bool `json:"correct"`
}

type LessonSummaryResponse struct {
	Total                  int      `json:"total"`
	Correct                int      `json:"correct"`
	Wrong                  int      `json:"wrong"`
	Accuracy               float64  `json:"accuracy"`
	ScheduledForRepractice []string `json:"scheduledForRepractice"`
}
