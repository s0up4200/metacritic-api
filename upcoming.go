package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Album struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

type Response struct {
	Title  string  `json:"title"`
	Albums []Album `json:"albums"`
}

var albumsJSONStr string
var lastFetchTime time.Time

func fetchAlbums() {
	// If the last fetch was less than 24 hours ago, return the cached data
	if time.Since(lastFetchTime) < 24*time.Hour {
		return
	}

	// Define the URL to scrape
	url := "https://www.metacritic.com/browse/albums/release-date/coming-soon"

	// Send a GET request to the specified URL
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching albums: %s\n", err)
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
	var albums []Album
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		artist := strings.TrimSpace(s.Find(".artistName a").Text())
		if artist == "" {
			artist = strings.TrimSpace(s.Find(".artistName").Text())
		}
		title := strings.TrimSpace(s.Find(".albumTitle").Text())
		// Remove square brackets and their content at the end of the title
		re := regexp.MustCompile(`\s*\[[^\]]*\]$`)
		title = re.ReplaceAllString(title, "")
		if artist != "" && title != "" && title != "[Title TBA]" {
			album := Album{artist, title}
			albums = append(albums, album)
		}
	})

	// Wrap the albums in a custom struct with a title
	responseStruct := Response{
		Title:  "UPCOMING ALBUM RELEASES - METACRITIC",
		Albums: albums,
	}

	// Convert the response to a JSON string with two-space indents
	albumsJSON, err := json.MarshalIndent(responseStruct, "", "  ")
	if err != nil {
		log.Printf("Error encoding albums to JSON: %s\n", err)
		return
	}
	albumsJSONStr = string(albumsJSON)
	albumsJSONStr = strings.Replace(albumsJSONStr, "\\u0026", "&", -1) // Convert HTML escape sequence to "&"
	//albumsJSONStr = strings.ReplaceAll(albumsJSONStr, "'", "\"")

	// Update the last fetch time to the current time
	lastFetchTime = time.Now()

	log.Println("Albums cache updated.")
}

func startCacheUpdater() {
	for {
		fetchAlbums()
		time.Sleep(24 * time.Hour)
	}
}

func handleAlbumsRequest(w http.ResponseWriter, r *http.Request) {
	// Lock the mutex to prevent race conditions
	mutex.Lock()
	defer mutex.Unlock()

	// Log the incoming request details
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the HTTP response writer
	w.Write([]byte(albumsJSONStr))

	// Log the response details
	log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, http.StatusOK)
}
