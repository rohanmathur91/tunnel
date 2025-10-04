package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all
	},
}

func (s *Server) HandleEcho(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal("Cannot create new connection", err)
		return
	}

	defer connection.Close()

	log.Println("Client connected, new connection created!")

	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			log.Fatal("Cannot read message from echo", err)
			return
		}

		log.Println("echo::", messageType)
		err = connection.WriteMessage(messageType, message)
		if err != nil {
			log.Fatal("Cannot write message from echo", err)
			return
		}
	}

}

func (s *Server) HandleNewConnection(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal("Cannot create new connection", err)
		return
	}

	defer connection.Close()

	log.Println("Client connected, new connection created!")
}
