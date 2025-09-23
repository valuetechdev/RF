package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		entries, err := terms.ReadDir(termsPath)
		cobra.CheckErr(err)
		for _, entry := range entries {
			fmt.Println(strings.ReplaceAll(entry.Name(), ".json", ""))
		}
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available terms",
		Args:    cobra.NoArgs,
		Run:     run,
	}
	return cmd
}
