package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rohanmathur91/tunnel/server"
)

func main() {
	config := server.LoadConfig()
	wsServer := server.New(&config)

	http.HandleFunc("/", wsServer.HandleHttp) // http
	http.HandleFunc("/tunnel", wsServer.HandleNewConnection)
	http.HandleFunc("/echo", wsServer.HandleEcho)
	http.HandleFunc("/health", wsServer.HandleHealthCheck)

	httpServer := http.Server{
		Addr: fmt.Sprintf(":%d", config.Port),
	}

	fmt.Printf("Server running on port %d... \n", config.Port)

	err := httpServer.ListenAndServe()

	if err != nil {
		log.Fatal(err)
		panic("Could not start the server!")
	}
}
