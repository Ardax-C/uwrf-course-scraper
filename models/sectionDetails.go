package models

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
