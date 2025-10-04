package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendJSONResponse(w http.ResponseWriter, data any) {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		fmt.Println("Error in SendJSONResponse", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
