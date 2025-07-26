package main

import (
	"flag"
	"fmt"
	"log/slog"
)

func main() {
	// parse flags
	word := flag.String("word", "abc", "Word to be generated candidate strings from")
	minLength := flag.Int("min", 4, "Minimum length of candidate strings")
	maxLength := flag.Int("max", 4, "Maximum length of candidate strings")
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
		// WordScore options
	}
	myWord := NewWordScore(*word, options...)
	slog.Debug("Set Capital score")
	myWord.Score(1, scoreCapital)
	slog.Debug("Set Consonant score")
	myWord.Score(1, scoreConsonant)
	candidates := myWord.Candidates(*minLength, *maxLength)
	fmt.Printf("Word: %v\n", *word)
	fmt.Printf("Candidates: %v\n", candidates)
}

func scoreCapital(score int, chars []rune) []int {
	isCapital := func(char rune) bool {
		return char >= 'A' && char <= 'Z'
	}
	scores := make([]int, len(chars))
	for i, char := range chars {
		if isCapital(char) {
			scores[i] = score
		}
	}
	return scores
}

func scoreConsonant(score int, chars []rune) []int {
	isConsonant := func(char rune) bool {
		switch char {
		case 'a', 'e', 'i', 'o', 'u':
			return false
		default:
			return true
		}
	}
	scores := make([]int, len(chars))
	for i, char := range chars {
		if isConsonant(char) {
			scores[i] = score
		}
	}
	return scores
}
