package internal

import "time"

type Document struct {
	CreatedAt  time.Time `json:"_createdAt"`
	UpdatedAt  time.Time `json:"_updatedAt"`
	Body       []Block   `json:"body"`
	Relaterte  []Related `json:"relaterte"`
	Sammendrag []Block   `json:"sammendrag"`
	Tittel     string    `json:"tittel"`
	Slug       Slug      `json:"slug"`
	Synonymer  []string  `json:"synonymer"`
}

type Block struct {
	Children []Span    `json:"children"`
	MarkDefs []MarkDef `json:"markDefs"`
	Style    string    `json:"style"`
	Level    int       `json:"level,omitempty"`
	ListItem string    `json:"listItem,omitempty"`
}

type Span struct {
	Marks []string `json:"marks"`
	Text  string   `json:"text"`
}

type MarkDef struct {
	Reference Reference `json:"reference,omitzero"`
	Href      string    `json:"href,omitempty"`
	Slug      Slug      `json:"slug,omitzero"`
}

type Reference struct {
	Ref string `json:"_ref"`
}

type Related struct {
	Slug   Slug   `json:"slug"`
	Tittel string `json:"tittel"`
}

type Slug struct {
	Current string `json:"current"`
}

type AccountGroup struct {
	Nummer  int       `json:"nummer"`
	Navn    string    `json:"navn"`
	Kontoer []Account `json:"kontoer"`
}

type Account struct {
	Kontonummer   int         `json:"kontonummer"`
	Navn          string      `json:"navn"`
	KunForOrgForm []string    `json:"kunForOrgForm"`
	MetaData      AccountMeta `json:"metaData"`
}

type AccountMeta struct {
	Hjelpetekst           string   `json:"hjelpetekst"`
	Advarsel              string   `json:"advarsel"`
	Eksempler             string   `json:"eksempler"`
	GyldigeMvakoder       []string `json:"gyldigeMvakoder"`
	DefaultMvakode        *string  `json:"defaultMvakode"`
	DefaultMvakodeErIngen bool     `json:"defaultMvakodeErIngen"`
	Labels                []string `json:"labels"`
	Statusfarge           string   `json:"statusfarge"`
	Sokeord               []string `json:"sokeord"`
}
