package utils

import (
	"fmt"
	"strings"

	"github.com/Ardax-C/uwrf-course-scraper/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func scrapeCoursePage(link string) (models.CourseListing, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.uwrf.edu"),
	)

	var course models.CourseListing
	var currentSection models.Section
	var isCollectingSectionData bool

	c.OnHTML("table", func(e *colly.HTMLElement) {
		// Define a map for section data fields
		fieldMap := map[string]*string{
			"Section":      &currentSection.SectionNum,
			"Class Number": &currentSection.ClassNumber,
			"Term":         &currentSection.Term,
			"Status":       &currentSection.Status,
			"Dates":        &currentSection.Dates,
			"Topic":        &currentSection.Topic,
			"Time":         &currentSection.Time,
			"Instructor":   &currentSection.Instructor,
			"Enrolled":     &currentSection.Enrolled,
			"Room":         &currentSection.Room,
			"Notes":        &currentSection.Notes,
		}

		e.ForEach("table tr", func(i int, el *colly.HTMLElement) {
			if el.Text == "" {
				if isCollectingSectionData {
					course.Sections = append(course.Sections, currentSection)
					currentSection = models.Section{}
					isCollectingSectionData = false
				}
			} else {
				el.DOM.Find("td.text-right.bold").Each(func(_ int, s *goquery.Selection) {
					label := strings.TrimSpace(s.Text())
					value := strings.TrimSpace(s.Next().Text())

					if label == "Section" {
						if isCollectingSectionData {
							course.Sections = append(course.Sections, currentSection)
							currentSection = models.Section{}
						}
						isCollectingSectionData = true
					}

					if ptr, ok := fieldMap[label]; ok {
						*ptr = CleanString(value)
					}
				})
			}
		})

		// Append the last section if it exists
		if isCollectingSectionData {
			course.Sections = append(course.Sections, currentSection)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error scraping:", r.Request.URL, "\nError:", err)
	})

	c.Visit(link)

	return course, nil
}
