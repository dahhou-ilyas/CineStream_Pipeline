package main

import (
	"context"
	"go-films-pipline/cleaner"
	"go-films-pipline/model"
	"go-films-pipline/natsProducers"
	"go-films-pipline/scrape"
	"os"
	"os/signal"
	"sync"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigChannel := make(chan os.Signal, 1)

		signal.Notify(sigChannel, os.Interrupt)

		<-sigChannel
		close(sigChannel)
		cancel()
	}()

	wg := new(sync.WaitGroup)

	filmsChannel := make(chan model.Movie)

	wg.Add(1)
	go scrape.ScrapeIMDB(wg, filmsChannel)

	go func() {
		for film := range filmsChannel {
			cleaner.CleanMovieData(&film)
			enrichedFilm := cleaner.ProcessMovie(film)

			natsProducers.Producer(ctx, enrichedFilm)
		}
	}()
	wg.Wait()

}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
