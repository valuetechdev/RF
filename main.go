package main

import (
	"embed"

	"github.com/valuetechtev/rf/cmd"
)

//go:embed terms
var terms embed.FS

func main() {
	cmd.Execute(terms)
}
