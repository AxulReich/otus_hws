package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(inText string) []string {
	const (
		cutSet       = `,.!@#$%^&*()-"'[]\/+=|`
		regNonCutSet = `[^-\p{L}]+`
	)
	var (
		words           = strings.Fields(inText)
		wordToFrequency = make(map[string]int64)
		arrangedWords   []string
		reg             = regexp.MustCompile(regNonCutSet)
		lessFunc        = func(i, j int) bool {
			if wordToFrequency[arrangedWords[i]] == wordToFrequency[arrangedWords[j]] {
				return arrangedWords[i] < arrangedWords[j]
			}
			return wordToFrequency[arrangedWords[i]] > wordToFrequency[arrangedWords[j]]
		}
	)

	for _, word := range words {
		word = strings.ToLower(strings.Trim(word, cutSet))
		if len(word) < 1 {
			continue
		}
		word = reg.ReplaceAllString(word, "")

		if _, ok := wordToFrequency[word]; !ok {
			wordToFrequency[word] = 1
			arrangedWords = append(arrangedWords, word)
			sort.Slice(arrangedWords, lessFunc)
			continue
		}
		wordToFrequency[word]++
		sort.Slice(arrangedWords, lessFunc)
	}
	if len(arrangedWords) < 10 {
		return arrangedWords
	}
	return arrangedWords[:10]
}
