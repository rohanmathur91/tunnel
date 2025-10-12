package dto

import (
	"io"
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

func CreateRequest(r *http.Request) (Request, error) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return Request{}, err
	}

	return Request{
		Id:     utils.GenerateID(),
		Method: r.Method,
		Header: r.Header,
		Path:   r.URL.Path,
		Query:  r.URL.RawQuery,
		Body:   body,
	}, nil
}

func CreateResponse(requestId string, response *http.Response) (Response, error) {
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return Response{}, err
	}

	return Response{
		RequestId: requestId,
		Header:    response.Header,
		Status:    response.StatusCode,
		Body:      body,
	}, nil
}
