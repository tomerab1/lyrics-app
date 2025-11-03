package utils

import "strings"

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
