/*
Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gen",
	Short: "A collection of small dev utilities for your terminal",
	Long: `gen bundles the things you'd normally google or visit a website for
into a single CLI — UUIDs, passwords, QR codes, OTPs, tokens, lorem
ipsum, timestamps and more.

Run 'gen --help' to see all available commands.`,
}

func Execute(version string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
