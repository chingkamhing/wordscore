package main

import (
	"log/slog"
	"slices"
	"strings"
)

// CommonAbbreviation map the full name to it's corresponding abbreviation
var CommonAbbreviation = map[string]string{
	// top 100 countries ranked by their financial assets
	"UNITED STATES":        "USA",
	"CHINA":                "CHN",
	"JAPAN":                "JPN",
	"GERMANY":              "DEU",
	"UNITED KINGDOM":       "GBR",
	"FRANCE":               "FRA",
	"INDIA":                "IND",
	"CANADA":               "CAN",
	"ITALY":                "ITA",
	"AUSTRALIA":            "AUS",
	"SOUTH KOREA":          "KOR",
	"SWITZERLAND":          "CHE",
	"NETHERLANDS":          "NLD",
	"BRAZIL":               "BRA",
	"SPAIN":                "ESP",
	"RUSSIA":               "RUS",
	"MEXICO":               "MEX",
	"INDONESIA":            "IDN",
	"SAUDI ARABIA":         "SAU",
	"TURKEY":               "TUR",
	"TAIWAN":               "TWN",
	"SWEDEN":               "SWE",
	"BELGIUM":              "BEL",
	"POLAND":               "POL",
	"THAILAND":             "THA",
	"AUSTRIA":              "AUT",
	"NORWAY":               "NOR",
	"UNITED ARAB EMIRATES": "ARE",
	"IRAN":                 "IRN",
	"SINGAPORE":            "SGP",
	"ISRAEL":               "ISR",
	"HONG KONG":            "HKG",
	"MALAYSIA":             "MYS",
	"DENMARK":              "DNK",
	"SOUTH AFRICA":         "ZAF",
	"PHILIPPINES":          "PHL",
	"EGYPT":                "EGY",
	"VIETNAM":              "VNM",
	"PAKISTAN":             "PAK",
	"ARGENTINA":            "ARG",
	"COLOMBIA":             "COL",
	"CHILE":                "CHL",
	"BANGLADESH":           "BGD",
	"FINLAND":              "FIN",
	"NIGERIA":              "NGA",
	"IRELAND":              "IRL",
	"PORTUGAL":             "PRT",
	"GREECE":               "GRC",
	"CZECH REPUBLIC":       "CZE",
	"ROMANIA":              "ROU",
	"PERU":                 "PER",
	"NEW ZEALAND":          "NZL",
	"IRAQ":                 "IRQ",
	"QATAR":                "QAT",
	"KAZAKHSTAN":           "KAZ",
	"HUNGARY":              "HUN",
	"UKRAINE":              "UKR",
	"KUWAIT":               "KWT",
	"MOROCCO":              "MAR",
	"SLOVAKIA":             "SVK",
	"SRI LANKA":            "LKA",
	"ECUADOR":              "ECU",
	"ANGOLA":               "AGO",
	"OMAN":                 "OMN",
	"CUBA":                 "CUB",
	"BELARUS":              "BLR",
	"AZERBAIJAN":           "AZE",
	"SUDAN":                "SDN",
	"DOMINICAN REPUBLIC":   "DOM",
	"LUXEMBOURG":           "LUX",
	"MYANMAR":              "MMR",
	"UZBEKISTAN":           "UZB",
	"KENYA":                "KEN",
	"GUATEMALA":            "GTM",
	"BULGARIA":             "BGR",
	"TUNISIA":              "TUN",
	"SERBIA":               "SRB",
	"ETHIOPIA":             "ETH",
	"CROATIA":              "HRV",
	"LEBANON":              "LBN",
	"LITHUANIA":            "LTU",
	"SLOVENIA":             "SVN",
	"GHANA":                "GHA",
	"TANZANIA":             "TZA",
	"PANAMA":               "PAN",
	"COSTA RICA":           "CRI",
	"JORDAN":               "JOR",
	"BOLIVIA":              "BOL",
	"PARAGUAY":             "PRY",
	"URUGUAY":              "URY",
	"CAMEROON":             "CMR",
	"EL SALVADOR":          "SLV",
	"UGANDA":               "UGA",
	"NEPAL":                "NPL",
	"HONDURAS":             "HND",
	"CYPRUS":               "CYP",
	"ICELAND":              "ISL",
	"ZAMBIA":               "ZMB",
	"CAMBODIA":             "KHM",
	"SENEGAL":              "SEN",

	// common abbreviations for some common words
	"INTERNATIONAL": "INTL",
}

type ScoreFunc func(score int, chars string) []int

type AbbreviateFunc func(chars string) (int, int)

type WordScore struct {
	Length         int
	Chars          []rune
	CharScores     []*CharScore
	MaxCombination int // Maximum number of combinations to generate
}

type CharScore struct {
	Char  rune
	Score int
	Index int
}

type WordScoreOption func(*WordScore)

func WordScoreOptionMaxCombinations(maxCombinations int) WordScoreOption {
	return func(wordScore *WordScore) {
		wordScore.MaxCombination = maxCombinations
	}
}

// NewWordScore creates a new WordScore instance that generate abbreviate combination words which is easily readable.
//
// Typical flow to get the abbreviation combinations:
// - Remove(): remove common words which hard to make the abbreviation unique (e.g. "Venture" for vessel, "Limited" for company, etc.)
// - Abbreviate(): abbreviate the word (e.g. "International" to "Intl", "Hong Kong" to "HKG", etc.)
// - Score(): calculate the char scores (e.g. acronyms, consonants, etc.)
// - Remove(): remove unused words or char which might still usefull for score calculation (e.g. space, special characters, etc.)
// - Transform(): transform the word (e.g. to uppercase or lowercase)
// - Combinations(): generate the combination words based on the char scores and the specified min and max length
func NewWordScore(word string, options ...WordScoreOption) *WordScore {
	lenWord := len(word)
	wordScore := &WordScore{
		Length:         lenWord,
		Chars:          []rune(word),
		CharScores:     make([]*CharScore, lenWord),
		MaxCombination: 10, // Default maximum combinations
	}
	for _, opt := range options {
		opt(wordScore)
	}
	for i := range lenWord {
		wordScore.CharScores[i] = &CharScore{
			Char:  wordScore.Chars[i],
			Index: i,
		}
	}
	return wordScore
}

// Abbreviate abbreviate a map of full names to their corresponding abbreviations. Expect the abbreviation full name to be uppercase. If found, the abbreviation string will have the specified score.
func (ws *WordScore) Abbreviate(score int, abbreviationMap map[string]string) {
	for full, abbr := range abbreviationMap {
		chars := strings.ToUpper(string(ws.Chars))
		base := strings.Index(chars, full)
		if base < 0 {
			continue
		}
		slog.Debug("Abbreviate", "full", full, "abbr", abbr, "base", base)
		charScores := make([]*CharScore, len(abbr))
		for i, c := range abbr {
			charScores[i] = &CharScore{Char: c, Score: score}
		}
		length := len(full)
		ws.Chars = slices.Concat(ws.Chars[:base], []rune(abbr), ws.Chars[base+length:])
		ws.CharScores = slices.Concat(ws.CharScores[:base], charScores, ws.CharScores[base+length:])
		ws.Length = len(ws.Chars)
		for i := range ws.Length {
			ws.CharScores[i].Index = i
		}
	}
}

// Transform applies a transformation function to the word, such as converting it to uppercase or lowercase. Expect the transformation function to return a string that has the same length as the original word.
func (ws *WordScore) Transform(transform func(word string) string) string {
	result := transform(string(ws.Chars))
	for i, c := range result {
		ws.Chars[i] = c
		ws.CharScores[i].Char = c
	}
	return result
}

func (ws *WordScore) Remove(removeMap map[string]struct{}) {
removeLoop:
	for remove := range removeMap {
		for range ws.Length {
			base := strings.Index(string(ws.Chars), remove)
			if base < 0 {
				continue removeLoop
			}
			length := len(remove)
			ws.Chars = slices.Concat(ws.Chars[:base], ws.Chars[base+length:])
			ws.CharScores = slices.Concat(ws.CharScores[:base], ws.CharScores[base+length:])
			ws.Length = len(ws.Chars)
			for i := range ws.Length {
				ws.CharScores[i].Index = i
			}
		}
	}
}

func (ws *WordScore) Score(score int, scoreFunc ScoreFunc) {
	scores := scoreFunc(score, string(ws.Chars))
	for i := range ws.Length {
		ws.CharScores[i].Score += scores[i]
		slog.Debug("Score", "i", i, "char", string(ws.CharScores[i].Char), "score", ws.CharScores[i].Score)
	}
}

// Combinations returns a list of combination strings base on different combination of Chars while meeting the specified min and max number of characters.
func (ws *WordScore) Combinations(length int) []string {
	// generate all combinations
	combinations := ws.generateCombinations(length)
	// cap the combinations to the maximum number of combinations
	combinations = combinations[:min(len(combinations), ws.MaxCombination)]
	// pad the combinations with zeroes to meet the minimum length
	combinations = padZero(combinations, length)
	return combinations
}

func (ws *WordScore) generateCombinations(length int) []string {
	// create a map of score to list of CharScore
	maxScore := ws.maxScore()
	scoreChars := make(map[int][]*CharScore, maxScore)
	for i := range maxScore {
		scoreChars[i] = make([]*CharScore, 0, ws.Length)
	}
	for _, cs := range ws.CharScores {
		scoreChars[cs.Score] = append(scoreChars[cs.Score], cs)
	}

	// find all the CharScore that just meet the length sorting from hightest to lowest score
	candidateChars := []*CharScore{}
	const minScore = 1
	for score := maxScore; score >= minScore; score-- {
		for _, cs := range scoreChars[score] {
			candidateChars = append(candidateChars, cs)
		}
	}
	slices.SortFunc(candidateChars, func(a, b *CharScore) int {
		return int(a.Index) - int(b.Index)
	})
	slog.Debug("candidateChars", "chars", charsToString(candidateChars))

	// create a map of length-of-chars to list of CharScore
	results := make([]string, 0, ws.MaxCombination)
	accumulateChars := make([]*CharScore, 0, ws.Length)
	lengthCharsMap := make(map[int][]*CharScore, maxScore)
forLoop:
	for score := maxScore; score >= minScore; score-- {
		prevChars := slices.Clone(accumulateChars)
		currChars := scoreChars[score]
		accumulateChars = append(accumulateChars, currChars...)
		slices.SortFunc(accumulateChars, func(a, b *CharScore) int {
			return int(a.Index) - int(b.Index)
		})
		lengthAccumulate := len(accumulateChars)
		lengthCharsMap[lengthAccumulate] = accumulateChars
		slog.Debug("LengthChars", "score", score, "length", lengthAccumulate, "curr", charsToString(currChars), "prev", charsToString(prevChars), "accumulate", charsToString(accumulateChars))
		switch {
		case lengthAccumulate < length:
			// does not meet the min length, skip
		case lengthAccumulate == length:
			results = append(results, charsToString(accumulateChars))
			break forLoop
		case lengthAccumulate > length:
			// find a combination of chars that have the exact length of length while sorting by highest score
			remainLength := length - len(prevChars)
			combinations := ws.combinationsChars(currChars, remainLength)
			for _, chars := range combinations {
				candidateChars := append(prevChars, chars...)
				slices.SortFunc(candidateChars, func(a, b *CharScore) int {
					return int(a.Index) - int(b.Index)
				})
				results = append(results, charsToString(candidateChars))
			}
			break forLoop
		}
	}

	// in case found no combinations, return the original chars
	if len(results) == 0 {
		results = append(results, charsToString(ws.CharScores))
	}
	return results
}

// combinationsChars generates all possible combinations of CharScores that have the exact length of length.
func (ws *WordScore) combinationsChars(chars []*CharScore, length int) [][]*CharScore {
	var results [][]*CharScore
	var backtrack func(start int, current []*CharScore)
	backtrack = func(start int, current []*CharScore) {
		if len(current) == length {
			results = append(results, append([]*CharScore(nil), current...))
			return
		}
		for i := start; i < len(chars); i++ {
			current = append(current, chars[i])
			backtrack(i+1, current)
			current = current[:len(current)-1]
		}
	}
	backtrack(0, []*CharScore{})
	return results
}

func (ws *WordScore) maxScore() int {
	maxScore := 0
	for _, cs := range ws.CharScores {
		if cs.Score > maxScore {
			maxScore = cs.Score
		}
	}
	return maxScore
}

func padZero(words []string, length int) []string {
	results := make([]string, 0, len(words))
	for _, word := range words {
		if len(word) >= length {
			results = append(results, word)
		} else {
			results = append(results, word+strings.Repeat("0", length-len(word)))
		}
	}
	return results
}

func charsToString(chars []*CharScore) string {
	var sb strings.Builder
	for _, c := range chars {
		sb.WriteRune(c.Char)
	}
	return sb.String()
}
