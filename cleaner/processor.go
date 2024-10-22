package cleaner

import "go-films-pipline/model"

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
