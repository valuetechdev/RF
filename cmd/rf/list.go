package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		for _, term := range allTerms() {
			fmt.Println(term)
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

func allTerms() []string {
	var s []string
	entries, err := terms.ReadDir(termsPath)
	cobra.CheckErr(err)
	for _, entry := range entries {
		s = append(s, strings.ReplaceAll(entry.Name(), ".json", ""))
	}
	return s
}

func completeList(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return allTerms(), cobra.ShellCompDirectiveNoFileComp
}
