package words

import (
	"crypto/rand"
	"embed"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strings"
)

//go:embed data/words.json
var embeddedWords embed.FS

type WordData struct {
	Words []string `json:"words"`
}

var wordList []string
var diceWords map[string]string
var wordDice map[string]string

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

	if len(wd.Words) == 0 {
		return fmt.Errorf("embedded EFF word list is empty")
	}

	nextWords := make([]string, 0, len(wd.Words))
	for _, word := range wd.Words {
		if word == "end" || word == "" {
			continue
		}
		if word != strings.ToLower(word) || strings.TrimSpace(word) != word {
			return fmt.Errorf("invalid word entry %q", word)
		}
		nextWords = append(nextWords, word)
	}
	if len(nextWords) != 7776 {
		return fmt.Errorf("EFF large word list has %d words; expected 7776", len(nextWords))
	}

	nextDiceWords := make(map[string]string, len(nextWords))
	nextWordDice := make(map[string]string, len(nextWords))
	for index, word := range nextWords {
		if index > 0 && nextWords[index-1] >= word {
			return fmt.Errorf("EFF word list is not strictly sorted at index %d", index)
		}
		code := diceCode(index)
		nextDiceWords[code] = word
		nextWordDice[word] = code
	}
	wordList = nextWords
	diceWords = nextDiceWords
	wordDice = nextWordDice
	return nil
}

// Count returns the number of words loaded
func Count() int {
	return len(wordList)
}

// All returns all words
func All() []string {
	return append([]string(nil), wordList...)
}

// ByDiceCode returns the EFF word mapped to a five-dice code from 11111 to 66666.
func ByDiceCode(code string) (string, bool) {
	if len(code) != 5 {
		return "", false
	}
	for _, digit := range code {
		if digit < '1' || digit > '6' {
			return "", false
		}
	}
	word, ok := diceWords[code]
	return word, ok
}

// DiceCodeByWord returns the EFF dice code for an exact word-list entry.
func DiceCodeByWord(word string) (string, bool) {
	code, ok := wordDice[strings.ToLower(word)]
	return code, ok
}

func diceCode(index int) string {
	var code [5]byte
	for position := len(code) - 1; position >= 0; position-- {
		code[position] = byte(index%6) + '1'
		index /= 6
	}
	return string(code[:])
}

// Random returns n words selected independently with the operating-system CSPRNG.
func Random(n int) []string {
	result, _ := RandomSecure(n)
	return result
}

// RandomSecure returns n words using crypto/rand with replacement.
func RandomSecure(n int) ([]string, error) {
	if n <= 0 {
		n = 1
	}
	if len(wordList) == 0 {
		return nil, fmt.Errorf("word list is not loaded")
	}
	if n > len(wordList) {
		n = len(wordList)
	}
	result := make([]string, n)
	max := big.NewInt(int64(len(wordList)))
	for i := range result {
		index, err := rand.Int(rand.Reader, max)
		if err != nil {
			return nil, fmt.Errorf("read secure random word index: %w", err)
		}
		result[i] = wordList[index.Int64()]
	}
	return result, nil
}

// Passphrase generates a passphrase using the OS CSPRNG. Empty indicates a random source error.
func Passphrase(n int, separator string, capitalize bool) string {
	if separator == "" {
		separator = "-"
	}
	result, _ := PassphraseSecure(n, separator, capitalize)
	return result
}

// PassphraseSecure generates a passphrase and reports CSPRNG failures to its caller.
func PassphraseSecure(n int, separator string, capitalize bool) (string, error) {
	if n <= 0 {
		n = 6
	}
	if n > 20 {
		return "", fmt.Errorf("passphrase word count %d exceeds maximum 20", n)
	}
	selected, err := RandomSecure(n)
	if err != nil {
		return "", err
	}
	if capitalize {
		for i, word := range selected {
			selected[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(selected, separator), nil
}

// EntropyBits reports the theoretical entropy of independent selections from the loaded list.
func EntropyBits(n int) float64 {
	if n <= 0 || len(wordList) == 0 {
		return 0
	}
	return float64(n) * math.Log2(float64(len(wordList)))
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
