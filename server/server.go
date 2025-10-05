package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rohanmathur91/tunnel/dto"
	"github.com/rohanmathur91/tunnel/utils"
)

type ResponseChannels map[string]chan dto.Response

type Tunnel struct {
	connection       *websocket.Conn
	responseChannels ResponseChannels
}

type Server struct {
	mutex   sync.RWMutex
	tunnels map[string]*Tunnel
}

func New() *Server {
	return &Server{
		tunnels: make(map[string]*Tunnel),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all
	},
}

func (s *Server) HandleHttp(w http.ResponseWriter, r *http.Request) {
	tunnelId := r.URL.Query().Get("tunnelId")
	if len(tunnelId) == 0 {
		fmt.Println("Invalid tunnel id", tunnelId)
		http.Error(w, "Invalid tunnel id", http.StatusNotFound)
		return
	}

	s.mutex.RLock()
	tunnelDetails, exists := s.tunnels[tunnelId]
	s.mutex.RUnlock()

	if !exists {
		fmt.Println("Invalid tunnel URL", tunnelId)
		http.Error(w, "Forward the port first", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Tunnel id from http handler %v\n", tunnelId)

	request := dto.ToRequest(r)
	if request == nil {
		fmt.Println("Something went wrong!", request)
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	fmt.Println("Forwarding request...")

	err := tunnelDetails.connection.WriteJSON(request)
	if err != nil {
		log.Fatal("Failed to forward request", err)
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}

	responseChannel := make(chan dto.Response, 1)

	s.mutex.Lock()
	tunnelDetails.responseChannels[request.Id] = responseChannel
	s.mutex.Unlock()

	defer func() {
		// cleanup
		s.mutex.Lock()
		close(responseChannel)
		delete(tunnelDetails.responseChannels, request.Id)
		s.mutex.Unlock()
	}()

	fmt.Printf("Tunnel details: %+v\n", tunnelDetails)
	fmt.Println("Waiting for response...")
	response := <-responseChannel

	fmt.Printf("Response received from tunnel %+v\n", response)
	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(response.Status)
	w.Write(response.Body)

	prettyJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("Response: \n%s\n", string(prettyJSON))
	fmt.Fprintln(w, string(prettyJSON))
}

func (s *Server) HandleNewConnection(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Cannot create new connection", err)
		http.Error(w, "Failed to upgrade", http.StatusInternalServerError)
		return
	}

	defer connection.Close()

	// init connection/tunnel
	tunnelId := utils.GenerateID()
	tunnel := &Tunnel{
		connection:       connection,
		responseChannels: make(ResponseChannels),
	}

	s.mutex.Lock()
	s.tunnels[tunnelId] = tunnel
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		delete(s.tunnels, tunnelId)
		s.mutex.Unlock()
		log.Printf("Tunnel %s disconnected", tunnelId)
	}()

	fmt.Println("Client connected, new connection created!", tunnelId)

	tunnelInfo := dto.ClientTunnelInfo{
		Id:  tunnelId,
		Url: fmt.Sprintf("http://localhost:8000/?tunnelId=%s", tunnelId),
	}

	fmt.Printf("TunnelInfo %+v\n", tunnelInfo)

	err = connection.WriteJSON(tunnelInfo)
	if err != nil {
		fmt.Println("Cannot send message from", err)
		return
	}

	for {
		var response dto.Response
		err := connection.ReadJSON(&response)
		if err != nil {
			fmt.Println("Cannot read message from client ", err)
		}

		responseChannels := s.tunnels[tunnelId].responseChannels
		pendingResponseChannel := responseChannels[response.RequestId]
		pendingResponseChannel <- response

		prettyJSON, _ := json.MarshalIndent(response, "", "  ")
		fmt.Printf("Incomming response: %+v\n", string(prettyJSON))
	}
}

func (s *Server) HandleEcho(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Cannot create new connection", err)
		return
	}

	defer connection.Close()

	fmt.Println("Echo: new connection created!")

	for {
		var message any
		err := connection.ReadJSON(&message)
		if err != nil {
			fmt.Println("Cannot read message from echo ", err)
			return
		}

		err = connection.WriteJSON(message)
		if err != nil {
			fmt.Println("Cannot write message from echo ", err)
			return
		}
	}

}

func (s *Server) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	res := map[string]string{
		"status": "healthy",
	}

	utils.SendJSONResponse(w, res)
}
