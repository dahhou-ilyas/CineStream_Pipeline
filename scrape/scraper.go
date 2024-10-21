package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tebeka/selenium"
	"go-films-pipline/model"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Film struct {
	Title  string
	URL    string
	Rating string
	Year   string
}

func ScrapeIMDB(wg *sync.WaitGroup, films chan<- Film) {

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

			url := movieMap["url"].(string)

			if err := driver.Get(url); err != nil {
				log.Fatalf("Erreur dans le lien de la film spécifique: %v", err)
			}
			time.Sleep(1 * time.Second)

			execScript := `
				const elemts=document.querySelector("ul.title-pc-list").querySelectorAll("li");
				const genre = Array.from(document.querySelector("div.ipc-chip-list__scroller").querySelectorAll("a span")).map(e => e.textContent);
				const plot = document.querySelector('p[data-testid="plot"] span').textContent
				return {
					"director":elemts[0].textContent.split(" "),
					"writers":elemts[1].textContent.split(" "),
					"stars":elemts[2].textContent.split(" "),
					"genre":genre,
					"plot":plot,
				}
			`

			result1, err := driver.ExecuteScript(execScript, nil)
			if err != nil {
				log.Fatalf("Erreur lors de l'exécution du script de info de films specifique: %v", err)
			}

			souInfoMovieMap, ok := result1.(map[string]interface{})
			if !ok {
				log.Fatalf("Erreur: résultat inattendu")
			}

			ratingNumb, err := strconv.ParseFloat(strings.Fields(movieMap["rating"].(string))[0], 64)
			if err != nil {
				fmt.Println("Erreur de conversion xxxxxxx:", err)
			}

			movie := model.Movie{
				ID:          uuid.New().String(),
				Title:       movieMap["title"].(string),
				Rating:      ratingNumb,
				ReleaseDate: movieMap["year"].(string),
				Director: func(directorInterfaces []interface{}) []string {
					strSlice := make([]string, len(directorInterfaces))
					for i, v := range directorInterfaces {
						strSlice[i] = v.(string)
					}
					return strSlice
				}(souInfoMovieMap["director"].([]interface{})),
				Genre: func(genreInterfaces []interface{}) []string {
					strSlice := make([]string, len(genreInterfaces))
					for i, v := range genreInterfaces {
						strSlice[i] = v.(string)
					}
					return strSlice
				}(souInfoMovieMap["genre"].([]interface{})),
				Plot: souInfoMovieMap["plot"].(string),
				Stars: func(starInterfaces []interface{}) []string {
					strSlice := make([]string, len(starInterfaces))
					for i, v := range starInterfaces {
						strSlice[i] = v.(string)
					}
					return strSlice
				}(souInfoMovieMap["stars"].([]interface{})),
				Writers: func(writerInterfaces []interface{}) []string {
					strSlice := make([]string, len(writerInterfaces))
					for i, v := range writerInterfaces {
						strSlice[i] = v.(string)
					}
					return strSlice
				}(souInfoMovieMap["writers"].([]interface{})),
			}
			fmt.Println(movie.Writers)

		}
	}
}

func main() {
	wg := new(sync.WaitGroup)

	filmsChannel := make(chan Film)

	wg.Add(1)
	go ScrapeIMDB(wg, filmsChannel)

	wg.Wait()

}
