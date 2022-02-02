package hw03frequencyanalysis

import (
	"bufio"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

func Top10(inText string) ([]string, error) {
	const (
		splitSet = `[^-\p{L}]+`
	)
	var (
		scanner    = bufio.NewScanner(strings.NewReader(inText))
		regToSplit = regexp.MustCompile(splitSet)

		wordToFrequency = make(map[string]int64)
		arrangedWords   []string
		lessFunc        = func(i, j int) bool {
			if wordToFrequency[arrangedWords[i]] == wordToFrequency[arrangedWords[j]] {
				return arrangedWords[i] < arrangedWords[j]
			}
			return wordToFrequency[arrangedWords[i]] > wordToFrequency[arrangedWords[j]]
		}
		trimFunc = func(r rune) bool { return !unicode.IsLetter(r) }
	)

	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		token := scanner.Text()
		for _, word := range regToSplit.Split(token, -1) {
			word = strings.ToLower(strings.TrimFunc(word, trimFunc))
			if len(word) < 1 {
				continue
			}
			if _, ok := wordToFrequency[word]; !ok {
				wordToFrequency[word] = 1
				arrangedWords = append(arrangedWords, word)
				continue
			}
			wordToFrequency[word]++
		}
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("top10 err: %w", scanner.Err())
	}
	sort.Slice(arrangedWords, lessFunc)

	if len(arrangedWords) < 10 {
		return arrangedWords, nil
	}
	return arrangedWords[:10], nil
}
