package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rohanmathur91/tunnel/server"
)

const port = 8000

func main() {
	config := server.LoadConfig()
	wsServer := server.New(&config)

	http.HandleFunc("/", wsServer.HandleHttp) // http
	http.HandleFunc("/tunnel", wsServer.HandleNewConnection)
	http.HandleFunc("/echo", wsServer.HandleEcho)
	http.HandleFunc("/health", wsServer.HandleHealthCheck)

	httpServer := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	fmt.Printf("Server running on port %d... \n", port)

	err := httpServer.ListenAndServe()

	if err != nil {
		log.Fatal(err)
		panic("Could not start the server!")
	}
}
