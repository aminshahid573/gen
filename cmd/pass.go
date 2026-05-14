// Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
package cmd

import (
	"crypto/rand"
	"fmt"
	"gen/internal/ui"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var passCmd = &cobra.Command{
	Use:   "pass",
	Short: "Generate a strong random password",
	Long: `Generate a cryptographically secure password with guaranteed character class coverage.

Unlike gen str, gen pass ensures at least one character from every enabled class
appears in the output — making it suitable for password policy requirements.

Examples:
  gen pass
  gen pass --length 24
  gen pass --length 32 --symbols
  gen pass --no-ambiguous
  gen pass --count 5
  gen pass --length 20 --symbols --no-ambiguous --count 3`,

	Run: func(cmd *cobra.Command, args []string) {
		length, _ := cmd.Flags().GetInt("length")
		upper, _ := cmd.Flags().GetBool("uppercase")
		lower, _ := cmd.Flags().GetBool("lowercase")
		digits, _ := cmd.Flags().GetBool("digits")
		symbols, _ := cmd.Flags().GetBool("symbols")
		noAmbig, _ := cmd.Flags().GetBool("no-ambiguous")
		count, _ := cmd.Flags().GetInt("count")

		if count < 1 {
			fatalf("--count must be at least 1")
		}
		if length < 1 {
			fatalf("--length must be at least 1")
		}

		// at least one class must be enabled
		if !upper && !lower && !digits && !symbols {
			fatalf("at least one character class must be enabled")
		}

		// minimum length must fit one char per enabled class
		minLen := 0
		if upper {
			minLen++
		}
		if lower {
			minLen++
		}
		if digits {
			minLen++
		}
		if symbols {
			minLen++
		}
		if length < minLen {
			fatalf("--length %d is too short to satisfy all enabled character classes (min %d)", length, minLen)
		}

		charset := buildCharset("", upper, lower, digits, symbols, noAmbig)

		rows := make([][]string, 0, count)
		for range count {
			pass, err := generatePassword(length, charset, upper, lower, digits, symbols, noAmbig)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error generating password: %v\n", err)
				os.Exit(1)
			}
			rows = append(rows, []string{strconv.Itoa(length), pass})
		}

		fmt.Println(ui.RenderTable([]string{"Length", "Password"}, rows))
	},
}

// generatePassword guarantees at least one char from each enabled class,
// then fills the rest from the full charset and shuffles with crypto/rand.
func generatePassword(length int, charset string, upper, lower, digits, symbols, noAmbig bool) (string, error) {
	const (
		upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowerChars  = "abcdefghijklmnopqrstuvwxyz"
		digitChars  = "0123456789"
		symbolChars = "!@#$%^&*()-_=+[]{}|;:,.<>?"
		ambiguous   = "0O1lIB8"
	)

	filter := func(s string) string {
		if !noAmbig {
			return s
		}
		var b strings.Builder
		for _, c := range s {
			if !strings.ContainsRune(ambiguous, c) {
				b.WriteRune(c)
			}
		}
		return b.String()
	}

	buf := make([]byte, 0, length)

	// guarantee one from each enabled class first
	if upper {
		c, err := randomChar(filter(upperChars))
		if err != nil {
			return "", err
		}
		buf = append(buf, c)
	}
	if lower {
		c, err := randomChar(filter(lowerChars))
		if err != nil {
			return "", err
		}
		buf = append(buf, c)
	}
	if digits {
		c, err := randomChar(filter(digitChars))
		if err != nil {
			return "", err
		}
		buf = append(buf, c)
	}
	if symbols {
		c, err := randomChar(filter(symbolChars))
		if err != nil {
			return "", err
		}
		buf = append(buf, c)
	}

	// fill remaining slots from the full charset
	for len(buf) < length {
		c, err := randomChar(charset)
		if err != nil {
			return "", err
		}
		buf = append(buf, c)
	}

	// shuffle so the guaranteed chars aren't always at the front
	if err := shuffleBytes(buf); err != nil {
		return "", err
	}

	return string(buf), nil
}

// shuffleBytes performs a Fisher-Yates shuffle using crypto/rand.
func shuffleBytes(s []byte) error {
	for i := len(s) - 1; i > 0; i-- {
		j, err := cryptoRandInt(i + 1)
		if err != nil {
			return err
		}
		s[i], s[j] = s[j], s[i]
	}
	return nil
}

// cryptoRandInt returns a cryptographically random int in [0, max).
func cryptoRandInt(max int) (int, error) {
	limit := 256 - (256 % max)
	b := make([]byte, 1)
	for {
		if _, err := rand.Read(b); err != nil {
			return 0, err
		}
		if int(b[0]) < limit {
			return int(b[0]) % max, nil
		}
	}
}

func init() {
	rootCmd.AddCommand(passCmd)
	passCmd.Flags().IntP("length", "l", 16, "Password length")
	passCmd.Flags().IntP("count", "c", 1, "Number of passwords to generate")
	passCmd.Flags().Bool("uppercase", true, "Include uppercase letters A-Z")
	passCmd.Flags().Bool("lowercase", true, "Include lowercase letters a-z")
	passCmd.Flags().Bool("digits", true, "Include digits 0-9")
	passCmd.Flags().Bool("symbols", false, "Include symbols !@#$%^&*()-_=+[]{}|;:,.<>?")
	passCmd.Flags().Bool("no-ambiguous", false, "Exclude visually ambiguous characters (0O 1lI B8)")
}
