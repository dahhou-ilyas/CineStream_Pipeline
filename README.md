## 📝 Description

CineStream Pipeline is an innovative data streaming solution dedicated to real-time processing of cinematographic information. This platform combines the power of automated scraping with Selenium, event-driven communication via NATS, and real-time streaming via Server-Sent Events (SSE) with Express.js.

### Overview
- **Data Collection**: Automated web scraping with Selenium WebDriver
- **ETL Processing**: Sophisticated transformation pipeline
- **Distribution**: High-performance messaging with NATS
- **Streaming**: Real-time data flow via SSE
- **API**: Express.js streaming endpoint

## 🌟 Key Features

### Advanced Scraping with Selenium
- Dynamic IMDB page navigation
- Support for complex JavaScript interactions
- Robust data extraction with WebDriverWait
- Support for dynamic elements and iframes
- Intelligent loading time management

### ETL Pipeline
- Data validation and cleaning
- Automatic metadata enrichment
- Format normalization
- Duplicate and conflict management

### Real-time Messaging with NATS
- Ultra-fast and lightweight communication
- Support for Pub/Sub and Request/Reply patterns
- Native horizontal scalability
- Message persistence with NATS JetStream

### Streaming API with Express.js and SSE
- Single streaming endpoint `/api/movies/stream`
- Server-Sent Events for real-time communication
- Persistent and efficient connection
- Automatic reconnection handling
- Backpressure support

## 💻 Technical Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm
- Docker & Docker Compose
- NATS 2.9+
- Selenium WebDriver
- Chrome/Firefox Driver

## 🚀 Installation

1. **Clone the repository**
```bash
git clone https://github.com/your-username/cinestream-pipeline.git
cd cinestream-pipeline
```

2. **WebDrivers Installation**
```bash
# ChromeDriver Installation
wget https://chromedriver.storage.googleapis.com/[VERSION]/chromedriver_linux64.zip
unzip chromedriver_linux64.zip
sudo mv chromedriver /usr/local/bin/
```

3. **Launch Infrastructure**
```bash
docker-compose up -d
```

4. **Start Services**
```bash
# Scraper
cd scraper
go run main.go

# API
cd ../api
npm install
npm run start
```

## 📊 Selenium Configuration Example

```go
// WebDriver configuration example
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

## 📊 SSE Stream Connection

```javascript
// JavaScript client example
const eventSource = new EventSource('http://localhost:3000/api/movies/stream');

eventSource.onmessage = (event) => {
    const movie = JSON.parse(event.data);
    console.log('New movie received:', movie);
};

eventSource.onerror = (error) => {
    console.error('SSE Error:', error);
    eventSource.close();
};
```

### SSE Event Structure
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
