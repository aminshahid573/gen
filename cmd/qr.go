// Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
	"rsc.io/qr"
)

var qrCmd = &cobra.Command{
	Use:   "qr [input]",
	Short: "Render a QR code in the terminal",
	Long: `Render a QR code directly in the terminal using Unicode block characters.

Input can be passed as a positional argument or via --input.
Error correction levels: L (7%), M (15%), Q (25%), H (30%).
Higher levels produce larger QR codes but tolerate more damage.

Examples:
  gen qr "https://example.com"
  gen qr --input "https://example.com"
  gen qr --input "Hello, World!" --level H
  gen qr --input "otpauth://totp/Example?secret=JBSWY3DP" --level M
  gen qr --input "mailto:aminShahid5515@gmail.com"`,

	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		level, _ := cmd.Flags().GetString("level")
		noHalf, _ := cmd.Flags().GetBool("no-half-blocks")
		quietZone, _ := cmd.Flags().GetInt("quiet-zone")

		if len(args) > 0 {
			input = strings.Join(args, " ")
		}

		if input == "" {
			fatalf("input is required — pass as argument or --input")
		}

		cfg := qrterminal.Config{
			Level:          resolveQRLevel(level), // now correctly qrterminal.Level
			Writer:         os.Stdout,
			BlackChar:      qrterminal.BLACK,
			WhiteChar:      qrterminal.WHITE,
			BlackWhiteChar: qrterminal.BLACK_WHITE,
			WhiteBlackChar: qrterminal.WHITE_BLACK,
			QuietZone:      quietZone,
			HalfBlocks:     !noHalf,
		}

		fmt.Println()
		qrterminal.GenerateWithConfig(input, cfg)
		fmt.Println()
		fmt.Printf("  Input : %s\n", input)
		fmt.Printf("  Level : %s\n", strings.ToUpper(level))
	},
}

// resolveQRLevel returns qrterminal.Level — NOT rsc.io/qr.Level
func resolveQRLevel(level string) qr.Level {
	switch strings.ToUpper(level) {
	case "M":
		return qr.M
	case "Q":
		return qr.Q
	case "H":
		return qr.H
	default:
		return qr.L
	}
}

func init() {
	rootCmd.AddCommand(qrCmd)
	qrCmd.Flags().StringP("input", "i", "", "String to encode into the QR code")
	qrCmd.Flags().StringP("level", "l", "L", "Error correction level: L, M, Q, H")
	qrCmd.Flags().Bool("no-half-blocks", false, "Use full blocks instead of half-block Unicode characters")
	qrCmd.Flags().Int("quiet-zone", 2, "Quiet zone size (blank border around QR code)")
}
