package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize database connection and run migrations
	if err := InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/courses", getCoursesHandler)

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "IronSnake API is running!")
}
