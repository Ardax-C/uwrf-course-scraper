package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Ardax-C/uwrf-course-scraper/models"
	"github.com/Ardax-C/uwrf-course-scraper/utils"
	"github.com/gocolly/colly"
)

func Init() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.uwrf.edu"),
	)

	var classes []models.CourseListing
	var wg sync.WaitGroup
	var mu sync.Mutex // For thread-safe append to classes slice

	// Limit the number of concurrent goroutines
	semaphore := make(chan struct{}, 10) // Adjust the concurrency level

	c.OnHTML("a.colorbox[href]", func(e *colly.HTMLElement) {
		rawLink := e.Request.AbsoluteURL(e.Attr("href"))

		// Clean the URL
		cleanedLink, err := utils.CleanURL(rawLink)
		if err != nil {
			fmt.Println("Error cleaning URL:", rawLink, "Error:", err)
			return
		}

		if utils.IsValidLink(cleanedLink) {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(link string) {
				defer wg.Done()
				defer func() { <-semaphore }()

				course, err := utils.ScrapeCoursePage(link)
				if err != nil {
					fmt.Println("Error scraping:", link, "Error:", err)
					return
				}

				mu.Lock()
				classes = append(classes, course)
				mu.Unlock()
			}(cleanedLink)
		}
	})

	c.Visit("https://www.uwrf.edu/ClassSchedule/DepartmentCourses.cfm?subject=CIDS")

	wg.Wait()

	jsonData, err := json.MarshalIndent(classes, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	file, err := os.Create("classes.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	if _, err := file.Write(jsonData); err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("Scraping completed and data saved to classes.json")
}
