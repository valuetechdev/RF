package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"github.com/valuetechtev/rf/internal"
)

func kontoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "konto [account number]",
		Short: "Look up Norwegian account details",
		Long:  `Look up details for a specific Norwegian accounting account number.`,
		Args:  cobra.ExactArgs(1),
		Example: `  rf konto 4900
  rf konto 1920`,
		ValidArgsFunction: completeAccounts,
		RunE: func(cmd *cobra.Command, args []string) error {
			accountNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid account number: %s", args[0])
			}

			account, err := lookupAccount(accountNum)
			if err != nil {
				return err
			}

			return renderAccount(account)
		},
	}
	return cmd
}

func lookupAccount(accountNum int) (*internal.Account, error) {
	f, err := accounts.ReadFile("accounts.json")
	if err != nil {
		return nil, err
	}

	var groups []internal.AccountGroup
	if err := json.Unmarshal(f, &groups); err != nil {
		return nil, fmt.Errorf("failed to parse accounts JSON: %w", err)
	}

	var closestMatch *internal.Account

	for _, group := range groups {
		for _, account := range group.Kontoer {
			if account.Kontonummer == accountNum {
				return &account, nil
			}

			if account.Kontonummer < accountNum {
				if closestMatch == nil || account.Kontonummer > closestMatch.Kontonummer {
					closestMatch = &account
				}
			}
		}
	}

	if closestMatch != nil {
		return closestMatch, nil
	}

	return nil, fmt.Errorf("account %d not found", accountNum)
}

func renderAccount(account *internal.Account) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %d - %s\n\n", account.Kontonummer, account.Navn))

	if account.MetaData.Hjelpetekst != "" {
		sb.WriteString(account.MetaData.Hjelpetekst)
		sb.WriteString("\n\n")
	}

	if account.MetaData.Advarsel != "" {
		sb.WriteString("⚠️  **Advarsel:**\n")
		sb.WriteString(account.MetaData.Advarsel)
		sb.WriteString("\n\n")
	}

	if account.MetaData.Eksempler != "" {
		sb.WriteString("**Eksempler:**\n")
		sb.WriteString(account.MetaData.Eksempler)
		sb.WriteString("\n\n")
	}

	if verbose {
		if len(account.MetaData.GyldigeMvakoder) > 0 {
			sb.WriteString("**Gyldige MVA-koder:** ")
			sb.WriteString(strings.Join(account.MetaData.GyldigeMvakoder, ", "))
			sb.WriteString("\n\n")
		}

		if account.MetaData.DefaultMvakode != nil {
			sb.WriteString(fmt.Sprintf("**Standard MVA-kode:** %s\n\n", *account.MetaData.DefaultMvakode))
		}

		if len(account.MetaData.Sokeord) > 0 {
			sb.WriteString("**Søkeord:** ")
			sb.WriteString(strings.Join(account.MetaData.Sokeord, ", "))
			sb.WriteString("\n\n")
		}
	}

	markdown := sb.String()

	return renderMarkdown(markdown)
}

func renderMarkdown(markdown string) error {
	if plain {
		fmt.Print(markdown)
		return nil
	}

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

func allAccounts() []string {
	f, err := accounts.ReadFile("accounts.json")
	if err != nil {
		return nil
	}

	var groups []internal.AccountGroup
	if err := json.Unmarshal(f, &groups); err != nil {
		return nil
	}

	var accountNumbers []string
	for _, group := range groups {
		for _, account := range group.Kontoer {
			accountNumbers = append(accountNumbers, fmt.Sprintf("%d\t%s", account.Kontonummer, account.Navn))
		}
	}

	return accountNumbers
}

func completeAccounts(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return allAccounts(), cobra.ShellCompDirectiveNoFileComp
}
