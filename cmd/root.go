package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/valuetechtev/rf/internal"

	"github.com/spf13/cobra"
)

var (
	baseAPI string = "https://fiken.no/forklarer/api/forklarer"
	timeout time.Duration
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "rf [term]",
	Short: "Lookup Norwegian financial terms and display as Markdown",
	Long: `RF (Regnskapsfaglig) is a CLI tool for looking up Norwegian financial and accounting terms.
It fetches data from an API endpoint and displays the result as formatted Markdown.

Examples:
  rf bokforing    # Look up "bokføring" 
  rf regnskap     # Look up "regnskap"
  rf balanse      # Look up "balanse"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		term := args[0]
		return lookupTerm(term)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Second, "HTTP request timeout")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "More than just the summary")
}

func lookupTerm(term string) error {
	url := fmt.Sprintf("%s/%s", baseAPI, term)

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var doc internal.Document
	if err := json.Unmarshal(body, &doc); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}
	if !verbose {
		var result strings.Builder
		fmt.Fprintln(&result, doc.Sammendrag[0].Children[0].Text)
		if len(doc.Relaterte) > 0 {
			result.WriteString("\n## Related\n\n")
			for _, rel := range doc.Relaterte {
				result.WriteString(fmt.Sprintf("- %s\n", rel.Tittel))
			}
			result.WriteString("\n")
		}
		fmt.Print(result.String())
		return nil
	}

	transformer := internal.NewMarkdownTransformer()

	transformer.SetLinkResolver(func(ref string) string {
		return fmt.Sprintf("%s/%s", baseAPI, ref)
	})

	markdown := transformer.Transform(doc)
	fmt.Print(markdown)

	return nil
}
