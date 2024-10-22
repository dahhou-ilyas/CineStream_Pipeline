package cleaner

import (
	"go-films-pipline/model"
	"time"
)

type MovieEnriched struct {
	model.Movie
	DecadeCategory string `json:"decade_category"`

	GenreCategories []string `json:"genre_categories"`

	DirectorStats Statistics `json:"director_stats"`

	SentimentAnalysis struct {
		PlotTone     string  `json:"plot_tone"`
		KeywordScore float64 `json:"keyword_score"`
	} `json:"sentiment_analysis"`

	Recommendations []string `json:"recommendations"`
}

type Statistics struct {
	TotalMovies int `json:"total_movies"`

	AverageRating float64 `json:"average_rating"`
}

func ProcessMovie(movie model.Movie) MovieEnriched {
	return MovieEnriched{}
}

func categorizeDecade(year string) string {
	if y, err := time.Parse("2006", year); err != nil {
		decad := (y.Year() / 10) * 10
		return string(decad) + "s"
	}

	return "Unknown"
}
