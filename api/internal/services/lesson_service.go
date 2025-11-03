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
)

type LessonService struct {
	songRepo   repositories.SongRepoIface
	lessonRepo repositories.LessonRepoIface
	logger     *slog.Logger
}

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

	// 2) Build vocabulary and line indexes from song.Lyrics ([][]string)
	lines := song.Lyrics
	if len(lines) == 0 {
		return nil, errors.New("chosen song has no lines")
	}

	vocab := uniqueLower(flatten(lines)) // []string of unique, lower-cased words for distractors
	indexes := randDistinct(r, len(lines), min(6, len(lines)))

	// 3) Build 6 lesson items (mix of types, fix at creation time)
	items := make([]models.LessonItem, 0, 6)
	for i := 0; len(items) < 6; i++ {
		idx := indexes[i%len(indexes)]
		words := slices.Clone(lines[idx])

		// Skip empty lines
		if len(words) == 0 {
			continue
		}

		// Alternate types to include both kinds
		if len(items)%2 == 0 && len(words) >= 2 {
			// Fill Blanks
			hidden := r.IntN(len(words))
			correct := words[hidden]

			options := buildOptions(r, correct, vocab)
			rendered := renderBlank(words, hidden)

			items = append(items, models.LessonItem{
				Type:         models.LessonTypeFillBlanks,
				LineIndex:    idx,
				RenderedLine: rendered,
				Words:        options, // 4 options: 1 correct + 3 distractors
				CorrectWord:  correct,
			})
		} else {
			// Arrange (UI will shuffle; we send correct order)
			items = append(items, models.LessonItem{
				Type:      models.LessonTypeArrange,
				LineIndex: idx,
				Words:     words,
			})
		}
	}

	// 4) Persist lesson
	lesson := &models.Lesson{
		UserId: dto.UserId,
		SongId: song.Id,
		Items:  items,
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
		out = append(out, contracts.LessonItem{
			Type:         it.Type,
			LineIndex:    it.LineIndex,
			RenderedLine: it.RenderedLine,
			Words:        slices.Clone(it.Words),
		})
	}
	return out
}
