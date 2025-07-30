package main

import (
	"fmt"
	"log/slog"
	"reflect"
	"testing"
)

// test WordScore
func Test_WordScore(t *testing.T) {
	logLevel := slog.LevelInfo
	tests := []struct {
		word      *WordScore
		minLength int
		maxLength int
		expect    []string
	}{
		// 0
		{&WordScore{MaxCombination: 100}, 4, 4, []string{"0000"}},
		// 1
		{&WordScore{
			Chars:          []rune{'A', 'B', 'C'},
			CharScores:     []*CharScore{{Char: 'A', Index: 0, Score: 1}, {Char: 'B', Index: 1, Score: 1}, {Char: 'C', Index: 2, Score: 1}},
			MaxCombination: 100,
		}, 4, 4, []string{"ABC0"}},
		// 2
		{&WordScore{
			Chars:          []rune{'A', 'B', 'C', 'D'},
			CharScores:     []*CharScore{{Char: 'A', Index: 0, Score: 1}, {Char: 'B', Index: 1, Score: 1}, {Char: 'C', Index: 2, Score: 1}, {Char: 'D', Index: 3, Score: 1}},
			MaxCombination: 100,
		}, 4, 4, []string{"ABCD"}},
		// 3
		{&WordScore{
			Chars:          []rune{'A', 'B', 'C', 'D', 'E'},
			CharScores:     []*CharScore{{Char: 'A', Index: 0, Score: 1}, {Char: 'B', Index: 1, Score: 1}, {Char: 'C', Index: 2, Score: 1}, {Char: 'D', Index: 3, Score: 1}, {Char: 'E', Index: 4, Score: 1}},
			MaxCombination: 100,
		}, 4, 4, []string{"ABCD", "ABCE", "ABDE", "ACDE", "BCDE"}},
		// 4
		{&WordScore{
			Chars:          []rune{'A', 'B', 'C', 'D', 'E', 'F'},
			CharScores:     []*CharScore{{Char: 'A', Index: 0, Score: 1}, {Char: 'B', Index: 1, Score: 1}, {Char: 'C', Index: 2, Score: 1}, {Char: 'D', Index: 3, Score: 1}, {Char: 'E', Index: 4, Score: 1}, {Char: 'F', Index: 5, Score: 1}},
			MaxCombination: 100,
		}, 4, 4, []string{"ABCD", "ABCE", "ABCF", "ABDE", "ABDF", "ABEF", "ACDE", "ACDF", "ACEF", "ADEF", "BCDE", "BCDF", "BCEF", "BDEF", "CDEF"}},
		// 5
		{&WordScore{
			Chars:          []rune{'A', 'B', 'C', 'D', 'E', 'F'},
			CharScores:     []*CharScore{{Char: 'A', Index: 0, Score: 2}, {Char: 'B', Index: 1, Score: 1}, {Char: 'C', Index: 2, Score: 1}, {Char: 'D', Index: 3, Score: 1}, {Char: 'E', Index: 4, Score: 1}, {Char: 'F', Index: 5, Score: 1}},
			MaxCombination: 100,
		}, 4, 4, []string{"ABCD", "ABCE", "ABCF", "ABDE", "ABDF", "ABEF", "ACDE", "ACDF", "ACEF", "ADEF"}},
		// 6
		{&WordScore{
			Chars:          []rune{'A', 'B', 'C', 'D', 'E', 'F'},
			CharScores:     []*CharScore{{Char: 'A', Index: 0, Score: 3}, {Char: 'B', Index: 1, Score: 1}, {Char: 'C', Index: 2, Score: 2}, {Char: 'D', Index: 3, Score: 1}, {Char: 'E', Index: 4, Score: 1}, {Char: 'F', Index: 5, Score: 1}},
			MaxCombination: 100,
		}, 4, 4, []string{"ABCD", "ABCE", "ABCF", "ACDE", "ACDF", "ACEF"}},
	}
	slog.SetLogLoggerLevel(logLevel)
	for i, test := range tests {
		t.Run(fmt.Sprintf("Test split() [%v]:", i), func(t *testing.T) {
			test.word.Length = len(test.word.Chars)
			result := test.word.Combinations(test.minLength, test.maxLength)
			if !reflect.DeepEqual(result, test.expect) {
				t.Fatalf("  Test %v mismatch: %q", i, result)
			}
			t.Logf("  Test %v succeed.", i)
		})
	}
}
