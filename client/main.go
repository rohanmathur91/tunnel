package main

import (
	"flag"
	"fmt"

	"github.com/gorilla/websocket"
)

/*

Client
- [DONE] Connect to websocket server
- [DONE] Test it
- Get unique connection tunnelId (later a subdomain)
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

	fmt.Println("Client port", *port)

	const server = "ws://localhost:8000/tunnel"
	connection, _, err := websocket.DefaultDialer.Dial(server, nil)

	if err != nil {
		fmt.Println("Client cannot be connected to websocket server!", err)
		return
	}

	defer connection.Close()

	var connectionInfo map[string]string
	err = connection.ReadJSON(&connectionInfo)
	if err != nil {
		fmt.Println("Client cannot read connection info!", err)
		return
	}

	// Talk to websocket
	// hit localhost
	// send response back to websocket
}
