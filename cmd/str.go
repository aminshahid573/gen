/*
Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
*/
package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gen/internal/ui"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// strCmd represents the str command
var strCmd = &cobra.Command{
	Use:   "str",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		length, _ := cmd.Flags().GetInt("length")
		enc, _ := cmd.Flags().GetString("encoding")

		str, err := generateStr(length, enc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating string: %w", err)
			os.Exit(1)
		}

		row := [][]string{{enc, strconv.Itoa(length), str}}
		fmt.Println(ui.RenderTable([]string{"Encoding", "Length", "Output"}, row))

	},
}

func generateStr(length int, enc string) (string, error) {
	bytes := make([]byte, length)
	rand.Read(bytes) // Practice using crypto/rand for security!

	var output string
	var charSet string = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789!@#$%^&*()-_=+[]{}|;:,.<>?"
	switch enc {
	case "raw":
		for range length {
			randomIndex, _ := rand.Read(bytes)
			output += string(charSet[randomIndex])
		}
		return output, nil

	case "base64":
		output = base64.StdEncoding.EncodeToString(bytes)
		return output, nil
	case "hex":
		output = fmt.Sprintf("%x", bytes)
		return output, nil
	default:
		return "", fmt.Errorf("Unknown encoding")
	}

}

func init() {
	rootCmd.AddCommand(strCmd)
	strCmd.Flags().IntP("length", "l", 16, "Length of the generated string")
	strCmd.Flags().StringP("encoding", "e", "hex", "Encoding type (hex, base64)")

}
