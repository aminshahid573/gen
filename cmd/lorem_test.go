package cmd

import (
	"math/rand"
	"strings"
	"testing"
)

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"HELLO", "HELLO"},
		{"", ""},
		{"a", "A"},
	}
	for _, tt := range tests {
		if got := capitalize(tt.input); got != tt.want {
			t.Errorf("capitalize(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGenerateLoremSentence(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for range 20 {
		s := generateLoremSentence(rng)
		if !strings.HasSuffix(s, ".") {
			t.Errorf("sentence %q does not end with period", s)
		}
		if s[0] != strings.ToUpper(s[:1])[0] {
			t.Errorf("sentence %q does not start with uppercase", s)
		}
		words := strings.Fields(s)
		if len(words) < 6 {
			t.Errorf("sentence has only %d words, want at least 6", len(words))
		}
	}
}

func TestGenerateLoremWords_Count(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for _, n := range []int{1, 5, 10, 50} {
		out := generateLoremWords(rng, n)
		words := strings.Fields(out)
		if len(words) != n {
			t.Errorf("generateLoremWords(%d) returned %d words", n, len(words))
		}
	}
}

func TestGenerateLoremParagraph_SentenceCount(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for _, n := range []int{1, 3, 5} {
		p := generateLoremParagraph(rng, n)
		// count periods as sentence terminators
		count := strings.Count(p, ".")
		if count != n {
			t.Errorf("generateLoremParagraph(%d) has %d sentences", n, count)
		}
	}
}

func TestGenerateLoremWords_UsesWordPool(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	out := generateLoremWords(rng, 100)
	// first word is capitalized — lowercase it for pool check
	words := strings.Fields(strings.ToLower(out))
	poolSet := map[string]bool{}
	for _, w := range loremWords {
		poolSet[w] = true
	}
	for _, w := range words {
		if !poolSet[w] {
			t.Errorf("word %q not found in lorem word pool", w)
		}
	}
}
