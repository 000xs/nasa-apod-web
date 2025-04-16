package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ApodResponse struct {
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Url         string `json:"url"`
}

func FetchApod() (*ApodResponse, error) {
	// Load API key from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("NASA_API_KEY")
	url := fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apiKey)

	// Send GET request to NASA's API
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var data ApodResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	// Fetch APOD data from NASA API
	apodData, err := FetchApod()
	if err != nil {
		http.Error(w, "Failed to fetch APOD data", http.StatusInternalServerError)
		return
	}

	// Parse and render the HTML template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to parse HTML template", http.StatusInternalServerError)
		return
	}

	// Pass APOD data to the template and render
	tmpl.Execute(w, apodData)
}

func main() {
	// Set up HTTP routes
	http.HandleFunc("/", HomePage)

	// Start the web server
	log.Println("Starting server on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
