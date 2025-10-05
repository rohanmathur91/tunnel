package dto

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type ClientTunnelInfo struct {
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
	RequestId string              `json:"request_id"`
	Header    map[string][]string `json:"header"`
	Body      []byte              `json:"body"`
	Status    int                 `json:"status"`
}

func ToRequest(r *http.Request) *Request {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Fatal("Cannot read the request body", string(body))
		return nil
	}

	return &Request{
		Id:     uuid.NewString(),
		Method: r.Method,
		Header: r.Header,
		Path:   r.URL.Path,
		Query:  r.URL.RawQuery,
		Body:   body,
	}
}

func ToJSONRequest(r *http.Request) ([]byte, *Request) {
	request := ToRequest(r)
	jsonBytes, err := json.Marshal(request)

	if err != nil {
		log.Fatal("Cannot parse request into json", err)
		return nil, request
	}

	return jsonBytes, request
}
