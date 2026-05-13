/*
Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
*/
package cmd

import (
	"os"

	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
)

// qrCmd represents the qr command
var qrCmd = &cobra.Command{
	Use:   "qr",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		textData, _ := cmd.Flags().GetString("text")
		if textData != "" {
			generateQR(textData)
		}
	},
}

func init() {
	rootCmd.AddCommand(qrCmd)
	qrCmd.Flags().StringP("text", "t", "", "Text to genrate QR Code")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// qrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// qrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generateQR(rawData string) {
	config := qrterminal.Config{
		HalfBlocks: true,
		Level:      qrterminal.M,
		Writer:     os.Stdout,
	}
	qrterminal.GenerateWithConfig(rawData, config)
}
