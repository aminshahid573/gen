// Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gen/internal/ui"
	"os"

	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate an opaque bearer token",
	Long: `Generate a cryptographically secure opaque bearer token.

Output is always URL-safe base64 (no padding) derived from raw random bytes,
making it safe for use in HTTP Authorization headers, cookies, and URLs.

Examples:
  gen token
  gen token --bytes 48
  gen token --prefix "tok_"
  gen token --bearer
  gen token --count 5
  gen token --bytes 64 --prefix "sess_"`,

	Run: func(cmd *cobra.Command, args []string) {
		bytes, _  := cmd.Flags().GetInt("bytes")
		prefix, _ := cmd.Flags().GetString("prefix")
		bearer, _ := cmd.Flags().GetBool("bearer")
		count, _  := cmd.Flags().GetInt("count")

		if count < 1 {
			fatalf("--count must be at least 1")
		}
		if bytes < 1 {
			fatalf("--bytes must be at least 1")
		}

		rows := make([][]string, 0, count)
		for range count {
			token, err := generateToken(bytes, prefix, bearer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error generating token: %v\n", err)
				os.Exit(1)
			}
			rows = append(rows, []string{fmt.Sprintf("%d", bytes), token})
		}

		fmt.Println(ui.RenderTable([]string{"Bytes", "Token"}, rows))
	},
}

func generateToken(byteCount int, prefix string, bearer bool) (string, error) {
	b := make([]byte, byteCount)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := prefix + base64.RawURLEncoding.EncodeToString(b)

	if bearer {
		token = "Bearer " + token
	}

	return token, nil
}

func init() {
	rootCmd.AddCommand(tokenCmd)
	tokenCmd.Flags().IntP("bytes", "b", 32, "Number of random bytes (more bytes = longer token)")
	tokenCmd.Flags().StringP("prefix", "p", "", "Prepend a fixed string (e.g. tok_, sess_)")
	tokenCmd.Flags().BoolP("bearer", "B", false, `Prepend "Bearer " for direct use in Authorization headers`)
	tokenCmd.Flags().IntP("count", "c", 1, "Number of tokens to generate")
}
