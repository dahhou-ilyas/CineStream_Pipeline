package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"go-films-pipline/model"
	"log"
	"sync"
	"time"
)

type Scraper struct {
	baseURL   string
	maxMovies int
	collector *colly.Collector
	movies    []model.Movie
	mu        sync.Mutex
}

func NewScraper(maxMovies int) *Scraper {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(2),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 14_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Safari/605.1.15"),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	return &Scraper{
		baseURL:   "https://www.imdb.com",
		maxMovies: maxMovies,
		collector: c,
		movies:    make([]model.Movie, 0),
	}
}

func main() {
	c := colly.NewCollector()

	url := "https://www.imdb.com/chart/top/"

	c.OnHTML(".cli-parent", func(e *colly.HTMLElement) {
		fmt.Println(e.ChildText("h3.ipc-title__text"))
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
}
