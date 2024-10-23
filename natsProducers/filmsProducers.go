package main

import (
	"fmt"
	"go-films-pipline/model"
	"go-films-pipline/scrape"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)

	filmsChannel := make(chan model.Movie)

	wg.Add(1)
	go scrape.ScrapeIMDB(wg, filmsChannel)

	go func() {
		for film := range filmsChannel {
			fmt.Printf("Title: %s, Year: %s, Rating: %s\n", film.Title, film.Year, film.Rating)
		}
	}()

	wg.Wait()
}
