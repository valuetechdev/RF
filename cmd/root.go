package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
		doc, err := lookupTerm(term)
		cobra.CheckErr(err)
		return RenderDoc(doc)
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

func lookupTerm(term string) (*internal.Document, error) {
	f, err := terms.ReadFile(filepath.Join(termsPath, fmt.Sprintf("%s.json", term)))
	if err != nil {
		return nil, err
	}

	var doc internal.Document

	if err := json.Unmarshal(f, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if !verbose {
		// in non-verbose mode, we modify the doc by removing the stuff
		// we dont want to render

		// remove the body, and keep the summary
		doc.Body = nil

		// remove synonyms
		doc.Synonymer = nil

		// remove the slug at the end of the related terms
		for idx, related := range doc.Relaterte {
			related.Slug = internal.Slug{}
			doc.Relaterte[idx] = related
		}
	}
	return &doc, nil
}

func RenderDoc(doc *internal.Document) error {
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
