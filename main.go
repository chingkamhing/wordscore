package main

import (
	"flag"
	"fmt"
	"log/slog"
	"strings"
	"unicode"
)

var removeCommonWord = map[string]struct{}{
	"Venture": {},
}

func main() {
	// parse flags
	word := flag.String("word", "abc", "Word to generate combination strings from")
	length := flag.Int("length", 4, "Length of combination strings")
	count := flag.Int("count", 10, "Number of candidates to generate")
	debugLevel := flag.String("debug", "info", "Enable debug printf")
	flag.Parse()

	// set log level
	switch *debugLevel {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	case "info":
		fallthrough
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	options := []WordScoreOption{
		WordScoreOptionMaxCombinations(*count),
	}
	myWord := NewWordScore(*word, options...)
	slog.Debug("Remove common words")
	myWord.Remove(removeCommonWord)
	slog.Debug("Set Capital score")
	myWord.Score(1, scoreCapital)
	slog.Debug("Set Acronym score")
	myWord.Score(1, scoreAcronym)
	slog.Debug("Abbreviate word")
	myWord.Abbreviate(1, CommonAbbreviation)
	slog.Debug("Transform to uppercase")
	myWord.Transform(transformUppercase)
	slog.Debug("Set Consonant score")
	myWord.Score(1, scoreConsonant)
	slog.Debug("Set Letter score")
	myWord.Score(1, scoreLetter)
	combinations := myWord.Combinations(*length)
	quoted := make([]string, len(combinations))
	for i, s := range combinations {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	fmt.Printf("%v\n", strings.Join(quoted, ", "))
}

func scoreCapital(score int, chars string) []int {
	isAllUpper := func(s string) bool {
		for _, c := range s {
			if unicode.IsLower(c) {
				return false
			}
		}
		return true
	}
	scores := make([]int, len(chars))
	if isAllUpper(chars) {
		// string is all uppercase, set score for the first letter
		isFirstLetter := true
		for i, char := range chars {
			if isFirstLetter {
				if unicode.IsLetter(char) {
					scores[i] = score
				}
				isFirstLetter = false
			} else {
				if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
					isFirstLetter = true
				}
			}
		}
	} else {
		// string is mix of uppercase and lowercase, set score for the uppercase letters only
		for i, char := range chars {
			if unicode.IsUpper(char) {
				scores[i] = score
			}
		}
	}
	return scores
}

func scoreAcronym(score int, chars string) []int {
	scores := make([]int, len(chars))
	isSplit := func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r)
	}
	upper := strings.ToUpper(chars)
	var word strings.Builder
	for i, char := range upper {
		if isSplit(char) {
			lastWord := word.String()
			if isAllConsonant(lastWord) {
				for j := i - len(word.String()); j < i; j++ {
					scores[j] = score
				}
			}
		} else {
			word.WriteRune(char)
		}
	}
	return scores
}

func scoreConsonant(score int, chars string) []int {
	scores := make([]int, len(chars))
	for i, char := range chars {
		if isConsonant(char) {
			scores[i] = score
		}
	}
	return scores
}

func scoreLetter(score int, chars string) []int {
	scores := make([]int, len(chars))
	for i, char := range chars {
		if unicode.IsLetter(char) {
			scores[i] = score
		}
	}
	return scores
}

func isConsonant(char rune) bool {
	if unicode.IsLetter(char) {
		char = unicode.ToUpper(char)
		switch char {
		case 'A', 'E', 'I', 'O', 'U':
			return false
		default:
			return true
		}
	}
	return false
}

func isAllConsonant(word string) bool {
	for _, char := range word {
		if !isConsonant(char) {
			return false
		}
	}
	return true
}

func transformUppercase(word string) string {
	return strings.ToUpper(word)
}
