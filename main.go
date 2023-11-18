package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

// ClassDetails holds information about a class
type ClassDetails struct {
	Subject     string
	CatalogNum  string
	Title       string
	Credits     string
	Description string
	// Additional fields like Term, Instructor, etc. can be added as needed
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.uwrf.edu"),
	)

	var classes []ClassDetails

	c.OnHTML("a.colorbox[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "courseLightbox.cfm?subject=CIDS") {
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.OnHTML("div#classSchedule", func(e *colly.HTMLElement) {
		var classDetails ClassDetails
		e.ForEach("table tr", func(_ int, el *colly.HTMLElement) {
			switch el.ChildText("td:nth-child(1)") {
			case "Subject":
				classDetails.Subject = el.ChildText("td:nth-child(2)")
			case "Catalog Number":
				classDetails.CatalogNum = el.ChildText("td:nth-child(2)")
			case "Title":
				classDetails.Title = el.ChildText("td:nth-child(2)")
			case "Credits":
				classDetails.Credits = el.ChildText("td:nth-child(2)")
			case "Details:":
				classDetails.Description = el.ChildText("td:nth-child(1)")
			}
		})
		classes = append(classes, classDetails)
	})

	c.Visit("https://www.uwrf.edu/ClassSchedule/DepartmentCourses.cfm?subject=CIDS")

	jsonData, err := json.Marshal(classes)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Create("classes.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.Write(jsonData)

	fmt.Println("Scraping completed and data saved to classes.json")
}
