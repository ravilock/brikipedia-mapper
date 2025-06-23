package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	visitedUrls := sync.Map{}

	brazilianWikipediaDomains := colly.AllowedDomains("pt.wikipedia.org")
	indexCollector := colly.NewCollector(brazilianWikipediaDomains)
	filterCollector := indexCollector.Clone()
	pageCollector := indexCollector.Clone()

	indexCollector.OnHTML("#toc a[href]", func(e *colly.HTMLElement) {
		filterCollector.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	indexCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Index Collector Visiting", r.URL)
	})

	filterCollector.OnHTML(".mw-allpages-chunk a[href]", func(e *colly.HTMLElement) {
		// get class ofelement
		classAttr := e.Attr("class")
		// skip if link is a redirect to another page
		if strings.Contains(classAttr, "mw-redirect") {
			return
		}
		if _, visited := visitedUrls.Load(e.Request.AbsoluteURL(e.Attr("href"))); visited {
			fmt.Println("Already visited:", e.Request.AbsoluteURL(e.Attr("href")))
			os.Exit(0) // Exit if we revisit a URL
			return
		}
		pageCollector.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	filterCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Filter Collector Visiting", r.URL)
	})

	pageCollector.OnRequest(func(r *colly.Request) {
		visitedUrls.Store(r.URL.String(), true)
		fmt.Println("Page Collector Visiting", r.URL)
	})

	pageCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error visiting", r.Request.URL, ":", err)
		os.Exit(1) // Exit on error
	})

	indexCollector.Visit("https://pt.wikipedia.org/wiki/Portal:Conte%C3%BAdo/%C3%8Dndice_alfab%C3%A9tico")
	for {
		time.Sleep(100 * time.Millisecond) // Keep the program running to allow asynchronous processing
	}
}
