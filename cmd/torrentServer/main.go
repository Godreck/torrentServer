package main

import (
	"log"
	"net/http"
	"torrentServer/http_server/handlers/search"
)

func main() {
	http.HandleFunc("/search", search.SearchHandler)
	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
