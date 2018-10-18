package football

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// SendMatches will send the scores based on date
func SendMatches(date string) string {
	i := 0
	url := fmt.Sprintf("https://www.goal.com/en-in/fixtures/%s/", date)
	fmt.Println(url)
	c := colly.NewCollector()
	var buffer bytes.Buffer

	// Find and visit all links
	c.OnHTML("div.match-row", func(e *colly.HTMLElement) {
		i++
		// Extract the link from the anchor HTML element
		buffer.WriteString(standardizeSpaces(e.Text))
		if i%3 == 0 {
			buffer.WriteString("|")
		} else {
			buffer.WriteString("\u000A")
		}
	})

	c.Visit(url)
	return buffer.String()
}
