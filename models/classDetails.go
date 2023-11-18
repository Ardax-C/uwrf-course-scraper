package models

type Course struct {
	Subject     string    `json:"subject"`
	CatalogNum  string    `json:"catalog_number"`
	Title       string    `json:"title"`
	Credits     string    `json:"credits"`
	Description string    `json:"description"`
	Sections    []Section `json:"sections,omitempty"`
}
