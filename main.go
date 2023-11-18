package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Ardax-C/uwrf-course-scraper/cmd"
	"github.com/Ardax-C/uwrf-course-scraper/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.uwrf.edu"),
	)

	var classes []models.Course

	c.OnHTML("a.colorbox[href]", func(e *colly.HTMLElement) {
		originalLink := e.Attr("href")

		params := strings.Split(originalLink, "&")
		cleanedParams := make([]string, len(params))

		for i, param := range params {
			parts := strings.Split(param, "=")
			if len(parts) == 2 {
				parts[1] = strings.TrimSpace(parts[1])
			}
			cleanedParams[i] = strings.Join(parts, "=")
		}

		cleanedLink := strings.Join(cleanedParams, "&")

		if strings.Contains(cleanedLink, "courseLightbox.cfm?subject=CIDS") {
			fmt.Println("Visiting:", e.Request.AbsoluteURL(cleanedLink))
			e.Request.Visit(e.Request.AbsoluteURL(cleanedLink))
		}
	})

	c.OnHTML("div#classSchedule", func(e *colly.HTMLElement) {
		var course models.Course

		e.DOM.Find("table").First().Find("tr").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				return
			} else if i == 1 {
				course.Subject = s.Find("td").Eq(0).Text()
				course.CatalogNum = s.Find("td").Eq(1).Text()
				course.Title = s.Find("td").Eq(2).Text()
				course.Credits = cmd.CleanString(s.Find("td").Eq(3).Text())
			} else if i == 2 {
				course.Description = s.Find("td").Eq(0).Text()
			}
		})

		var currentSection models.Section
		var isCollectingSectionData bool

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
						*ptr = cmd.CleanString(value)
					}
				})
			}
		})

		if isCollectingSectionData {
			course.Sections = append(course.Sections, currentSection)
		}

		if course.Subject != "" || course.CatalogNum != "" {
			classes = append(classes, course)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://www.uwrf.edu/ClassSchedule/DepartmentCourses.cfm?subject=CIDS")

	c.Wait()

	jsonData, err := json.MarshalIndent(classes, "", "    ")
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

	if _, err := file.Write(jsonData); err != nil {
		fmt.Println("Error writing JSON to file:", err)
	} else {
		fmt.Println("Scraping completed and data saved to classes.json")
	}
}
