package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Album struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

func main() {
	// Define the URL to scrape
	url := "https://www.metacritic.com/browse/albums/release-date/coming-soon"

	// Send a GET request to the specified URL
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		panic(err)
	}

	// Find all album rows and store the artist and title in a slice of structs
	var albums []Album
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		artist := strings.TrimSpace(s.Find(".artistName a").Text())
		if artist == "" {
			artist = strings.TrimSpace(s.Find(".artistName").Text())
		}
		title := strings.TrimSpace(s.Find(".albumTitle").Text())
		if artist != "" && title != "" && title != "[Title TBA]" {
			album := Album{artist, title}
			albums = append(albums, album)
		}
	})

	// Convert the slice of structs to a JSON string with single quotes
	albumsJSON, err := json.Marshal(albums)
	if err != nil {
		panic(err)
	}
	albumsJSONStr := string(albumsJSON)
	albumsJSONStr = strings.Replace(albumsJSONStr, "\\u0026", "&", -1) // Convert HTML escape sequence to "&"
	albumsJSONStr = strings.ReplaceAll(albumsJSONStr, "\"", "'")

	// Print the JSON string
	fmt.Println(albumsJSONStr)
}
