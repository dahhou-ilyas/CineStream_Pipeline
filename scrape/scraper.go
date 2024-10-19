package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/tebeka/selenium"

	"go-films-pipline/model"
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

func main() {
	// Configuration du client Selenium
	const (
		port = 4444 // Le port où Selenium écoute
	)

	// Connexion au WebDriver
	wd, err := selenium.NewRemote(selenium.Capabilities{"browserName": "chrome"}, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit() // Fermer le navigateur à la fin

	// Accéder à une page
	if err := wd.Get("http://google.com"); err != nil {
		panic(err)
	}

	// Exemple d'attente pour s'assurer que la page est chargée
	time.Sleep(2 * time.Second)

	// Extraire le titre de la page
	title, err := wd.Title()

	if err != nil {
		panic(err)
	}
	fmt.Println("Title:", title)
}
