package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var lettersInWord = regexp.MustCompile(`[a-zA-Zа-яА-Я]`)

func Top10(st string) []string {
	wordCouter := make(map[string]int)
	formatedString := strings.ReplaceAll(st, "\n", " ")
	stSlice := strings.Split(formatedString, " ")

	for _, item := range stSlice {
		itemWord := strings.ToLower(item)
		if lettersInWord.MatchString(itemWord) {
			itemWord = strings.TrimFunc(itemWord, func(r rune) bool {
				return !unicode.IsLetter(r) && !unicode.IsNumber(r)
			})
		} else if len(itemWord) <= 1 {
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
