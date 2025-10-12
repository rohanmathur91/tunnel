package dto

import (
	"io"
	"log"
	"net/http"

	"github.com/rohanmathur91/tunnel/utils"
)

type TunnelInfo struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type Request struct {
	Id     string              `json:"id"`
	Method string              `json:"method"`
	Header map[string][]string `json:"header"`
	Path   string              `json:"path"`
	Query  string              `json:"query"`
	Body   []byte              `json:"body"`
}

type Response struct {
	RequestId string              `json:"requestId"`
	Header    map[string][]string `json:"header"`
	Body      []byte              `json:"body"`
	Status    int                 `json:"status"`
}

func CreateRequest(r *http.Request) *Request {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Fatal("Cannot read the request body", string(body))
		return nil
	}

	return &Request{
		Id:     utils.GenerateID(),
		Method: r.Method,
		Header: r.Header,
		Path:   r.URL.Path,
		Query:  r.URL.RawQuery,
		Body:   body,
	}
}

func CreateResponse(requestId string, response *http.Response) *Response {
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		log.Fatal("Error while reading reponse body ", err)
		return nil
	}

	return &Response{
		RequestId: requestId,
		Header:    response.Header,
		Status:    response.StatusCode,
		Body:      body,
	}
}
