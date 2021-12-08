package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func Top10(text string) []string {
	wordCounters := make(map[string]int)
	uniqueWords := []string{}

	for _, word := range strings.Fields(text) {
		word = strings.Trim(word, ".?!,:;-[]{}()\"'")
		word = strings.ToLower(word)
		if len(word) == 0 {
			continue
		}

		if _, found := wordCounters[word]; !found {
			uniqueWords = append(uniqueWords, word)
		}
		wordCounters[word]++
	}

	sort.Slice(uniqueWords, func(i, j int) bool {
		wordI := uniqueWords[i]
		wordJ := uniqueWords[j]
		if wordCounters[wordI] == wordCounters[wordJ] {
			return wordI < wordJ
		}
		return wordCounters[wordI] > wordCounters[wordJ]
	})

	return uniqueWords[:min(len(uniqueWords), 10)]
}
