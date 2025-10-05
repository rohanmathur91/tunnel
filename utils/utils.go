package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func SendJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		fmt.Println("Error in SendJSONResponse", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GenerateID() string {
	return uuid.NewString()
}
