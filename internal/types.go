package internal

import "time"

type Document struct {
	CreatedAt  time.Time `json:"_createdAt"`
	ID         string    `json:"_id"`
	Rev        string    `json:"_rev"`
	Type       string    `json:"_type"`
	UpdatedAt  time.Time `json:"_updatedAt"`
	Body       []Block   `json:"body"`
	BrodTekst  []Block   `json:"brodtekst"`
	Ord        string    `json:"ord"`
	Relaterte  []Related `json:"relaterte"`
	Sammendrag []Block   `json:"sammendrag"`
	Tittel     string    `json:"tittel"`
	Slug       Slug      `json:"slug"`
	Synonymer  []string  `json:"synonymer"`
}

type Block struct {
	Key      string    `json:"_key"`
	Type     string    `json:"_type"`
	Children []Span    `json:"children"`
	MarkDefs []MarkDef `json:"markDefs"`
	Style    string    `json:"style"`
	Level    int       `json:"level,omitempty"`
	ListItem string    `json:"listItem,omitempty"`
}

type Span struct {
	Key   string   `json:"_key"`
	Type  string   `json:"_type"`
	Marks []string `json:"marks"`
	Text  string   `json:"text"`
}

type MarkDef struct {
	Key       string    `json:"_key"`
	Type      string    `json:"_type"`
	Reference Reference `json:"reference,omitempty"`
	Href      string    `json:"href,omitempty"`
	Slug      Slug      `json:"slug,omitempty"`
}

type Reference struct {
	Ref  string `json:"_ref"`
	Type string `json:"_type"`
}

type Related struct {
	ID     string `json:"id"`
	Slug   Slug   `json:"slug"`
	Tittel string `json:"tittel"`
}

type Slug struct {
	Type    string `json:"_type"`
	Current string `json:"current"`
}
