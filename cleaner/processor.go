package cleaner

import (
	"go-films-pipline/model"
	"strings"
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

func categorizeGenres(genres []string) []string {
	categories := make([]string, 0)

	mainCategories := map[string][]string{
		"Action-Adventure": {"Action", "Adventure", "Thriller"},
		"Drama-Romance":    {"Drama", "Romance"},
		"Comedy-Family":    {"Comedy", "Family", "Animation"},
		"Sci-Fi-Fantasy":   {"Science Fiction", "Fantasy", "Sci-Fi"},
	}

	genreSet := make(map[string]bool)
	for _, genre := range genres {
		genreSet[genre] = true
	}

	for category, relatedGenres := range mainCategories {
		for _, genre := range relatedGenres {
			if genreSet[genre] {
				categories = append(categories, category)
				break
			}
		}
	}

	return categories
}

func analyzePlotTone(plot string) string {
	plot = strings.ToLower(plot)

	toneKeywords := map[string][]string{
		"Dark": {
			// Violence et Crime
			"death", "murder", "kill", "violent", "blood", "crime", "criminal",
			"deadly", "lethal", "massacre", "assassin", "sinister",
			// Atmosphère sombre
			"dark", "darkness", "grim", "bleak", "noir", "shadow", "haunting",
			// Tragédie
			"tragedy", "tragic", "suffering", "doom", "despair", "misery",
			// Horreur
			"horror", "terrifying", "nightmare", "evil", "demon", "monster",
			// Psychologique
			"paranoid", "insanity", "madness", "twisted", "disturbing", "psychological",
		},

		"Light": {
			// Comédie
			"comedy", "funny", "humor", "laugh", "hilarious", "witty", "joke",
			"amusing", "comedy", "comedic", "goofy", "silly",
			// Émotions positives
			"happy", "joy", "delight", "cheerful", "uplifting", "optimistic",
			"playful", "fun", "light-hearted", "whimsical",
			// Romance légère
			"romantic", "charming", "sweet", "lovely", "heartwarming",
			// Aventure positive
			"adventure", "exciting", "wonderful", "magical", "enchanting",
			// Family-friendly
			"family", "friendship", "innocent", "gentle", "wholesome",
		},

		"Dramatic": {
			// Drame émotionnel
			"drama", "emotional", "intense", "powerful", "moving", "touching",
			"profound", "deep", "meaningful", "serious",
			// Conflits personnels
			"struggle", "conflict", "challenge", "overcome", "perseverance",
			"determination", "ambition", "rivalry",
			// Relations
			"relationship", "family", "love", "betrayal", "reconciliation",
			"sacrifice", "loyalty", "redemption",
			// Transformation
			"journey", "change", "growth", "transformation", "discovery",
			// Société
			"social", "political", "cultural", "historical", "revolution",
			"war", "justice", "inequality", "prejudice",
		},

		"Suspense": {
			// Mystère
			"mystery", "enigma", "secret", "clue", "investigation", "detective",
			"puzzle", "unsolved", "mysterious", "hidden",
			// Tension
			"suspense", "tension", "thriller", "anticipation", "anxiety",
			"nervous", "paranoia", "suspicious", "uncertain",
			// Action
			"chase", "escape", "pursuit", "race", "hunt", "dangerous",
			"risk", "threat", "deadline", "countdown",
			// Intrigue
			"conspiracy", "plot", "scheme", "deception", "betrayal", "spy",
			"espionage", "infiltration", "sabotage",
			// Atmosphère
			"tense", "gripping", "edge", "breathtaking", "shocking", "twist",
		},

		"Epic": {
			// Grande échelle
			"epic", "grand", "massive", "spectacular", "monumental",
			"legendary", "mythical", "saga", "empire",
			// Bataille et guerre
			"battle", "war", "conquest", "victory", "defeat", "warrior",
			"hero", "army", "kingdom", "throne",
			// Aventure épique
			"quest", "journey", "expedition", "discovery", "exploration",
			"adventure", "mission", "destiny",
			// Fantasy/Science-Fiction
			"magic", "fantasy", "supernatural", "alien", "futuristic",
			"space", "cosmic", "mythological", "dragon",
			// Thèmes épiques
			"destiny", "fate", "power", "glory", "honor", "legacy",
		},

		"Philosophical": {
			// Questions existentielles
			"existential", "philosophical", "metaphysical", "consciousness",
			"reality", "existence", "truth", "meaning", "purpose",
			// Réflexion
			"contemplative", "thoughtful", "meditation", "introspection",
			"questioning", "wisdom", "understanding",
			// Concepts abstraits
			"time", "memory", "dream", "perception", "identity", "soul",
			"mind", "reality", "illusion", "nature",
			// Thèmes sociaux profonds
			"humanity", "society", "civilization", "progress", "decay",
			"morality", "ethics", "ideology", "belief",
			// Exploration psychologique
			"psychological", "subconscious", "mind", "awareness", "self",
		},

		"Satirical": {
			// Satire sociale
			"satire", "parody", "mockery", "irony", "sarcasm", "criticism",
			"cynical", "satirical", "ridicule", "exaggeration",
			// Commentaire social
			"commentary", "critique", "absurd", "bizarre", "eccentric",
			"unconventional", "controversial", "provocative",
			// Humour noir
			"dark humor", "black comedy", "gallows humor", "wit", "clever",
			"sharp", "biting", "caustic", "sardonic",
			// Critique culturelle
			"stereotype", "cliché", "convention", "norm", "tradition",
			"society", "culture", "media", "politics",
		},
	}

	toneScores := make(map[string]float64)
	wordCount := len(strings.Fields(plot))

	for tone, keywords := range toneKeywords {
		score := 0.0
		for _, keyword := range keywords {
			occurrences := strings.Count(plot, keyword)
			if occurrences > 0 {
				score += float64(occurrences) / float64(wordCount) * 100
			}
		}
		toneScores[tone] = score
	}

	var (
		primaryTone    string
		primaryScore   float64
		secondaryTone  string
		secondaryScore float64
	)

	for ton, score := range toneScores {
		if score > primaryScore {
			secondaryTone = primaryTone
			primaryTone = ton
			secondaryScore = primaryScore
			primaryScore = score
		} else if score > secondaryScore {
			secondaryTone = ton
			secondaryScore = score
		}
	}

	if secondaryScore >= primaryScore*0.5 {
		return primaryTone + "-" + secondaryTone
	}

	return primaryTone
}
