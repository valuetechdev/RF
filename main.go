package main

import (
	"embed"

	"github.com/valuetechtev/rf/cmd"
)

//go:embed terms
var terms embed.FS

//go:embed accounts.json
var accounts embed.FS

func main() {
	cmd.Execute(terms, accounts)
}
