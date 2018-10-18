package football

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// SendMatches will send the scores based on date
func SendMatches(date string) string {
	url := fmt.Sprintf("https://www.goal.com/en-in/fixtures/%s/", date)
	c := colly.NewCollector()
	var sb strings.Builder
	// Find and visit all links
	c.OnHTML("div.match-row", func(e *colly.HTMLElement) {
		// Extract the link from the anchor HTML element
		sb.WriteString(e.Text)
		sb.WriteString("\n")
	})

	c.Visit(url)
	return sb.String()
}
