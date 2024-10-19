package model

type Movie struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Rating      float64 `json:"rating"`
	ReleaseDate string  `json:"release_date"`
	Duration    int     `json:"duration"`
	Director    string  `json:"director"`
	Genre       string  `json:"genre"`
	Plot        string  `json:"plot"`
	PosterURL   string  `json:"posterUrl"`
}
