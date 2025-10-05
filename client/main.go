package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/rohanmathur91/tunnel/dto"
)

/*

Client
- [DONE] Connect to websocket server
- [DONE] Test it
- [DONE] Get unique connection tunnelId (later a subdomain)
- Try forwarding requests and print response

Server
- Maintain client to tunnelId map in memory in server
	- handle duplicate connections
- Test all

Client
- Prepare client to run in CLI, add ability to do something like
	tunnel --port 3000 (this should give a server URL)

*/

func main() {
	port := flag.Int("port", 3000, "local port")
	flag.Parse()

	log.Println("Client port", *port)

	const server = "ws://localhost:8000/tunnel"
	connection, _, err := websocket.DefaultDialer.Dial(server, nil)

	if err != nil {
		log.Println("Client cannot be connected to websocket server!", err)
		return
	}

	defer connection.Close()

	var info dto.ClientTunnelInfo
	err = connection.ReadJSON(&info)
	if err != nil {
		log.Println("Client cannot read connection info! ", err)
		return
	}

	fmt.Println("-----------------------------------------")
	fmt.Printf("Tunnel ID:  %s\n", info.Id)
	fmt.Printf("Public URL: %s\n", info.Url)
	fmt.Printf("Forwarding: http://localhost:%d\n", *port)
	fmt.Println("-----------------------------------------")

	for {
		var request dto.Request
		err := connection.ReadJSON(&request)
		if err != nil {
			log.Fatal("Cannot read message from echo ", err)
			return
		}

		prettyJSON, _ := json.MarshalIndent(request, "", "  ")
		fmt.Printf("Incomming request:\n%s\n", string(prettyJSON))

		// TODO: execute the request and send response
		response := dto.Response{
			RequestId: request.Id,
			Status:    200,
		}

		err = connection.WriteJSON(response)
		if err != nil {
			log.Fatal("Cannot write message from echo ", err)
			return
		}
	}
}
