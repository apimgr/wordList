package words

import (
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	if err := Load(); err != nil {
		panic(err)
	}
	m.Run()
}

func TestLoad(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if Count() == 0 {
		t.Fatal("expected words to be loaded")
	}
	for _, w := range All() {
		if w == "end" || w == "" {
			t.Fatalf("placeholder word %q was not filtered out", w)
		}
	}
}

func TestCountAndAll(t *testing.T) {
	all := All()
	if len(all) != Count() {
		t.Fatalf("Count() = %d, len(All()) = %d", Count(), len(all))
	}
}

func TestRandom(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"positive", 5, 5},
		{"zero", 0, 1},
		{"negative", -3, 1},
		{"more than available", Count() + 1000, Count()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Random(tt.n)
			if len(got) != tt.want {
				t.Errorf("Random(%d) len = %d, want %d", tt.n, len(got), tt.want)
			}
		})
	}
}

func TestPassphrase(t *testing.T) {
	p := Passphrase(4, "-", false)
	parts := strings.Split(p, "-")
	if len(parts) != 4 {
		t.Fatalf("expected 4 words, got %d: %q", len(parts), p)
	}

	// defaults
	def := Passphrase(0, "", false)
	if strings.Count(def, "-") != 5 {
		t.Fatalf("expected default 6 words joined by '-', got %q", def)
	}

	// capitalization
	capped := Passphrase(3, "-", true)
	for _, w := range strings.Split(capped, "-") {
		if len(w) > 0 && !(w[0] >= 'A' && w[0] <= 'Z') {
			t.Errorf("expected capitalized word, got %q", w)
		}
	}
}

func TestSearch(t *testing.T) {
	all := All()
	if len(all) == 0 {
		t.Skip("no words loaded")
	}
	prefix := all[0][:1]
	results := Search(prefix)
	if len(results) == 0 {
		t.Fatalf("expected matches for prefix %q", prefix)
	}
	for _, w := range results {
		if !strings.HasPrefix(w, strings.ToLower(prefix)) {
			t.Errorf("word %q does not have prefix %q", w, prefix)
		}
	}
}

func TestSearchContains(t *testing.T) {
	all := All()
	if len(all) == 0 {
		t.Skip("no words loaded")
	}
	sub := all[0]
	if len(sub) > 3 {
		sub = sub[1:3]
	}
	results := SearchContains(sub)
	if len(results) == 0 {
		t.Fatalf("expected matches for substring %q", sub)
	}
	for _, w := range results {
		if !strings.Contains(w, strings.ToLower(sub)) {
			t.Errorf("word %q does not contain %q", w, sub)
		}
	}
}

func TestByLetter(t *testing.T) {
	if got := ByLetter(""); got != nil {
		t.Errorf("ByLetter(\"\") = %v, want nil", got)
	}

	results := ByLetter("a")
	for _, w := range results {
		if !strings.HasPrefix(w, "a") {
			t.Errorf("word %q does not start with 'a'", w)
		}
	}
}

func TestByLength(t *testing.T) {
	results := ByLength(5)
	for _, w := range results {
		if len(w) != 5 {
			t.Errorf("word %q does not have length 5", w)
		}
	}
}

func TestStats(t *testing.T) {
	stats := Stats()
	if stats["total_words"] != Count() {
		t.Errorf("total_words = %v, want %d", stats["total_words"], Count())
	}
	if stats["min_length"].(int) > stats["max_length"].(int) {
		t.Errorf("min_length %v > max_length %v", stats["min_length"], stats["max_length"])
	}
}

func TestLetters(t *testing.T) {
	letters := Letters()
	if len(letters) == 0 {
		t.Fatal("expected at least one letter")
	}
	seen := make(map[string]bool)
	for _, l := range letters {
		if seen[l] {
			t.Errorf("duplicate letter %q", l)
		}
		seen[l] = true
	}
}


func TestDiceCodeMapping(t *testing.T) {
	if word, ok := ByDiceCode("11111"); !ok || word != "abacus" {
		t.Fatalf("ByDiceCode(11111) = %q, %v", word, ok)
	}
	if word, ok := ByDiceCode("66666"); !ok || word != "zoom" {
		t.Fatalf("ByDiceCode(66666) = %q, %v", word, ok)
	}
	if code, ok := DiceCodeByWord("abacus"); !ok || code != "11111" {
		t.Fatalf("DiceCodeByWord(abacus) = %q, %v", code, ok)
	}
	for _, invalid := range []string{"", "00000", "1111", "77777"} {
		if _, ok := ByDiceCode(invalid); ok {
			t.Errorf("ByDiceCode(%q) unexpectedly found a mapping", invalid)
		}
	}
}

func TestSecureGenerationAndEntropy(t *testing.T) {
	generated, err := RandomSecure(8)
	if err != nil || len(generated) != 8 {
		t.Fatalf("RandomSecure(8) = %d words, %v", len(generated), err)
	}
	passphrase, err := PassphraseSecure(3, "", false)
	if err != nil || strings.Contains(passphrase, "-") || len(strings.Fields(passphrase)) != 1 {
		t.Fatalf("PassphraseSecure with no separator = %q, %v", passphrase, err)
	}
	if got := EntropyBits(6); got < 77 || got > 78 {
		t.Errorf("EntropyBits(6) = %.2f, expected approximately 77.5", got)
	}
}
