package utils

import (
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"

	"github.tomerab1/todo-api/internal/contracts"
	"github.tomerab1/todo-api/internal/models"
)

func LyricsToSlices(lyrics string) [][]string {
	// Split the lyrics by lines
	lines := strings.Split(strings.TrimSpace(lyrics), "\n")

	var result [][]string
	for _, line := range lines {
		// Split each line into words
		words := strings.Fields(line)
		result = append(result, words)
	}
	return result
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Flatten(lines [][]string) []string {
	var out []string
	for _, ln := range lines {
		out = append(out, ln...)
	}
	return out
}

func UniqueLower(xs []string) []string {
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

func RenderBlank(words []string, hiddenIdx int) string {
	cp := slices.Clone(words)
	cp[hiddenIdx] = "___"
	return strings.Join(cp, " ")
}

func BuildOptions(r *rand.Rand, correct string, vocab []string) []string {
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
func ToContractItems(items []models.LessonItem) []contracts.LessonItem {
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
func ItemSignature(it models.LessonItem) string {
	if it.Type == models.LessonTypeFillBlanks {
		return "F:" + it.RenderedLine
	}
	// arrange uniqueness by type + line index
	return "A:" + strconv.Itoa(it.LineIndex)
}
