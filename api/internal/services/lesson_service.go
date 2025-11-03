package services

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"slices"
	"strings"
	"time"

	"github.tomerab1/todo-api/internal/contracts"
	"github.tomerab1/todo-api/internal/models"
	"github.tomerab1/todo-api/internal/repositories"
    "strconv"
)

type LessonService struct {
	songRepo   repositories.SongRepoIface
	lessonRepo repositories.LessonRepoIface
	logger     *slog.Logger
}

var ErrDuplicateAnswer = errors.New("duplicate answer")

func NewLessonService(
	songRepo repositories.SongRepoIface,
	lessonRepo repositories.LessonRepoIface,
	logger *slog.Logger,
) *LessonService {
	return &LessonService{
		songRepo:   songRepo,
		lessonRepo: lessonRepo,
		logger:     logger,
	}
}

// CreateLesson generates a 6-item lesson and persists it.
func (svc *LessonService) CreateLesson(
	ctx context.Context,
	dto contracts.CreateLessonDto,
) (*contracts.CreateLessonResponse, error) {
	if strings.TrimSpace(dto.UserId) == "" {
		return nil, errors.New("userId is required")
	}

	// 1) Pick a song (random for now)
	songs, err := svc.songRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(songs) == 0 {
		return nil, errors.New("no songs available")
	}

	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	song := songs[r.IntN(len(songs))]

    // 2) Build vocabulary and candidate line indexes from song.Lyrics ([][]string)
    lines := song.Lyrics
    if len(lines) == 0 {
        return nil, errors.New("chosen song has no lines")
    }

    vocab := uniqueLower(flatten(lines)) // []string of unique, lower-cased words for distractors

    // Prepare distinct candidates for each type
    fillCands := make([]int, 0, len(lines)) // len(words) >= 2
    arrCands := make([]int, 0, len(lines))  // len(words) >= 1
    for i, ln := range lines {
        if len(ln) >= 2 {
            fillCands = append(fillCands, i)
        }
        if len(ln) >= 1 {
            arrCands = append(arrCands, i)
        }
    }
    // Shuffle candidates
    r.Shuffle(len(fillCands), func(i, j int) { fillCands[i], fillCands[j] = fillCands[j], fillCands[i] })
    r.Shuffle(len(arrCands), func(i, j int) { arrCands[i], arrCands[j] = arrCands[j], arrCands[i] })

    // Target: 3 fillblanks + 3 arrange, without repeating the same line
    used := make(map[int]struct{})
    items := make([]models.LessonItem, 0, 6)

    // helper to consume from a candidate list ensuring unique line usage
    take := func(cands *[]int) (int, bool) {
        for len(*cands) > 0 {
            idx := (*cands)[0]
            *cands = (*cands)[1:]
            if _, ok := used[idx]; ok {
                continue
            }
            used[idx] = struct{}{}
            return idx, true
        }
        return 0, false
    }

    // Build fills
    for len(items) < 3 {
        idx, ok := take(&fillCands)
        if !ok {
            break
        }
        words := slices.Clone(lines[idx])
        hidden := r.IntN(len(words))
        correct := words[hidden]
        options := buildOptions(r, correct, vocab)
        rendered := renderBlank(words, hidden)
        items = append(items, models.LessonItem{
            Type:         models.LessonTypeFillBlanks,
            LineIndex:    idx,
            RenderedLine: rendered,
            Words:        options,
            CorrectWord:  correct,
        })
    }

    // Build arrange
    for len(items) < 6 {
        idx, ok := take(&arrCands)
        if !ok {
            break
        }
        words := slices.Clone(lines[idx])
        items = append(items, models.LessonItem{
            Type:      models.LessonTypeArrange,
            LineIndex: idx,
            Words:     words,
        })
    }

    // If we still don't have 6 items (e.g., not enough distinct lines), allow reuse but avoid exact duplicates
    if len(items) < 6 {
        // fallback pool of all indexes
        pool := r.Perm(len(lines))
        seen := make(map[string]struct{})
        for _, it := range items {
            sig := itemSignature(it)
            seen[sig] = struct{}{}
        }
        for _, idx := range pool {
            if len(items) >= 6 {
                break
            }
            words := slices.Clone(lines[idx])
            if len(words) == 0 {
                continue
            }
            // alternate types while creating distinct signatures
            if len(items)%2 == 0 && len(words) >= 2 {
                hidden := r.IntN(len(words))
                correct := words[hidden]
                options := buildOptions(r, correct, vocab)
                rendered := renderBlank(words, hidden)
                cand := models.LessonItem{Type: models.LessonTypeFillBlanks, LineIndex: idx, RenderedLine: rendered, Words: options, CorrectWord: correct}
                if _, ok := seen[itemSignature(cand)]; ok {
                    continue
                }
                seen[itemSignature(cand)] = struct{}{}
                items = append(items, cand)
            } else {
                cand := models.LessonItem{Type: models.LessonTypeArrange, LineIndex: idx, Words: words}
                if _, ok := seen[itemSignature(cand)]; ok {
                    continue
                }
                seen[itemSignature(cand)] = struct{}{}
                items = append(items, cand)
            }
        }
    }

	// 4) Persist lesson
	lesson := &models.Lesson{
		UserId: dto.UserId,
		SongId: song.Id,
		Items:  items,
		Answers: make([]models.LessonAnswer, 0),
	}
	lesson, err = svc.lessonRepo.Create(ctx, dto.UserId, lesson)
	if err != nil {
		return nil, err
	}

	return &contracts.CreateLessonResponse{
		LessonId: lesson.Id,
		Items:    toContractItems(lesson.Items),
	}, nil
}

// SubmitAnswer persists an answer only for fillblanks; returns correctness and 409 on duplicate.
func (svc *LessonService) SubmitAnswer(
    ctx context.Context,
    lessonId string,
    itemIndex int,
    ansType string,
    userInput string,
) (bool, error) {
    // Only fillblanks are persisted; correctness computed against stored lesson item
    lesson, err := svc.lessonRepo.GetById(ctx, lessonId)
    if err != nil {
        return false, err
    }
    // reject duplicate submissions for same item
    for _, a := range lesson.Answers {
        if a.ItemIndex == itemIndex {
            return false, ErrDuplicateAnswer
        }
    }
    if itemIndex < 0 || itemIndex >= len(lesson.Items) {
        return false, errors.New("invalid item index")
    }
    item := lesson.Items[itemIndex]
    if ansType != string(models.LessonTypeFillBlanks) {
        // Ignore persistence for arrange; compute correctness locally if possible
        // For arrange, UI checks correctness itself; we reply ok without persisting
        return true, nil
    }

    correct := strings.EqualFold(userInput, item.CorrectWord)
    // Try to push answer; repo enforces single submission per item
    err = svc.lessonRepo.AddAnswer(ctx, lessonId, models.LessonAnswer{
        ItemIndex: itemIndex,
        Type:      ansType,
        UserInput: userInput,
        Correct:   correct,
    })
    if err != nil {
        return false, err
    }
    return correct, nil
}

func (svc *LessonService) GetSummary(
	ctx context.Context,
	lessonId string,
) (total int, correct int, wrong int, accuracy float64, scheduled []string, err error) {
	lesson, err := svc.lessonRepo.GetById(ctx, lessonId)
	if err != nil {
		return 0, 0, 0, 0, nil, err
	}
	total = len(lesson.Items)
	
	// Count fillblanks answers
	fillblanksCorrect := 0
	fillblanksWrong := 0
	for _, a := range lesson.Answers {
		if a.Type == string(models.LessonTypeFillBlanks) {
			if a.Correct {
				fillblanksCorrect++
			} else {
				fillblanksWrong++
				scheduled = append(scheduled, a.UserInput)
			}
		}
	}
	
	// Count arrange items (they're not persisted, so we count them as correct for accuracy)
	arrangeCount := 0
	for _, item := range lesson.Items {
		if item.Type == models.LessonTypeArrange {
			arrangeCount++
		}
	}
	
	// Total correct = fillblanks correct + arrange items (assumed correct since not tracked)
	correct = fillblanksCorrect + arrangeCount
	// Wrong = only fillblanks wrong
	wrong = fillblanksWrong
	
	// Calculate accuracy based on all items
	if total > 0 {
		accuracy = float64(correct) / float64(total) * 100
	}
	if scheduled == nil {
		scheduled = []string{}
	}
	return
}

// --- helpers ---

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// randDistinct returns k distinct indexes from [0..n)
func randDistinct(r *rand.Rand, n, k int) []int {
	if k > n {
		k = n
	}
	idxs := r.Perm(n)
	return idxs[:k]
}

func flatten(lines [][]string) []string {
	var out []string
	for _, ln := range lines {
		out = append(out, ln...)
	}
	return out
}

func uniqueLower(xs []string) []string {
	seen := make(map[string]struct{}, len(xs))
	out := make([]string, 0, len(xs))
	for _, w := range xs {
		lw := strings.ToLower(strings.TrimSpace(w))
		if lw == "" {
			continue
		}
		if _, ok := seen[lw]; !ok {
			seen[lw] = struct{}{}
			out = append(out, lw)
		}
	}
	return out
}

func renderBlank(words []string, hiddenIdx int) string {
	cp := slices.Clone(words)
	cp[hiddenIdx] = "___"
	return strings.Join(cp, " ")
}

func buildOptions(r *rand.Rand, correct string, vocab []string) []string {
	// Keep correct in its original case; distractors come from lower-case vocab.
	correctLower := strings.ToLower(correct)

	// Build candidate distractors from vocab (exclude the correct word)
	cands := make([]string, 0, len(vocab))
	for _, w := range vocab {
		if w == correctLower {
			continue
		}
		cands = append(cands, w)
	}

	// Shuffle and take up to 3
	r.Shuffle(len(cands), func(i, j int) { cands[i], cands[j] = cands[j], cands[i] })
	d := min(3, len(cands))
	opts := make([]string, 0, 4)
	opts = append(opts, correct)
	for i := 0; i < d; i++ {
		opts = append(opts, cands[i])
	}
	// If we don't have enough distractors, duplicate correct (simple fallback)
	for len(opts) < 4 {
		opts = append(opts, correct)
	}
	// Final shuffle so correct isn't always first
	r.Shuffle(len(opts), func(i, j int) { opts[i], opts[j] = opts[j], opts[i] })
	return opts
}

// map models -> contracts for the UI
func toContractItems(items []models.LessonItem) []contracts.LessonItem {
	out := make([]contracts.LessonItem, 0, len(items))
	for _, it := range items {
        ci := contracts.LessonItem{
            Type:         it.Type,
            LineIndex:    it.LineIndex,
            RenderedLine: it.RenderedLine,
            Words:        slices.Clone(it.Words),
        }
        if it.Type == models.LessonTypeFillBlanks {
            ci.CorrectWord = it.CorrectWord
        }
        out = append(out, ci)
	}
	return out
}

// create a uniqueness signature for a lesson item to avoid exact duplicates
func itemSignature(it models.LessonItem) string {
    if it.Type == models.LessonTypeFillBlanks {
        return "F:" + it.RenderedLine
    }
    // arrange uniqueness by type + line index
    return "A:" + strconv.Itoa(it.LineIndex)
}
