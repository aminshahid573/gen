package cmd

import (
	"strings"
	"testing"
)

func TestBuildCharset(t *testing.T) {
	tests := []struct {
		name       string
		custom     string
		upper      bool
		lower      bool
		digits     bool
		symbols    bool
		noAmbig    bool
		wantSubstr string
		wantEmpty  bool
	}{
		{"all enabled", "", true, true, true, false, false, "A", false},
		{"custom overrides flags", "XYZ", false, false, false, false, false, "XYZ", false},
		{"only digits", "", false, false, true, false, false, "0", false},
		{"no ambiguous removes 0", "", false, false, true, false, true, "", false},
		{"nothing enabled", "", false, false, false, false, false, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := buildCharset(tt.custom, tt.upper, tt.lower, tt.digits, tt.symbols, tt.noAmbig)
			if tt.wantEmpty && cs != "" {
				t.Errorf("expected empty charset, got %q", cs)
			}
			if tt.wantSubstr != "" && !strings.Contains(cs, tt.wantSubstr) {
				t.Errorf("charset %q missing expected substring %q", cs, tt.wantSubstr)
			}
			if tt.noAmbig {
				for _, c := range "0O1lIB8" {
					if strings.ContainsRune(cs, c) {
						t.Errorf("no-ambiguous charset still contains %c", c)
					}
				}
			}
		})
	}
}

func TestGenerateStr_Lengths(t *testing.T) {
	charset := "abcdefghijklmnopqrstuvwxyz"

	tests := []struct {
		enc      string
		length   int
		exactLen bool // exact output length matches input length
	}{
		{"raw", 16, true},
		{"raw", 32, true},
		{"hex", 16, false}, // hex doubles byte count
		{"base64", 16, false},
		{"base64url", 16, false},
		{"base58", 16, false},
	}

	for _, tt := range tests {
		t.Run(tt.enc, func(t *testing.T) {
			out, err := generateStr(tt.length, tt.enc, charset)
			if err != nil {
				t.Fatalf("generateStr() error: %v", err)
			}
			if out == "" {
				t.Fatal("output should not be empty")
			}
			if tt.exactLen && len(out) != tt.length {
				t.Errorf("len(%q) = %d, want %d", out, len(out), tt.length)
			}
		})
	}
}

func TestGenerateStr_InvalidEncoding(t *testing.T) {
	_, err := generateStr(16, "invalid", "abc")
	if err == nil {
		t.Fatal("expected error for invalid encoding")
	}
}

func TestGenerateStr_RawUsesCharset(t *testing.T) {
	charset := "XYZ"
	out, err := generateStr(100, "raw", charset)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range out {
		if !strings.ContainsRune(charset, c) {
			t.Errorf("output char %q not in charset %q", c, charset)
		}
	}
}

func TestGeneratePattern(t *testing.T) {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	tests := []struct {
		name    string
		pattern string
		length  int // expected output length
	}{
		{"uppercase only", "AAA", 3},
		{"lowercase only", "aaa", 3},
		{"digits only", "000", 3},
		{"alphanumeric", "***", 3},
		{"literals pass through", "A-0-a", 5},
		{"mixed", "AA-00-aa", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := generatePattern(tt.pattern, charset)
			if err != nil {
				t.Fatalf("generatePattern() error: %v", err)
			}
			if len(out) != tt.length {
				t.Errorf("len(%q) = %d, want %d", out, len(out), tt.length)
			}
		})
	}
}

func TestRandomChar(t *testing.T) {
	charset := "abc"
	for range 100 {
		c, err := randomChar(charset)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.ContainsRune(charset, rune(c)) {
			t.Errorf("randomChar returned %q not in charset %q", c, charset)
		}
	}
}

func TestRandomFromCharset(t *testing.T) {
	charset := "abcdef"
	out, err := randomFromCharset(charset, 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 20 {
		t.Errorf("len = %d, want 20", len(out))
	}
	for _, c := range out {
		if !strings.ContainsRune(charset, c) {
			t.Errorf("char %q not in charset", c)
		}
	}
}
