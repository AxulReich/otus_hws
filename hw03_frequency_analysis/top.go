package hw03frequencyanalysis

import (
	"bufio"
	"regexp"
	"sort"
	"strings"
)

func Top10(inText string) []string {
	const (
		cutSet       = `,.!@#$%^&*()-"'[]\/+=|:;*`
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
		//word = strings.ToLower(strings.Trim(word, cutSet))
		word = strings.ToLower(reg.ReplaceAllString(word, ""))
		if len(word) < 1 || word == "-" {
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

func MyTop10(inText string) []string {
	const (
		cutSet   = `[^\p{L}]+`
		splitSet = `[^-\p{L}]+`
	)
	var (
		wordToFrequency = make(map[string]int64)
		arrangedWords []string
		regToSplit = regexp.MustCompile(splitSet)
		regToTrim  = regexp.MustCompile(cutSet)
		lessFunc   = func(i, j int) bool {
			if wordToFrequency[arrangedWords[i]] == wordToFrequency[arrangedWords[j]] {
				return arrangedWords[i] < arrangedWords[j]
			}
			return wordToFrequency[arrangedWords[i]] > wordToFrequency[arrangedWords[j]]
		}
	)
	scanner := bufio.NewScanner(strings.NewReader(inText))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		token := scanner.Text()
		for _, word := range regToSplit.Split(token, -1) {
			word = strings.ToLower(regToTrim.ReplaceAllString(word, ""))
			if len(word) < 1 {
				continue
			}
			if _, ok := wordToFrequency[word]; !ok {
				wordToFrequency[word] = 1
				arrangedWords = append(arrangedWords, word)
				sort.Slice(arrangedWords, lessFunc)
				continue
			}
			wordToFrequency[word]++
			sort.Slice(arrangedWords, lessFunc)
		}
	}

	if len(arrangedWords) < 10 {
		return arrangedWords
	}
	return arrangedWords[:10]
}