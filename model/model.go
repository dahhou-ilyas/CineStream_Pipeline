package model

type Movie struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Rating      float64  `json:"rating"`
	ReleaseDate string   `json:"release_date"`
	Director    []string `json:"director"`
	Genre       []string `json:"genre"`
	Plot        string   `json:"plot"`
	Writers     []string `json:"writers"`
	Stars       []string `json:"stars"`
}
