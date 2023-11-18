package utils

import (
	"fmt"
	"strings"

	"github.com/Ardax-C/uwrf-course-scraper/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func ScrapeCoursePage(link string) (models.CourseListing, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.uwrf.edu"),
	)

	var course models.CourseListing
	var currentSection models.Section
	var isCollectingSectionData bool

	// Parsing main course information
	c.OnHTML("div#classSchedule", func(e *colly.HTMLElement) {
		// Assuming the first row contains subject, catalog number, title, and credits
		e.DOM.Find("table").First().Find("tr").Each(func(i int, s *goquery.Selection) {
			if i == 1 { // Skip the header, start with second row
				course.Subject = s.Find("td").Eq(0).Text()
				course.CatalogNum = s.Find("td").Eq(1).Text()
				course.Title = s.Find("td").Eq(2).Text()
				course.Credits = CleanString(s.Find("td").Eq(3).Text())
			} else if i == 2 {
				// Extract course description
				course.Description = s.Find("td").Eq(0).Text()
			}
		})
	})

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
