// Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
package cmd

import (
	"fmt"
	"github.com/aminshahid573/gen/internal/ui"
	"os"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
)

var otpCmd = &cobra.Command{
	Use:   "otp",
	Short: "Generate a TOTP or HOTP one-time password",
	Long: `Generate a one-time password from a base32 secret for testing 2FA flows.

TOTP (default) uses the current time as input (RFC 6238).
HOTP uses an incrementing counter as input (RFC 4226).

Algorithms: sha1 (default), sha256, sha512
Digits:     6 (default), 8

Examples:
  gen otp --secret JBSWY3DPEHPK3PXP
  gen otp --secret JBSWY3DPEHPK3PXP --type hotp --counter 42
  gen otp --secret JBSWY3DPEHPK3PXP --digits 8 --algorithm sha256
  gen otp --secret JBSWY3DPEHPK3PXP --remaining
  gen otp --secret JBSWY3DPEHPK3PXP --validate 123456`,

	Run: func(cmd *cobra.Command, args []string) {
		secret, _ := cmd.Flags().GetString("secret")
		otpType, _ := cmd.Flags().GetString("type")
		digits, _ := cmd.Flags().GetInt("digits")
		period, _ := cmd.Flags().GetUint("period")
		counter, _ := cmd.Flags().GetUint64("counter")
		algo, _ := cmd.Flags().GetString("algorithm")
		remaining, _ := cmd.Flags().GetBool("remaining")
		validate, _ := cmd.Flags().GetString("validate")

		if secret == "" {
			fatalf("--secret is required (base32-encoded, e.g. JBSWY3DPEHPK3PXP)")
		}

		secret = strings.ToUpper(strings.TrimSpace(secret))

		otpDigits := otp.DigitsSix
		if digits == 8 {
			otpDigits = otp.DigitsEight
		} else if digits != 6 {
			fatalf("--digits must be 6 or 8")
		}

		otpAlgo := resolveOTPAlgorithm(algo)

		// validate mode
		if validate != "" {
			var valid bool
			var err error
			switch strings.ToLower(otpType) {
			case "totp":
				valid, err = totp.ValidateCustom(validate, secret, time.Now().UTC(), totp.ValidateOpts{
					Period:    period,
					Skew:      1,
					Digits:    otpDigits,
					Algorithm: otpAlgo,
				})
			case "hotp":
				valid, err = hotp.ValidateCustom(validate, counter, secret, hotp.ValidateOpts{
					Digits:    otpDigits,
					Algorithm: otpAlgo,
				})
			default:
				fatalf("--type must be totp or hotp")
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "error validating OTP: %v\n", err)
				os.Exit(1)
			}

			result := "invalid ✗"
			if valid {
				result = "valid ✓"
			}
			fmt.Println(ui.RenderTable(
				[]string{"Code", "Type", "Result"},
				[][]string{{validate, strings.ToUpper(otpType), result}},
			))
			return
		}

		// generation mode
		var code string
		var err error

		switch strings.ToLower(otpType) {
		case "totp":
			code, err = totp.GenerateCodeCustom(secret, time.Now().UTC(), totp.ValidateOpts{
				Period:    period,
				Digits:    otpDigits,
				Algorithm: otpAlgo,
			})
		case "hotp":
			code, err = hotp.GenerateCodeCustom(secret, counter, hotp.ValidateOpts{
				Digits:    otpDigits,
				Algorithm: otpAlgo,
			})
		default:
			fatalf("--type must be totp or hotp")
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "error generating OTP: %v\n", err)
			os.Exit(1)
		}

		rows := [][]string{{
			strings.ToUpper(otpType),
			fmt.Sprintf("%d", digits),
			strings.ToUpper(algo),
			code,
		}}

		if remaining && strings.ToLower(otpType) == "totp" {
			elapsed := time.Now().Unix() % int64(period)
			left := int64(period) - elapsed
			rows[0] = append(rows[0], fmt.Sprintf("%ds", left))
			fmt.Println(ui.RenderTable([]string{"Type", "Digits", "Algorithm", "Code", "Remaining"}, rows))
		} else {
			fmt.Println(ui.RenderTable([]string{"Type", "Digits", "Algorithm", "Code"}, rows))
		}
	},
}

func resolveOTPAlgorithm(algo string) otp.Algorithm {
	switch strings.ToLower(algo) {
	case "sha256":
		return otp.AlgorithmSHA256
	case "sha512":
		return otp.AlgorithmSHA512
	default:
		return otp.AlgorithmSHA1
	}
}

func init() {
	rootCmd.AddCommand(otpCmd)
	otpCmd.Flags().StringP("secret", "s", "", "Base32-encoded shared secret (required)")
	otpCmd.Flags().StringP("type", "t", "totp", "OTP type: totp or hotp")
	otpCmd.Flags().IntP("digits", "d", 6, "Code length: 6 or 8")
	otpCmd.Flags().UintP("period", "p", 30, "TOTP time step in seconds")
	otpCmd.Flags().Uint64P("counter", "C", 0, "HOTP counter value")
	otpCmd.Flags().StringP("algorithm", "a", "sha1", "Hash algorithm: sha1, sha256, sha512")
	otpCmd.Flags().BoolP("remaining", "r", false, "Show seconds remaining until next TOTP code")
	otpCmd.Flags().StringP("validate", "v", "", "Validate this code against the secret instead of generating")
}
