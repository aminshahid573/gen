package cmd

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	fixed := time.Date(2026, 5, 13, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		format  string
		layout  string
		check   func(string) bool
		wantErr bool
	}{
		{"unix", "", func(s string) bool {
			n, err := strconv.ParseInt(s, 10, 64)
			return err == nil && n == fixed.Unix()
		}, false},
		{"unixmilli", "", func(s string) bool {
			n, err := strconv.ParseInt(s, 10, 64)
			return err == nil && n == fixed.UnixMilli()
		}, false},
		{"unixnano", "", func(s string) bool {
			n, err := strconv.ParseInt(s, 10, 64)
			return err == nil && n == fixed.UnixNano()
		}, false},
		{"rfc3339", "", func(s string) bool {
			_, err := time.Parse(time.RFC3339, s)
			return err == nil
		}, false},
		{"rfc3339nano", "", func(s string) bool {
			return strings.Contains(s, "T")
		}, false},
		{"iso8601", "", func(s string) bool {
			_, err := time.Parse(time.RFC3339, s)
			return err == nil
		}, false},
		{"http", "", func(s string) bool {
			_, err := time.Parse(time.RFC1123, s)
			return err == nil
		}, false},
		{"date", "", func(s string) bool {
			return s == "2026-05-13"
		}, false},
		{"datetime", "", func(s string) bool {
			return strings.Contains(s, "2026-05-13")
		}, false},
		{"kitchen", "", func(s string) bool {
			return strings.Contains(s, "AM") || strings.Contains(s, "PM")
		}, false},
		{"custom", "02 Jan 2006", func(s string) bool {
			return s == "13 May 2026"
		}, false},
		{"custom no layout", "", nil, true},
		{"invalid", "", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			out, err := formatTime(fixed, tt.format, tt.layout)
			if (err != nil) != tt.wantErr {
				t.Fatalf("formatTime(%q) error = %v, wantErr = %v", tt.format, err, tt.wantErr)
			}
			if err == nil && !tt.check(out) {
				t.Errorf("formatTime(%q) = %q, failed check", tt.format, out)
			}
		})
	}
}

func TestParseTimeRange_Defaults(t *testing.T) {
	from, to := parseTimeRange("", "")
	if from.Unix() != 0 {
		t.Errorf("default from should be unix epoch, got %v", from)
	}
	if to.Unix() <= 0 {
		t.Errorf("default to should be approximately now, got %v", to)
	}
}

func TestParseTimeRange_DateFormat(t *testing.T) {
	from, to := parseTimeRange("2020-01-01", "2025-12-31")
	if from.Year() != 2020 {
		t.Errorf("from year = %d, want 2020", from.Year())
	}
	if to.Year() != 2025 {
		t.Errorf("to year = %d, want 2025", to.Year())
	}
}

func TestParseTimeRange_UnixTimestamp(t *testing.T) {
	from, _ := parseTimeRange("0", "")
	if from.Unix() != 0 {
		t.Errorf("unix 0 should parse to epoch, got %v", from)
	}
}
