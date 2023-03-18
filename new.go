package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var cacheMutex sync.Mutex
var requestMutex sync.Mutex

type NewAlbum struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

var newAlbumsJSONStr string
var newLastFetchTime time.Time

func fetchNewAlbums() {
	// Lock the mutex to prevent race conditions
	cacheMutex.Lock()

	// If the last fetch was less than 24 hours ago, release the lock and return the cached data
	if time.Since(newLastFetchTime) < 24*time.Hour {
		cacheMutex.Unlock()
		return
	}

	// Define the URL to scrape
	url := "https://www.metacritic.com/browse/albums/release-date/new-releases/date"

	// Send a GET request to the specified URL
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching new albums: %s\n", err)
		return
	}
	defer response.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Printf("Error parsing HTML: %s\n", err)
		return
	}

	// Find all album rows and store the artist and title in a slice of structs
	var newAlbums []NewAlbum
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		artist := strings.TrimSpace(s.Find(".clamp-details .artist").Text())
		artist = strings.TrimPrefix(artist, "by ") // Remove the "by " prefix
		title := strings.TrimSpace(s.Find(".title h3").Text())
		// Remove square brackets and their content at the end of the title
		re := regexp.MustCompile(`\s*\[[^\]]*\]$`)
		title = re.ReplaceAllString(title, "")
		if artist != "" && title != "" && title != "[Title TBA]" {
			newAlbum := NewAlbum{artist, title}
			newAlbums = append(newAlbums, newAlbum)
		}
	})

	// Define a custom struct with the desired JSON order
	type NewResponse struct {
		Title     string     `json:"title"`
		NewAlbums []NewAlbum `json:"albums"`
	}

	// Wrap the newAlbums in a custom struct with a title
	newResponseStruct := NewResponse{
		Title:     "NEW ALBUM RELEASES - METACRITIC",
		NewAlbums: newAlbums,
	}

	// Convert the response to a JSON string with two-space indents
	newAlbumsJSON, err := json.MarshalIndent(newResponseStruct, "", "  ")
	if err != nil {
		log.Printf("Error encoding new albums to JSON: %s\n", err)
		return
	}
	newAlbumsJSONStr = string(newAlbumsJSON)
	newAlbumsJSONStr = strings.Replace(newAlbumsJSONStr, "\\u0026", "&", -1) // Convert HTML escape sequence to "&"
	//newAlbumsJSONStr = strings.ReplaceAll(newAlbumsJSONStr, "'", "\"")

	// Update the last fetch time to the current time
	newLastFetchTime = time.Now()

	log.Println("New albums cache updated.")

	// Release the lock
	mutex.Unlock()
}

func startNewCacheUpdater() {
	for {
		fetchNewAlbums()
		time.Sleep(24 * time.Hour)
	}
}

func handleNewAlbumsRequest(w http.ResponseWriter, r *http.Request) {
	// Lock the request mutex to prevent race conditions
	requestMutex.Lock()
	defer requestMutex.Unlock()

	// Retrieve the client IP address from the X-Forwarded-For header
	remoteAddr := r.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr = r.RemoteAddr
	}

	// Log the incoming request details
	log.Printf("%s %s %s", remoteAddr, r.Method, r.URL)

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the HTTP response writer
	w.Write([]byte(newAlbumsJSONStr))

	// Log the response details
	log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, http.StatusOK)
}
