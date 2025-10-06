package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	port := 8080

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/data", dataHandler)
	http.HandleFunc("/api/users", usersHandler)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	fmt.Println("Available endpoints:")
	fmt.Println("  - http://localhost:8080/")
	fmt.Println("  - http://localhost:8080/api/data")
	fmt.Println("  - http://localhost:8080/api/users")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
		<h1>Welcome to Go Server</h1>
		<p>Available endpoints:</p>
		<ul>
			<li><a href="/api/data">/api/data</a> - Get dummy data</li>
			<li><a href="/api/users">/api/users</a> - Get user list</li>
			<li><a href="/health">/health</a> - Health check</li>
		</ul>
	`
	fmt.Fprint(w, html)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Message:   "This is dummy data from the server",
		Timestamp: time.Now(),
		Status:    "success",
	}

	json.NewEncoder(w).Encode(response)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := []User{
		{ID: 1, Name: "Alice Johnson", Email: "alice@example.com"},
		{ID: 2, Name: "Bob Smith", Email: "bob@example.com"},
		{ID: 3, Name: "Charlie Brown", Email: "charlie@example.com"},
	}

	json.NewEncoder(w).Encode(users)
}
