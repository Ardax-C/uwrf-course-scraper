package models

type CourseListing struct {
	Subject     string    `json:"subject"`
	CatalogNum  string    `json:"catalog_number"`
	Title       string    `json:"title"`
	Credits     string    `json:"credits"`
	Description string    `json:"description"`
	Sections    []Section `json:"sections,omitempty"`
}

type Section struct {
	SectionNum  string `json:"section"`
	ClassNumber string `json:"class_number"`
	Term        string `json:"term"`
	Status      string `json:"status"`
	Dates       string `json:"dates"`
	Topic       string `json:"topic"`
	Time        string `json:"time"`
	Instructor  string `json:"instructor"`
	Enrolled    string `json:"enrolled"`
	Room        string `json:"room"`
	Notes       string `json:"notes"`
}
