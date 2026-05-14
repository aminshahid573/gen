package cmd

import (
	"strings"
	"testing"
)

func TestGeneratePassword_Length(t *testing.T) {
	charset := buildCharset("", true, true, true, false, false)
	pass, err := generatePassword(16, charset, true, true, true, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(pass) != 16 {
		t.Errorf("len = %d, want 16", len(pass))
	}
}

func TestGeneratePassword_GuaranteedClasses(t *testing.T) {
	tests := []struct {
		name    string
		upper   bool
		lower   bool
		digits  bool
		symbols bool
		checkFn func(string) bool
	}{
		{
			name: "uppercase guaranteed",
			upper: true, lower: false, digits: false, symbols: false,
			checkFn: func(s string) bool {
				return strings.ContainsAny(s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
			},
		},
		{
			name: "digit guaranteed",
			upper: false, lower: false, digits: true, symbols: false,
			checkFn: func(s string) bool {
				return strings.ContainsAny(s, "0123456789")
			},
		},
		{
			name: "symbol guaranteed",
			upper: false, lower: false, digits: false, symbols: true,
			checkFn: func(s string) bool {
				return strings.ContainsAny(s, "!@#$%^&*()-_=+[]{}|;:,.<>?")
			},
		},
		{
			name: "all classes guaranteed",
			upper: true, lower: true, digits: true, symbols: true,
			checkFn: func(s string) bool {
				hasUpper   := strings.ContainsAny(s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
				hasLower   := strings.ContainsAny(s, "abcdefghijklmnopqrstuvwxyz")
				hasDigit   := strings.ContainsAny(s, "0123456789")
				hasSymbol  := strings.ContainsAny(s, "!@#$%^&*()-_=+[]{}|;:,.<>?")
				return hasUpper && hasLower && hasDigit && hasSymbol
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charset := buildCharset("", tt.upper, tt.lower, tt.digits, tt.symbols, false)
			// run multiple times to catch shuffle issues
			for i := range 20 {
				pass, err := generatePassword(16, charset, tt.upper, tt.lower, tt.digits, tt.symbols, false)
				if err != nil {
					t.Fatalf("iteration %d: %v", i, err)
				}
				if !tt.checkFn(pass) {
					t.Errorf("iteration %d: password %q failed character class check", i, pass)
				}
			}
		})
	}
}

func TestGeneratePassword_NoAmbiguous(t *testing.T) {
	charset := buildCharset("", true, true, true, false, true)
	for range 20 {
		pass, err := generatePassword(32, charset, true, true, true, false, true)
		if err != nil {
			t.Fatal(err)
		}
		for _, c := range "0O1lIB8" {
			if strings.ContainsRune(pass, c) {
				t.Errorf("password %q contains ambiguous char %c", pass, c)
			}
		}
	}
}

func TestShuffleBytes_PreservesContent(t *testing.T) {
	original := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	buf := make([]byte, len(original))
	copy(buf, original)

	if err := shuffleBytes(buf); err != nil {
		t.Fatal(err)
	}

	if len(buf) != len(original) {
		t.Errorf("length changed after shuffle: %d != %d", len(buf), len(original))
	}

	// same bytes, different order (check frequency)
	freq := map[byte]int{}
	for _, b := range original {
		freq[b]++
	}
	for _, b := range buf {
		freq[b]--
	}
	for b, count := range freq {
		if count != 0 {
			t.Errorf("byte %q frequency changed after shuffle", b)
		}
	}
}
