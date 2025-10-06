package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/rohanmathur91/tunnel/dto"
)

type Client struct {
	port int
}

func NewClient(port int) *Client {
	return &Client{
		port: port,
	}
}

const baseUrl = "http://localhost"
const serverUrl = "ws://localhost:8000/tunnel"

func (c *Client) Start() {
	port := c.port
	log.Println("Client port", port)

	connection, _, err := websocket.DefaultDialer.Dial(serverUrl, nil)

	if err != nil {
		log.Println("Client cannot be connected to websocket server!", err)
		return
	}

	defer connection.Close()

	var tunnelInfo dto.ClientTunnelInfo
	err = connection.ReadJSON(&tunnelInfo)
	if err != nil {
		log.Println("Client cannot read connection tunnelInfo! ", err)
		return
	}

	fmt.Println("-----------------------------------------")
	fmt.Printf("Tunnel ID:  %s\n", tunnelInfo.Id)
	fmt.Printf("Public URL: %s\n", tunnelInfo.Url)
	fmt.Printf("Forwarding: %s:%d\n", baseUrl, port)
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

		// this is blocking
		c.handleRequest(connection, tunnelInfo, request)
	}
}

func (c *Client) prepareUrl(tunnelInfo dto.ClientTunnelInfo, request dto.Request) string {
	baseURL, err := url.Parse(fmt.Sprintf("%s:%d", baseUrl, c.port))
	if err != nil {
		fmt.Println("Invalid base URL: %w", err)
		return ""
	}

	baseURL.Path = request.Path
	query := baseURL.Query()
	if request.Query != "" {
		existingParams, err := url.ParseQuery(request.Query)
		if err != nil {
			fmt.Println("Invalid query string: %w", err)
			return ""
		}

		for key, values := range existingParams {
			for _, value := range values {
				query.Add(key, value)
			}
		}
	}

	query.Set("tunnelId", tunnelInfo.Id)
	baseURL.RawQuery = query.Encode()
	localURL := baseURL.String()
	log.Printf("Forwarding to: %s", localURL)
	return localURL
}

func (c *Client) sendRequest(conn *websocket.Conn, tunnelInfo dto.ClientTunnelInfo, request dto.Request) (*http.Response, error) {
	localURL := c.prepareUrl(tunnelInfo, request)
	httpReq, err := http.NewRequest(request.Method, localURL, bytes.NewReader(request.Body))

	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.sendErrorResponseToTunnel(conn, request.Id, http.StatusInternalServerError, err.Error())
		return nil, err
	}

	// Copy headers
	for key, values := range request.Header {
		for _, value := range values {
			httpReq.Header.Add(key, value)
		}
	}

	httpReq.Header.Add("x-tunnel-id", tunnelInfo.Id)
	httpClient := &http.Client{}
	return httpClient.Do(httpReq)
}

func (c *Client) sendErrorResponseToTunnel(conn *websocket.Conn, requestId string, status int, message string) error {
	res := &dto.Response{
		RequestId: requestId,
		Status:    status,
		Body:      []byte(message),
	}

	err := conn.WriteJSON(res)
	return err
}

func (c *Client) sendResponseToTunnel(conn *websocket.Conn, requestId string, rawResponse http.Response) error {

	body, err := io.ReadAll(rawResponse.Body)
	if err != nil {
		fmt.Println("Error while reading body ", err)
		return err
	}

	defer rawResponse.Body.Close()

	res := &dto.Response{
		RequestId: requestId,
		Header:    rawResponse.Header,
		Status:    rawResponse.StatusCode,
		Body:      body,
	}

	err = conn.WriteJSON(res)
	return err
}

func (c *Client) handleRequest(conn *websocket.Conn, tunnelInfo dto.ClientTunnelInfo, request dto.Request) {
	response, err := c.sendRequest(conn, tunnelInfo, request)

	if err != nil {
		c.sendErrorResponseToTunnel(conn, request.Id, http.StatusInternalServerError, err.Error())
	} else {
		c.sendResponseToTunnel(conn, request.Id, *response)
	}
}
