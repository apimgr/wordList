package words

import (
	"embed"
	"encoding/json"
	"math/rand"
	"sort"
	"strings"
	"time"
)

//go:embed data/words.json
var embeddedWords embed.FS

type WordData struct {
	Words []string `json:"words"`
}

var wordList []string
var rng *rand.Rand

// Load loads the word list from embedded data
func Load() error {
	data, err := embeddedWords.ReadFile("data/words.json")
	if err != nil {
		return err
	}

	var wd WordData
	if err := json.Unmarshal(data, &wd); err != nil {
		return err
	}

	// Filter out the "end" placeholder
	wordList = make([]string, 0, len(wd.Words))
	for _, w := range wd.Words {
		if w != "end" && w != "" {
			wordList = append(wordList, w)
		}
	}

	// Initialize random with time seed
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	return nil
}

// Count returns the number of words loaded
func Count() int {
	return len(wordList)
}

// All returns all words
func All() []string {
	return wordList
}

// Random returns n random words
func Random(n int) []string {
	if n <= 0 {
		n = 1
	}
	if n > len(wordList) {
		n = len(wordList)
	}

	// Create a copy and shuffle
	shuffled := make([]string, len(wordList))
	copy(shuffled, wordList)

	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:n]
}

// Passphrase generates a passphrase with n words
func Passphrase(n int, separator string, capitalize bool) string {
	if n <= 0 {
		n = 6
	}
	if separator == "" {
		separator = "-"
	}

	words := Random(n)

	if capitalize {
		for i, w := range words {
			if len(w) > 0 {
				words[i] = strings.ToUpper(w[:1]) + w[1:]
			}
		}
	}

	return strings.Join(words, separator)
}

// Search returns words matching a prefix
func Search(prefix string) []string {
	prefix = strings.ToLower(prefix)
	var matches []string

	for _, w := range wordList {
		if strings.HasPrefix(w, prefix) {
			matches = append(matches, w)
		}
	}

	return matches
}

// SearchContains returns words containing a substring
func SearchContains(substring string) []string {
	substring = strings.ToLower(substring)
	var matches []string

	for _, w := range wordList {
		if strings.Contains(w, substring) {
			matches = append(matches, w)
		}
	}

	return matches
}

// ByLetter returns words starting with a specific letter
func ByLetter(letter string) []string {
	letter = strings.ToLower(letter)
	if len(letter) == 0 {
		return nil
	}

	var matches []string
	for _, w := range wordList {
		if strings.HasPrefix(w, letter[:1]) {
			matches = append(matches, w)
		}
	}

	return matches
}

// ByLength returns words of a specific length
func ByLength(length int) []string {
	var matches []string

	for _, w := range wordList {
		if len(w) == length {
			matches = append(matches, w)
		}
	}

	return matches
}

// Stats returns statistics about the word list
func Stats() map[string]interface{} {
	letterCounts := make(map[string]int)
	lengthCounts := make(map[int]int)
	minLen := 999
	maxLen := 0
	totalLen := 0

	for _, w := range wordList {
		if len(w) > 0 {
			letter := strings.ToLower(w[:1])
			letterCounts[letter]++
		}

		length := len(w)
		lengthCounts[length]++

		if length < minLen {
			minLen = length
		}
		if length > maxLen {
			maxLen = length
		}
		totalLen += length
	}

	avgLen := float64(totalLen) / float64(len(wordList))

	// Find most common letter
	mostCommonLetter := ""
	maxCount := 0
	for letter, count := range letterCounts {
		if count > maxCount {
			mostCommonLetter = letter
			maxCount = count
		}
	}

	// Find most common length
	mostCommonLength := 0
	maxLenCount := 0
	for length, count := range lengthCounts {
		if count > maxLenCount {
			mostCommonLength = length
			maxLenCount = count
		}
	}

	return map[string]interface{}{
		"total_words":        len(wordList),
		"min_length":         minLen,
		"max_length":         maxLen,
		"average_length":     avgLen,
		"most_common_letter": mostCommonLetter,
		"most_common_length": mostCommonLength,
		"letters_count":      letterCounts,
		"lengths_count":      lengthCounts,
	}
}

// Letters returns all unique starting letters
func Letters() []string {
	seen := make(map[string]bool)
	var letters []string

	for _, w := range wordList {
		if len(w) > 0 {
			letter := strings.ToLower(w[:1])
			if !seen[letter] {
				seen[letter] = true
				letters = append(letters, letter)
			}
		}
	}

	sort.Strings(letters)
	return letters
}
