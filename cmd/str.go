// Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/aminshahid573/gen/internal/ui"
	"os"
	"strconv"
	"strings"

	"github.com/mr-tron/base58"
	"github.com/spf13/cobra"
)

var strCmd = &cobra.Command{
	Use:   "str",
	Short: "Generate a cryptographically secure random string",
	Long: `Generate a cryptographically secure random string using crypto/rand.

Character set flags (combinable):
  --uppercase     Include A-Z (default true)
  --lowercase     Include a-z (default true)
  --digits        Include 0-9 (default true)
  --symbols       Include special characters: !@#$%^&*()-_=+[]{}|;:,.<>?
  --no-ambiguous  Exclude visually confusing characters: 0O 1lI B8

Encoding transforms the raw bytes directly (ignores charset flags):
  raw       Random characters from the charset (default)
  hex       Lowercase hexadecimal
  base64    Standard base64
  base64url URL-safe base64 (no +/=)
  base58    Base58 (no ambiguous characters)

Pattern mode (overrides --length and charset flags):
  A  uppercase letter    a  lowercase letter
  0  digit               *  any alphanumeric
  !  symbol              -  literal passthrough

Examples:
  gen str
  gen str --length 32 --symbols
  gen str --length 16 --no-ambiguous
  gen str --encoding base64url
  gen str --encoding hex --count 5
  gen str --prefix "sk_live_" --length 24
  gen str --pattern "AAA-000-aaa"
  gen str --charset "abcdef0123456789" --length 20`,

	Run: func(cmd *cobra.Command, args []string) {
		length, _ := cmd.Flags().GetInt("length")
		enc, _ := cmd.Flags().GetString("encoding")
		count, _ := cmd.Flags().GetInt("count")
		prefix, _ := cmd.Flags().GetString("prefix")
		suffix, _ := cmd.Flags().GetString("suffix")
		pattern, _ := cmd.Flags().GetString("pattern")
		customCS, _ := cmd.Flags().GetString("charset")
		upper, _ := cmd.Flags().GetBool("uppercase")
		lower, _ := cmd.Flags().GetBool("lowercase")
		digits, _ := cmd.Flags().GetBool("digits")
		symbols, _ := cmd.Flags().GetBool("symbols")
		noAmbig, _ := cmd.Flags().GetBool("no-ambiguous")

		if count < 1 {
			fatalf("--count must be at least 1")
		}

		// build charset (only used for raw and pattern modes)
		charset := buildCharset(customCS, upper, lower, digits, symbols, noAmbig)
		if charset == "" {
			fatalf("charset is empty — enable at least one of: --uppercase, --lowercase, --digits, --symbols")
		}

		rows := make([][]string, 0, count)

		for range count {
			var out string
			var err error

			if pattern != "" {
				out, err = generatePattern(pattern, charset)
			} else {
				out, err = generateStr(length, enc, charset)
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "error generating string: %v\n", err)
				os.Exit(1)
			}

			out = prefix + out + suffix

			displayLen := strconv.Itoa(len(out))
			displayEnc := enc
			if pattern != "" {
				displayEnc = "pattern"
			}
			rows = append(rows, []string{displayEnc, displayLen, out})
		}

		fmt.Println(ui.RenderTable([]string{"Encoding", "Length", "Output"}, rows))
	},
}

// buildCharset assembles the character pool from flags.
func buildCharset(custom string, upper, lower, digits, symbols, noAmbiguous bool) string {
	if custom != "" {
		return custom
	}

	const (
		upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowerChars  = "abcdefghijklmnopqrstuvwxyz"
		digitChars  = "0123456789"
		symbolChars = "!@#$%^&*()-_=+[]{}|;:,.<>?"
		ambiguous   = "0O1lIB8"
	)

	var sb strings.Builder
	if upper {
		sb.WriteString(upperChars)
	}
	if lower {
		sb.WriteString(lowerChars)
	}
	if digits {
		sb.WriteString(digitChars)
	}
	if symbols {
		sb.WriteString(symbolChars)
	}

	cs := sb.String()

	if noAmbiguous {
		filtered := strings.Builder{}
		for _, c := range cs {
			if !strings.ContainsRune(ambiguous, c) {
				filtered.WriteRune(c)
			}
		}
		return filtered.String()
	}

	return cs
}

// generateStr produces a random string of the given length in the requested encoding.
func generateStr(length int, enc, charset string) (string, error) {
	switch strings.ToLower(enc) {
	case "raw":
		return randomFromCharset(charset, length)

	case "hex":
		b := make([]byte, length)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		return hex.EncodeToString(b)[:length], nil

	case "base64":
		b := make([]byte, length)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(b), nil

	case "base64url":
		b := make([]byte, length)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		return base64.RawURLEncoding.EncodeToString(b), nil

	case "base58":
		b := make([]byte, length)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		return base58.Encode(b), nil

	default:
		return "", fmt.Errorf("unknown encoding %q (choose: raw, hex, base64, base64url, base58)", enc)
	}
}

// generatePattern builds a string matching a pattern template.
//
//	A = uppercase  a = lowercase  0 = digit  * = alphanumeric  ! = symbol
//	any other character is passed through as a literal
func generatePattern(pattern, charset string) (string, error) {
	const (
		upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowerChars  = "abcdefghijklmnopqrstuvwxyz"
		digitChars  = "0123456789"
		symbolChars = "!@#$%^&*()-_=+[]{}|;:,.<>?"
		alphaNum    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	)

	var sb strings.Builder
	for _, ch := range pattern {
		var pool string
		switch ch {
		case 'A':
			pool = upperChars
		case 'a':
			pool = lowerChars
		case '0':
			pool = digitChars
		case '*':
			pool = alphaNum
		case '!':
			pool = symbolChars
		default:
			sb.WriteRune(ch) // literal
			continue
		}
		c, err := randomChar(pool)
		if err != nil {
			return "", err
		}
		sb.WriteByte(c)
	}
	return sb.String(), nil
}

// randomFromCharset picks `n` random characters from the given pool
// using rejection sampling to avoid modulo bias.
func randomFromCharset(charset string, n int) (string, error) {
	if len(charset) == 0 {
		return "", fmt.Errorf("charset is empty")
	}
	sb := strings.Builder{}
	sb.Grow(n)
	for range n {
		c, err := randomChar(charset)
		if err != nil {
			return "", err
		}
		sb.WriteByte(c)
	}
	return sb.String(), nil
}

// randomChar picks one random byte from the charset using rejection sampling.
func randomChar(charset string) (byte, error) {
	size := len(charset)
	limit := 256 - (256 % size) // rejection ceiling to eliminate modulo bias
	b := make([]byte, 1)
	for {
		if _, err := rand.Read(b); err != nil {
			return 0, err
		}
		if int(b[0]) < limit {
			return charset[int(b[0])%size], nil
		}
	}
}

func init() {
	rootCmd.AddCommand(strCmd)
	strCmd.Flags().IntP("length", "l", 32, "Length of the generated string")
	strCmd.Flags().StringP("encoding", "e", "raw", "Encoding: raw, hex, base64, base64url, base58")
	strCmd.Flags().IntP("count", "c", 1, "Number of strings to generate")
	strCmd.Flags().StringP("prefix", "p", "", "Prepend a fixed string to output")
	strCmd.Flags().StringP("suffix", "x", "", "Append a fixed string to output")
	strCmd.Flags().StringP("pattern", "P", "", `Pattern template (A=upper a=lower 0=digit *=alnum !=symbol)`)
	strCmd.Flags().StringP("charset", "C", "", "Custom character pool (overrides charset flags)")
	strCmd.Flags().Bool("uppercase", true, "Include uppercase letters A-Z")
	strCmd.Flags().Bool("lowercase", true, "Include lowercase letters a-z")
	strCmd.Flags().Bool("digits", true, "Include digits 0-9")
	strCmd.Flags().Bool("symbols", false, "Include symbols !@#$%^&*()-_=+[]{}|;:,.<>?")
	strCmd.Flags().Bool("no-ambiguous", false, "Exclude visually ambiguous characters (0O 1lI B8)")
}
