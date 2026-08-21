//go:generate curl -sSfL -o accounts.json https://kontohjelp.fiken.no/data/kontoGruppeInfo
//go:generate go run cmd/download/main.go

package main
