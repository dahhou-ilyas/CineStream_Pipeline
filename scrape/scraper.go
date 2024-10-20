package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/tebeka/selenium"
	"go-films-pipline/model"
	"log"
	"os"
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

type Movie struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Rating string `json:"rating"`
	Year   string `json:"year"`
}

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
	// Configuration pour ARM
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

	// Contexte avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Println(ctx)

	// Connexion à Selenium
	seleniumURL := "http://localhost:4444/wd/hub"
	driver, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		log.Fatalf("Erreur de connexion au WebDriver: %v", err)
	}
	defer driver.Quit()

	// Chargement de la page
	if err := driver.Get("https://www.imdb.com/chart/top"); err != nil {
		log.Fatalf("Erreur lors du chargement de la page: %v", err)
	}

	// Attendre que la page soit chargée
	time.Sleep(5 * time.Second) // Peut être amélioré avec une attente explicite

	// Script d'extraction
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

	// Exécution du script
	result, err := driver.ExecuteScript(extractScript, nil)
	if err != nil {
		log.Fatalf("Erreur lors de l'exécution du script: %v", err)
	}

	// Vérification du type de résultat
	resultSlice, ok := result.([]interface{})
	if !ok {
		log.Fatalf("Erreur de type dans le résultat")
	}

	// Conversion des données
	var movies []Movie
	for _, item := range resultSlice {
		if movieMap, ok := item.(map[string]interface{}); ok {
			movie := Movie{
				Title:  fmt.Sprintf("%v", movieMap["title"]),
				URL:    fmt.Sprintf("%v", movieMap["url"]),
				Rating: fmt.Sprintf("%v", movieMap["rating"]),
				Year:   fmt.Sprintf("%v", movieMap["year"]),
			}
			movies = append(movies, movie)
		}
	}

	// Création du fichier JSON
	file, err := os.Create("movies.json")
	if err != nil {
		log.Fatalf("Erreur lors de la création du fichier: %v", err)
	}
	defer file.Close()

	// Encodage en JSON avec indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(movies); err != nil {
		log.Fatalf("Erreur lors de l'encodage JSON: %v", err)
	}

	// Affichage du résumé
	fmt.Printf("Extraction terminée. %d films ont été extraits et sauvegardés dans movies.json\n", len(movies))

	// Affichage des 5 premiers films pour vérification
	fmt.Println("\nAperçu des 5 premiers films :")
	for i, movie := range movies {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s (%s) - Rating: %s\n", i+1, movie.Title, movie.Year, movie.Rating)
	}
}
