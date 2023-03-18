package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/s0up4200/metacritic-api/music"
)

func main() {
	// Start the cache updater goroutines
	go music.StartCacheUpdater()
	go music.StartNewCacheUpdater()

	// Register the handler function for the "/metacritic/upcoming-albums" endpoint
	http.HandleFunc("/metacritic/upcoming-albums", music.HandleAlbumsRequest)

	// Register the handler function for the "/metacritic/new-albums" endpoint
	http.HandleFunc("/metacritic/new-albums", music.HandleNewAlbumsRequest)

	// Start the HTTP server
	fmt.Println("Server listening on port 45323...")
	log.Println("Server started successfully.")
	log.Fatal(http.ListenAndServe(":45323", nil))
}
