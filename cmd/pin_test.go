package cmd

import (
	"strings"
	"testing"
)

func TestGeneratePIN_Length(t *testing.T) {
	for _, length := range []int{4, 6, 8, 12} {
		pin, err := generatePIN(length, 0, "")
		if err != nil {
			t.Fatalf("length %d: %v", length, err)
		}
		if len(pin) != length {
			t.Errorf("length %d: got pin %q with len %d", length, pin, len(pin))
		}
	}
}

func TestGeneratePIN_OnlyDigits(t *testing.T) {
	pin, err := generatePIN(20, 0, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range pin {
		if c < '0' || c > '9' {
			t.Errorf("non-digit character %q in PIN %q", c, pin)
		}
	}
}

func TestGeneratePIN_Grouping(t *testing.T) {
	tests := []struct {
		length    int
		group     int
		separator string
		wantParts int
	}{
		{6, 3, "-", 2}, // 123-456
		{9, 3, "-", 3}, // 123-456-789
		{8, 4, " ", 2}, // 1234 5678
		{6, 2, ".", 3}, // 12.34.56
	}

	for _, tt := range tests {
		pin, err := generatePIN(tt.length, tt.group, tt.separator)
		if err != nil {
			t.Fatal(err)
		}
		parts := strings.Split(pin, tt.separator)
		if len(parts) != tt.wantParts {
			t.Errorf("grouping(%d,%d,%q): got %d parts in %q, want %d",
				tt.length, tt.group, tt.separator, len(parts), pin, tt.wantParts)
		}
		for _, p := range parts {
			if len(p) != tt.group {
				t.Errorf("part %q has len %d, want %d", p, len(p), tt.group)
			}
		}
	}
}

func TestGeneratePIN_NoGroupingWhenGroupZero(t *testing.T) {
	pin, err := generatePIN(6, 0, "-")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(pin, "-") {
		t.Errorf("expected no separator when group=0, got %q", pin)
	}
}
