package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Ardax-C/uwrf-course-scraper/models"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.uwrf.edu"),
	)

	var classes []models.ClassDetails

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
		var classDetails models.ClassDetails

		e.ForEach("table tr.tr-background", func(i int, el *colly.HTMLElement) {
			if i == 1 {
				classDetails.Subject = el.ChildText("td:nth-child(1)")
				classDetails.CatalogNum = strings.TrimSpace(el.ChildText("td:nth-child(2)"))
				classDetails.Title = el.ChildText("td:nth-child(3)")
			}
		})

		classDetails.Description = e.ChildText("table tr:nth-child(3) td[colspan='3']")

		classes = append(classes, classDetails)
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
