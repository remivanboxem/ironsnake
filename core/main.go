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

	// Load configuration
	config := LoadConfig()
	log.Println("Configuration loaded successfully")

	// Initialize authentication services
	InitAuthServices(config)
	log.Println("Authentication services initialized")

	// Public routes
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/auth/login", loginHandler)
	http.HandleFunc("/auth/logout", logoutHandler)

	// Protected routes (require authentication)
	http.HandleFunc("/auth/me", AuthMiddleware(getMeHandler))
	http.HandleFunc("/courses", AuthMiddleware(getCoursesHandler))
	http.HandleFunc("/run", AuthMiddleware(runCodeHandler))

	// Task routes need to be registered before course routes due to path matching
	http.HandleFunc("/courses/", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a task request
		if len(r.URL.Path) > len("/courses/") {
			remaining := r.URL.Path[len("/courses/"):]
			for i := 0; i < len(remaining)-6; i++ {
				if remaining[i:i+7] == "/tasks/" {
					getTaskByIDHandler(w, r)
					return
				}
			}
		}
		// Otherwise, it's a course request
		getCourseByIDHandler(w, r)
	}))

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "IronSnake API is running!")
}
