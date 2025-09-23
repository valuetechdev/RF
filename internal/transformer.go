package internal

import (
	"fmt"
	"sort"
	"strings"
)

type MarkdownTransformer struct {
	linkResolver func(string) string
}

func NewMarkdownTransformer() *MarkdownTransformer {
	return &MarkdownTransformer{
		linkResolver: func(ref string) string {
			return fmt.Sprintf("#%s", ref)
		},
	}
}

func (t *MarkdownTransformer) Transform(doc Document) string {
	var result strings.Builder

	if doc.Tittel != "" {
		result.WriteString(fmt.Sprintf("# %s\n\n", doc.Tittel))
	}

	if len(doc.Sammendrag) > 0 {
		result.WriteString("## Summary\n\n")
		for _, block := range doc.Sammendrag {
			result.WriteString(t.transformBlock(block))
		}
		result.WriteString("\n")
	}

	if len(doc.Body) > 0 {
		for _, block := range doc.Body {
			result.WriteString(t.transformBlock(block))
		}
	}

	if len(doc.Synonymer) > 0 {
		result.WriteString("\n## Synonyms\n\n")
		result.WriteString(strings.Join(doc.Synonymer, ", "))
		result.WriteString("\n\n")
	}

	if len(doc.Relaterte) > 0 {
		result.WriteString("## Related Terms\n\n")
		for _, rel := range doc.Relaterte {
			result.WriteString(fmt.Sprintf("- [%s](%s)\n", rel.Tittel, rel.Slug.Current))
		}
		result.WriteString("\n")
	}

	return result.String()
}

func (t *MarkdownTransformer) transformBlock(block Block) string {
	content := t.transformSpans(block.Children, block.MarkDefs)

	switch block.Style {
	case "h1":
		return fmt.Sprintf("# %s\n\n", content)
	case "h2":
		return fmt.Sprintf("## %s\n\n", content)
	case "h3":
		return fmt.Sprintf("### %s\n\n", content)
	case "h4":
		return fmt.Sprintf("#### %s\n\n", content)
	case "h5":
		return fmt.Sprintf("##### %s\n\n", content)
	case "h6":
		return fmt.Sprintf("###### %s\n\n", content)
	case "normal":
		if block.ListItem != "" {
			indent := strings.Repeat("  ", block.Level-1)
			if block.ListItem == "number" {
				return fmt.Sprintf("%s1. %s\n", indent, content)
			}
			return fmt.Sprintf("%s- %s\n", indent, content)
		}
		return fmt.Sprintf("%s\n\n", content)
	default:
		return fmt.Sprintf("%s\n\n", content)
	}
}

func (t *MarkdownTransformer) transformSpans(spans []Span, markDefs []MarkDef) string {
	markDefMap := make(map[string]MarkDef)
	for i, md := range markDefs {
		// Use index as key since we removed the Key field
		markDefMap[fmt.Sprintf("%d", i)] = md
	}

	var result strings.Builder

	for _, span := range spans {
		text := span.Text

		appliedMarks := make(map[string]bool)

		sort.Slice(span.Marks, func(i, j int) bool {
			return getMarkPriority(span.Marks[i]) > getMarkPriority(span.Marks[j])
		})

		for _, mark := range span.Marks {
			if appliedMarks[mark] {
				continue
			}

			if markDef, exists := markDefMap[mark]; exists {
				// Determine type based on content since we removed the Type field
				if markDef.Href != "" {
					// This is a regular link
					text = fmt.Sprintf("[%s](%s)", text, markDef.Href)
				} else if markDef.Reference.Ref != "" || markDef.Slug.Current != "" {
					// This is an ordLink (term reference)
					if markDef.Slug.Current != "" {
						text = fmt.Sprintf("[%s](%s)", text, t.linkResolver(markDef.Slug.Current))
					} else {
						text = fmt.Sprintf("[%s](%s)", text, t.linkResolver(markDef.Reference.Ref))
					}
				}
				appliedMarks[mark] = true
			} else {
				switch mark {
				case "strong":
					text = fmt.Sprintf("**%s**", text)
				case "em":
					text = fmt.Sprintf("*%s*", text)
				case "code":
					text = fmt.Sprintf("`%s`", text)
				case "underline":
					text = fmt.Sprintf("<u>%s</u>", text)
				case "strike-through":
					text = fmt.Sprintf("~~%s~~", text)
				}
				appliedMarks[mark] = true
			}
		}

		result.WriteString(text)
	}

	return result.String()
}

func getMarkPriority(mark string) int {
	priorities := map[string]int{
		"link":           100,
		"ordLink":        100,
		"strong":         50,
		"em":             40,
		"code":           30,
		"underline":      20,
		"strike-through": 10,
	}

	if priority, exists := priorities[mark]; exists {
		return priority
	}
	return 0
}
