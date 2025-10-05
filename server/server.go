package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rohanmathur91/tunnel/dto"
)

type Tunnel struct {
	connection       *websocket.Conn
	responseChannels map[string]chan dto.Response
}

type Server struct {
	mutex   sync.RWMutex
	tunnels map[string]Tunnel
}

func New() *Server {
	return &Server{}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all
	},
}

func (s *Server) HandleHttp(w http.ResponseWriter, r *http.Request) {
	tunnelId := r.URL.Query().Get("tunnelId")
	if len(tunnelId) == 0 {
		log.Fatal("Invalid tunnel id", tunnelId)
		http.Error(w, "Invalid tunnel id", http.StatusNotFound)
	}

	tunnelDetails, exists := s.tunnels[tunnelId]

	if !exists {
		log.Fatal("Invalid tunnel URL", tunnelId)
		http.Error(w, "Forward the port first", http.StatusInternalServerError)

	}

	request, rawRequest := dto.ToJSONRequest(r)

	if request == nil {
		log.Fatal("Something went wrong!", rawRequest)
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
	}

	err := tunnelDetails.connection.WriteJSON(request)
	if err != nil {
		log.Fatal("Failed to forward request", err)
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
	}

	responseChannel := make(chan dto.Response, 1)

	s.mutex.Lock()
	tunnelDetails.responseChannels[rawRequest.Id] = responseChannel
	s.mutex.Unlock()

	defer func() {
		// cleanup
		s.mutex.Lock()
		close(responseChannel)
		delete(tunnelDetails.responseChannels, rawRequest.Id)
		s.mutex.Unlock()
	}()

	log.Print("Waiting for response...")
	response := <-responseChannel

	log.Print("response received from tunnel", response)
	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(response.Status)
	w.Write(response.Body)
}

func (s *Server) HandleNewConnection(w http.ResponseWriter, r *http.Request) {
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
		log.Println("message::", string(message))
	}
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

func (s *Server) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	res := map[string]string{
		"status": "healthy",
	}

	SendJSONResponse(w, res)
}
