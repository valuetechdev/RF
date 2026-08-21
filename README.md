# RF - Norwegian Financial Terms CLI

RF (Regnskapsfører) is a command-line tool for looking up Norwegian financial
and accounting terms from [Fiken forklarer](https://fiken.no/forklarer), plus
the Norwegian standard chart of accounts from
[Kontohjelp](https://kontohjelp.fiken.no). All data is embedded in the binary at
build time, so lookups work offline and instantly.

## Installation

```bash
go install github.com/valuetechdev/rf@latest
```

### Build from source

From a checkout of this repo:

```bash
go install .
# or
mise run install
```

Either way the `rf` binary lands in `$GOBIN` (`$(go env GOPATH)/bin` by default).

## Usage

### Look up a term

```bash
# Look up "bokføring" (accounting) — terms are addressed by their slug
rf bokforing

# The full article, not just the summary
rf -v bokforing

# Raw Markdown instead of the rendered terminal output
rf --plain bokforing
```

By default only the summary and related terms are shown. `--verbose` adds the
body text and synonyms.

### List available terms

```bash
rf list
rf ls      # alias
```

### Look up an account number

```bash
rf konto 1920      # Bankinnskudd
rf konto 4900
rf -v konto 4900   # adds valid VAT codes, default VAT code and search words
```

If the exact account number does not exist, the closest lower-numbered account
is shown instead — `rf konto 4905` resolves to 4900.

### Shell completion

Both terms and account numbers complete from the embedded data:

```bash
rf completion fish | source          # fish, current session
rf completion bash > /etc/bash_completion.d/rf
rf completion zsh > "${fpath[1]}/_rf"
```

Run `rf completion --help` for the full instructions per shell.

## Configuration

### Flags

- `-v`, `--verbose`: show the full article / extra account metadata instead of
  just the summary
- `--plain`: print raw Markdown instead of the styled terminal output
- `-h`, `--help`: show help information

Both flags are global and work on every subcommand.

## Output Format

Terms are converted from the source rich-text blocks to Markdown and rendered
with [glamour](https://github.com/charmbracelet/glamour):

- Document title as H1
- Summary section
- Structured content with proper headings
- Rich text formatting (bold, italic, code, underline, strikethrough)
- Links and cross-references to other terms
- Numbered and bulleted lists with nesting
- Synonyms and related terms sections (verbose mode)

### Example

```console
$ rf --plain bokforing
# bokføring

## Summary

Bokføring er å føre regnskap. Det er å registrere *posteringer* i regnskapet.
En postering er én enkelt registrering på én av kontoene i regnskapet ditt.

## Related Terms

- [regnskap]()
- [postering]()
- [bilag]()
```

## Development

### Project Structure

```
.
├── cmd/
│   ├── rf/                  # CLI command definitions
│   │   ├── root.go          # root command, term lookup, rendering
│   │   ├── list.go          # `list` command + term completion
│   │   └── konto.go         # `konto` command + account completion
│   └── download/main.go     # fetches and cleans the term data (go generate)
├── internal/
│   ├── types.go             # JSON structure definitions
│   └── transformer.go       # Markdown transformation logic
├── terms/                   # embedded term data, one JSON file per term
├── accounts.json            # embedded chart of accounts
├── generate.go              # go:generate directives
├── main.go                  # entry point, embeds terms/ and accounts.json
└── README.md
```

### Tasks

Tooling is managed with [mise](https://mise.jdx.dev):

```bash
mise run install    # go install .
mise run generate   # refresh accounts.json and terms/ from upstream
mise run fmt        # go fmt ./...
mise run tidy       # go mod tidy
mise run test       # go test ./...
mise run check      # golangci-lint run
```

### Refreshing the data

`mise run generate` runs the directives in `generate.go`:

1. Downloads `accounts.json` from `kontohjelp.fiken.no/data/kontoGruppeInfo`
2. Runs `cmd/download` to fetch every term from
   `fiken.no/forklarer/api/forklarer/`, strip the fields the CLI does not need,
   and write one JSON file per term into `terms/`

Rebuild afterwards so the new data is embedded in the binary.

### Adding Features

The transformer supports:

- Multiple heading levels (h1-h6)
- Rich text marks (strong, em, code, underline, strikethrough)
- Links and custom ordLinks
- Numbered and bulleted lists with nesting
- Custom link resolvers
