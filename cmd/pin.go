// Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
package cmd

import (
	"fmt"
	"github.com/aminshahid573/gen/internal/ui"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Generate a numeric PIN",
	Long: `Generate a cryptographically secure numeric PIN.

Optionally group digits with a separator for readability.

Examples:
  gen pin
  gen pin --length 8
  gen pin --length 6 --group 3 --separator "-"
  gen pin --count 5
  gen pin --length 8 --group 4 --separator " "`,

	Run: func(cmd *cobra.Command, args []string) {
		length, _ := cmd.Flags().GetInt("length")
		count, _ := cmd.Flags().GetInt("count")
		group, _ := cmd.Flags().GetInt("group")
		separator, _ := cmd.Flags().GetString("separator")

		if count < 1 {
			fatalf("--count must be at least 1")
		}
		if length < 1 {
			fatalf("--length must be at least 1")
		}
		if group < 0 {
			fatalf("--group must be 0 or greater")
		}

		rows := make([][]string, 0, count)
		for range count {
			pin, err := generatePIN(length, group, separator)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error generating PIN: %v\n", err)
				os.Exit(1)
			}
			rows = append(rows, []string{fmt.Sprintf("%d", length), pin})
		}

		fmt.Println(ui.RenderTable([]string{"Length", "PIN"}, rows))
	},
}

func generatePIN(length, group int, separator string) (string, error) {
	const digitChars = "0123456789"

	digits := make([]byte, length)
	for i := range length {
		c, err := randomChar(digitChars)
		if err != nil {
			return "", err
		}
		digits[i] = c
	}

	raw := string(digits)

	// group digits if requested
	if group > 0 && separator != "" {
		var parts []string
		for i := 0; i < len(raw); i += group {
			end := i + group
			if end > len(raw) {
				end = len(raw)
			}
			parts = append(parts, raw[i:end])
		}
		return strings.Join(parts, separator), nil
	}

	return raw, nil
}

func init() {
	rootCmd.AddCommand(pinCmd)
	pinCmd.Flags().IntP("length", "l", 6, "Number of digits")
	pinCmd.Flags().IntP("count", "c", 1, "Number of PINs to generate")
	pinCmd.Flags().IntP("group", "g", 0, "Group digits into chunks of N (0 = no grouping)")
	pinCmd.Flags().StringP("separator", "s", "-", "Separator character between groups")
}
