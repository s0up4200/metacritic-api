package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mutex sync.Mutex

func Main() {
	// Start the cache updater goroutines
	go StartCacheUpdater()
	go StartNewCacheUpdater()

	// Register the handler function for the "/metacritic/upcoming-albums" endpoint
	http.HandleFunc("/metacritic/upcoming-albums", HandleAlbumsRequest)

	// Register the handler function for the "/metacritic/new-albums" endpoint
	http.HandleFunc("/metacritic/new-albums", HandleNewAlbumsRequest)

	// Start the HTTP server
	fmt.Println("Server listening on port 45323...")
	log.Println("Server started successfully.")
	log.Fatal(http.ListenAndServe(":45323", nil))
}
