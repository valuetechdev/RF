# RF - Norwegian Financial Terms CLI

A command-line tool for looking up Norwegian financial and accounting terms. Fetches data from an API endpoint and displays formatted Markdown output.

## Installation

### Build from source
```bash
git clone <this-repo>
cd rf-lookup
chmod +x install.sh
./install.sh
```

### Manual installation
```bash
go build -o rf
sudo mv rf /usr/local/bin/
```

## Usage

### Basic Usage
Set your API endpoint and look up terms:

```bash
export BASE_API="https://your-api-endpoint.com"
rf bokforing
```

### Using API flag
```bash
rf --api https://your-api-endpoint.com bokforing
```

### Examples
```bash
# Look up "bokføring" (accounting)
rf bokforing

# Look up "regnskap" (accounting records)  
rf regnskap

# Look up "balanse" (balance)
rf balanse

# With custom timeout
rf --timeout 10s bokforing
```

## Configuration

### Environment Variables
- `BASE_API`: Base URL for the API endpoint (required)

### Flags
- `--api`: Base API URL (overrides BASE_API env var)
- `--timeout`: HTTP request timeout (default: 30s)
- `--help`: Show help information

## API Requirements

The CLI expects the API to:
1. Accept GET requests to `{BASE_API}/{term}`
2. Return JSON in the specific format with rich text blocks
3. Include fields like `tittel`, `sammendrag`, `body`, `synonymer`, etc.

## Output Format

The tool converts JSON responses to clean Markdown with:
- Document title as H1
- Summary section
- Structured content with proper headings
- Rich text formatting (bold, italic, links)
- Lists and nested content
- Synonyms and related terms sections

## Example Output

```markdown
# bokføring

## Summary
Bokføring er å føre regnskap. Det er å registrere *posteringer* i regnskapet.

## Hvorfor må man føre bokføre
Du må føre bokføre – altså føre regnskap – fordi det å tjene penger...

## Synonyms
postering, transaksjon, regnskapsføring, bokføre

## Related Terms
- [regnskap](regnskap)
- [postering](postering)
```

## Development

### Project Structure
```
.
├── cmd/root.go          # CLI command definitions
├── internal/
│   ├── types.go         # JSON structure definitions
│   └── transformer.go   # Markdown transformation logic
├── main.go              # Entry point
└── README.md
```

### Adding Features
The transformer supports:
- Multiple heading levels (h1-h6)
- Rich text marks (strong, em, code, underline, strikethrough)
- Links and custom ordLinks
- Numbered and bulleted lists with nesting
- Custom link resolvers

## License

MIT License