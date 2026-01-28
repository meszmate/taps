package typing

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/meszmate/taps/internal/words"
)

type Quote struct {
	Text   string `json:"text"`
	Source string `json:"source"`
	Length string `json:"length"`
}

var (
	englishWords   []string
	english1kWords []string
	quotes         []Quote
)

func init() {
	_ = json.Unmarshal(words.EnglishJSON, &englishWords)
	_ = json.Unmarshal(words.English1kJSON, &english1kWords)
	_ = json.Unmarshal(words.QuotesJSON, &quotes)
}

func GetWordList(language string) []string {
	switch language {
	case "english_1k":
		return english1kWords
	default:
		return englishWords
	}
}

var punctuationMarks = []string{".", ",", ";", ":", "!", "?"}

func GenerateWords(count int, language string, addPunctuation, addNumbers bool) string {
	wordList := GetWordList(language)
	if len(wordList) == 0 {
		return ""
	}

	result := make([]string, 0, count)
	for i := 0; i < count; i++ {
		if addNumbers && rand.Float64() < 0.1 {
			result = append(result, fmt.Sprintf("%d", rand.Intn(100)))
			continue
		}

		word := wordList[rand.Intn(len(wordList))]

		if addPunctuation && rand.Float64() < 0.15 {
			p := punctuationMarks[rand.Intn(len(punctuationMarks))]
			if rand.Float64() < 0.5 {
				word = word + p
			} else {
				word = strings.ToUpper(word[:1]) + word[1:]
			}
		}

		result = append(result, word)
	}

	return strings.Join(result, " ")
}

func GenerateWordsForTime(language string, addPunctuation, addNumbers bool) string {
	return GenerateWords(200, language, addPunctuation, addNumbers)
}

func GetRandomQuote(length string) Quote {
	var filtered []Quote
	for _, q := range quotes {
		if q.Length == length {
			filtered = append(filtered, q)
		}
	}
	if len(filtered) == 0 {
		if len(quotes) == 0 {
			return Quote{Text: "No quotes available.", Source: "System", Length: "short"}
		}
		return quotes[rand.Intn(len(quotes))]
	}
	return filtered[rand.Intn(len(filtered))]
}
