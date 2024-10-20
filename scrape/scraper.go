package main

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/tebeka/selenium"
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

func (s *Scraper) ScrapeMovies() ([]model.Movie, error) {
	c := s.collector.Clone()

	count := 0

	c.OnHTML("td.titleColumn a", func(e *colly.HTMLElement) {
		if count >= s.maxMovies {
			return
		}

		movieURL := fmt.Sprintf("%s%s", s.baseURL, e.Attr("href"))
		s.scrapeMovieDetails(movieURL)
		count++
	})

	err := c.Visit(s.baseURL + "/chart/top/")
	if err != nil {
		return nil, fmt.Errorf("error visiting IMDB: %w", err)
	}
	c.Wait()
	return s.movies, nil

}

func (s *Scraper) scrapeMovieDetails(url string) {

}

type Film struct {
	Title  string
	URL    string
	Rating string
	Year   string
}

func scrapeIMDB(wg *sync.WaitGroup, films chan<- Film) {

	defer wg.Done()

	caps := selenium.Capabilities{
		"browserName":  "chrome",
		"platformName": "Linux",
		"goog:chromeOptions": map[string]interface{}{
			"args": []string{
				"--no-sandbox",
				"--disable-dev-shm-usage",
				"--disable-gpu",
			},
		},
	}

	_, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	seleniumURL := "http://localhost:4444/wd/hub"
	driver, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		log.Fatalf("Erreur de connexion au WebDriver: %v", err)
	}
	defer driver.Quit()

	if err := driver.Get("https://www.imdb.com/chart/top"); err != nil {
		log.Fatalf("Erreur lors du chargement de la page: %v", err)
	}

	time.Sleep(5 * time.Second)

	extractScript := `
        return Array.from(document.querySelectorAll('li.ipc-metadata-list-summary-item'))
        .map(el => {
            const titleElement = el.querySelector('h3.ipc-title__text');
            const linkElement = el.querySelector('a.ipc-title-link-wrapper');
            const ratingElement = el.querySelector('span.ipc-rating-star--imdb');
            const yearElement = el.querySelector('span.cli-title-metadata-item');
            
            let title = '';
            if (titleElement) {
                // Suppression du numéro de classement au début du titre
                title = titleElement.textContent.replace(/^\d+\.\s*/, '').trim();
            }

            return {
                title: title,
                url: linkElement ? linkElement.href : '',
                rating: ratingElement ? ratingElement.textContent.trim() : '',
                year: yearElement ? yearElement.textContent.trim() : '',
            };
        });
    `

	result, err := driver.ExecuteScript(extractScript, nil)
	if err != nil {
		log.Fatalf("Erreur lors de l'exécution du script: %v", err)
	}

	resultSlice, ok := result.([]interface{})
	if !ok {
		log.Fatalf("Erreur de type dans le résultat")
	}

	for _, item := range resultSlice {
		if movieMap, ok := item.(map[string]interface{}); ok {
			fmt.Println(movieMap["title"].(string))
			films <- Film{
				Title:  fmt.Sprintf("%v", movieMap["title"]),
				URL:    fmt.Sprintf("%v", movieMap["url"]),
				Rating: fmt.Sprintf("%v", movieMap["rating"]),
				Year:   fmt.Sprintf("%v", movieMap["year"]),
			}
		}
	}
}

func main() {
	wg := new(sync.WaitGroup)

	filmsChannel := make(chan Film)

	wg.Add(1)
	go scrapeIMDB(wg, filmsChannel)

	go func() {
		for film := range filmsChannel {
			fmt.Printf("Title: %s, Year: %s, Rating: %s\n", film.Title, film.Year, film.Rating)
		}
	}()

	wg.Wait()

}
