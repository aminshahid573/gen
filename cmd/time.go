// Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
package cmd

import (
	"fmt"
	"gen/internal/ui"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Generate or convert timestamps in various formats",
	Long: `Output the current time (or a random timestamp) in a variety of formats.

Formats:
  unix          Seconds since Unix epoch (integer)
  unixmilli     Milliseconds since Unix epoch
  unixnano      Nanoseconds since Unix epoch
  rfc3339       2006-01-02T15:04:05Z07:00
  rfc3339nano   2006-01-02T15:04:05.999999999Z07:00
  iso8601       Same as rfc3339
  http          Mon, 02 Jan 2006 15:04:05 GMT  (HTTP-date)
  date          2006-01-02
  datetime      2006-01-02 15:04:05
  kitchen       3:04PM
  custom        Use Go time layout via --layout

Examples:
  gen time
  gen time --format rfc3339
  gen time --format unix
  gen time --format date --utc
  gen time --format all
  gen time --random
  gen time --random --from 2020-01-01 --to 2025-12-31
  gen time --count 5 --random
  gen time --format custom --layout "02 Jan 2006"`,

	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		layout, _ := cmd.Flags().GetString("layout")
		utc, _ := cmd.Flags().GetBool("utc")
		random, _ := cmd.Flags().GetBool("random")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		count, _ := cmd.Flags().GetInt("count")

		if count < 1 {
			fatalf("--count must be at least 1")
		}

		fromTime, toTime := parseTimeRange(from, to)

		// "all" mode — show every format for current/random time, ignores count
		if strings.ToLower(format) == "all" {
			t := pickTime(random, fromTime, toTime, utc)
			printAllFormats(t)
			return
		}

		rows := make([][]string, 0, count)
		for range count {
			t := pickTime(random, fromTime, toTime, utc)
			out, err := formatTime(t, format, layout)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			rows = append(rows, []string{strings.ToLower(format), out})
		}

		fmt.Println(ui.RenderTable([]string{"Format", "Value"}, rows))
	},
}

func pickTime(random bool, from, to time.Time, utc bool) time.Time {
	var t time.Time
	if random {
		delta := to.Unix() - from.Unix()
		if delta <= 0 {
			delta = 1
		}
		//nolint:gosec // math/rand is fine for non-crypto timestamp generation
		t = time.Unix(from.Unix()+rand.Int63n(delta), 0)
	} else {
		t = time.Now()
	}
	if utc {
		return t.UTC()
	}
	return t
}

func formatTime(t time.Time, format, customLayout string) (string, error) {
	switch strings.ToLower(format) {
	case "unix":
		return strconv.FormatInt(t.Unix(), 10), nil
	case "unixmilli":
		return strconv.FormatInt(t.UnixMilli(), 10), nil
	case "unixnano":
		return strconv.FormatInt(t.UnixNano(), 10), nil
	case "rfc3339", "iso8601":
		return t.Format(time.RFC3339), nil
	case "rfc3339nano":
		return t.Format(time.RFC3339Nano), nil
	case "http":
		return t.UTC().Format(time.RFC1123), nil
	case "date":
		return t.Format("2006-01-02"), nil
	case "datetime":
		return t.Format("2006-01-02 15:04:05"), nil
	case "kitchen":
		return t.Format(time.Kitchen), nil
	case "custom":
		if customLayout == "" {
			return "", fmt.Errorf("--layout is required with --format custom")
		}
		return t.Format(customLayout), nil
	default:
		return "", fmt.Errorf("unknown format %q — run 'gen time --format all' to see options", format)
	}
}

func printAllFormats(t time.Time) {
	rows := [][]string{
		{"unix", strconv.FormatInt(t.Unix(), 10)},
		{"unixmilli", strconv.FormatInt(t.UnixMilli(), 10)},
		{"unixnano", strconv.FormatInt(t.UnixNano(), 10)},
		{"rfc3339", t.Format(time.RFC3339)},
		{"rfc3339nano", t.Format(time.RFC3339Nano)},
		{"iso8601", t.Format(time.RFC3339)},
		{"http", t.UTC().Format(time.RFC1123)},
		{"date", t.Format("2006-01-02")},
		{"datetime", t.Format("2006-01-02 15:04:05")},
		{"kitchen", t.Format(time.Kitchen)},
	}
	fmt.Println(ui.RenderTable([]string{"Format", "Value"}, rows))
}

// parseTimeRange parses --from and --to into time.Time.
// Accepts: YYYY-MM-DD, YYYY-MM-DDTHH:MM:SS, or unix timestamps.
func parseTimeRange(from, to string) (time.Time, time.Time) {
	parseOne := func(s string, fallback time.Time) time.Time {
		if s == "" {
			return fallback
		}
		// try unix timestamp first
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			return time.Unix(n, 0)
		}
		// try date formats
		for _, layout := range []string{"2006-01-02", "2006-01-02T15:04:05", time.RFC3339} {
			if t, err := time.Parse(layout, s); err == nil {
				return t
			}
		}
		fatalf("could not parse time %q — use YYYY-MM-DD, RFC3339, or a unix timestamp", s)
		return fallback
	}

	defaultFrom := time.Unix(0, 0)
	defaultTo := time.Now()

	return parseOne(from, defaultFrom), parseOne(to, defaultTo)
}

func init() {
	rootCmd.AddCommand(timeCmd)
	timeCmd.Flags().StringP("format", "f", "rfc3339", "Output format (use 'all' to show every format)")
	timeCmd.Flags().StringP("layout", "l", "", "Custom Go time layout string (used with --format custom)")
	timeCmd.Flags().BoolP("utc", "u", false, "Force UTC output")
	timeCmd.Flags().BoolP("random", "r", false, "Generate a random timestamp instead of current time")
	timeCmd.Flags().String("from", "", "Random range start: YYYY-MM-DD, RFC3339, or unix timestamp")
	timeCmd.Flags().String("to", "", "Random range end:   YYYY-MM-DD, RFC3339, or unix timestamp")
	timeCmd.Flags().IntP("count", "c", 1, "Number of timestamps to generate")
}
