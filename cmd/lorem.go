// Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
package cmd

import (
	"fmt"
	"gen/internal/ui"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// classic lorem ipsum word pool
var loremWords = strings.Fields(`lorem ipsum dolor sit amet consectetur adipiscing elit
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua ut enim ad minim veniam
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat
duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat
nulla pariatur excepteur sint occaecat cupidatat non proident sunt in culpa qui officia
deserunt mollit anim id est laborum curabitur pretium tincidunt lacus nulla gravida orci
a odio tempus ullamcorper volutpat ultricies mi justo etiam bibendum egestas facilisis
posuere vulputate aliquet porta nunc faucibus risus pretium accumsan phasellus`)

var loremCmd = &cobra.Command{
	Use:   "lorem",
	Short: "Generate lorem ipsum placeholder text",
	Long: `Generate lorem ipsum placeholder text for mockups, testing, and UI development.

Output modes (pick one):
  --words N       Generate exactly N words
  --sentences N   Generate N sentences  (default: 1)
  --paragraphs N  Generate N paragraphs

Examples:
  gen lorem
  gen lorem --words 20
  gen lorem --sentences 3
  gen lorem --paragraphs 2
  gen lorem --paragraphs 3 --min-sentences 4 --max-sentences 8`,

	Run: func(cmd *cobra.Command, args []string) {
		words, _       := cmd.Flags().GetInt("words")
		sentences, _   := cmd.Flags().GetInt("sentences")
		paragraphs, _  := cmd.Flags().GetInt("paragraphs")
		minSent, _     := cmd.Flags().GetInt("min-sentences")
		maxSent, _     := cmd.Flags().GetInt("max-sentences")

		// validate
		if minSent > maxSent {
			fatalf("--min-sentences cannot be greater than --max-sentences")
		}

		wordsFlagSet      := cmd.Flags().Changed("words")
		sentencesFlagSet  := cmd.Flags().Changed("sentences")
		paragraphsFlagSet := cmd.Flags().Changed("paragraphs")

		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		var output string

		switch {
		case wordsFlagSet:
			if words < 1 {
				fatalf("--words must be at least 1")
			}
			output = generateLoremWords(rng, words)

		case paragraphsFlagSet:
			if paragraphs < 1 {
				fatalf("--paragraphs must be at least 1")
			}
			var parts []string
			for range paragraphs {
				n := minSent + rng.Intn(maxSent-minSent+1)
				parts = append(parts, generateLoremParagraph(rng, n))
			}
			output = strings.Join(parts, "\n\n")

		case sentencesFlagSet:
			fallthrough
		default:
			if sentences < 1 {
				fatalf("--sentences must be at least 1")
			}
			var sents []string
			for range sentences {
				sents = append(sents, generateLoremSentence(rng))
			}
			output = strings.Join(sents, " ")
		}

		// for words/sentences: table output; for paragraphs: plain (too long for table)
		if paragraphsFlagSet {
			fmt.Println(output)
		} else {
			fmt.Println(ui.RenderTable(
				[]string{"Length", "Output"},
				[][]string{{fmt.Sprintf("%d chars", len(output)), output}},
			))
		}
	},
}

func generateLoremWords(rng *rand.Rand, n int) string {
	picked := make([]string, n)
	for i := range n {
		picked[i] = loremWords[rng.Intn(len(loremWords))]
	}
	sentence := strings.Join(picked, " ")
	return capitalize(sentence)
}

func generateLoremSentence(rng *rand.Rand) string {
	wordCount := 6 + rng.Intn(10) // 6–15 words per sentence
	words := make([]string, wordCount)
	for i := range wordCount {
		words[i] = loremWords[rng.Intn(len(loremWords))]
	}
	return capitalize(strings.Join(words, " ")) + "."
}

func generateLoremParagraph(rng *rand.Rand, sentenceCount int) string {
	sents := make([]string, sentenceCount)
	for i := range sentenceCount {
		sents[i] = generateLoremSentence(rng)
	}
	return strings.Join(sents, " ")
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func init() {
	rootCmd.AddCommand(loremCmd)
	loremCmd.Flags().IntP("words", "w", 0, "Generate exactly N words")
	loremCmd.Flags().IntP("sentences", "s", 1, "Generate N sentences")
	loremCmd.Flags().IntP("paragraphs", "p", 0, "Generate N paragraphs")
	loremCmd.Flags().Int("min-sentences", 3, "Min sentences per paragraph")
	loremCmd.Flags().Int("max-sentences", 7, "Max sentences per paragraph")
}
