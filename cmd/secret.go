// Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/aminshahid573/gen/internal/ui"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Generate a cryptographic secret from raw random bytes",
	Long: `Generate a cryptographic secret suitable for API keys, JWT signing secrets,
and HMAC keys. Unlike gen str, output is derived from raw random bytes —
not a character set — giving maximum entropy per byte.

Types and their defaults:
  generic   32 bytes, hex encoded
  jwt       64 bytes, base64url encoded  (HS256/HS512 signing secret)
  hmac      32 bytes, hex encoded        (HMAC-SHA256 key)
  api       32 bytes, base64url encoded  (API key secret component)

Examples:
  gen secret
  gen secret --type jwt
  gen secret --type hmac
  gen secret --type api --prefix "sk_live_"
  gen secret --bytes 64 --encoding hex
  gen secret --count 3`,

	Run: func(cmd *cobra.Command, args []string) {
		secretType, _ := cmd.Flags().GetString("type")
		encoding, _ := cmd.Flags().GetString("encoding")
		bytes, _ := cmd.Flags().GetInt("bytes")
		prefix, _ := cmd.Flags().GetString("prefix")
		count, _ := cmd.Flags().GetInt("count")

		if count < 1 {
			fatalf("--count must be at least 1")
		}

		// apply type defaults if the user hasn't overridden bytes/encoding
		bytesFlagSet := cmd.Flags().Changed("bytes")
		encodingFlagSet := cmd.Flags().Changed("encoding")

		byteCount, enc := resolveSecretDefaults(secretType, bytes, encoding, bytesFlagSet, encodingFlagSet)

		rows := make([][]string, 0, count)
		for range count {
			secret, err := generateSecret(byteCount, enc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error generating secret: %v\n", err)
				os.Exit(1)
			}
			out := prefix + secret
			rows = append(rows, []string{secretType, enc, fmt.Sprintf("%d", byteCount), out})
		}

		fmt.Println(ui.RenderTable([]string{"Type", "Encoding", "Bytes", "Secret"}, rows))
	},
}

// resolveSecretDefaults applies sensible per-type defaults
// unless the user explicitly set --bytes or --encoding.
func resolveSecretDefaults(secretType string, bytes int, encoding string, bytesFlagSet, encodingFlagSet bool) (int, string) {
	type defaults struct {
		bytes    int
		encoding string
	}

	typeDefaults := map[string]defaults{
		"generic": {32, "hex"},
		"jwt":     {64, "base64url"},
		"hmac":    {32, "hex"},
		"api":     {32, "base64url"},
	}

	d, ok := typeDefaults[strings.ToLower(secretType)]
	if !ok {
		fatalf("unknown secret type %q (choose: generic, jwt, hmac, api)", secretType)
	}

	if !bytesFlagSet {
		bytes = d.bytes
	}
	if !encodingFlagSet {
		encoding = d.encoding
	}
	return bytes, encoding
}

func generateSecret(byteCount int, encoding string) (string, error) {
	if byteCount < 1 {
		return "", fmt.Errorf("--bytes must be at least 1")
	}

	b := make([]byte, byteCount)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	switch strings.ToLower(encoding) {
	case "hex":
		return hex.EncodeToString(b), nil
	case "base64":
		return base64.StdEncoding.EncodeToString(b), nil
	case "base64url":
		return base64.RawURLEncoding.EncodeToString(b), nil
	default:
		return "", fmt.Errorf("unknown encoding %q (choose: hex, base64, base64url)", encoding)
	}
}

func init() {
	rootCmd.AddCommand(secretCmd)
	secretCmd.Flags().StringP("type", "t", "generic", "Secret type: generic, jwt, hmac, api")
	secretCmd.Flags().IntP("bytes", "b", 32, "Number of random bytes (overrides type default)")
	secretCmd.Flags().StringP("encoding", "e", "hex", "Encoding: hex, base64, base64url (overrides type default)")
	secretCmd.Flags().StringP("prefix", "p", "", "Prepend a fixed string (e.g. sk_live_)")
	secretCmd.Flags().IntP("count", "c", 1, "Number of secrets to generate")
}
