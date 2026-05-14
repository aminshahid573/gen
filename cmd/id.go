// Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
package cmd

import (
	"fmt"
	"github.com/aminshahid573/gen/internal/ui"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"github.com/spf13/cobra"
)

var idCmd = &cobra.Command{
	Use:   "id",
	Short: "Generate a UUID",
	Long: `Generate a Universally Unique Identifier (UUID) of a specified version.

Supported versions:
  1  Time-based + MAC address
  3  MD5 hash of a namespace and name (requires --name and --namespace)
  4  Random (default)
  5  SHA-1 hash of a namespace and name (requires --name and --namespace)
  6  Sortable time-based (reordered v1)
  7  Sortable Unix timestamp-based

Examples:
  gen id
  gen id --version 7
  gen id --count 10 --version 7
  gen id --version 5 --namespace url --name https://example.com
  gen id --version 3 --namespace dns --name example.com
  gen id --version 4 --format base58
  gen id --version 7 --format urn
  gen id --format random --count 5
  gen id --decode 019681a2-f3c8-7000-a7e1-5b9f4c3d2e1a`,

	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetInt("version")
		name, _ := cmd.Flags().GetString("name")
		ns, _ := cmd.Flags().GetString("namespace")
		format, _ := cmd.Flags().GetString("format")
		raw, _ := cmd.Flags().GetString("decode")
		count, _ := cmd.Flags().GetInt("count")

		// decode mode — skip generation
		if raw != "" {
			if count > 1 {
				fatalf("--count cannot be used with --decode")
			}
			decodeUUID(raw)
			return
		}

		if count < 1 {
			fatalf("--count must be at least 1")
		}

		rows := make([][]string, 0, count)

		for range count {
			id, err := generateUUID(version, name, ns)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error generating UUID v%d: %v\n", version, err)
				os.Exit(1)
			}

			// resolve "random" to an actual format name before formatting
			actualFormat := resolveFormat(format)

			out, err := formatUUID(id, actualFormat)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}

			rows = append(rows, []string{fmt.Sprintf("v%d", version), actualFormat, out})
		}

		fmt.Println(ui.RenderTable([]string{"Version", "Format", "UUID"}, rows))
	},
}

// generateUUID returns a UUID for the given version, name, and namespace.
func generateUUID(version int, name, ns string) (uuid.UUID, error) {
	switch version {
	case 1:
		return uuid.NewUUID()
	case 3:
		if name == "" {
			return uuid.UUID{}, fmt.Errorf("--name is required for version 3")
		}
		return uuid.NewMD5(resolveNamespace(ns), []byte(name)), nil
	case 4:
		return uuid.NewRandom()
	case 5:
		if name == "" {
			return uuid.UUID{}, fmt.Errorf("--name is required for version 5")
		}
		return uuid.NewSHA1(resolveNamespace(ns), []byte(name)), nil
	case 6:
		return uuid.NewV6()
	case 7:
		return uuid.NewV7()
	default:
		return uuid.UUID{}, fmt.Errorf("unsupported version %d (choose from 1,3,4,5,6,7)", version)
	}
}

// resolveFormat maps "random" to a concrete format name.
// All other values pass through unchanged.
func resolveFormat(format string) string {
	if strings.ToLower(format) != "random" {
		return format
	}
	formats := []string{"standard", "compact", "upper", "urn", "base58"}
	return formats[rand.Intn(len(formats))]
}

// formatUUID returns the UUID string in the requested format.
func formatUUID(id uuid.UUID, format string) (string, error) {
	switch strings.ToLower(format) {
	case "standard":
		return id.String(), nil
	case "compact":
		return strings.ReplaceAll(id.String(), "-", ""), nil
	case "upper":
		return strings.ToUpper(id.String()), nil
	case "urn":
		return id.URN(), nil
	case "base58":
		return base58.Encode(id[:]), nil
	default:
		return "", fmt.Errorf("invalid format %q (choose: standard, compact, upper, urn, base58, random)", format)
	}
}

// decodeUUID parses a UUID string and prints its metadata.
func decodeUUID(raw string) {
	id, err := uuid.Parse(raw)
	if err != nil {
		fatalf("invalid UUID %q: %v", raw, err)
	}

	rows := [][]string{
		{"String", id.String()},
		{"Compact", strings.ReplaceAll(id.String(), "-", "")},
		{"URN", id.URN()},
		{"Version", strconv.Itoa(int(id.Version()))},
		{"Variant", id.Variant().String()},
		{"Timestamp", extractTimestamp(id)},
	}

	fmt.Println(ui.RenderTable([]string{"Field", "Value"}, rows))
}

// extractTimestamp returns a human-readable timestamp for versions that embed one.
func extractTimestamp(id uuid.UUID) string {
	switch id.Version() {
	case 1, 6:
		sec, nsec := id.Time().UnixTime()
		return time.Unix(sec, nsec).UTC().Format(time.RFC3339Nano)
	case 7:
		ms := int64(id[0])<<40 | int64(id[1])<<32 | int64(id[2])<<24 |
			int64(id[3])<<16 | int64(id[4])<<8 | int64(id[5])
		return time.UnixMilli(ms).UTC().Format(time.RFC3339Nano)
	case 3, 5:
		return "n/a (deterministic hash)"
	default:
		return "n/a (random)"
	}
}

// resolveNamespace maps a string name to a UUID namespace.
func resolveNamespace(ns string) uuid.UUID {
	switch ns {
	case "dns":
		return uuid.NameSpaceDNS
	case "url":
		return uuid.NameSpaceURL
	case "oid":
		return uuid.NameSpaceOID
	case "x500":
		return uuid.NameSpaceX500
	default:
		parsed, err := uuid.Parse(ns)
		if err != nil {
			return uuid.NameSpaceDNS
		}
		return parsed
	}
}

// fatalf prints an error to stderr and exits.
func fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(idCmd)
	idCmd.Flags().IntP("version", "v", 4, "UUID version to generate (1,3,4,5,6,7)")
	idCmd.Flags().StringP("name", "n", "", "Name string for v3/v5 hashed UUIDs")
	idCmd.Flags().StringP("namespace", "s", "dns", "Namespace for v3/v5 (dns, url, oid, x500, or a UUID)")
	idCmd.Flags().StringP("format", "f", "standard", "Output format: standard, compact, upper, urn, base58, random")
	idCmd.Flags().StringP("decode", "d", "", "Decode and inspect an existing UUID")
	idCmd.Flags().IntP("count", "c", 1, "Number of UUIDs to generate")
}
