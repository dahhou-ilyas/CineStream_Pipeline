## 📝 Description

CineStream Pipeline est une solution innovante de streaming de données dédiée au traitement en temps réel des informations cinématographiques. Cette plateforme combine la puissance du scraping automatisé avec Selenium, une diffusion événementielle via NATS, et un streaming temps réel via Server-Sent Events (SSE) avec Express.js.

### Vue d'ensemble
- **Collecte de données** : Web scraping automatisé avec Selenium WebDriver
- **Traitement ETL** : Pipeline de transformation sophistiqué
- **Diffusion** : Messaging haute performance avec NATS
- **Streaming** : Flux de données temps réel via SSE
- **API** : Endpoint de streaming Express.js

## 🌟 Fonctionnalités clés

### Scraping avancé avec Selenium
- Navigation dynamique des pages IMDB
- Support des interactions JavaScript complexes
- Extraction robuste des données avec WebDriverWait
- Support des éléments dynamiques et des iframes
- Gestion intelligente des temps de chargement

### Pipeline ETL
- Validation et nettoyage des données
- Enrichissement automatique des métadonnées
- Normalisation des formats
- Gestion des doublons et des conflits

### Messaging temps réel avec NATS
- Communication ultra-rapide et légère
- Support des patterns Pub/Sub et Request/Reply
- Scalabilité horizontale native
- Persistance des messages avec NATS JetStream

### Streaming API avec Express.js et SSE
- Endpoint de streaming unique `/api/movies/stream`
- Server-Sent Events pour une communication temps réel
- Connexion persistante et efficace
- Gestion automatique des reconnexions
- Support des backpressure

## 💻 Prérequis techniques

- Go 1.21 ou supérieur
- Node.js 18+ et npm
- Docker & Docker Compose
- NATS 2.9+
- Selenium WebDriver
- Chrome/Firefox Driver

## 🚀 Installation

1. **Cloner le repository**
```bash
git clone https://github.com/votre-username/cinestream-pipeline.git
cd cinestream-pipeline
```

2. **Installation des WebDrivers**
```bash
# Installation de ChromeDriver
wget https://chromedriver.storage.googleapis.com/[VERSION]/chromedriver_linux64.zip
unzip chromedriver_linux64.zip
sudo mv chromedriver /usr/local/bin/
```

3. **Lancement de l'infrastructure**
```bash
docker-compose up -d
```

4. **Démarrage des services**
```bash
# Scraper
cd scraper
go run main.go

# API
cd ../api
npm install
npm run start
```

## 📊 Exemple de configuration Selenium

```go
// Exemple de configuration du WebDriver
func setupSelenium() (*selenium.WebDriver, error) {
    caps := selenium.Capabilities{
        "browserName": "chrome",
    }
    
    chromeCaps := chrome.Capabilities{
        Args: []string{
            "--headless",
            "--no-sandbox",
            "--disable-dev-shm-usage",
            "--disable-gpu",
            "--window-size=1920,1080",
        },
    }
    
    caps.AddChrome(chromeCaps)
    
    driver, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
    if err != nil {
        return nil, fmt.Errorf("failed to create selenium driver: %w", err)
    }
    
    return driver, nil
}
```

## 📊 Connexion au flux SSE

```javascript
// Exemple de client JavaScript
const eventSource = new EventSource('http://localhost:3000/api/movies/stream');

eventSource.onmessage = (event) => {
    const movie = JSON.parse(event.data);
    console.log('Nouveau film reçu:', movie);
};

eventSource.onerror = (error) => {
    console.error('Erreur SSE:', error);
    eventSource.close();
};
```

### Structure des événements SSE
```json
{
    "id": "msg_123",
    "event": "movie",
    "data": {
        "id": "tt0111161",
        "title": "The Shawshank Redemption",
        "rating": 9.3,
        "year": "1994",
        "release_date": "1994-09-23",
        "director": ["Frank Darabont"],
        "genre": ["Drama", "Crime"],
        "plot": "Two imprisoned men bond over a number of years...",
        "writers": ["Stephen King", "Frank Darabont"],
        "stars": ["Tim Robbins", "Morgan Freeman"],
        "decade_category": "1990s",
        "genre_categories": ["Prison", "Drama", "Adaptation"],
        "director_stats": {
            "total_movies": 4,
            "average_rating": 8.2
        },
        "sentiment_analysis": {
            "plot_tone": "Hopeful",
            "keyword_score": 0.85
        },
        "recommendations": [
        ]
    }
}
```

