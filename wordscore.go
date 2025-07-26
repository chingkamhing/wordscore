package main

import (
	"log/slog"
	"sort"
	"strconv"
	"strings"
)

type ScoreFunc func(score int, chars []rune) []int

type WordScore struct {
	Length        int
	Chars         []rune
	CharScores    []CharScore
	MaxCandidates int // Maximum number of candidates to generate
}

type CharScore struct {
	Char  rune
	Score int
	Index int
}

type WordScoreOption func(*WordScore)

func WordScoreOptionMaxCandidates(maxCandidates int) WordScoreOption {
	return func(wordScore *WordScore) {
		wordScore.MaxCandidates = maxCandidates
	}
}

func NewWordScore(word string, options ...WordScoreOption) *WordScore {
	lenWord := len(word)
	wordScore := &WordScore{
		Length:        lenWord,
		Chars:         []rune(word),
		CharScores:    make([]CharScore, lenWord),
		MaxCandidates: 10, // Default maximum candidates
	}
	for _, opt := range options {
		opt(wordScore)
	}
	for i := range lenWord {
		wordScore.CharScores[i] = CharScore{
			Char:  wordScore.Chars[i],
			Index: i,
		}
	}
	return wordScore
}

// Candidates returns a list of candidate strings base on different combination of Chars while meeting the specified min and max number of characters.
// - base on the highest in Scores of corresponding Chars and return list of candidate strings (e.g. Scores{0, 1, 1, 1} and Chars{'a', 'b', 'c', 'd'}, Candidates(3, 3) return "bcd")
// - the order of the candidate strings must be based on the order of Chars
// - the candidate strings contain all the combination that meet the specified min and max number of characters (i.e. Scores{0, 2, 2, 1} and Chars{'a', 'b', 'c', 'd'}, Candidates(2, 3) return all strings that has min 2 chars and max 3 chars e.g. "bc", "bd", "cd", "bcd")
// - the order of the candidate strings must be sorted from highest to lowest based on the sum of Scores of corresponding Chars
// - if the length of Chars is less than min, append zero-leading number up to max 10 candidate strings (e.g. Chars{'a'}, Candidates(3, 3) return "a00", "a01", "a02", "a03", etc.)
func (ws *WordScore) Candidates(min, max int) []string {
	// Handle case where word length is less than min
	if ws.Length < min {
		return ws.generateZeroPaddedCandidates(min, max)
	}

	// Generate all possible combinations within min-max length
	var candidates []string
	for l := min; l <= max && l <= ws.Length; l++ {
		combinations := ws.generateCombinations(l)
		candidates = append(candidates, combinations...)
	}

	// Sort candidates by total score (descending)
	sort.Slice(candidates, func(i, j int) bool {
		return ws.calculateCandidateScore(candidates[i]) > ws.calculateCandidateScore(candidates[j])
	})

	return candidates
}

func (ws *WordScore) generateZeroPaddedCandidates(min, max int) []string {
	var candidates []string
	base := string(ws.Chars)

	// We'll generate up to 10 candidates as specified
	count := ws.MaxCandidates
	digits := min - ws.Length

	for i := range count {
		// Format the number with leading zeros
		numStr := strconv.Itoa(i)
		if len(numStr) < digits {
			numStr = strings.Repeat("0", digits-len(numStr)) + numStr
		} else if len(numStr) > digits {
			numStr = numStr[:digits]
		}
		candidates = append(candidates, base+numStr)
	}

	return candidates
}

func (ws *WordScore) generateCombinations(length int) []string {
	var result []string
	var backtrack func(start int, current []rune)

	backtrack = func(start int, current []rune) {
		if len(current) == length {
			result = append(result, string(current))
			return
		}

		for i := start; i < ws.Length; i++ {
			current = append(current, ws.CharScores[i].Char)
			backtrack(i+1, current)
			current = current[:len(current)-1]
		}
	}

	backtrack(0, []rune{})
	return result
}

func (ws *WordScore) calculateCandidateScore(candidate string) int {
	score := 0
	chars := []rune(candidate)
	charMap := make(map[rune]int, ws.Length)

	// Create a map of char to score for quick lookup
	for _, cs := range ws.CharScores {
		charMap[cs.Char] = cs.Score
	}

	for _, c := range chars {
		score += charMap[c]
	}
	return score
}

func (ws *WordScore) Score(score int, scoreFunc ScoreFunc) {
	scores := scoreFunc(score, ws.Chars)
	for i := range ws.Length {
		ws.CharScores[i].Score += scores[i]
		slog.Debug("Score", "i", i, "char", ws.CharScores[i].Char, "score", ws.CharScores[i].Score)
	}
}
