package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rohanmathur91/tunnel/dto"
	"github.com/rohanmathur91/tunnel/utils"
)

type ResponseChannels map[string]chan dto.Response

type Tunnel struct {
	mutex            sync.RWMutex
	connection       *websocket.Conn
	responseChannels ResponseChannels
}

type Server struct {
	mutex   sync.RWMutex
	config  *Config
	tunnels map[string]*Tunnel
}

func New(config *Config) *Server {
	return &Server{
		tunnels: make(map[string]*Tunnel),
		config:  config,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all
	},
}

func (s *Server) extractTunnelID(r *http.Request) string {
	host := r.Host

	portIndex := strings.Index(host, ":")
	if portIndex != -1 {
		host = host[:portIndex]
	}

	// check subdomain pattern
	if !strings.HasSuffix(host, fmt.Sprintf(".%s", s.config.Domain)) {
		return ""
	}

	tunnelId := strings.TrimSuffix(host, fmt.Sprintf(".%s", s.config.Domain))
	return tunnelId
}

func (s *Server) HandleHttp(w http.ResponseWriter, r *http.Request) {
	tunnelId := s.extractTunnelID(r)
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
		http.Error(w, "Invalid tunnel", http.StatusInternalServerError)
		return
	}

	request := dto.CreateRequest(r)
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

	tunnelDetails.mutex.Lock()
	tunnelDetails.responseChannels[request.Id] = responseChannel
	tunnelDetails.mutex.Unlock()

	defer func() {
		tunnelDetails.mutex.Lock()
		close(responseChannel)
		delete(tunnelDetails.responseChannels, request.Id)
		tunnelDetails.mutex.Unlock()
	}()

	fmt.Printf("Tunnel details: %+v\n", tunnelDetails)
	fmt.Println("Waiting for response...")
	response := <-responseChannel // blocks here

	fmt.Println("Response received from tunnel!")
	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Add("x-tunnel-id", tunnelId)
	w.WriteHeader(response.Status)
	w.Write(response.Body)
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

	fmt.Println("New tunnel created with id:", tunnelId)

	tunnelInfo := dto.TunnelInfo{
		Id:  tunnelId,
		Url: fmt.Sprintf("http://%s.%s:%d", tunnelId, s.config.Domain, s.config.Port),
	}

	fmt.Println("-----------------------------------------")
	fmt.Printf("Tunnel ID:  %s\n", tunnelInfo.Id)
	fmt.Printf("Public URL: %s\n", tunnelInfo.Url)
	fmt.Println("-----------------------------------------")

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

	utils.SendJSONResponse(w, http.StatusOK, res)
}
