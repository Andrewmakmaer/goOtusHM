// package main

package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

func Top10(st string) []string {
	wordCouter := make(map[string]int)
	formatedString := strings.ReplaceAll(st, "\n", " ")
	stSlice := strings.Split(formatedString, " ")
	for _, item := range stSlice {
		itemWord := strings.TrimSpace(item)
		itemWord = strings.ToLower(itemWord)
		itemWord = strings.TrimFunc(itemWord, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		if itemWord == "" {
			continue
		}
		value, ok := wordCouter[itemWord]
		if ok {
			wordCouter[itemWord] = value + 1
		} else {
			wordCouter[itemWord] = 1
		}
	}
	words := make([]string, 0, len(wordCouter))

	for word := range wordCouter {
		words = append(words, word)
	}

	sort.SliceStable(words, func(i, j int) bool {
		if wordCouter[words[i]] != wordCouter[words[j]] {
			return wordCouter[words[i]] > wordCouter[words[j]]
		}
		return words[i] < words[j]
	})

	if len(words) > 10 {
		return words[:10]
	}

	return words
}
