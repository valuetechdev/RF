//go:generate curl https://kontohjelp.fiken.no/data/kontoGruppeInfo > accounts.json
//go:generate go run cmd/download/main.go

package main
