package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/tebeka/selenium"
	"log"

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

type Film struct {
	Title  string
	URL    string
	Rating string
	Year   string
}

const (
	seleniumURL = "http://172.17.0.2:4444/wd/hub"
	maxRetries  = 3
	waitTime    = 5 * time.Second
)

func scrapeIMDB(wg *sync.WaitGroup, films chan<- Film) {

	defer wg.Done()

	driver, err := selenium.NewRemote(selenium.Capabilities{"browserName": "chrome"}, "http://localhost:4444/wd/hub")

	if err != nil {
		log.Printf("Erreur de connexion au WebDriver: %v", err)
		return
	}

	defer driver.Quit()

	// Charger une page
	if err := driver.Get("https://www.imdb.com/chart/top/"); err != nil {
		log.Printf("Erreur lors du chargement de la page: %v", err)
		return
	}

	script := `
        window.scrollTo(0, document.body.scrollHeight);
        return document.body.scrollHeight;
    `

	// Faire défiler plusieurs fois pour charger tout le contenu
	for i := 0; i < 5; i++ {
		if _, err := driver.ExecuteScript(script, nil); err != nil {
			log.Printf("Erreur lors du défilement: %v", err)
			return
		}
		time.Sleep(2 * time.Second) // Attendre un peu plus longtemps
	}

	extractScript := `
        return Array.from(document.querySelectorAll('li.ipc-metadata-list-summary-item'))
    	.map(el => {
        	const titleElement = el.querySelector('h3.ipc-title__text');
        	const linkElement = el.querySelector('a.ipc-title-link-wrapper');
        	const ratingElement = el.querySelector('span.ipc-rating-star--imdb');
        	const yearElement = el.querySelector("span.cli-title-metadata-item");
        	return {
            	title: titleElement ? titleElement.textContent.trim() : '',
            	url: linkElement ? linkElement.href : '',
            	rating: ratingElement ? ratingElement.textContent.trim() : '',
            	year: yearElement ? yearElement.textContent.trim() : '',
			};
    	});
    `

	result, err := driver.ExecuteScript(extractScript, nil)
	fmt.Println(result)
	if err != nil {
		log.Printf("Erreur lors de l'extraction: %v", err)
		return
	}
	fmt.Println(result)

	// Traiter les résultats
	if movies, ok := result.([]interface{}); ok {
		for _, movie := range movies {
			if movieMap, ok := movie.(map[string]interface{}); ok {
				films <- Film{
					Title:  movieMap["title"].(string),
					URL:    movieMap["url"].(string),
					Rating: movieMap["rating"].(string),
					Year:   movieMap["year"].(string),
				}
			}
		}
	} else {
		log.Println("Le résultat de l'extraction n'est pas un tableau.")
	}
}

func main() {
	start := time.Now()

	// Créer un canal pour les films
	films := make(chan Film, 253)
	var wg sync.WaitGroup

	// Lancer le scraping
	wg.Add(1)
	go scrapeIMDB(&wg, films)

	// Goroutine pour collecter les résultats
	go func() {
		wg.Wait()
		close(films)
	}()

	// Collecter et afficher les résultats
	var count int
	for film := range films {
		count++
		fmt.Printf("%d. %s (%s) - Note: %s\n    URL: %s\n\n",
			count, film.Title, film.Year, film.Rating, film.URL)
	}

	duration := time.Since(start)
	fmt.Printf("\nScraping terminé en %s\nNombre total de films : %d\n", duration, count)
}
