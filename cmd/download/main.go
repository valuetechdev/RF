package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	baseAPI    = "https://fiken.no/forklarer/api/forklarer/"
	termsDir   = "terms"
	maxRetries = 3
)

type TermListItem struct {
	ID    string `json:"id"`
	Slug  Slug   `json:"slug"`
	Title string `json:"tittel"`
}

type Slug struct {
	Current string `json:"current"`
}

// Raw types for parsing the full API response
type RawDocument struct {
	CreatedAt       time.Time    `json:"_createdAt"`
	ID              string       `json:"_id"`
	Rev             string       `json:"_rev"`
	Type            string       `json:"_type"`
	UpdatedAt       time.Time    `json:"_updatedAt"`
	Body            []RawBlock   `json:"body"`
	BrodTekst       []RawBlock   `json:"brodtekst"`
	Ord             string       `json:"ord"`
	Relaterte       []RawRelated `json:"relaterte"`
	Sammendrag      []RawBlock   `json:"sammendrag"`
	SammendragTekst string       `json:"sammendragTekst"`
	SistPublisert   string       `json:"sistPublisert"`
	Tittel          string       `json:"tittel"`
	Slug            RawSlug      `json:"slug"`
	Synonymer       []string     `json:"synonymer"`
}

type RawBlock struct {
	Key      string       `json:"_key"`
	Type     string       `json:"_type"`
	Children []RawSpan    `json:"children"`
	MarkDefs []RawMarkDef `json:"markDefs"`
	Style    string       `json:"style"`
	Level    int          `json:"level,omitempty"`
	ListItem string       `json:"listItem,omitempty"`
}

type RawSpan struct {
	Key   string   `json:"_key"`
	Type  string   `json:"_type"`
	Marks []string `json:"marks"`
	Text  string   `json:"text"`
}

type RawMarkDef struct {
	Key       string       `json:"_key"`
	Type      string       `json:"_type"`
	Reference RawReference `json:"reference,omitzero"`
	Href      string       `json:"href,omitempty"`
	Slug      RawSlug      `json:"slug,omitzero"`
}

type RawReference struct {
	Ref  string `json:"_ref"`
	Type string `json:"_type"`
}

type RawRelated struct {
	ID     string  `json:"id"`
	Slug   RawSlug `json:"slug"`
	Tittel string  `json:"tittel"`
}

type RawSlug struct {
	Type    string `json:"_type"`
	Current string `json:"current"`
}

// Clean types for saving
type CleanDocument struct {
	CreatedAt  time.Time      `json:"_createdAt"`
	UpdatedAt  time.Time      `json:"_updatedAt"`
	Body       []CleanBlock   `json:"body"`
	Relaterte  []CleanRelated `json:"relaterte"`
	Sammendrag []CleanBlock   `json:"sammendrag"`
	Tittel     string         `json:"tittel"`
	Slug       CleanSlug      `json:"slug"`
	Synonymer  []string       `json:"synonymer"`
}

type CleanBlock struct {
	Children []CleanSpan    `json:"children"`
	MarkDefs []CleanMarkDef `json:"markDefs"`
	Style    string         `json:"style"`
	Level    int            `json:"level,omitempty"`
	ListItem string         `json:"listItem,omitempty"`
}

type CleanSpan struct {
	Marks []string `json:"marks"`
	Text  string   `json:"text"`
}

type CleanMarkDef struct {
	Reference CleanReference `json:"reference,omitzero"`
	Href      string         `json:"href,omitempty"`
	Slug      CleanSlug      `json:"slug,omitzero"`
}

type CleanReference struct {
	Ref string `json:"_ref"`
}

type CleanRelated struct {
	Slug   CleanSlug `json:"slug"`
	Tittel string    `json:"tittel"`
}

type CleanSlug struct {
	Current string `json:"current"`
}

// Helper functions to clean up the data
func cleanDocument(raw RawDocument) CleanDocument {
	return CleanDocument{
		CreatedAt:  raw.CreatedAt,
		UpdatedAt:  raw.UpdatedAt,
		Body:       cleanBlocks(raw.Body),
		Relaterte:  cleanRelated(raw.Relaterte),
		Sammendrag: cleanBlocks(raw.Sammendrag),
		Tittel:     raw.Tittel,
		Slug:       cleanSlug(raw.Slug),
		Synonymer:  raw.Synonymer,
	}
}

func cleanBlocks(rawBlocks []RawBlock) []CleanBlock {
	result := make([]CleanBlock, len(rawBlocks))
	for i, block := range rawBlocks {
		result[i] = CleanBlock{
			Children: cleanSpansWithMarkDefs(block.Children, block.MarkDefs),
			MarkDefs: cleanMarkDefs(block.MarkDefs),
			Style:    block.Style,
			Level:    block.Level,
			ListItem: block.ListItem,
		}
	}
	return result
}

func cleanSpansWithMarkDefs(rawSpans []RawSpan, rawMarkDefs []RawMarkDef) []CleanSpan {
	// Create mapping from old keys to new indices
	keyToIndex := make(map[string]string)
	for i, markDef := range rawMarkDefs {
		keyToIndex[markDef.Key] = fmt.Sprintf("%d", i)
	}

	result := make([]CleanSpan, len(rawSpans))
	for i, span := range rawSpans {
		cleanMarks := make([]string, len(span.Marks))
		for j, mark := range span.Marks {
			if newIndex, exists := keyToIndex[mark]; exists {
				cleanMarks[j] = newIndex
			} else {
				// Keep non-markDef marks (like "strong", "em") as-is
				cleanMarks[j] = mark
			}
		}
		result[i] = CleanSpan{
			Marks: cleanMarks,
			Text:  span.Text,
		}
	}
	return result
}

func cleanMarkDefs(rawMarkDefs []RawMarkDef) []CleanMarkDef {
	result := make([]CleanMarkDef, len(rawMarkDefs))
	for i, markDef := range rawMarkDefs {
		result[i] = CleanMarkDef{
			Reference: CleanReference{Ref: markDef.Reference.Ref},
			Href:      markDef.Href,
			Slug:      cleanSlug(markDef.Slug),
		}
	}
	return result
}

func cleanRelated(rawRelated []RawRelated) []CleanRelated {
	result := make([]CleanRelated, len(rawRelated))
	for i, related := range rawRelated {
		result[i] = CleanRelated{
			Slug:   cleanSlug(related.Slug),
			Tittel: related.Tittel,
		}
	}
	return result
}

func cleanSlug(rawSlug RawSlug) CleanSlug {
	return CleanSlug{
		Current: rawSlug.Current,
	}
}

func main() {
	log.Println("Starting download of all terms...")

	// Create terms directory
	if err := os.MkdirAll(termsDir, 0750); err != nil {
		log.Fatalf("Failed to create terms directory: %v", err)
	}

	// Get list of all terms
	terms, err := fetchTermsList()
	if err != nil {
		log.Fatalf("Failed to fetch terms list: %v", err)
	}

	log.Printf("Found %d terms to download", len(terms))

	// Download each term
	for i, term := range terms {
		log.Printf("Downloading %d/%d: %s (%s)", i+1, len(terms), term.Title, term.Slug.Current)

		if err := downloadTerm(term); err != nil {
			log.Printf("Failed to download term %s: %v", term.Slug.Current, err)
			continue
		}

		// Small delay to be respectful
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Download completed!")
}

func fetchTermsList() ([]TermListItem, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(baseAPI)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch terms list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var terms []TermListItem
	if err := json.Unmarshal(body, &terms); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return terms, nil
}

func downloadTerm(term TermListItem) error {
	url := fmt.Sprintf("%s%s", baseAPI, term.Slug.Current)

	client := &http.Client{Timeout: 30 * time.Second}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("API returned status %d", resp.StatusCode)
			if attempt < maxRetries && resp.StatusCode >= 500 {
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			return lastErr
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			return fmt.Errorf("failed to read response body after %d attempts: %w", maxRetries, err)
		}

		// Parse raw response
		var rawDoc RawDocument
		if err := json.Unmarshal(body, &rawDoc); err != nil {
			return fmt.Errorf("failed to parse JSON response: %w", err)
		}

		// Clean up the data
		cleanDoc := cleanDocument(rawDoc)

		// Marshal clean data
		cleanBody, err := json.Marshal(cleanDoc)
		if err != nil {
			return fmt.Errorf("failed to marshal clean data: %w", err)
		}

		// Save to file
		filename := filepath.Join(termsDir, fmt.Sprintf("%s.json", term.Slug.Current))
		if err := os.WriteFile(filename, cleanBody, 0600); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filename, err)
		}

		return nil
	}

	return lastErr
}
