package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/valuetechtev/rf/internal"

	"github.com/spf13/cobra"
)

var (
	timeout time.Duration
	verbose bool
	plain   bool

	terms     embed.FS
	termsPath = filepath.Clean("terms")
)

var rootCmd = &cobra.Command{
	Use:   "rf [term]",
	Short: "Lookup Norwegian financial terms and display as Markdown",
	Long: `RF (Regnskapsfaglig) is a CLI tool for looking up Norwegian financial and accounting terms.
It fetches data from an API endpoint and displays the result as formatted Markdown.`,
	Args: cobra.ExactArgs(1),
	Example: `  rf bokforing    # Look up "bokføring" 
  rf regnskap     # Look up "regnskap"
  rf balanse      # Look up "balanse"`, //nolint:misspell
	ValidArgsFunction: completeList,
	RunE: func(cmd *cobra.Command, args []string) error {
		term := args[0]
		return lookupTerm(term)
	},
}

func Execute(t embed.FS) {
	terms = t
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Second, "HTTP request timeout")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "More than just the summary")
	rootCmd.PersistentFlags().BoolVar(&plain, "plain", false, "print plain output instead of pretty")
	rootCmd.AddCommand(listCmd())
}

func lookupTerm(term string) error {
	f, err := terms.ReadFile(filepath.Join(termsPath, fmt.Sprintf("%s.json", term)))
	if err != nil {
		return err
	}

	var doc internal.Document
	if err := json.Unmarshal(f, &doc); err != nil {
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

	markdown := transformer.Transform(doc)

	if !plain {
		r, err := glamour.NewTermRenderer(
			glamour.WithWordWrap(80),
			glamour.WithStylePath("dark"),
		)
		if err != nil {
			return err
		}

		out, err := r.Render(markdown)
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	}

	fmt.Print(markdown)

	return nil
}
