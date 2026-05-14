package cmd

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateUUID(t *testing.T) {
	tests := []struct {
		name    string
		version int
		uName   string
		ns      string
		wantErr bool
	}{
		{"v1", 1, "", "", false},
		{"v3 with name", 3, "example.com", "dns", false},
		{"v3 missing name", 3, "", "dns", true},
		{"v4", 4, "", "", false},
		{"v5 with name", 5, "https://example.com", "url", false},
		{"v5 missing name", 5, "", "url", true},
		{"v6", 6, "", "", false},
		{"v7", 7, "", "", false},
		{"unsupported version", 99, "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := generateUUID(tt.version, tt.uName, tt.ns)
			if (err != nil) != tt.wantErr {
				t.Fatalf("generateUUID() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && id == (uuid.UUID{}) {
				t.Fatal("expected non-zero UUID")
			}
		})
	}
}

func TestGenerateUUID_V3Deterministic(t *testing.T) {
	a, _ := generateUUID(3, "example.com", "dns")
	b, _ := generateUUID(3, "example.com", "dns")
	if a != b {
		t.Fatal("v3 UUID should be deterministic for same input")
	}
}

func TestGenerateUUID_V5Deterministic(t *testing.T) {
	a, _ := generateUUID(5, "https://example.com", "url")
	b, _ := generateUUID(5, "https://example.com", "url")
	if a != b {
		t.Fatal("v5 UUID should be deterministic for same input")
	}
}

func TestFormatUUID(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	tests := []struct {
		format  string
		check   func(string) bool
		wantErr bool
	}{
		{"standard", func(s string) bool { return strings.Count(s, "-") == 4 }, false},
		{"compact", func(s string) bool { return len(s) == 32 && !strings.Contains(s, "-") }, false},
		{"upper", func(s string) bool { return s == strings.ToUpper(s) }, false},
		{"urn", func(s string) bool { return strings.HasPrefix(s, "urn:uuid:") }, false},
		{"base58", func(s string) bool { return len(s) > 0 }, false},
		{"invalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			out, err := formatUUID(id, tt.format)
			if (err != nil) != tt.wantErr {
				t.Fatalf("formatUUID() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && !tt.check(out) {
				t.Fatalf("formatUUID(%q) = %q, failed check", tt.format, out)
			}
		})
	}
}

func TestResolveFormat(t *testing.T) {
	validFormats := map[string]bool{
		"standard": true, "compact": true, "upper": true,
		"urn": true, "base58": true,
	}

	// non-random passthrough
	for f := range validFormats {
		if got := resolveFormat(f); got != f {
			t.Errorf("resolveFormat(%q) = %q, want passthrough", f, got)
		}
	}

	// random always resolves to a known format
	for range 50 {
		got := resolveFormat("random")
		if !validFormats[got] {
			t.Errorf("resolveFormat(random) returned unknown format %q", got)
		}
	}
}

func TestResolveNamespace(t *testing.T) {
	tests := []struct {
		input string
		want  uuid.UUID
	}{
		{"dns", uuid.NameSpaceDNS},
		{"url", uuid.NameSpaceURL},
		{"oid", uuid.NameSpaceOID},
		{"x500", uuid.NameSpaceX500},
		{"unknown", uuid.NameSpaceDNS}, // fallback
		{"550e8400-e29b-41d4-a716-446655440000",
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := resolveNamespace(tt.input); got != tt.want {
				t.Errorf("resolveNamespace(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractTimestamp(t *testing.T) {
	tests := []struct {
		name    string
		version int
		want    string // substring to check
	}{
		{"v1", 1, "T"},   // RFC3339 contains T
		{"v6", 6, "T"},
		{"v7", 7, "T"},
		{"v3", 3, "n/a"},
		{"v5", 5, "n/a"},
		{"v4", 4, "n/a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id uuid.UUID
			switch tt.version {
			case 1:
				id, _ = uuid.NewUUID()
			case 3:
				id = uuid.NewMD5(uuid.NameSpaceDNS, []byte("test"))
			case 4:
				id, _ = uuid.NewRandom()
			case 5:
				id = uuid.NewSHA1(uuid.NameSpaceDNS, []byte("test"))
			case 6:
				id, _ = uuid.NewV6()
			case 7:
				id, _ = uuid.NewV7()
			}

			ts := extractTimestamp(id)
			if !strings.Contains(ts, tt.want) {
				t.Errorf("extractTimestamp() = %q, want substring %q", ts, tt.want)
			}
		})
	}
}
